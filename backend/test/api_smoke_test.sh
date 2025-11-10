#!/usr/bin/env bash
# Strict curl smoke test for ask.zip backend routes (Gin)
# - Covers Auth, Users, Memberships, Orders, Web, Chat
# - JWT required for most endpoints
# - No hard dependency on jq (falls back to python3); needs at least one of {jq, python3}
set -euo pipefail

# ----------------------------
# Config (override via env)
# ----------------------------
BASE_URL="${BASE_URL:-http://127.0.0.1:5000}"
USERNAME="${USERNAME:-testuser}"
PASSWORD="${PASSWORD:-Test@12345}"
NEW_PASSWORD="${NEW_PASSWORD:-Test@12345_new}"
FULL_NAME="${FULL_NAME:-Test User}"
EMAIL="${EMAIL:-testuser@example.com}"
PHONE="${PHONE:-13900001234}"
SQ1="${SQ1:-Your first pet?}"
SA1="${SA1:-cat}"
SQ2="${SQ2:-Your favorite teacher?}"
SA2="${SA2:-alice}"

# ----------------------------
# Tiny utils
# ----------------------------
die() { echo "ERROR: $*" >&2; exit 1; }
need() { command -v "$1" >/dev/null 2>&1 || die "Missing dependency: $1"; }
need curl

HAS_JQ=0; command -v jq >/dev/null 2>&1 && HAS_JQ=1 || true
HAS_PY=0; command -v python3 >/dev/null 2>&1 && HAS_PY=1 || true
[ $HAS_JQ -eq 1 -o $HAS_PY -eq 1 ] || die "Need either jq or python3 for JSON parsing"

color() { printf "\033[%sm%s\033[0m" "$1" "$2"; }
h1() { printf "\n%s %s\n" "$(color '1;36' '#')" "$(color '1;36' "$*")"; }
step() { printf "%s %s\n" "$(color '0;33' '-')" "$*"; }
auth_header() { printf 'Authorization: Bearer %s' "$TOKEN"; }

pp() {
  if [ $HAS_JQ -eq 1 ]; then
    jq -C . 2>/dev/null || cat
  else
    python3 - <<'PY' 2>/dev/null || cat
import sys, json
try:
    obj = json.load(sys.stdin)
    print(json.dumps(obj, ensure_ascii=False, indent=2))
except Exception:
    print(sys.stdin.read())
PY
  fi
}

# Read JSON from stdin; try keys in order; print first non-empty scalar
# Supports dotted paths (e.g., data.user_id)
json_pick() {
  if [ $HAS_PY -eq 1 ]; then
    python3 - "$@" <<'PY'
import sys, json
try:
    data = json.load(sys.stdin)
except Exception:
    print("")
    sys.exit(0)

def get_path(d, path):
    cur = d
    for p in path.split("."):
        if isinstance(cur, dict) and p in cur:
            cur = cur[p]
        else:
            return None
    return cur

for key in sys.argv[1:]:
    v = get_path(data, key)
    if v not in (None, "", []):
        if isinstance(v, (dict, list)):
            continue
        print(v)
        sys.exit(0)
print("")
PY
  else
    # crude last-resort
    tr -d '\n' | sed 's/ //g' | awk -v k="$1" -F'"' '
      { for(i=1;i<NF;i++){ if($i==k){ print $(i+2); exit } } }
    '
  fi
}

# ----------------------------
# 0) Auth: register, login, me
# ----------------------------
h1 "Auth: register + login + me"

step "Register (idempotent) POST /api/auth/register"
set +e
register_resp=$(curl -sS -X POST "$BASE_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d @- <<JSON
{
  "username": "$USERNAME",
  "password": "$PASSWORD",
  "full_name": "$FULL_NAME",
  "email": "$EMAIL",
  "phone_number": "$PHONE",
  "security_question1": "$SQ1",
  "security_answer1": "$SA1",
  "security_question2": "$SQ2",
  "security_answer2": "$SA2"
}
JSON
) || true
set -e
echo "$register_resp" | pp || true

step "Login POST /api/auth/login"
login_resp=$(curl -sS -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}")
echo "$login_resp" | pp
TOKEN=$(echo "$login_resp" | json_pick token)
[ -n "$TOKEN" ] || die "No token in login response"

step "Me GET /api/auth/me"
me_resp=$(curl -sS -X GET "$BASE_URL/api/auth/me" -H "$(auth_header)")
echo "$me_resp" | pp
USER_ID=$(echo "$me_resp" | json_pick user_id data.user_id id data.id)
[ -n "$USER_ID" ] || die "Failed to parse user_id from /api/auth/me"

# ----------------------------
# 1) Users: update (strict route)
# ----------------------------
h1 "Users: update self (PUT /api/users/:user_id)"

step "Update user PUT /api/users/$USER_ID"
update_resp=$(curl -sS -X PUT "$BASE_URL/api/users/$USER_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"full_name\":\"${FULL_NAME}_updated\",\"email\":\"updated_${EMAIL}\",\"phone_number\":\"${PHONE}\"}")
echo "$update_resp" | pp

# ----------------------------
# 2) Auth: verify-security + reset-password
# ----------------------------
h1 "Auth: verify-security + reset-password"

step "Verify POST /api/auth/verify-security"
verify_resp=$(curl -sS -X POST "$BASE_URL/api/auth/verify-security" \
  -H "Content-Type: application/json" \
  -d @- <<JSON
{"username":"$USERNAME","security_answer1":"$SA1","security_answer2":"$SA2"}
JSON
)
echo "$verify_resp" | pp
RESET_TOKEN=$(echo "$verify_resp" | json_pick reset_token data.reset_token)

if [ -n "${RESET_TOKEN:-}" ]; then
  step "Reset password POST /api/auth/reset-password"
  reset_resp=$(curl -sS -X POST "$BASE_URL/api/auth/reset-password" \
    -H "Content-Type: application/json" \
    -d "{\"reset_token\":\"$RESET_TOKEN\",\"new_password\":\"$NEW_PASSWORD\"}")
  echo "$reset_resp" | pp

  step "Re-login with NEW password"
  login2_resp=$(curl -sS -X POST "$BASE_URL/api/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$USERNAME\",\"password\":\"$NEW_PASSWORD\"}")
  echo "$login2_resp" | pp
  TOKEN_NEW=$(echo "$login2_resp" | json_pick token)
  if [ -n "$TOKEN_NEW" ]; then TOKEN="$TOKEN_NEW"; fi
else
  step "No reset_token returned, skip password reset"
fi

# ----------------------------
# 3) Membership CRUD
# ----------------------------
h1 "Membership: list/create/get/update"

step "List GET /api/membership"
curl -sS -X GET "$BASE_URL/api/membership" -H "$(auth_header)" | pp

step "Create POST /api/membership"
create_member_resp=$(curl -sS -X POST "$BASE_URL/api/membership" \
  -H "Content-Type: application/json" -H "$(auth_header)" \
  -d @- <<JSON
{
  "user_id": $USER_ID,
  "start_date": "2025-01-01",
  "expire_date": "2026-01-01",
  "status": "active"
}
JSON
)
echo "$create_member_resp" | pp
MEMBERSHIP_ID=$(echo "$create_member_resp" | json_pick membership_id data.membership_id id data.id)

step "By user GET /api/membership/$USER_ID"
curl -sS -X GET "$BASE_URL/api/membership/$USER_ID" -H "$(auth_header)" | pp

if [ -n "${MEMBERSHIP_ID:-}" ]; then
  step "Update PUT /api/membership/$MEMBERSHIP_ID"
  curl -sS -X PUT "$BASE_URL/api/membership/$MEMBERSHIP_ID" \
    -H "Content-Type: application/json" -H "$(auth_header)" \
    -d '{"expire_date":"2027-01-01","status":"active"}' | pp
fi

# ----------------------------
# 4) Orders
# ----------------------------
h1 "Orders: create/list/latest/recent"

step "Create order POST /api/membership/orders"
order_resp=$(curl -sS -X POST "$BASE_URL/api/membership/orders" \
  -H "Content-Type: application/json" -H "$(auth_header)" \
  -d @- <<JSON
{
  "user_id": $USER_ID,
  "duration_months": 12,
  "amount": 199.00,
  "payment_method": "other"
}
JSON
)
echo "$order_resp" | pp

step "List GET /api/membership/orders/$USER_ID"
curl -sS -X GET "$BASE_URL/api/membership/orders/$USER_ID" -H "$(auth_header)" | pp

step "Latest GET /api/membership/orders/$USER_ID/latest"
curl -sS -X GET "$BASE_URL/api/membership/orders/$USER_ID/latest" -H "$(auth_header)" | pp

step "Recent GET /api/membership/orders/$USER_ID/recent?n=3"
curl -sS -X GET "$BASE_URL/api/membership/orders/$USER_ID/recent?n=3" -H "$(auth_header)" | pp

# ----------------------------
# 5) Web: items/search/ingest/chunk
# ----------------------------
h1 "Web: items + search + ingest/chunk"

step "Create page POST /web/items (fetch=true)"
create_page_resp=$(curl -sS -X POST "$BASE_URL/web/items" \
  -H "Content-Type: application/json" -H "$(auth_header)" \
  -d '{"url":"https://example.com","title":"Example Home","fetch":true}')
echo "$create_page_resp" | pp
PAGE_ID=$(echo "$create_page_resp" | json_pick id data.id)

step "List pages GET /web/items"
curl -sS -X GET "$BASE_URL/web/items" -H "$(auth_header)" | pp

if [ -n "${PAGE_ID:-}" ]; then
  step "Get page GET /web/page/$PAGE_ID"
  curl -sS -X GET "$BASE_URL/web/page/$PAGE_ID" -H "$(auth_header)" | pp

  step "Update page PUT /web/items/$PAGE_ID"
  curl -sS -X PUT "$BASE_URL/web/items/$PAGE_ID" \
    -H "Content-Type: application/json" -H "$(auth_header)" \
    -d '{"title":"Example Updated","content":"Hello world"}' | pp
fi

step "Search by query POST /web/search"
curl -sS -X POST "$BASE_URL/web/search" \
  -H "Content-Type: application/json" -H "$(auth_header)" \
  -d '{"q":"Example","top_k":5}' | pp

step "Ingest by urls POST /web/search"
curl -sS -X POST "$BASE_URL/web/search" \
  -H "Content-Type: application/json" -H "$(auth_header)" \
  -d '{"urls":["https://www.rfc-editor.org/"],"top_k":3}' | pp

step "Ingest POST /web/ingest"
curl -sS -X POST "$BASE_URL/web/ingest" -H "$(auth_header)" | pp

step "Chunk POST /web/chunk"
curl -sS -X POST "$BASE_URL/web/chunk"  -H "$(auth_header)" | pp

if [ -n "${PAGE_ID:-}" ]; then
  step "Delete page DELETE /web/items/$PAGE_ID"
  curl -sS -X DELETE "$BASE_URL/web/items/$PAGE_ID" -H "$(auth_header)" | pp
fi

# ----------------------------
# 6) Chat: sessions/messages/completion/stream
# ----------------------------
h1 "Chat: sessions + messages + completion + stream"

step "Create session POST /chat/sessions"
create_sess_resp=$(curl -sS -X POST "$BASE_URL/chat/sessions" \
  -H "Content-Type: application/json" -H "$(auth_header)" \
  -d '{"title":"Smoke Test Session"}')
echo "$create_sess_resp" | pp
SESSION_ID=$(echo "$create_sess_resp" | json_pick id data.id session_id data.session_id)
[ -n "$SESSION_ID" ] || die "Failed to create session"

step "List sessions GET /chat/sessions"
curl -sS -X GET "$BASE_URL/chat/sessions" -H "$(auth_header)" | pp

step "Append message POST /chat/sessions/:id/messages"
curl -sS -X POST "$BASE_URL/chat/sessions/$SESSION_ID/messages" \
  -H "Content-Type: application/json" -H "$(auth_header)" \
  -d '{"role":"user","content":"Hello from curl"}' | pp

step "Alias POST /chat/messages (with session_id)"
curl -sS -X POST "$BASE_URL/chat/messages" \
  -H "Content-Type: application/json" -H "$(auth_header)" \
  -d "{\"session_id\":\"$SESSION_ID\",\"content\":\"Second message via alias\"}" | pp

step "Get messages GET /chat/sessions/:id/messages"
curl -sS -X GET "$BASE_URL/chat/sessions/$SESSION_ID/messages" -H "$(auth_header)" | pp

step "Complete POST /chat/sessions/:id/complete (tolerant)"
curl -sS -X POST "$BASE_URL/chat/sessions/$SESSION_ID/complete" \
  -H "Content-Type: application/json" -H "$(auth_header)" \
  -d '{"content":"Briefly echo hello","model":""}' | pp || true

step "Stream POST /chat/sessions/:id/stream (tolerant)"
curl -sS -N -X POST "$BASE_URL/chat/sessions/$SESSION_ID/stream" \
  -H "Content-Type: application/json" -H "$(auth_header)" \
  -d '{"content":"streaming hello","model":""}' || true

# ----------------------------
# 7) Dangerous delete (disabled)
# ----------------------------
h1 "Auth: delete user (skipped)"
echo "# To test account deletion:"
echo "# curl -sS -X DELETE \"$BASE_URL/api/auth/delete\" -H \"$(auth_header)\" |"
echo "#   { if [ $HAS_JQ -eq 1 ]; then jq -C .; else cat; fi; }"

echo
h1 "All requests completed."
