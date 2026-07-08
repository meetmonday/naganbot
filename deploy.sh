#!/bin/bash
set -e

SERVER_USER=${SERVER_USER:-stas}
SERVER_HOST=${SERVER_HOST:-randomshit.icu}
SERVER_DIR=${SERVER_DIR:-/home/stas/naganbot}
SSH_KEY=${SSH_KEY:-$HOME/.ssh/id_ed25519}

# Запустить ssh-agent, если не запущен, и добавить ключ (passphrase спросят один раз)
if [ -z "$SSH_AUTH_SOCK" ]; then
  eval "$(ssh-agent -s)"
fi
ssh-add -l &>/dev/null || ssh-add "$SSH_KEY"
LOCAL_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "==> Deploying to $SERVER_USER@$SERVER_HOST:$SERVER_DIR"

rsync -avz \
    --exclude '.git' \
    --exclude 'data' \
    --exclude 'deploy.sh' \
    "$LOCAL_DIR/" \
    "$SERVER_USER@$SERVER_HOST:$SERVER_DIR/"

echo "==> Restarting containers on server"
ssh -i "$SSH_KEY" "$SERVER_USER@$SERVER_HOST" \
    "cd $SERVER_DIR && \
    DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 \
    docker compose -f docker-compose.yaml -f docker-compose.dev.yaml -f docker-compose.build.yaml up -d --build --force-recreate && \
    docker compose logs -f app"
