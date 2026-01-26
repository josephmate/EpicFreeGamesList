#!/usr/bin/env bash
set -euo pipefail

# Usage: ./serve.sh [port] [dir]
# Defaults: port=8000, dir=.

PORT="${1:-8000}"
DIR="${2:-.}"
# URL prefix used by GitHub Pages for this repo
PREFIX="/EpicFreeGamesList"

# Find Python 3 (fallback to python if that's all that's available)
if command -v python3 >/dev/null 2>&1; then
  PY_EXEC=python3
elif command -v python >/dev/null 2>&1; then
  PY_EXEC=python
else
  echo "Python 3 is required but was not found on PATH. Install Python and retry." >&2
  exit 1
fi

echo "Serving '${DIR}' at http://localhost:${PORT} (handling prefix ${PREFIX})"
cd "$DIR"

# Try to open default browser (best-effort; will not fail the script)
if command -v xdg-open >/dev/null 2>&1; then
  xdg-open "http://localhost:${PORT}${PREFIX}/" >/dev/null 2>&1 || true
elif command -v start >/dev/null 2>&1; then
  start "http://localhost:${PORT}${PREFIX}/" || true
elif command -v open >/dev/null 2>&1; then
  open "http://localhost:${PORT}${PREFIX}/" || true
fi

# Run an embedded Python server that strips the GitHub Pages prefix before resolving files.
${PY_EXEC} - "${PORT}" "${PREFIX}" <<'PY'
import sys, os, http.server, socketserver

PORT = int(sys.argv[1])
PREFIX = sys.argv[2]

class PrefixHandler(http.server.SimpleHTTPRequestHandler):
    def translate_path(self, path):
        # If request starts with the repo prefix, strip it so files are served from the repo root
        if path.startswith(PREFIX + '/'):
            path = path[len(PREFIX):]
        elif path == PREFIX:
            path = '/'
        return super().translate_path(path)

Handler = PrefixHandler
os.chdir(os.getcwd())
print(f"Serving '{os.getcwd()}' on http://0.0.0.0:{PORT} (prefix={PREFIX})")
socketserver.TCPServer.allow_reuse_address = True
with socketserver.TCPServer(('', PORT), Handler) as httpd:
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        pass
PY
