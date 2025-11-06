package handlers

import (
	"net/http"
	"time"

	"backend/database"
	"backend/models"
	"backend/services"
	"backend/utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// 初始化用户服务
var userService = services.NewUserService()

// Register 用户注册
func Register(c *gin.Context) {
	log.Info("开始处理用户注册请求")

	var req models.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Error("解析注册请求失败")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	log.WithFields(log.Fields{
		"username":    req.Username,
		"email":       req.Email,
		"phoneNumber": req.PhoneNumber,
	}).Debug("用户注册请求参数")

	// 调用用户服务处理注册逻辑
	user, err := userService.RegisterUser(req)
	if err != nil {
		log.WithError(err).Error("用户注册失败")
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": err.Error()})
		return
	}

	log.WithField("userID", user.UserID).Info("用户注册成功")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user_id": user.UserID,
		"message": "注册成功",
	})
}

// Login 用户登录
func Login(c *gin.Context) {
	log.Info("开始处理用户登录请求")

	var req models.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Error("解析登录请求失败")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	log.WithField("username", req.Username).Debug("用户登录请求参数")

	// 调用用户服务处理登录逻辑
	token, err := userService.LoginUser(req)
	if err != nil {
		log.WithError(err).Error("用户登录失败")
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
		return
	}

	// 计算令牌过期时间
	expireAt := time.Now().Add(24 * time.Hour).Format(time.RFC3339)

	log.WithField("username", req.Username).Info("用户登录成功")

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"token":     token,
		"expire_at": expireAt,
	})
}

// GetProfile 获取用户信息
func GetProfile(c *gin.Context) {
	log.Info("开始处理获取用户信息请求")

	// 从JWT令牌中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		log.Warn("未授权访问，缺少user_id")
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "未授权访问"})
		return
	}

	log.WithField("userID", userID).Debug("获取用户信息请求参数")

	// 调用用户服务获取用户信息
	userProfile, err := userService.GetUserProfile(userID.(uint))
	if err != nil {
		log.WithError(err).WithField("userID", userID).Error("获取用户信息失败")
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户不存在"})
		return
	}

	log.WithField("userID", userID).Info("成功获取用户信息")
	c.JSON(http.StatusOK, userProfile)
}

// UpdateUser 更新用户信息
func UpdateUser(c *gin.Context) {
	log.Info("开始处理更新用户信息请求")

	// 从JWT令牌中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		log.Warn("未授权访问，缺少user_id")
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "未授权访问"})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Error("解析更新用户信息请求失败")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	log.WithFields(log.Fields{
		"userID":      userID,
		"fullName":    req.FullName,
		"email":       req.Email,
		"phoneNumber": req.PhoneNumber,
	}).Debug("更新用户信息请求参数")

	// 调用用户服务更新用户信息
	if err := userService.UpdateUser(userID.(uint), req); err != nil {
		log.WithError(err).WithField("userID", userID).Error("更新用户信息失败")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	log.WithField("userID", userID).Info("用户信息更新成功")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "用户信息已更新",
	})
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	log.Info("开始处理删除用户请求")

	// 从JWT令牌中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		log.Warn("未授权访问，缺少user_id")
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "未授权访问"})
		return
	}

	log.WithField("userID", userID).Debug("删除用户请求参数")

	// 调用用户服务删除用户
	if err := userService.DeleteUser(userID.(uint)); err != nil {
		log.WithError(err).WithField("userID", userID).Error("删除用户失败")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除失败"})
		return
	}

	log.WithField("userID", userID).Info("用户删除成功")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "用户已删除",
	})
}

// VerifySecurity 验证密保问题
func VerifySecurity(c *gin.Context) {
	log.Info("开始处理验证密保问题请求")

	var req models.SecurityVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Error("解析验证密保问题请求失败")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	log.WithField("username", req.Username).Debug("验证密保问题请求参数")

	// 1. 调用用户服务验证密保问题（这一步只是判断答案对不对）
	valid, err := userService.VerifySecurity(req)
	if err != nil {
		log.WithError(err).Error("验证密保问题失败")
		// 基本是业务错误，比如用户不存在
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if !valid {
		log.WithField("username", req.Username).Warn("密保答案验证失败")
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "密保答案错误"})
		return
	}

	// 2. 答案正确，生成重置令牌
	resetToken := utils.GenerateResetToken()
	expiresAt := time.Now().Add(15 * time.Minute)

	// 3. 把令牌写回这个用户，跟 ResetPassword 的查法对上
	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		log.WithError(err).Error("密保验证成功但查询用户失败")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户不存在"})
		return
	}

	if err := database.DB.Model(&user).Updates(map[string]interface{}{
		"reset_token":            resetToken,
		"reset_token_expires_at": expiresAt,
	}).Error; err != nil {
		log.WithError(err).Error("保存重置令牌失败")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "保存重置令牌失败"})
		return
	}

	log.WithFields(log.Fields{
		"username":   req.Username,
		"resetToken": resetToken,
	}).Info("密保问题验证成功并生成重置令牌")

	// 4. 返回给前端 / pytest
	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"reset_token": resetToken,
	})
}

// ResetPassword 重置密码
func ResetPassword(c *gin.Context) {
	log.Info("开始处理重置密码请求")

	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Error("解析重置密码请求失败")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "参数错误",
		})
		return
	}

	log.WithField("resetToken", req.ResetToken).Debug("重置密码请求参数")

	// 调用用户服务重置密码
	err := userService.ResetPassword(req)
	if err != nil {
		msg := err.Error()
		switch msg {
		case "无效的重置令牌", "重置令牌已过期":
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": msg,
			})
		case "密码处理失败", "更新密码失败":
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "密码重置失败",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "服务器内部错误",
			})
		}
		return
	}

	log.Info("密码重置成功")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "密码已更新",
	})
}
