#!/bin/sh
# Charites Universal Installer for Linux and macOS
# Usage: curl -fsSL https://raw.githubusercontent.com/will2469/charites/main/scripts/install.sh | sh
# Options:
#   CHARITES_VERSION="v1.0.0"  Specify exact version to install
#   CHARITES_INSTALL_DIR="..." Specify custom installation directory

set -eu

REPO="will2469/charites"
GITHUB_URL="${GITHUB_URL:-https://github.com/${REPO}}"

# Colors for output
BOLD="\033[1m"
GREEN="\033[32m"
BLUE="\033[34m"
YELLOW="\033[33m"
RED="\033[31m"
NC="\033[0m"

info() {
    printf "%b==>%b %b%s%b\n" "${BLUE}" "${NC}" "${BOLD}" "$1" "${NC}"
}

success() {
    printf "%b==>%b %b%s%b\n" "${GREEN}" "${NC}" "${BOLD}" "$1" "${NC}"
}

warn() {
    printf "%bwarning:%b %s\n" "${YELLOW}" "${NC}" "$1"
}

error() {
    printf "%berror:%b %s\n" "${RED}" "${NC}" "$1" >&2
    exit 1
}

# 1. Detect Operating System
OS="$(uname -s)"
case "${OS}" in
    Linux*)   TARGET_OS="linux" ;;
    Darwin*)  TARGET_OS="darwin" ;;
    *)        error "Unsupported operating system: ${OS}. Charites supports Linux and macOS via this script." ;;
esac

# 2. Detect Architecture
ARCH="$(uname -m)"
case "${ARCH}" in
    x86_64|amd64)   TARGET_ARCH="amd64" ;;
    arm64|aarch64)  TARGET_ARCH="arm64" ;;
    *)              error "Unsupported architecture: ${ARCH}. Charites supports amd64 (x86_64) and arm64 (Apple Silicon / ARM64)." ;;
esac

if [ "${TARGET_OS}" = "darwin" ] && [ "${TARGET_ARCH}" = "amd64" ]; then
    error "macOS Intel (x86_64) prebuilt binaries are not published. Charites officially supports macOS Apple Silicon (arm64). You can build from source with: go install github.com/${REPO}/cmd/charites@latest"
fi

# 3. Determine Release Tag (User specified or Latest)
REQUESTED_TAG="${1:-${CHARITES_VERSION:-${VERSION:-}}}"
TAG="${REQUESTED_TAG}"

if [ -z "${TAG}" ]; then
    info "Resolving latest release tag from GitHub..."
    if command -v curl >/dev/null 2>&1; then
        TAG="$(curl -sSfL "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null | grep '"tag_name":' | head -n1 | sed -E 's/.*"tag_name":[[:space:]]*"([^"]+)".*/\1/' || true)"
    fi

    if [ -z "${TAG}" ] && command -v curl >/dev/null 2>&1; then
        LATEST_URL="$(curl -sSfLI -o /dev/null -w '%{url_effective}' "${GITHUB_URL}/releases/latest" 2>/dev/null || true)"
        TAG="${LATEST_URL##*/}"
    fi

    if [ -z "${TAG}" ] || [ "${TAG}" = "latest" ] || [ "${TAG}" = "releases" ]; then
        TAG="v1.0.0"
        warn "Could not resolve latest release tag via GitHub API, defaulting to ${TAG}"
    fi
fi

# 4. Prepare Package Metadata
ARCHIVE_NAME="charites_${TAG}_${TARGET_OS}_${TARGET_ARCH}.tar.gz"
DOWNLOAD_URL="${GITHUB_URL}/releases/download/${TAG}/${ARCHIVE_NAME}"
CHECKSUMS_URL="${GITHUB_URL}/releases/download/${TAG}/checksums.txt"

# Display Detection Summary
printf "\n%b==>%b %bCharites Universal Installer%b\n" "${BLUE}" "${NC}" "${BOLD}" "${NC}"
printf "  • Platform:        %s / %s (%s)\n" "${TARGET_OS}" "${TARGET_ARCH}" "${ARCH}"
printf "  • Release Tag:     %s\n" "${TAG}"
printf "  • Release Package: %s\n\n" "${ARCHIVE_NAME}"

# 5. Create Temporary Directory
TMP_DIR="$(mktemp -d)"
cleanup() {
    rm -rf "${TMP_DIR}"
}
trap cleanup EXIT INT TERM

# 6. Download Package
info "Downloading ${DOWNLOAD_URL}..."
if command -v curl >/dev/null 2>&1; then
    curl -fsSL "${DOWNLOAD_URL}" -o "${TMP_DIR}/${ARCHIVE_NAME}"
elif command -v wget >/dev/null 2>&1; then
    wget -q "${DOWNLOAD_URL}" -O "${TMP_DIR}/${ARCHIVE_NAME}"
else
    error "Neither curl nor wget is installed. Please install one of them to proceed."
fi

# 7. Verify Checksum (if checksums.txt is available)
info "Checking release signature & SHA256 integrity..."
CHECKSUM_FILE="${TMP_DIR}/checksums.txt"
if curl -fsSL "${CHECKSUMS_URL}" -o "${CHECKSUM_FILE}" 2>/dev/null || wget -q "${CHECKSUMS_URL}" -O "${CHECKSUM_FILE}" 2>/dev/null; then
    EXPECTED_SHA=$(grep "${ARCHIVE_NAME}" "${CHECKSUM_FILE}" | awk '{print $1}' || true)
    if [ -n "${EXPECTED_SHA}" ]; then
        if command -v sha256sum >/dev/null 2>&1; then
            ACTUAL_SHA=$(sha256sum "${TMP_DIR}/${ARCHIVE_NAME}" | awk '{print $1}')
        elif command -v shasum >/dev/null 2>&1; then
            ACTUAL_SHA=$(shasum -a 256 "${TMP_DIR}/${ARCHIVE_NAME}" | awk '{print $1}')
        else
            ACTUAL_SHA=""
        fi

        if [ -n "${ACTUAL_SHA}" ]; then
            if [ "${ACTUAL_SHA}" = "${EXPECTED_SHA}" ]; then
                SHORT_SHA=$(printf '%.16s' "${ACTUAL_SHA}")
                success "SHA256 checksum verified: ${SHORT_SHA}..."
            else
                error "Checksum mismatch! Expected ${EXPECTED_SHA}, got ${ACTUAL_SHA}."
            fi
        else
            error "Neither sha256sum nor shasum is installed to verify checksum."
        fi
    fi
else
    warn "Checksums file not found for release ${TAG}, skipping checksum verification."
fi

# 8. Extract Binary
info "Validating archive entries for security..."
if tar -tzf "${TMP_DIR}/${ARCHIVE_NAME}" 2>/dev/null | grep -E '(^|/)\.\.(/|$)' >/dev/null 2>&1; then
    error "Archive contains illegal relative path traversal entries (..)!"
fi

info "Extracting binary..."
tar -xzf "${TMP_DIR}/${ARCHIVE_NAME}" -C "${TMP_DIR}"

if [ ! -f "${TMP_DIR}/charites" ]; then
    error "Failed to locate 'charites' binary inside extracted archive."
fi

chmod +x "${TMP_DIR}/charites"

# 9. Select Destination Directory
USE_SUDO=false
if [ -n "${CHARITES_INSTALL_DIR:-}" ]; then
    INSTALL_DIR="${CHARITES_INSTALL_DIR}"
    mkdir -p "${INSTALL_DIR}"
    DEST="${INSTALL_DIR}/charites"
elif [ -w "/usr/local/bin" ]; then
    INSTALL_DIR="/usr/local/bin"
    DEST="${INSTALL_DIR}/charites"
elif command -v sudo >/dev/null 2>&1 && [ -n "${SUDO_USER:-}" ] || (command -v sudo >/dev/null 2>&1 && sudo -n true 2>/dev/null); then
    INSTALL_DIR="/usr/local/bin"
    USE_SUDO=true
    DEST="${INSTALL_DIR}/charites"
else
    INSTALL_DIR="${HOME}/.local/bin"
    mkdir -p "${INSTALL_DIR}"
    DEST="${INSTALL_DIR}/charites"
fi

info "Installing to ${DEST}..."
if [ "${USE_SUDO}" = true ]; then
    sudo cp "${TMP_DIR}/charites" "${DEST}"
    sudo chmod 755 "${DEST}"
else
    cp "${TMP_DIR}/charites" "${DEST}"
    chmod 755 "${DEST}"
fi

# macOS quarantine removal
if [ "${TARGET_OS}" = "darwin" ]; then
    xattr -d com.apple.quarantine "${DEST}" 2>/dev/null || true
fi

# 10. Check and Configure PATH if installed to user directory
PATH_NEEDS_UPDATE=false
case ":${PATH}:" in
    *":${INSTALL_DIR}:"*) ;;
    *) PATH_NEEDS_UPDATE=true ;;
esac

if [ "${PATH_NEEDS_UPDATE}" = true ]; then
    SHELL_PROFILE=""

    if [ -n "${ZSH_VERSION:-}" ] || [ -f "${HOME}/.zshrc" ]; then
        SHELL_PROFILE="${HOME}/.zshrc"
    elif [ -f "${HOME}/.bashrc" ]; then
        SHELL_PROFILE="${HOME}/.bashrc"
    elif [ -f "${HOME}/.profile" ]; then
        SHELL_PROFILE="${HOME}/.profile"
    fi

    if [ -n "${SHELL_PROFILE}" ]; then
        EXPORT_CMD="export PATH=\"${INSTALL_DIR}:\$PATH\""
        if ! grep -qs "${INSTALL_DIR}" "${SHELL_PROFILE}"; then
            {
                echo ""
                echo "# Charites Compiler & Static Linter"
                echo "${EXPORT_CMD}"
            } >> "${SHELL_PROFILE}"
            info "Added ${INSTALL_DIR} to ${SHELL_PROFILE}."
        fi
        export PATH="${INSTALL_DIR}:${PATH}"
    fi
fi

# 11. Verify Installation
success "Charites ${TAG} successfully installed to ${DEST}!"
printf "\n"
if command -v charites >/dev/null 2>&1; then
    charites --version || true
else
    "${DEST}" --version || true
fi

printf "\n%bQuick Start:%b\n" "${BOLD}" "${NC}"
printf "  charites --help\n"
printf "  charites scan .\n\n"
