#!/usr/bin/env bash
set -euo pipefail
cd "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if [[ ! -f nginx/ssl/fullchain.pem ]]; then
  chmod +x scripts/gen-ssl-selfsigned.sh 2>/dev/null || true
  ./scripts/gen-ssl-selfsigned.sh
fi

if docker compose version &>/dev/null; then
  DC=(docker compose)
elif docker-compose version &>/dev/null; then
  DC=(docker-compose)
else
  echo "Need Docker Compose: install plugin (docker compose) or docker-compose"
  exit 1
fi

"${DC[@]}" up -d --build
echo "Running: :80→:443 (nginx HTTPS), :9000 (tunnel). Stop: ${DC[*]} down"
