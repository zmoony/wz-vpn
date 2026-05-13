#!/usr/bin/env bash
set -euo pipefail

PREFIX="${PI_GATEWAY_INSTALL_PREFIX:-/opt/pi-gateway}"
SYSTEMD_DIR="${PI_GATEWAY_SYSTEMD_DIR:-/etc/systemd/system}"
SERVICE_NAME="pi-gateway.service"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

BIN_SOURCE="${PI_GATEWAY_BIN_SOURCE:-${REPO_ROOT}/pi-gateway}"
WEB_DIST_SOURCE="${PI_GATEWAY_WEB_DIST_SOURCE:-${REPO_ROOT}/web/dist}"
ENV_SOURCE="${PI_GATEWAY_ENV_SOURCE:-${REPO_ROOT}/.env.example}"
SERVICE_SOURCE="${PI_GATEWAY_SERVICE_SOURCE:-${REPO_ROOT}/deploy/${SERVICE_NAME}}"

BIN_DIR="${PREFIX}/bin"
ETC_DIR="${PREFIX}/etc"
DATA_DIR="${PREFIX}/data"
WEB_DIR="${PREFIX}/web"
WEB_DIST_DIR="${WEB_DIR}/dist"
TARGET_ENV_FILE="${ETC_DIR}/pi-gateway.env"

require_command() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "missing required command: $1" >&2
    exit 1
  fi
}

copy_tree() {
  local source="$1"
  local target="$2"

  mkdir -p "${target}"
  cp -R "${source}/." "${target}/"
}

echo "[1/5] Checking required commands"
for cmd in install mkdir cp systemctl; do
  require_command "${cmd}"
done

if [[ ! -f "${BIN_SOURCE}" ]]; then
  echo "pi-gateway binary not found: ${BIN_SOURCE}" >&2
  echo "set PI_GATEWAY_BIN_SOURCE or place the built binary at the repo root" >&2
  exit 1
fi

if [[ ! -f "${ENV_SOURCE}" ]]; then
  echo "env template not found: ${ENV_SOURCE}" >&2
  exit 1
fi

if [[ ! -f "${SERVICE_SOURCE}" ]]; then
  echo "systemd service template not found: ${SERVICE_SOURCE}" >&2
  exit 1
fi

echo "[2/5] Creating target directories under ${PREFIX}"
mkdir -p \
  "${BIN_DIR}" \
  "${ETC_DIR}" \
  "${DATA_DIR}" \
  "${DATA_DIR}/nginx/sites-enabled" \
  "${DATA_DIR}/ddns-go" \
  "${DATA_DIR}/certs" \
  "${DATA_DIR}/nftables/backups" \
  "${WEB_DIST_DIR}"

echo "[3/5] Installing binary and optional frontend assets"
install -m 0755 "${BIN_SOURCE}" "${BIN_DIR}/pi-gateway"
if [[ -d "${WEB_DIST_SOURCE}" ]]; then
  copy_tree "${WEB_DIST_SOURCE}" "${WEB_DIST_DIR}"
else
  echo "warning: web dist directory not found, skipping static assets: ${WEB_DIST_SOURCE}" >&2
fi

echo "[4/5] Installing environment file and systemd unit"
if [[ ! -f "${TARGET_ENV_FILE}" ]]; then
  install -m 0644 "${ENV_SOURCE}" "${TARGET_ENV_FILE}"
else
  echo "keeping existing env file: ${TARGET_ENV_FILE}"
fi
install -m 0644 "${SERVICE_SOURCE}" "${SYSTEMD_DIR}/${SERVICE_NAME}"

echo "[5/5] Reloading and enabling systemd service"
systemctl daemon-reload
systemctl enable "${SERVICE_NAME}"

if [[ "${PI_GATEWAY_RESTART_SERVICE:-1}" == "1" ]]; then
  systemctl restart "${SERVICE_NAME}"
  systemctl --no-pager --full status "${SERVICE_NAME}" || true
else
  echo "skipping service restart because PI_GATEWAY_RESTART_SERVICE=${PI_GATEWAY_RESTART_SERVICE:-0}"
fi

cat <<EOF
Pi-Gateway installed.

Prefix: ${PREFIX}
Binary: ${BIN_DIR}/pi-gateway
Env: ${TARGET_ENV_FILE}
Service: ${SYSTEMD_DIR}/${SERVICE_NAME}

Next steps:
1. Review ${TARGET_ENV_FILE}
2. Ensure host dependencies like wg/nginx/nft/acme.sh/ddns-go are installed
3. Restart the service after any env changes: systemctl restart ${SERVICE_NAME}
EOF
