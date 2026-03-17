#!/usr/bin/env bash
set -euo pipefail

TAGS_FILE="docker/.tags"
IMAGE="mikespook/shopify-admin-mcp-server"

# Auto-increment patch version unless VERSION is overridden
if [ -n "${VERSION:-}" ]; then
  FULL_VERSION="$VERSION"
else
  if [ -f "$TAGS_FILE" ]; then
    PREV=$(cat "$TAGS_FILE")
    MAJOR=$(echo "$PREV" | cut -d. -f1 | tr -d 'v')
    MINOR=$(echo "$PREV" | cut -d. -f2)
    PATCH=$(echo "$PREV" | cut -d. -f3)
    FULL_VERSION="v${MAJOR}.${MINOR}.$((PATCH + 1))"
  else
    FULL_VERSION="v0.1.0"
  fi
fi

MINOR_VERSION=$(echo "$FULL_VERSION" | cut -d. -f1-2)

echo "Building $IMAGE with version $FULL_VERSION"
echo "$FULL_VERSION" > "$TAGS_FILE"

docker build \
  -f docker/Dockerfile \
  -t "${IMAGE}:latest" \
  -t "${IMAGE}:${MINOR_VERSION}" \
  -t "${IMAGE}:${FULL_VERSION}" \
  .

echo "Built tags: latest, ${MINOR_VERSION}, ${FULL_VERSION}"
