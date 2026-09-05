#!/usr/bin/env bash
set -euo pipefail

# Charites Rule Scaffolding Generator (1-SSOT Tri-Corpus Standard)
# Usage: ./scaffold_rule.sh <CATEGORY> <SLUG> [SEVERITY] [DESCRIPTION]
# Example: ./scaffold_rule.sh theme hardcode-color HIGH "Detects raw un-tokenized colors"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SKILL_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
ASSETS_DIR="${SKILL_ROOT}/assets"

if [ "$#" -lt 2 ]; then
    echo "Usage: $0 <CATEGORY> <SLUG> [SEVERITY] [DESCRIPTION]"
    echo "Example: $0 theme hardcode-color HIGH \"Detects hardcoded hex/rgb color literals\""
    exit 1
fi

CATEGORY="$1"
SLUG="$2"
SEVERITY="${3:-HIGH}"
DESCRIPTION="${4:-Detects violations of ${CATEGORY}.${SLUG}}"

RULE_ID="${CATEGORY}.${SLUG}"
SNAKE_SLUG="${SLUG//-/_}"
RULE_FILE_NAME="${SNAKE_SLUG}.go"

# Convert to PascalCase for Go struct name
to_pascal() {
    echo "$1" | sed -r 's/(^|[-_])([a-z])/\U\2/g'
}

CAT_PASCAL=$(to_pascal "$CATEGORY")
SLUG_PASCAL=$(to_pascal "$SLUG")
STRUCT_NAME="${CAT_PASCAL}${SLUG_PASCAL}Rule"
SEVERITY_TITLE=$(tr '[:lower:]' '[:upper:]' <<< "${SEVERITY:0:1}")$(tr '[:upper:]' '[:lower:]' <<< "${SEVERITY:1}")
# shellcheck disable=SC2001
TITLE_CASE_NAME=$(echo "${SLUG//-/ }" | sed -e 's/\b\(.\)/\u\1/g')

RULES_DIR="${REPO_ROOT}/internal/rules/${CATEGORY}"
CORPUS_DIR="${REPO_ROOT}/tests/correctness/${CATEGORY}/${SLUG}"
WIKI_FILE="${REPO_ROOT}/wiki/${CATEGORY}.${SLUG}.md"

echo "=== Scaffolding Charites Rule: ${RULE_ID} ==="
echo "Implementation: ${RULES_DIR}/${RULE_FILE_NAME}"
echo "Tri-Corpus:     ${CORPUS_DIR}/"
echo "Wiki Doc:       ${WIKI_FILE}"

mkdir -p "${RULES_DIR}" "${REPO_ROOT}/wiki"
mkdir -p "${CORPUS_DIR}/positive" "${CORPUS_DIR}/negative" "${CORPUS_DIR}/adversarial"

# 1. Generate internal/rules/<category>_<slug>.go
sed \
    -e "s|{{STRUCT_NAME}}|${STRUCT_NAME}|g" \
    -e "s|{{RULE_ID}}|${RULE_ID}|g" \
    -e "s|{{CATEGORY}}|${CATEGORY}|g" \
    -e "s|{{SEVERITY_TITLE}}|${SEVERITY_TITLE}|g" \
    -e "s|{{DESCRIPTION}}|${DESCRIPTION}|g" \
    "${ASSETS_DIR}/rule.go.tmpl" > "${RULES_DIR}/${RULE_FILE_NAME}"

# 2. Generate tests/correctness/<category>.<slug>/rule_test.go
sed \
    -e "s|{{STRUCT_NAME}}|${STRUCT_NAME}|g" \
    -e "s|{{RULE_ID}}|${RULE_ID}|g" \
    "${ASSETS_DIR}/rule_test.go.tmpl" > "${CORPUS_DIR}/rule_test.go"

# 3. Generate Tri-Corpus Fixture skeletons
cat <<EOF > "${CORPUS_DIR}/positive/direct.astro"
---
// P1: Direct inline violation
---
<!-- want "${RULE_ID}" -->
<div style="color: #ff0000;">Inline violation</div>
EOF

cat <<EOF > "${CORPUS_DIR}/positive/indirect.tsx"
// P2: Indirect / dynamic class concatenation
export function Card() {
  // want "${RULE_ID}"
  return <div className={"bg-[#ff0000]"}>Indirect violation</div>;
}
EOF

cat <<EOF > "${CORPUS_DIR}/negative/tokens.astro"
---
// N1: Valid design tokens
---
<div class="bg-primary text-primary-foreground">Safe token usage</div>
EOF

cat <<EOF > "${CORPUS_DIR}/negative/ignore.tsx"
// N2: Explicit ignore directive
export function Ignored() {
  // charites:ignore ${RULE_ID} intentional test
  return <div style={{ color: "#123456" }}>Exempted</div>;
}
EOF

cat <<EOF > "${CORPUS_DIR}/adversarial/template_interp.astro"
---
// A1: Template string interpolation stress case
const dynamic = "foo";
---
<div class={\`prefix-\${dynamic}\`}>Stress testing</div>
EOF

cat <<EOF > "${CORPUS_DIR}/adversarial/ternary.tsx"
// A2: Ternary expression evasion attempt
export function Evasion(props: { active: boolean }) {
  return <div className={props.active ? "text-[#000]" : "text-primary"}>Ternary test</div>;
}
EOF

# 4. Generate 8-Pillars Wiki Documentation
if [ ! -f "${WIKI_FILE}" ]; then
    sed \
        -e "s|{{RULE_ID}}|${RULE_ID}|g" \
        -e "s|{{TITLE_CASE_NAME}}|${TITLE_CASE_NAME}|g" \
        -e "s|{{SEVERITY}}|${SEVERITY}|g" \
        -e "s|{{CATEGORY}}|${CATEGORY}|g" \
        -e "s|{{DESCRIPTION}}|${DESCRIPTION}|g" \
        "${ASSETS_DIR}/wiki_rule.md.tmpl" > "${WIKI_FILE}"
fi

echo "=== Scaffolding Complete! ==="
echo ""
echo "Next Steps Checklist:"
echo "1. Register in internal/rules/builtin.go:"
echo "     reg.Register(${CATEGORY}.New${STRUCT_NAME}())"
echo "2. Implement pattern inspection in internal/rules/${CATEGORY}/${RULE_FILE_NAME}"
echo "3. Run golden corpus test:"
echo "     go test -v ./tests/correctness/${CATEGORY}/${SLUG}/..."
echo "4. Regenerate wiki docs automatically:"
echo "     make wiki"

