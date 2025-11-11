## 0) 基础环境变量先设定

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
```
{"message":"注册成功","success":true,"user_id":1}
```

### 1.2 登录 → 提取 TOKEN
```bash
curl -sS -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}" | jq -r .token
```
```bash
TOKEN=$(curl -sS -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}" | jq -r .token) && echo "$TOKEN"
```
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjI4NzI3MjMsInVzZXJfaWQiOjF9.-kkDJEChejtDVUsolTkZ_JE3jy1TcCvrRtBR12CADb0
```

### 1.3 获取我的信息（拿 user_id）

```bash
curl -X GET "$BASE_URL/api/auth/me" \
  -H "Authorization: Bearer $TOKEN"
```
```
{"user_id":1,"username":"testuser","full_name":"Test User","email":"testuser@example.com","phone_number":"13900001234","created_at":"2025-11-10T22:51:48.721+08:00","updated_at":"2025-11-10T22:51:48.721+08:00"}
```
* 提取 `USER_ID`（jq）

```bash
USER_ID=$(curl -sS -X GET "$BASE_URL/api/auth/me" \
  -H "Authorization: Bearer $TOKEN" | jq -r '.user_id // .data.user_id') && echo "$USER_ID"
```
```
1
```


### 1.4 更新用户资料（注意**路由是 /api/users/:user_id**）

```bash
export USER_ID=1
curl -X PUT "$BASE_URL/api/users/$USER_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"full_name\":\"${FULL_NAME}_updated\",
    \"email\":\"updated_${EMAIL}\",
    \"phone_number\":\"$PHONE\"
  }"
```
```
{"message":"用户信息已更新","success":true}
```

### 1.5 校验密保 → 取 reset_token

```bash
curl -X POST "$BASE_URL/api/auth/verify-security" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"security_answer1\":\"$SA1\",\"security_answer2\":\"$SA2\"}"
```
```
{"reset_token":"02ba320129cbfa1da4b5795f1cf38a16","success":true}
```
* 提取 `RESET_TOKEN`

```bash
RESET_TOKEN=$(curl -sS -X POST "$BASE_URL/api/auth/verify-security" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"security_answer1\":\"$SA1\",\"security_answer2\":\"$SA2\"}" | jq -r '.reset_token') && echo "$RESET_TOKEN"
```
```
24d8f04d485bc07691b0b3f522bd47c4
```


### 1.6 重置密码（可选）

```bash
curl -X POST "$BASE_URL/api/auth/reset-password" \
  -H "Content-Type: application/json" \
  -d "{\"reset_token\":\"$RESET_TOKEN\",\"new_password\":\"$NEW_PASSWORD\"}"
```
```失败
{"message":"参数错误","success":false}
```
```成功
{"message":"密码已更新","success":true}
```

### 1.7 用新密码再登录（可选）

```bash
TOKEN=$(curl -sS -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$NEW_PASSWORD\"}" | jq -r .token) && echo "$TOKEN"
```
```bash
curl -sS -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$NEW_PASSWORD\"}"
```
{"expire_at":"2025-11-12T11:37:03+08:00","success":true,"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjI5MTg2MjMsInVzZXJfaWQiOjF9.jesg2jqCdpxvPkMSc-fAd4tsbUAZX32OjWyr0qKceog"}
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjI4NzYwNjAsInVzZXJfaWQiOjF9.OlwdYL0pWSnqr0d6BuqVShoiRSN9BO69ar9FKYV-qnU
```

## 2) Membership 相关

### 2.1 列表

```bash
curl -X GET "$BASE_URL/api/membership" \
  -H "Authorization: Bearer $TOKEN"
```
```
[]
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
```
{"membership_id":1,"message":"会员信息已创建","success":true}
```

### 2.3 按用户查询

```bash
curl -X GET "$BASE_URL/api/membership/$USER_ID" \
  -H "Authorization: Bearer $TOKEN"
```
```
{"membership_id":1,"user_id":1,"start_date":"2025-01-01T00:00:00+08:00","expire_date":"2026-01-01T00:00:00+08:00","status":"active"}
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
```
{"message":"会员信息已更新","success":true}
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
```
{"message":"订单已创建","order_id":1,"success":true}
```

### 3.2 订单列表

```bash
curl -X GET "$BASE_URL/api/membership/orders/$USER_ID" \
  -H "Authorization: Bearer $TOKEN"
```
```
[{"order_id":1,"user_id":1,"purchase_date":"2025-11-10T23:49:37.761+08:00","duration_months":12,"amount":199,"payment_method":"other"}]
```
### 3.3 最新订单

```bash
curl -X GET "$BASE_URL/api/membership/orders/$USER_ID/latest" \
  -H "Authorization: Bearer $TOKEN"
```
```
{"order_id":1,"user_id":1,"purchase_date":"2025-11-10T23:49:37.761+08:00","duration_months":12,"amount":199,"payment_method":"other"}
```

### 3.4 最近 N 条

```bash
curl -X GET "$BASE_URL/api/membership/orders/$USER_ID/recent?n=3" \
  -H "Authorization: Bearer $TOKEN"
```
```
[{"order_id":1,"user_id":1,"purchase_date":"2025-11-10T23:49:37.761+08:00","duration_months":12,"amount":199,"payment_method":"other"}]
```

## 4) Web 相关

### 4.1 创建页面（可抓取）
https://www.nowcoder.com/
https://www.njupt.edu.cn/
https://leetcode.cn/
```bash
curl -X POST "$BASE_URL/web/items" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"url":"https://leetcode.cn/","title":"Example Home","fetch":true}'
```
```
{"error":"db error"}
```
```
{"id":1,"user_id":1,"url":"https://www.njupt.edu.cn/","title":"Example Home","content":"旧版入口 English 电子邮件 首页 南邮概况 \u003e 学校简介 学校章程 南邮精神 校标校训 南邮校史 南邮校歌 现任领导 视频展播 校园实景漫... 校园景色 校区地图 内设机构 \u003e 党政群部门 教学机构 基层党的组......... 科研机构 直属单位和... 独立学院 学科建设 科学研究 \u003e 自然科学研... 社会科学研... 高等教育研... 科210023 三牌楼校区地址：南京市新模范马路66号 邮编：210003 锁金村校区地址：南京市龙蟠路177号 邮编：210042 联系电话:（86）-25-85866888 传真:（86）-25-85866999 邮箱:njupt@njupt.edu.cn 苏公网安备32011302320419号 |苏ICP备11073489号-1 版权所有：南京邮电大学","created_at":"2025-11-10T23:50:34.639+08:00","updated_at":"2025-11-10T23:50:34.639+08:00"}
```

### 4.2 列出页面

```bash
curl -X GET "$BASE_URL/web/items" \
  -H "Authorization: Bearer $TOKEN"
```
```
{"items":[{"id":1,"user_id":1,"url":"https://www.njupt.edu.cn/","title":"Example Home","content":"旧版入口 English 电子邮件 首页 南邮概况 \u003e 学校简介 学校章程 南邮精神 校标校训 南邮校史 南邮校歌 现任领导 视频展播 校园实景漫... 校园景色 校区地图 内设机构 \u003e 党政群部门 教学机构 基层党的组......... 科研机构 直属单位和... 独立学院 学科建设 科学研究 \u003e 自然科学研... 社会科学研... 高等教育研... 科210023 三牌楼校区地址：南京市新模范马路66号 邮编：210003 锁金村校区地址：南京市龙蟠路177号 邮编：210042 联系电话:（86）-25-85866888 传真:（86）-25-85866999 邮箱:njupt@njupt.edu.cn 苏公网安备32011302320419号 |苏ICP备11073489号-1 版权所有：南京邮电大学","created_at":"2025-11-10T23:50:34.639+08:00","updated_at":"2025-11-10T23:50:34.639+08:00"}]}
```
### 4.3 获取页面详情（需 PAGE_ID）

```bash
export PAGE_ID=1
curl -X GET "$BASE_URL/web/page/$PAGE_ID" \
  -H "Authorization: Bearer $TOKEN"
```
```
{"id":1,"user_id":1,"url":"https://www.njupt.edu.cn/","title":"Example Home","content":"旧版入口 English 电子邮件 首页 南邮概况 \u003e 学校简介 学校章程 南邮精神 校标校训 南邮校史 南邮校歌 现任领导 视频展播 校园实景漫... 校园景色 校区地图 内设机构 \u003e 党政群部门 教学机构 基层党的组......... 科研机构 直属单位和... 独立学院 学科建设 科学研究 \u003e 自然科学研... 社会科学研... 高等教育研... 科210023 三牌楼校区地址：南京市新模范马路66号 邮编：210003 锁金村校区地址：南京市龙蟠路177号 邮编：210042 联系电话:（86）-25-85866888 传真:（86）-25-85866999 邮箱:njupt@njupt.edu.cn 苏公网安备32011302320419号 |苏ICP备11073489号-1 版权所有：南京邮电大学","created_at":"2025-11-10T23:50:34.639+08:00","updated_at":"2025-11-10T23:50:34.639+08:00"}
```

### 4.4 更新页面

```bash
curl -X PUT "$BASE_URL/web/items/$PAGE_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Example Updated","content":"Hello world"}'
```
```
{"id":1,"user_id":1,"url":"https://www.njupt.edu.cn/","title":"Example Updated","content":"Hello world","created_at":"2025-11-10T23:50:34.639+08:00","updated_at":"2025-11-10T23:53:11.175+08:00"}
```

### 4.5 搜索（按关键词）

```bash
curl -X POST "$BASE_URL/web/search" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"q":"Example","top_k":5}'
```
```
{"results":[{"id":1,"url":"https://www.njupt.edu.cn/","title":"Example Updated","snippet":"Hello world","score":80}]}
```
### 4.6 搜索（传 URL，会触发抓取）

```bash
curl -X POST "$BASE_URL/web/search" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"urls":["https://www.rfc-editor.org/"],"top_k":3}'
```
```
{"results":[{"id":2,"url":"https://www.rfc-editor.org/","title":"» RFC Editor","snippet":"Search RFCs Advanced Search RFC Editor The RFC Series The RFC Series (ISSN 2070-1721) contains technical and organizational documents about the Internet, including the specifications and policy documents produced by five streams: the Intern","score":100}]}
```
### 4.7 兼容占位：ingest / chunk

```bash
curl -X POST "$BASE_URL/web/ingest" -H "Authorization: Bearer $TOKEN"
curl -X POST "$BASE_URL/web/chunk"  -H "Authorization: Bearer $TOKEN"
```
```
{"message":"Web content processed successfully","success":true}
{"message":"Content chunked successfully","success":true}
```
### 4.8 删除页面

```bash
curl -X DELETE "$BASE_URL/web/items/$PAGE_ID" \
  -H "Authorization: Bearer $TOKEN"
```
```
{"success":true}
```

## 5) Chat 相关

### 5.1 创建会话 → 取 SESSION_ID

```bash
curl -X POST "$BASE_URL/chat/sessions" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Smoke Test Session"}'
```
```
{"id":1,"user_id":1,"title":"Smoke Test Session","created_at":"2025-11-10T23:54:33.506+08:00","updated_at":"2025-11-10T23:54:33.506+08:00"}
```
* 提取 `SESSION_ID`（任选其一）：

```bash
SESSION_ID=$(curl -sS -X POST "$BASE_URL/chat/sessions" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Smoke Test Session"}' | jq -r '.id // .data.id // .session_id // .data.session_id') && echo "$SESSION_ID"
```
```
2
```


### 5.2 列出会话

```bash
curl -X GET "$BASE_URL/chat/sessions" \
  -H "Authorization: Bearer $TOKEN"
```
```
[{"id":2,"user_id":1,"title":"Smoke Test Session","created_at":"2025-11-10T23:54:56+08:00","updated_at":"2025-11-10T23:54:56+08:00"},{"id":1,"user_id":1,"title":"Smoke Test Session","created_at":"2025-11-10T23:54:34+08:00","updated_at":"2025-11-10T23:54:34+08:00"}]
```
### 5.3 追加消息（按会话）

```bash
curl -X POST "$BASE_URL/chat/sessions/$SESSION_ID/messages" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role":"user","content":"Hello from curl"}'
```
```
{"id":1,"user_id":1,"session_id":2,"content":"Hello from curl","role":"user","created_at":"2025-11-10T23:56:10.180238516+08:00"}
```

### 5.4 追加消息（别名 /chat/messages）

```bash
curl -X POST "$BASE_URL/chat/messages" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"session_id\":\"$SESSION_ID\",\"content\":\"Second message via alias\"}"
```
```
{"id":2,"user_id":1,"session_id":2,"content":"Second message via alias","role":"user","created_at":"2025-11-10T23:56:35.103559241+08:00"}
```
### 5.5 拉取消息

```bash
curl -X GET "$BASE_URL/chat/sessions/$SESSION_ID/messages" \
  -H "Authorization: Bearer $TOKEN"
```
```
[{"id":1,"user_id":1,"session_id":2,"content":"Hello from curl","role":"user","created_at":"2025-11-10T23:56:10+08:00"},{"id":2,"user_id":1,"session_id":2,"content":"Second message via alias","role":"user","created_at":"2025-11-10T23:56:35+08:00"}]
```
### 5.6 删除会话
```bash
curl -X DELETE "$BASE_URL/chat/sessions/11" \
  -H "Authorization: Bearer $TOKEN"
```
```bash
{"message":"会话已删除","success":true}
```

### 5.7 一次性补全（可能依赖 LLM 配置，失败忽略）

```bash
curl -X POST "$BASE_URL/chat/sessions/$SESSION_ID/complete" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content":"Briefly echo hello","model":""}'
```
```
{"content":"Hello!"}
```

### 5.8 流式补全（原样打印流）

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
```
{"delta": "H"}
{"delta": "-e"}
{"delta": "-l"}
{"delta": "-l"}
{"delta": "-o-!"}
{"event": "done"}
{"delta": "你好"}
{"delta": "！有什么"}
{"delta": "我可以帮助"}
{"delta": "你的吗？"}
{"event": "done"}
```
## 6) 危险操作：删除账号（默认不执行,暂不可用）

```bash
curl -X DELETE "$BASE_URL/api/auth/delete" \
  -H "Authorization: Bearer $TOKEN"
```
```
404 page not found
```