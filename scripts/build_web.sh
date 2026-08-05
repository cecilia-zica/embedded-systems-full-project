#!/usr/bin/env bash
# Builds the Flutter app for the web, injecting the backend URL and API key at
# build time (nothing hardcoded). Output lands in proj/app/build/web, ready to
# host on GitHub Pages / Netlify / Firebase Hosting.
#
# Usage:
#   API_BASE_URL=https://your-backend.fly.dev API_KEY=your-key ./scripts/build_web.sh
set -euo pipefail

API_BASE_URL="${API_BASE_URL:?set API_BASE_URL, e.g. https://your-backend.fly.dev}"
API_KEY="${API_KEY:?set API_KEY (must match the backend)}"

cd "$(dirname "$0")/../proj/app"

flutter pub get
flutter build web --release \
  --dart-define=API_BASE_URL="$API_BASE_URL" \
  --dart-define=API_KEY="$API_KEY"

echo "Done -> proj/app/build/web"
