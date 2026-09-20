#!/bin/sh
set -eu

: "${DEPLOY_PATH:?DEPLOY_PATH is required}"
: "${IMAGE_PREFIX:?IMAGE_PREFIX is required}"
: "${IMAGE_TAG:?IMAGE_TAG is required}"

cd "$DEPLOY_PATH"

compose() {
    docker compose --env-file .env -f docker-compose.yml "$@"
}

compose pull
compose up -d --remove-orphans
compose ps

docker image prune -f
