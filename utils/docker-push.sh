#!/usr/bin/env bash
set -euo pipefail

TAGS_FILE="docker/.tags"
IMAGE="mikespook/shopify-admin-mcp-server"

if [ ! -f "$TAGS_FILE" ]; then
  echo "Error: $TAGS_FILE not found. Run 'make docker-build-images' first." >&2
  exit 1
fi

FULL_VERSION=$(cat "$TAGS_FILE")
MINOR_VERSION=$(echo "$FULL_VERSION" | cut -d. -f1-2)

for TAG in latest "${MINOR_VERSION}" "${FULL_VERSION}"; do
  echo "Pushing ${IMAGE}:${TAG}..."
  docker push "${IMAGE}:${TAG}"
done

echo "All tags pushed."
