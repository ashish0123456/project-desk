#!/usr/bin/env bash
set -euo pipefail

SERVICE="${1:-}"
ACTION="${2:-up}"

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

if ! command -v migrate >/dev/null 2>&1; then
  echo "golang-migrate is not installed."
  echo "Install: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
  exit 1
fi

case "$SERVICE" in
  auth)
    DB_URL="postgres://postgres:postgres@localhost:5432/pm_auth?sslmode=disable"
    PATH_MIGRATIONS="services/auth-service/migrations"
    ;;
  project)
    DB_URL="postgres://postgres:postgres@localhost:5433/pm_project?sslmode=disable"
    PATH_MIGRATIONS="services/project-service/migrations"
    ;;
  issue)
    DB_URL="postgres://postgres:postgres@localhost:5434/pm_issue?sslmode=disable"
    PATH_MIGRATIONS="services/issue-service/migrations"
    ;;
  *)
    echo "Usage: $0 <auth|project|issue> [up|down|version]"
    exit 1
    ;;
esac

migrate -path "$PATH_MIGRATIONS" -database "$DB_URL" "$ACTION"
