#!/usr/bin/env bash
set -euo pipefail

# Requirements:
# - bash, openssl, curl
# - jq (optional but recommended). If not present, simple grep/sed fallbacks are used.

STREAM_HOST_DEFAULT="https://chat.stream-io-api.com"

err() { echo "[error] $*" >&2; }
info() { echo "[info]  $*" >&2; }

require() {
  command -v "$1" >/dev/null 2>&1 || { err "'$1' is required"; exit 1; }
}

require openssl
require curl

STREAM_KEY=${STREAM_KEY:-}
STREAM_SECRET=${STREAM_SECRET:-}
STREAM_HOST=${STREAM_HOST:-$STREAM_HOST_DEFAULT}

if [[ -z "$STREAM_KEY" || -z "$STREAM_SECRET" ]]; then
  err "Set STREAM_KEY and STREAM_SECRET env vars first"
  exit 1
fi

# base64url without padding
b64url() {
  openssl base64 -A | tr '+/' '-_' | tr -d '='
}

# Create a server-side JWT: header {alg:HS256, typ:JWT}, payload {server:true}
make_server_jwt() {
  local header payload unsigned signature
  header=$(printf '{"alg":"HS256","typ":"JWT"}' | b64url)
  payload=$(printf '{"server":true}' | b64url)
  unsigned="$header.$payload"
  signature=$(printf '%s' "$unsigned" | openssl dgst -binary -sha256 -hmac "$STREAM_SECRET" | b64url)
  printf '%s.%s\n' "$unsigned" "$signature"
}

AUTH_TOKEN=$(make_server_jwt)
info "Using STREAM_HOST: $STREAM_HOST"
info "Using STREAM_KEY: $STREAM_KEY"
info "Authorization (server JWT): $AUTH_TOKEN"

# Helpers
post() {
  local path="$1"; shift
  curl -sS -X POST "$STREAM_HOST/$path?api_key=$STREAM_KEY" \
    -H "Authorization: $AUTH_TOKEN" \
    -H "Stream-Auth-Type: jwt" \
    -H "Content-Type: application/json" \
    -d "$*"
}

del() {
  local path="$1"
  # Remaining args are query string (already encoded) and body json (optional)
  local qs="${2-}"
  local data="${3-}"
  if [[ -n "$data" ]]; then
    curl -sS -X DELETE "$STREAM_HOST/$path?api_key=$STREAM_KEY&$qs" \
      -H "Authorization: $AUTH_TOKEN" \
      -H "Stream-Auth-Type: jwt" \
      -H "Content-Type: application/json" \
      -d "$data"
  else
    curl -sS -X DELETE "$STREAM_HOST/$path?api_key=$STREAM_KEY&$qs" \
      -H "Authorization: $AUTH_TOKEN" \
      -H "Stream-Auth-Type: jwt" \
      -H "Content-Type: application/json"
  fi
}

# Generate IDs
rand() { LC_ALL=C tr -dc 'a-zA-Z0-9' </dev/urandom | head -c 10; echo; }
USER_ID="test-user-delete-for-me-$(rand)"
CHANNEL_ID="test-channel-delete-for-me-$(rand)"
CHANNEL_TYPE="messaging"

info "Creating user: $USER_ID"
upsert_payload=$(cat <<JSON
{"users": {"$USER_ID": {"id": "$USER_ID"}}}
JSON
)
post users "$upsert_payload" >/dev/null

info "Creating channel: $CHANNEL_TYPE:$CHANNEL_ID with member $USER_ID"
create_channel_payload=$(cat <<JSON
{"data": {"members": ["$USER_ID"], "created_by_id": "$USER_ID"}}
JSON
)
post channels/$CHANNEL_TYPE/$CHANNEL_ID/query "$create_channel_payload" >/dev/null

info "Sending message from $USER_ID"
send_msg_payload=$(cat <<JSON
{"message": {"text": "Test message for delete_for_me","user_id": "$USER_ID"}}
JSON
)
msg_resp=$(post channels/$CHANNEL_TYPE/$CHANNEL_ID/message "$send_msg_payload")
echo "$msg_resp" | sed 's/.*/[debug] send message response: &/' >&2

if command -v jq >/dev/null 2>&1; then
  MESSAGE_ID=$(printf '%s' "$msg_resp" | jq -r '.message.id')
else
  MESSAGE_ID=$(printf '%s' "$msg_resp" | sed -n 's/.*"id":"\([^"]\+\)".*/\1/p' | head -n1)
fi

if [[ -z "$MESSAGE_ID" || "$MESSAGE_ID" == "null" ]]; then
  err "Failed to parse message id"
  echo "$msg_resp" >&2
  exit 1
fi

# info "Deleting message for me (user_id) message_id=$MESSAGE_ID user_id=$USER_ID"
# del_resp=$(del messages/$MESSAGE_ID "delete_for_me=true&user_id=$USER_ID")
# echo "$del_resp" | sed 's/.*/[debug] delete-for-me (user_id) response: &/' >&2

# If server expects 'deleted_by' instead, try fallback
info "Retrying delete_for_me using deleted_by param"
del_resp2=$(del messages/$MESSAGE_ID "delete_for_me=true&deleted_by=$USER_ID")
echo "$del_resp2" | sed 's/.*/[debug] delete-for-me (deleted_by) response: &/' >&2

info "Done"

