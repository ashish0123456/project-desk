package models

import "time"

type Project struct {
    ID        string    `db:"id" json:"id"`
    Name      string    `db:"name" json:"name"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Member struct {
    ID        string    `db:"id" json:"id"`
    ProjectID string    `db:"project_id" json:"project_id"`
    UserID    string    `db:"user_id" json:"user_id"`
    Role      string    `db:"role" json:"role"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
}
