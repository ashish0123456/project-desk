package models

import "time"

type Issue struct {
    ID        string    `db:"id" json:"id"`
    ProjectID string    `db:"project_id" json:"project_id"`
    Title     string    `db:"title" json:"title"`
    Body      string    `db:"body" json:"body"`
    Status    string    `db:"status" json:"status"`
    Priority  string    `db:"priority" json:"priority"`
    Assignee  string    `db:"assignee" json:"assignee"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
}
