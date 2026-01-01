#!/usr/bin/env bash
set -e

INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
REPO="comethrusws/mirage"
BINARY_NAME="mirage"

echo "Installing Mirage..."

OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
    Linux*)     OS_TYPE=Linux;;
    Darwin*)    OS_TYPE=Darwin;;
    *)          echo "Unsupported OS: $OS"; exit 1;;
esac

case "$ARCH" in
    x86_64)     ARCH_TYPE=x86_64;;
    arm64|aarch64) ARCH_TYPE=arm64;;
    *)          echo "Unsupported architecture: $ARCH"; exit 1;;
esac

LATEST_RELEASE=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
    echo "Failed to fetch latest release"
    exit 1
fi

DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_RELEASE}/mirage_${OS_TYPE}_${ARCH_TYPE}.tar.gz"

echo "Downloading ${BINARY_NAME} ${LATEST_RELEASE}..."
curl -sL "$DOWNLOAD_URL" | tar xz "$BINARY_NAME"

echo "Installing to ${INSTALL_DIR}..."
sudo mv "$BINARY_NAME" "${INSTALL_DIR}/${BINARY_NAME}"
sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

echo "✓ ${BINARY_NAME} installed successfully!"
echo "Run '${BINARY_NAME} --help' to get started"
