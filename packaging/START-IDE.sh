#!/usr/bin/env sh
set -e
cd "$(dirname "$0")"
chmod +x moonbasic-ide moonbasic moonrun 2>/dev/null || true
if [ ! -f ./moonbasic-ide ]; then
  echo "moonbasic-ide not found in this folder." >&2
  exit 1
fi
exec ./moonbasic-ide
