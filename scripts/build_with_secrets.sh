#!/bin/sh
set -eu

# Usage:
#  export MQTT_BROKER_URL="mqtt://host:1883"
#  export MQTT_USERNAME="user"
#  export MQTT_PASSWORD="pass"
#  export MQTT_TOPIC_PREFIX="printer/"
#  export ID_MERCHANT="..."
#  export NAME_MERCHANT="..."
#  export PRINTER_NAME="EPSON TM-U220 Receipt"
#  export PAPER_WIDTH="42"
#  export OFFSITE="8"
#  ./scripts/build_with_secrets.sh

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
ENV_FILE="${1:-.env}"

FIRST=1
json='{'
add_pair() {
  key="$1"
  val="${2-}"
  if [ -n "$val" ]; then
    esc=$(printf '%s' "$val" | sed -e 's/\\/\\\\/g' -e 's/"/\\"/g')
    if [ "$FIRST" -eq 1 ]; then
      json="$json\"$key\":\"$esc\""
      FIRST=0
    else
      json="$json,\"$key\":\"$esc\""
    fi
  fi
}

if [ ! -f "$ENV_FILE" ]; then
  echo "ENV file not found: $ENV_FILE" >&2
  exit 1
fi

while IFS= read -r line; do
  case "$line" in
    ''|\#*) continue ;;
  esac
  case "$line" in
    *=*) ;;
    *) continue ;;
  esac
  key=${line%%=*}
  val=${line#*=}
  key=$(printf '%s' "$key" | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//')
  val=$(printf '%s' "$val" | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//')
  case "$val" in
    \"*\") val=$(printf '%s' "$val" | sed -e 's/^"//' -e 's/"$//') ;;
    \'*\') val=$(printf '%s' "$val" | sed -e "s/^'//" -e "s/'$//") ;;
  esac
  add_pair "$key" "$val"
done < "$ENV_FILE"
json="$json}"

# Base64 encode JSON to avoid ldflags quoting issues
if command -v base64 >/dev/null 2>&1; then
  b64=$(printf '%s' "$json" | base64 | tr -d '\n')
else
  echo "base64 command not found" >&2
  exit 1
fi

ldflags="-X other/go-printer-termal/infrastructure/config.EmbeddedEnvBase64=${b64}"

echo "Building clientMidasPrinter.exe with embedded secrets (WSL sh)..."
cd "$ROOT_DIR"

# Detect go command (WSL or Windows)
GO_CMD="go"
if ! command -v "$GO_CMD" >/dev/null 2>&1; then
  if command -v go.exe >/dev/null 2>&1; then
    GO_CMD="go.exe"
  else
    echo "Go toolchain not found in WSL or Windows PATH" >&2
    exit 1
  fi
fi

if [ "$GO_CMD" = "go" ]; then
  GOOS=windows GOARCH=amd64 "$GO_CMD" build -ldflags "$ldflags" -o clientMidasPrinter.exe ./cmd/printer
else
  "$GO_CMD" build -ldflags "$ldflags" -o clientMidasPrinter.exe ./cmd/printer
fi
echo "Output: $ROOT_DIR/clientMidasPrinter.exe"