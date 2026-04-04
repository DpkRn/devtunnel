#!/usr/bin/env bash
set -euo pipefail
cd "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
docker compose up -d --build
echo "Running: port 80 (nginx → HTTP), 9000 (tunnel). Stop: docker compose down"
