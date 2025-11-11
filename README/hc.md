根据main.go介绍一下chat和web的路由。

根据 `model` 层的 `chat.go`,`web.go`讲述使用到的模型（表），与MySQL表一一对应。原先MySQL初始化会建表，这边go有迁移机制，如果没有的话就会创建。

根据 `service`层的 `chat_service.go` 和 `web_service.go`大致讲解代码

根据 `handler`层的 `chat_handler.go` 和 `web_handler.go`和 `web_search_handler.go`大致讲解代码，并说明`web_handler.go`是为后续功能预留的接口（由于显存不能跑更多模型，没有实现）

测试端口,按下面的演示
## 在cmd设置一下环境变量
```cmd
set BASE_URL=http://120.79.25.184:5000
set USERNAME=hc
set PASSWORD=hc@12345
set NEW_PASSWORD=hc@12345_new
set FULL_NAME=hc
set EMAIL=hc@example.com
set PHONE=13900001235
set SQ1=where are you from?
set SA1=china
set SQ2=what's you professor?
set SA2=softengineering
```
## 注册（拿到 JWT）
```cmd
curl -X POST "%BASE_URL%/api/auth/register" ^
More?   -H "Content-Type: application/json" ^
More?   -d "{\"username\":\"%USERNAME%\",\"password\":\"%PASSWORD%\",\"full_name\":\"%FULL_NAME%\",\"email\":\"%EMAIL%\",\"phone_number\":\"%PHONE%\",\"security_question1\":\"%SQ1%\",\"security_answer1\":\"%SA1%\",\"security_question2\":\"%SQ2%\",\"security_answer2\":\"%SA2%\"}"
```
## 登录（拿到 JWT）
```cmd
curl -sS -X POST "%BASE_URL%/api/auth/login" ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"%USERNAME%\",\"password\":\"%PASSWORD%\"}"
```
```cmd
for /f "usebackq delims=" %T in (`powershell -NoProfile -Command "$r = Invoke-RestMethod -Method Post -Uri '%BASE_URL%/api/auth/login' -ContentType 'application/json' -Body '{\"username\":\"%USERNAME%\",\"password\":\"%PASSWORD%\"}'; $r.token"`) do set TOKEN=%T
```
```cmd
echo %TOKEN%
```
## 创建页面
```cmd
curl -X POST "%BASE_URL%/web/items" ^
  -H "Authorization: Bearer %TOKEN%" ^
  -H "Content-Type: application/json" ^
  -d "{\"url\":\"https://leetcode.cn/\",\"title\":\"Example Home\",\"fetch\":true}"
```

## 拿Page_id
```cmd
for /f "usebackq delims=" %P in (`powershell -NoProfile -Command "$r = Invoke-RestMethod -Method Post -Uri '%BASE_URL%/web/items' -Headers @{Authorization='Bearer %TOKEN%'} -ContentType 'application/json' -Body '{\"url\":\"https://leetcode.cn/\",\"title\":\"Example Home\",\"fetch\":true}'; $r.id"`) do set PAGE_ID=%P
echo %PAGE_ID%
```

## 列出页面
```cmd
curl -X GET "%BASE_URL%/web/items" -H "Authorization: Bearer %TOKEN%"
```

## 获取页面详情
```cmd
curl -X GET "%BASE_URL%/web/page/%PAGE_ID%" -H "Authorization: Bearer %TOKEN%"
```

## 更新页面
```cmd
curl -X PUT "%BASE_URL%/web/items/%PAGE_ID%" ^
  -H "Authorization: Bearer %TOKEN%" ^
  -H "Content-Type: application/json" ^
  -d "{\"title\":\"Example Updated\",\"content\":\"Hello world\"}"
```

## 搜索（按关键词）
```cmd
curl -X POST "%BASE_URL%/web/search" ^
  -H "Authorization: Bearer %TOKEN%" ^
  -H "Content-Type: application/json" ^
  -d "{\"q\":\"Example\",\"top_k\":5}"
```

## 搜索（传URL，会触发抓取）
```cmd
curl -X POST "%BASE_URL%/web/search" ^
  -H "Authorization: Bearer %TOKEN%" ^
  -H "Content-Type: application/json" ^
  -d "{\"urls\":[\"https://www.rfc-editor.org/\"],\"top_k\":3}"
```

## ingest / chunk 测试
```cmd
curl -X POST "%BASE_URL%/web/ingest" -H "Authorization: Bearer %TOKEN%"
curl -X POST "%BASE_URL%/web/chunk"  -H "Authorization: Bearer %TOKEN%"
```

## 删除页面
```cmd
curl -X DELETE "%BASE_URL%/web/items/%PAGE_ID%" -H "Authorization: Bearer %TOKEN%"
```

## 创建会话
```cmd
curl -X POST "%BASE_URL%/chat/sessions" ^
  -H "Authorization: Bearer %TOKEN%" ^
  -H "Content-Type: application/json" ^
  -d "{\"title\":\"Smoke Test Session\"}"
```

## 拿Session_id
```cmd
for /f "usebackq delims=" %S in (`powershell -NoProfile -Command "$r = Invoke-RestMethod -Method Post -Uri '%BASE_URL%/chat/sessions' -Headers @{Authorization='Bearer %TOKEN%'} -ContentType 'application/json' -Body '{\"title\":\"Smoke Test Session\"}'; $r.id"`) do set SESSION_ID=%S
echo %SESSION_ID%
```

## 列出会话
```cmd
curl -X GET "%BASE_URL%/chat/sessions" -H "Authorization: Bearer %TOKEN%"
```

## 追加消息（按会话）
```cmd
curl -X POST "%BASE_URL%/chat/sessions/%SESSION_ID%/messages" ^
  -H "Authorization: Bearer %TOKEN%" ^
  -H "Content-Type: application/json" ^
  -d "{\"role\":\"user\",\"content\":\"Hello from curl on Windows\"}"
```

## 另一种追加方式
```cmd
curl -X POST "%BASE_URL%/chat/messages" ^
  -H "Authorization: Bearer %TOKEN%" ^
  -H "Content-Type: application/json" ^
  -d "{\"session_id\":\"%SESSION_ID%\",\"content\":\"Second message via alias\"}"
```

## 拉取消息
```cmd
curl -X GET "%BASE_URL%/chat/sessions/%SESSION_ID%/messages" -H "Authorization: Bearer %TOKEN%"
```