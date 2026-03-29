#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BUF_IMAGE="bufbuild/buf:1.58.0"

cmd="${1:-lint}"
shift || true

case "$cmd" in
  lint)
    docker run --rm -v "$ROOT:/workspace" -w /workspace "$BUF_IMAGE" lint "$@"
    ;;
  build)
    docker run --rm -v "$ROOT:/workspace" -w /workspace "$BUF_IMAGE" build "$@"
    ;;
  generate)
    docker run --rm -v "$ROOT:/workspace" -w /workspace "$BUF_IMAGE" generate "$@"
    ;;
  format)
    docker run --rm -v "$ROOT:/workspace" -w /workspace "$BUF_IMAGE" format -w "$@"
    ;;
  *)
    echo "usage: $0 {lint|build|generate|format} [args...]" >&2
    exit 1
    ;;
esac
