## Go 模型把这两个字段屏蔽掉
`gorm:"-"` 的意思就是：**这个字段不要参与 ORM 映射**，也就是：

* AutoMigrate 不会帮你建这个字段
* 普通的 `db.Save(&user)`、`db.Model(&user).Update(...)` 也不会管它

但你后端的 handler 里又真的去执行了：

```sql
UPDATE `users` SET `reset_token` = 'xxx', `reset_token_expires_at` = '...' ...
```

也就是说：**代码层想要这个字段，模型层说我不要这个字段，数据库那边自然就没有这个字段**，然后你运行的时候 MySQL 说：`Unknown column 'reset_token'`，这就是你日志里那一条。

所以第一层原因就是：**模型和实际要用的字段不一致**。
要解决就很直接：把 `gorm:"-"` 去掉，给它正常映射上：

```go
ResetToken          string     `json:"-" gorm:"column:reset_token"`
ResetTokenExpiresAt *time.Time `json:"-" gorm:"column:reset_token_expires_at"`
```

然后让你的程序跑一遍迁移（如果你启了 AutoMigrate），或者自己写 `ALTER TABLE`。



这回不是前端代码写错，而是浏览器在当警察了：CORS。

你现在的情况是：

* 前端：`http://localhost:5173`
* 接口：`http://120.79.25.184:7207/api/chat/stream`
* 浏览器发现“域名+端口不一样”，于是先发了一个 **预检请求**（OPTIONS），结果后端返回里没有
  `Access-Control-Allow-Origin: http://localhost:5173`
  所以浏览器直接拦了。

一句话版：**后端没说“可以让 localhost:5173 来访问我”**，所以浏览器不让你用。

下面给你几个能落地的解法，你挑一个。

---

## 方案 1：在后端开 CORS（根治法，推荐）

你控制这个 `120.79.25.184:7207` 的接口的话，就在后端加 CORS 响应头。不同后端写法不一样，我给你几个常见的。

### 1) Flask

```python
from flask import Flask
from flask_cors import CORS

app = Flask(__name__)
CORS(app, origins=["http://localhost:5173"])
```

或者全开（开发期）：

```python
CORS(app)
```

### 2) FastAPI

```python
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI()

origins = [
    "http://localhost:5173",
]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
```

### 3) Django (django-cors-headers)

```python
INSTALLED_APPS = [
    ...
    "corsheaders",
    ...
]

MIDDLEWARE = [
    "corsheaders.middleware.CorsMiddleware",
    ...
]

CORS_ALLOWED_ORIGINS = [
    "http://localhost:5173",
]
```

重点就是让响应里带上：

```http
Access-Control-Allow-Origin: http://localhost:5173
Access-Control-Allow-Methods: GET,POST,OPTIONS
Access-Control-Allow-Headers: Content-Type,Authorization
```

这样浏览器就不吵了。

---

## 方案 2：用 Vite 开发代理，把跨域“藏”起来

如果这个接口你暂时改不了（比如线上同事写的，或你现在不想动后端），那可以让 Vite 在本地做代理，把跨域变成同源。

在你的项目根目录建/改 `vite.config.js`（如果是 TS 就是 `vite.config.ts`），加一个 `server.proxy`：

```js
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      // 你前端里写的是 /api/chat/stream
      '/api': {
        target: 'http://120.79.25.184:7207',
        changeOrigin: true,
        // 如果后端真实路径就是 /api 开头，就不用 rewrite
        // 如果后端没有 /api 前缀，再写：
        // rewrite: (path) => path.replace(/^\/api/, '')
      }
    }
  }
})
```

然后前端代码里请求写成：

```js
fetch('/api/chat/stream', { ... })
```

这样浏览器觉得你是在请求 `http://localhost:5173/api/...`（同源），Vite 再帮你转发到 `http://120.79.25.184:7207/...`，浏览器就不参与 CORS 审查了。

这个方法对“我只是在本机开发”特别好用。

---

## 方案 3：Nginx 做反向代理

如果你最后要部署到一台机器上，让 Nginx 当入口，也可以在 Nginx 上统一加 CORS：

```nginx
location /api/ {
    proxy_pass http://120.79.25.184:7207;
    add_header Access-Control-Allow-Origin *;
    add_header Access-Control-Allow-Methods GET,POST,OPTIONS;
    add_header Access-Control-Allow-Headers Content-Type,Authorization;
    if ($request_method = OPTIONS) {
        return 204;
    }
}
```

（生产别用 `*`，用具体域名。）

---

## 要点再说一遍

* CORS 是浏览器拦的，不是你写的 JS 错了。
* 真正能“说了算”的是**后端/代理**，前端自己是加不了那个响应头的。
* 你现在这个接口还是返回了内容，但是没带允许的源 → 浏览器直接: ❌
* 本地开发最快的是 **Vite 代理**；长期上线要么后端加 CORS，要么网关/Nginx 加。

你现在最少操作的路径就是：**加 vite.config.js 的 proxy → 前端把 URL 改成 `/api/...` → 重启 `npm run dev`**，再看控制台，就会安静很多。接着我们再看你那个 `/api/chat/stream` 的流式响应要不要做 `ReadableStream` 解析。
