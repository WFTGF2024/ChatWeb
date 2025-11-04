# 项目结构
```
backend/
├── config/           # 配置文件
│   └── config.go     # 配置加载和管理，使用全局变量存储配置
├── database/         # 数据库连接和初始化
│   └── database.go   # 数据库连接和模型迁移，使用logrus记录日志
├── docs/             # 项目文档
│   ├── API.md        # API接口文档
│   ├── API_Develop.md # API开发文档
│   └── MySQL.sql     # 数据库初始化脚本
├── handlers/         # HTTP请求处理层
│   └── user_handler.go # 用户相关HTTP处理函数
├── middleware/       # 中间件
│   └── auth.go       # JWT认证中间件
├── models/           # 数据模型层
│   └── user.go       # 用户模型定义
├── services/         # 业务逻辑层
│   └── user_service.go # 用户相关业务逻辑
├── utils/            # 工具函数
│   └── utils.go      # 密码哈希、JWT等工具函数
└── main.go           # 主程序入口，负责初始化配置、日志系统和数据库连接
```

# 接口的请求和回传格式，curl命令
[具体内容请点击这个](./api.md)
[curl命令请点击这个](./curl.md)
[curl命令及结果点击这个](./curl_and_response.md)

# 运行项目

1. 安装 `Go>1.20`,并把Go应用的bin文件夹的绝对路径加入到你的系统Path中
2. 确保 `MySQL` 数据库已启动并创建 `Qiniu_Project` 数据库
3. 根据实际情况修改 `config.yaml` 中的数据库连接信息和服务器设置
在 `backend`目录运行以下命令
```bash
go mod init backend
go mod tidy
go run main.go -r ./config.yaml
```
config.yaml若未给出，可以按下列格式新建给出：
```yaml
mysql:
  host: localhost
  port: 3306
  user: root
  password: yourpassword
  database: Qiniu_Project
app_log_file: logs/app.log
server_port: 8080
```