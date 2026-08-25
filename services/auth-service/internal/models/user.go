package models

import "time"

// User is the domain model for an authenticated user.
// PasswordHash is stored in the DB but omitted from JSON responses.
type User struct {
    ID           string    `db:"id" json:"id"`
    Email        string    `db:"email" json:"email"`
    PasswordHash string    `db:"password_hash" json:"-"`
    CreatedAt    time.Time `db:"created_at" json:"created_at"`
}
