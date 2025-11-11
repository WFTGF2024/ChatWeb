package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"

	"backend/config"
	"backend/database"
	"backend/handlers"
	"backend/middleware"
	"backend/models"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

func setupRouter() *gin.Engine {
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建路由
	r := gin.New()

	// 添加中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 设置CORS中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 创建API路由组
	api := r.Group("/api")

	// Web路由组（受 JWT 保护）
	web := r.Group("/web")
	web.Use(middleware.JWTAuthMiddleware())
	{
		web.POST("/search", handlers.HandleWebSearch)
		web.GET("/page/:id", handlers.HandleGetPage)

		web.GET("/items", handlers.ListPages)
		web.POST("/items", handlers.CreatePage)
		web.PUT("/items/:id", handlers.UpdatePage)
		web.DELETE("/items/:id", handlers.DeletePage)

		// 兼容你之前的两个占位接口（只注册一次，避免 panic）
		web.POST("/ingest", handlers.HandleWebIngest)
		web.POST("/chunk", handlers.HandleWebChunk)
	}

	// Chat路由组（受 JWT 保护）——仅在原基础上补齐缺失的路由
	chat := r.Group("/chat")
	chat.Use(middleware.JWTAuthMiddleware())
	{
		// 已有
		chat.POST("/sessions", handlers.HandleCreateSession)
		chat.GET("/sessions/:session_id/messages", handlers.HandleGetSessionMessages)
		chat.POST("/messages", handlers.HandleSaveMessage) // JSON 内带 {session_id, role?, content}

		// 新增：列出我的会话，便于前端获取 session_id 列表
		chat.GET("/sessions", handlers.HandleListSessions)

		// 新增：按路径会话追加消息（与 /chat/messages 二选一皆可用）
		chat.POST("/sessions/:session_id/messages", handlers.HandleAddMessage)
		chat.PUT("/sessions/:session_id", handlers.HandleUpdateSession)
		chat.DELETE("/sessions/:session_id", handlers.HandleDeleteSession)

		// 新增：基于上下文与 LLM 对话（一次性/流式）
		chat.POST("/sessions/:session_id/complete", handlers.HandleLLMCompleteOnce)
		chat.POST("/sessions/:session_id/stream", handlers.HandleLLMStream)
	}

	// 认证相关路由
	auth := api.Group("/auth")
	{
		// 用户注册
		auth.POST("/register", handlers.Register)

		// 用户登录
		auth.POST("/login", handlers.Login)

		// 验证密保问题
		auth.POST("/verify-security", handlers.VerifySecurity)

		// 重置密码
		auth.POST("/reset-password", handlers.ResetPassword)

		// 获取用户信息（需要认证）
		auth.GET("/me", middleware.JWTAuthMiddleware(), handlers.GetProfile)
	}

	// 用户管理路由（需要认证）
	users := api.Group("/users")
	users.Use(middleware.JWTAuthMiddleware())
	{
		// 更新用户信息
		users.PUT("/:user_id", handlers.UpdateUser)

		// 删除用户
		users.DELETE("/:user_id", handlers.DeleteUser)
	}

	// 会员管理路由（需要认证）
	membership := api.Group("/membership")
	membership.Use(middleware.JWTAuthMiddleware())
	{
		// 查询会员信息
		membership.GET("/:user_id", handlers.GetMembershipInfo)

		// 查询所有会员信息
		membership.GET("", handlers.GetAllMemberships)

		// 新增会员信息
		membership.POST("", handlers.CreateMembership)

		// 更新会员信息
		membership.PUT("/:membership_id", handlers.UpdateMembership)

		// 删除会员信息
		membership.DELETE("/:membership_id", handlers.DeleteMembership)
	}

	// 会员订单管理路由（需要认证）
	orders := api.Group("/membership/orders")
	orders.Use(middleware.JWTAuthMiddleware())
	{
		// 查询会员订单记录
		orders.GET("/:user_id", handlers.GetMembershipOrders)

		// 新增订单
		orders.POST("", handlers.CreateOrder)

		// 查询最近一条订单
		orders.GET("/:user_id/latest", handlers.GetLatestOrder)

		// 查询最近N条订单
		orders.GET("/:user_id/recent", handlers.GetRecentOrders)
	}

	return r
}

func main() {
	log.Info("应用程序启动...")

	// 解析命令行参数
	configPath := flag.String("r", "", "配置文件绝对路径 (必填)")
	flag.Parse()

	log.WithField("configPath", *configPath).Info("解析命令行参数")

	// 参数校验
	if *configPath == "" {
		log.Fatal("错误: 必须通过 -r 参数指定配置文件绝对路径")
	}

	// 初始化logrus
	log.Info("初始化日志系统...")
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	// 启动HTTP服务器以支持pprof性能分析
	go func() {
		log.Info("启动pprof性能分析服务器，监听地址: 127.0.0.1:6060")
		if err := http.ListenAndServe("127.0.0.1:6060", nil); err != nil {
			log.WithError(err).Fatal("启动pprof HTTP服务器失败")
		}
	}()

	// 1. 加载配置文件
	log.Info("开始加载配置文件...")
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.WithError(err).Fatal("加载配置文件失败")
	}

	log.WithFields(log.Fields{
		"mysql_host":     cfg.MySQL.Host,
		"mysql_port":     cfg.MySQL.Port,
		"mysql_user":     cfg.MySQL.User,
		"mysql_database": cfg.MySQL.Database,
		"app_log_file":   cfg.AppLogFile,
		"server_port":    cfg.ServerPort,
	}).Info("配置文件加载成功")

	// 设置运行时日志输出到指定的日志文件
	log.Info("设置日志输出到文件...")
	if err := os.MkdirAll(filepath.Dir(cfg.AppLogFile), 0755); err != nil {
		log.WithError(err).Fatal("创建日志目录失败")
	}

	f, err := os.OpenFile(cfg.AppLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.WithError(err).Fatal("打开日志文件失败")
	}
	defer f.Close()

	log.SetOutput(f)
	log.Info("日志系统配置完成")

	// 2. 初始化数据库连接
	log.Info("初始化数据库连接...")
	database.InitDB()
	// 自动迁移（确保加上 UserID 与唯一索引）
	if err := database.DB.AutoMigrate(&models.WebPage{}, &models.ContentChunk{}); err != nil {
		log.WithError(err).Fatal("AutoMigrate failed")
	}

	// 3. 设置路由
	log.Info("设置HTTP路由...")
	router := setupRouter()

	// 设置服务器端口，如果配置中没有设置则使用默认值8080
	serverPort := cfg.ServerPort
	if serverPort == 0 {
		serverPort = 8080
		log.Info("使用默认服务器端口: 8080")
	} else {
		log.WithField("serverPort", serverPort).Info("使用配置的服务器端口")
	}

	// 启动HTTP服务器
	serverAddr := fmt.Sprintf(":%d", serverPort)
	log.WithField("serverAddr", serverAddr).Info("启动HTTP服务器")

	if err := router.Run(serverAddr); err != nil {
		log.WithError(err).Fatal("启动服务器失败")
	}
}
