package migrate

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ApplyMigrations applies any *.up.sql files in directory order
// and records applied migrations in a simple schema_migrations table.
func ApplyMigrations(db *sql.DB, migrationsDir string) error {
    if err := ensureTable(db); err != nil {
        return err
    }

    entries, err := os.ReadDir(migrationsDir)
    if err != nil {
        // if directory doesn't exist, nothing to do
        if os.IsNotExist(err) {
            return nil
        }
        return fmt.Errorf("read migrations dir: %w", err)
    }

    var files []string
    for _, e := range entries {
        if e.IsDir() {
            continue
        }
        name := e.Name()
        if strings.HasSuffix(name, ".up.sql") {
            files = append(files, filepath.Join(migrationsDir, name))
        }
    }

    sort.Strings(files)

    for _, f := range files {
        applied, err := isApplied(db, f)
        if err != nil {
            return err
        }
        if applied {
            continue
        }
        content, err := os.ReadFile(f)
        if err != nil {
            return fmt.Errorf("read migration %s: %w", f, err)
        }
        if err := execSQL(db, string(content)); err != nil {
            return fmt.Errorf("exec migration %s: %w", f, err)
        }
        if err := record(db, f); err != nil {
            return fmt.Errorf("record migration %s: %w", f, err)
        }
    }

    return nil
}

func ensureTable(db *sql.DB) error {
    _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (name TEXT PRIMARY KEY, applied_at TIMESTAMP WITH TIME ZONE DEFAULT now())`)
    return err
}

func isApplied(db *sql.DB, path string) (bool, error) {
    var name string
    row := db.QueryRow(`SELECT name FROM schema_migrations WHERE name = $1`, path)
    if err := row.Scan(&name); err != nil {
        if err == sql.ErrNoRows {
            return false, nil
        }
        return false, err
    }
    return true, nil
}

func execSQL(db *sql.DB, sqlContent string) error {
    // split on ; statements conservatively. Many migrations are single statements.
    statements := splitSQL(sqlContent)
    for _, s := range statements {
        s = strings.TrimSpace(s)
        if s == "" {
            continue
        }
        if _, err := db.Exec(s); err != nil {
            return err
        }
    }
    return nil
}

func record(db *sql.DB, path string) error {
    _, err := db.Exec(`INSERT INTO schema_migrations (name) VALUES ($1)`, path)
    return err
}

func splitSQL(s string) []string {
    // naive split by semicolon that ignores file system newlines.
    // Good enough for simple migration files used here.
    var out []string
    parts := strings.Split(s, ";")
    for _, p := range parts {
        if strings.TrimSpace(p) == "" {
            continue
        }
        out = append(out, p)
    }
    return out
}
