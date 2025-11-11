## 0) 基础环境变量（建议先执行）

```bash
export BASE_URL="http://127.0.0.1:5000"
export USERNAME="testuser"
export PASSWORD="Test@12345"
export NEW_PASSWORD="Test@12345_new"
export FULL_NAME="Test User"
export EMAIL="testuser@example.com"
export PHONE="13900001234"
export SQ1="Your first pet?"
export SA1="cat"
export SQ2="Your favorite teacher?"
export SA2="alice"
```

## 1) Auth 相关

### 1.1 注册

```bash
curl -X POST "$BASE_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\":\"$USERNAME\",
    \"password\":\"$PASSWORD\",
    \"full_name\":\"$FULL_NAME\",
    \"email\":\"$EMAIL\",
    \"phone_number\":\"$PHONE\",
    \"security_question1\":\"$SQ1\",
    \"security_answer1\":\"$SA1\",
    \"security_question2\":\"$SQ2\",
    \"security_answer2\":\"$SA2\"
  }"
```

### 1.2 登录 → 提取 TOKEN
```bash
TOKEN=$(curl -sS -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}" | jq -r .token) && echo "$TOKEN"
```

### 1.3 获取我的信息（拿 user_id）

```bash
curl -X GET "$BASE_URL/api/auth/me" \
  -H "Authorization: Bearer $TOKEN"
```

* 提取 `USER_ID`（jq）

```bash
USER_ID=$(curl -sS -X GET "$BASE_URL/api/auth/me" \
  -H "Authorization: Bearer $TOKEN" | jq -r '.user_id // .data.user_id') && echo "$USER_ID"
```


### 1.4 更新用户资料（注意**路由是 /api/users/:user_id**）

```bash
curl -X PUT "$BASE_URL/api/users/$USER_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"full_name\":\"${FULL_NAME}_updated\",
    \"email\":\"updated_${EMAIL}\",
    \"phone_number\":\"$PHONE\"
  }"
```

### 1.5 校验密保 → 取 reset_token

```bash
curl -X POST "$BASE_URL/api/auth/verify-security" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"security_answer1\":\"$SA1\",\"security_answer2\":\"$SA2\"}"
```

* 提取 `RESET_TOKEN`（任选其一）

```bash
RESET_TOKEN=$(curl -sS -X POST "$BASE_URL/api/auth/verify-security" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"security_answer1\":\"$SA1\",\"security_answer2\":\"$SA2\"}" | jq -r '.reset_token') && echo "$RESET_TOKEN"


### 1.6 重置密码（可选）

```bash
curl -X POST "$BASE_URL/api/auth/reset-password" \
  -H "Content-Type: application/json" \
  -d "{\"reset_token\":\"$RESET_TOKEN\",\"new_password\":\"$NEW_PASSWORD\"}"
```

### 1.7 用新密码再登录（可选）

```bash
TOKEN=$(curl -sS -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$NEW_PASSWORD\"}" | jq -r .token) && echo "$TOKEN"
```

## 2) Membership 相关

### 2.1 列表

```bash
curl -X GET "$BASE_URL/api/membership" \
  -H "Authorization: Bearer $TOKEN"
```

### 2.2 创建

```bash
curl -X POST "$BASE_URL/api/membership" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"user_id\": $USER_ID,
    \"start_date\": \"2025-01-01\",
    \"expire_date\": \"2026-01-01\",
    \"status\": \"active\"
  }"
```

### 2.3 按用户查询

```bash
curl -X GET "$BASE_URL/api/membership/$USER_ID" \
  -H "Authorization: Bearer $TOKEN"
```

### 2.4 更新（需 membership_id）

* 先拿 `MEMBERSHIP_ID`（你可以从 2.2 的返回里提取，或查 2.3 的返回）：

```bash
export MEMBERSHIP_ID=1
```

* 更新：

```bash
curl -X PUT "$BASE_URL/api/membership/$MEMBERSHIP_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"expire_date":"2027-01-01","status":"active"}'
```


## 3) Orders 相关

### 3.1 创建订单

```bash
curl -X POST "$BASE_URL/api/membership/orders" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"user_id\": $USER_ID,
    \"duration_months\": 12,
    \"amount\": 199.00,
    \"payment_method\": \"other\"
  }"
```

### 3.2 订单列表

```bash
curl -X GET "$BASE_URL/api/membership/orders/$USER_ID" \
  -H "Authorization: Bearer $TOKEN"
```

### 3.3 最新订单

```bash
curl -X GET "$BASE_URL/api/membership/orders/$USER_ID/latest" \
  -H "Authorization: Bearer $TOKEN"
```

### 3.4 最近 N 条

```bash
curl -X GET "$BASE_URL/api/membership/orders/$USER_ID/recent?n=3" \
  -H "Authorization: Bearer $TOKEN"
```

---

## 4) Web 相关

### 4.1 创建页面（可抓取）

```bash
curl -X POST "$BASE_URL/web/items" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"url":"https://www.njupt.edu.cn/","title":"Example Home","fetch":true}'
```

### 4.2 列出页面

```bash
curl -X GET "$BASE_URL/web/items" \
  -H "Authorization: Bearer $TOKEN"
```

### 4.3 获取页面详情（需 PAGE_ID）

```bash
export PAGE_ID=1
curl -X GET "$BASE_URL/web/page/$PAGE_ID" \
  -H "Authorization: Bearer $TOKEN"
```

### 4.4 更新页面

```bash
curl -X PUT "$BASE_URL/web/items/$PAGE_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Example Updated","content":"Hello world"}'
```

### 4.5 搜索（按关键词）

```bash
curl -X POST "$BASE_URL/web/search" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"q":"Example","top_k":5}'
```

### 4.6 搜索（传 URL，会触发抓取）

```bash
curl -X POST "$BASE_URL/web/search" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"urls":["https://www.rfc-editor.org/"],"top_k":3}'
```

### 4.7 兼容占位：ingest / chunk

```bash
curl -X POST "$BASE_URL/web/ingest" -H "Authorization: Bearer $TOKEN"
curl -X POST "$BASE_URL/web/chunk"  -H "Authorization: Bearer $TOKEN"
```

### 4.8 删除页面

```bash
curl -X DELETE "$BASE_URL/web/items/$PAGE_ID" \
  -H "Authorization: Bearer $TOKEN"
```


## 5) Chat 相关

### 5.1 创建会话 → 取 SESSION_ID

```bash
curl -X POST "$BASE_URL/chat/sessions" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Smoke Test Session"}'
```

* 提取 `SESSION_ID`（任选其一）：

```bash
SESSION_ID=$(curl -sS -X POST "$BASE_URL/chat/sessions" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Smoke Test Session"}' | jq -r '.id // .data.id // .session_id // .data.session_id') && echo "$SESSION_ID"

```

### 5.2 列出会话

```bash
curl -X GET "$BASE_URL/chat/sessions" \
  -H "Authorization: Bearer $TOKEN"
```

### 5.3 追加消息（按会话）

```bash
curl -X POST "$BASE_URL/chat/sessions/$SESSION_ID/messages" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role":"user","content":"Hello from curl"}'
```

### 5.4 追加消息（别名 /chat/messages）

```bash
curl -X POST "$BASE_URL/chat/messages" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"session_id\":\"$SESSION_ID\",\"content\":\"Second message via alias\"}"
```

### 5.5 拉取消息

```bash
curl -X GET "$BASE_URL/chat/sessions/$SESSION_ID/messages" \
  -H "Authorization: Bearer $TOKEN"
```

### 5.6 一次性补全（可能依赖 LLM 配置，失败忽略）

```bash
curl -X POST "$BASE_URL/chat/sessions/$SESSION_ID/complete" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content":"Briefly echo hello","model":""}'
```

### 5.7 流式补全（原样打印流）

```bash
curl -N -X POST "$BASE_URL/chat/sessions/$SESSION_ID/stream" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content":"streaming hello","model":""}'

curl -N -X POST "$BASE_URL/chat/sessions/$SESSION_ID/stream" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content":"你好","model":""}'
```

---

## 6) 危险操作：删除账号（默认不执行）

```bash
curl -X DELETE "$BASE_URL/api/auth/delete" \
  -H "Authorization: Bearer $TOKEN"
```


