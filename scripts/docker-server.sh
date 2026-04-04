#!/usr/bin/env bash
# Build and run the tunnel server (mytunneld) in Docker on EC2 or any Linux host.
# Requires: docker, repo cloned (run from anywhere — script cds to repo root).
#
# Usage:
#   ./scripts/docker-server.sh
#   GO_VERSION=1.25 ./scripts/docker-server.sh
#
# Open EC2 security group: TCP 3000 (HTTP edge), TCP 9000 (tunnel control).

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

IMAGE_NAME="${IMAGE_NAME:-devtunnel-server}"
CONTAINER_NAME="${CONTAINER_NAME:-devtunnel-server}"

echo "==> Stopping/removing existing container (if any): ${CONTAINER_NAME}"
docker stop "${CONTAINER_NAME}" 2>/dev/null || true
docker rm "${CONTAINER_NAME}" 2>/dev/null || true

echo "==> Removing existing image (if any): ${IMAGE_NAME}"
docker rmi "${IMAGE_NAME}" 2>/dev/null || true

BUILD_ARGS=()
if [[ -n "${GO_VERSION:-}" ]]; then
  BUILD_ARGS+=(--build-arg "GO_VERSION=${GO_VERSION}")
fi

echo "==> Building image: ${IMAGE_NAME}"
docker build "${BUILD_ARGS[@]}" -t "${IMAGE_NAME}" .

echo "==> Starting container: ${CONTAINER_NAME}"
docker run -d \
  --name "${CONTAINER_NAME}" \
  --restart unless-stopped \
  -p 3000:3000 \
  -p 9000:9000 \
  "${IMAGE_NAME}"

echo ""
echo "Done. Ports: 3000 (HTTP edge), 9000 (control)."
echo "  docker logs -f ${CONTAINER_NAME}"
echo "  docker stop ${CONTAINER_NAME}"
