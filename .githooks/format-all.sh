#!/usr/bin/env bash
# format-all.sh
# Comprehensive code formatting & hygiene utility for Charites
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
cd "${REPO_DIR}"

echo "[Charites] Running formatting & hygiene pipeline..."

# Optional private developer hygiene hooks (runs if present locally)
if [ -f "${SCRIPT_DIR}/remove-emojis.sh" ] && [ -x "${SCRIPT_DIR}/remove-emojis.sh" ]; then
    echo "  [Hygiene] Stripping emojis..."
    "${SCRIPT_DIR}/remove-emojis.sh"
fi

if [ -f "${SCRIPT_DIR}/clean-ai-artifacts.sh" ] && [ -x "${SCRIPT_DIR}/clean-ai-artifacts.sh" ]; then
    echo "  [Hygiene] Cleaning AI artifacts & typography..."
    "${SCRIPT_DIR}/clean-ai-artifacts.sh"
fi

# 1. Format Go files
echo "  1/3 Formatting Go files (gofmt)..."
find . -name "*.go" -not -path "./vendor/*" -not -path "./.git/*" -exec gofmt -w {} + 2>/dev/null || true

# 2. Trim trailing whitespaces from markdown, yaml, and json
echo "  2/3 Cleaning trailing whitespaces (docs, yaml, json)..."
find . \( -name "*.md" -o -name "*.yaml" -o -name "*.yml" -o -name "*.json" \) \
    -not -path "./vendor/*" \
    -not -path "./.git/*" \
    -not -path "./node_modules/*" \
    -exec sed -i 's/[[:space:]]*$//' {} + 2>/dev/null || true

# 3. Format Web / Frontend templates if tools available
if command -v prettier >/dev/null 2>&1; then
    echo "  3/3 Formatting web files with prettier..."
    prettier --write "**/*.{astro,tsx,jsx,ts,js,css}" --ignore-path .charitesignore 2>/dev/null || true
elif command -v npx >/dev/null 2>&1 && [ -f "package.json" ]; then
    echo "  3/3 Formatting web files with npx prettier..."
    npx prettier --write "**/*.{astro,tsx,jsx,ts,js,css}" --ignore-path .charitesignore 2>/dev/null || true
else
    echo "  3/3 Skipping prettier (no npm/prettier in PATH or no web bundle yet)"
fi

echo "[Charites] All files formatted cleanly!"
