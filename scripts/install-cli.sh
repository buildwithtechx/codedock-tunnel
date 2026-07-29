#!/usr/bin/env bash

set -euo pipefail

repo="${CODEDOCK_TUNNEL_REPO:-codedock-tunnel/codedock-tunnel}"
version="${CODEDOCK_TUNNEL_VERSION:-latest}"
install_dir="${CODEDOCK_TUNNEL_INSTALL_DIR:-$HOME/.local/bin}"

case "$(uname -s)" in
  Linux) os="linux" ;;
  Darwin) os="darwin" ;;
  *) echo "unsupported operating system" >&2; exit 1 ;;
esac

case "$(uname -m)" in
  x86_64|amd64) arch="amd64" ;;
  arm64|aarch64) arch="arm64" ;;
  *) echo "unsupported architecture" >&2; exit 1 ;;
esac

asset="codedock-cli_${os}_${arch}.tar.gz"
base_url="https://github.com/${repo}/releases"
if [ "$version" = "latest" ]; then
  download_url="${base_url}/latest/download/${asset}"
else
  download_url="${base_url}/download/${version}/${asset}"
fi

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

curl --fail --location --silent --show-error "$download_url" -o "$tmp_dir/$asset"
tar -xzf "$tmp_dir/$asset" -C "$tmp_dir"

mkdir -p "$install_dir"
install "$tmp_dir/codedock-cli" "$install_dir/codedock-tunnel"
echo "installed codedock-tunnel to $install_dir/codedock-tunnel"
