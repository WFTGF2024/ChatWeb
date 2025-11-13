package services

import (
	"errors"
	"fmt"
	"time"

	"backend/database"
	"backend/models"

	log "github.com/sirupsen/logrus"
)

/* ---------- 小工具：解析 & 格式化日期 ---------- */

// 允许两种输入：RFC3339（2025-11-12T08:00:00Z）或日期（YYYY-MM-DD）
func parseAnyDate(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, errors.New("empty date")
	}
	// RFC3339
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	// 仅日期
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}
	// 常见本地时间
	if t, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("invalid date: %q (expect RFC3339 or YYYY-MM-DD)", s)
}

// 如果你的表是 DATE 类型，建议只返回 YYYY-MM-DD；若为 DATETIME 可改成 RFC3339
func fmtDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
	// 若需要完整时间戳：
	// return t.Format(time.RFC3339)
}

/* ---------- 接口定义 ---------- */

type MembershipService interface {
	GetMembershipInfo(userID uint) (*models.MembershipResponse, error)
	GetAllMemberships() ([]models.MembershipResponse, error)
	CreateMembership(req models.CreateMembershipRequest) (*models.MembershipResponse, error)
	UpdateMembership(membershipID uint, req models.UpdateMembershipRequest) error
	DeleteMembership(membershipID uint) error
	GetMembershipOrders(userID uint) ([]models.OrderResponse, error)
	CreateOrder(req models.CreateOrderRequest) (*models.OrderResponse, error)
	GetLatestOrder(userID uint) (*models.OrderResponse, error)
	GetRecentOrders(userID uint, n int) ([]models.OrderResponse, error)
}

type membershipService struct{}

func NewMembershipService() MembershipService { return &membershipService{} }

/* ---------- 会员 ---------- */

func (s *membershipService) GetMembershipInfo(userID uint) (*models.MembershipResponse, error) {
	log.Info("开始获取会员信息")
	log.WithField("userID", userID).Debug("获取会员信息请求参数")

	var membership models.MembershipInfo
	if err := database.DB.Where("user_id = ?", userID).First(&membership).Error; err != nil {
		log.WithError(err).WithField("userID", userID).Warn("会员信息不存在")
		return nil, errors.New("会员信息不存在")
	}

	resp := models.MembershipResponse{
		MembershipID: membership.MembershipID,
		UserID:       membership.UserID,
		StartDate:    fmtDate(membership.StartDate),
		ExpireDate:   fmtDate(membership.ExpireDate),
		Status:       membership.Status,
	}
	log.WithField("userID", userID).Info("获取会员信息成功")
	return &resp, nil
}

func (s *membershipService) GetAllMemberships() ([]models.MembershipResponse, error) {
	log.Info("开始获取所有会员信息")

	var memberships []models.MembershipInfo
	if err := database.DB.Find(&memberships).Error; err != nil {
		log.WithError(err).Error("查询会员信息失败")
		return nil, errors.New("查询会员信息失败")
	}
	if len(memberships) == 0 {
		return []models.MembershipResponse{}, nil
	}

	out := make([]models.MembershipResponse, len(memberships))
	for i, m := range memberships {
		out[i] = models.MembershipResponse{
			MembershipID: m.MembershipID,
			UserID:       m.UserID,
			StartDate:    fmtDate(m.StartDate),
			ExpireDate:   fmtDate(m.ExpireDate),
			Status:       m.Status,
		}
	}
	log.Info("获取所有会员信息成功")
	return out, nil
}

func (s *membershipService) CreateMembership(req models.CreateMembershipRequest) (*models.MembershipResponse, error) {
    log.Info("开始创建会员信息")
    log.WithFields(log.Fields{
        "userID":     req.UserID,
        "startDate":  req.StartDate,
        "expireDate": req.ExpireDate,
        "status":     req.Status,
    }).Debug("创建会员信息请求参数")

    // 用户存在校验
    var user models.User
    if err := database.DB.First(&user, req.UserID).Error; err != nil {
        log.WithError(err).WithField("userID", req.UserID).Warn("用户不存在")
        return nil, errors.New("用户不存在")
    }

    // 检查用户是否已有会员记录
    var existing models.MembershipInfo
    if err := database.DB.Where("user_id = ?", req.UserID).First(&existing).Error; err == nil {
        // 用户已有会员记录，延长1个月有效期
        log.WithField("userID", req.UserID).Info("用户已有会员信息，将延长1个月有效期")

        // 计算新的有效期，从现有有效期开始延长1个月
        var newExpireDate time.Time
        // 如果现有会员还未过期，从现有过期日期开始延长1个月
        if existing.ExpireDate.After(time.Now()) {
            newExpireDate = existing.ExpireDate.AddDate(0, 1, 0) // 加1个月
        } else {
            // 如果现有会员已过期，从当前时间开始延长1个月
            newExpireDate = time.Now().AddDate(0, 1, 0) // 从当前时间加1个月
        }

        // 更新会员信息
        updateData := map[string]interface{}{
            "expire_date": newExpireDate,
            "status":      "active",
        }

        if err := database.DB.Model(&existing).Updates(updateData).Error; err != nil {
            log.WithError(err).Error("延长会员有效期失败")
            return nil, errors.New("延长会员有效期失败")
        }

        // 返回更新后的会员信息
        resp := models.MembershipResponse{
            MembershipID: existing.MembershipID,
            UserID:       existing.UserID,
            StartDate:    fmtDate(existing.StartDate),
            ExpireDate:   fmtDate(newExpireDate),
            Status:       "active",
        }

        log.WithField("membershipID", existing.MembershipID).Info("延长会员有效期成功")
        return &resp, nil
    }

    // 解析日期字符串为 time.Time
    start, err := parseAnyDate(req.StartDate)
    if err != nil {
        return nil, fmt.Errorf("start_date 解析失败: %w", err)
    }
    expire, err := parseAnyDate(req.ExpireDate)
    if err != nil {
        return nil, fmt.Errorf("expire_date 解析失败: %w", err)
    }

    newMembership := models.MembershipInfo{
        UserID:     req.UserID,
        StartDate:  start,
        ExpireDate: expire,
        Status:     req.Status,
    }

    if err := database.DB.Create(&newMembership).Error; err != nil {
        log.WithError(err).Error("创建会员信息失败")
        return nil, errors.New("创建会员信息失败")
    }

    resp := models.MembershipResponse{
        MembershipID: newMembership.MembershipID,
        UserID:       newMembership.UserID,
        StartDate:    fmtDate(newMembership.StartDate),
        ExpireDate:   fmtDate(newMembership.ExpireDate),
        Status:       newMembership.Status,
    }
    log.WithField("membershipID", newMembership.MembershipID).Info("创建会员信息成功")
    return &resp, nil
}

func (s *membershipService) UpdateMembership(membershipID uint, req models.UpdateMembershipRequest) error {
	log.Info("开始更新会员信息")
	log.WithFields(log.Fields{
		"membershipID": membershipID,
		"expireDate":   req.ExpireDate,
		"status":       req.Status,
	}).Debug("更新会员信息请求参数")

	var membership models.MembershipInfo
	if err := database.DB.First(&membership, membershipID).Error; err != nil {
		log.WithError(err).WithField("membershipID", membershipID).Warn("会员信息不存在")
		return errors.New("会员信息不存在")
	}

	updateData := make(map[string]interface{})
	if req.ExpireDate != "" {
		if t, err := parseAnyDate(req.ExpireDate); err == nil {
			updateData["expire_date"] = t
		} else {
			return fmt.Errorf("expire_date 解析失败: %w", err)
		}
	}
	if req.Status != "" {
		updateData["status"] = req.Status
	}

	if len(updateData) > 0 {
		if err := database.DB.Model(&membership).Updates(updateData).Error; err != nil {
			log.WithError(err).Error("更新会员信息失败")
			return errors.New("更新会员信息失败")
		}
	}

	log.WithField("membershipID", membershipID).Info("更新会员信息成功")
	return nil
}

func (s *membershipService) DeleteMembership(membershipID uint) error {
	log.Info("开始删除会员信息")
	log.WithField("membershipID", membershipID).Debug("删除会员信息请求参数")

	var membership models.MembershipInfo
	if err := database.DB.First(&membership, membershipID).Error; err != nil {
		log.WithError(err).WithField("membershipID", membershipID).Warn("会员信息不存在")
		return errors.New("会员信息不存在")
	}

	if err := database.DB.Delete(&membership).Error; err != nil {
		log.WithError(err).Error("删除会员信息失败")
		return errors.New("删除会员信息失败")
	}

	log.WithField("membershipID", membershipID).Info("删除会员信息成功")
	return nil
}

/* ---------- 订单 ---------- */

func (s *membershipService) GetMembershipOrders(userID uint) ([]models.OrderResponse, error) {
	log.Info("开始获取会员订单记录")
	log.WithField("userID", userID).Debug("获取会员订单记录请求参数")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.WithError(err).WithField("userID", userID).Warn("用户不存在")
		return nil, errors.New("用户不存在")
	}

	var orders []models.MembershipOrder
	if err := database.DB.Where("user_id = ?", userID).Order("purchase_date DESC").Find(&orders).Error; err != nil {
		log.WithError(err).Error("查询会员订单失败")
		return nil, errors.New("查询会员订单失败")
	}

	res := make([]models.OrderResponse, 0, len(orders))
	for _, o := range orders {
		res = append(res, models.OrderResponse{
			OrderID:        o.OrderID,
			UserID:         o.UserID,
			PurchaseDate:   o.PurchaseDate,
			DurationMonths: o.DurationMonths,
			Amount:         o.Amount,
			PaymentMethod:  o.PaymentMethod,
		})
	}
	log.WithField("userID", userID).Info("获取会员订单记录成功")
	return res, nil
}

func (s *membershipService) CreateOrder(req models.CreateOrderRequest) (*models.OrderResponse, error) {
	log.Info("开始创建订单")
	log.WithFields(log.Fields{
		"userID":         req.UserID,
		"durationMonths": req.DurationMonths,
		"amount":         req.Amount,
		"paymentMethod":  req.PaymentMethod,
	}).Debug("创建订单请求参数")

	var user models.User
	if err := database.DB.First(&user, req.UserID).Error; err != nil {
		log.WithError(err).WithField("userID", req.UserID).Warn("用户不存在")
		return nil, errors.New("用户不存在")
	}

	newOrder := models.MembershipOrder{
		UserID:         req.UserID,
		PurchaseDate:   time.Now(),
		DurationMonths: req.DurationMonths,
		Amount:         req.Amount,
		PaymentMethod:  req.PaymentMethod,
	}
	if err := database.DB.Create(&newOrder).Error; err != nil {
		log.WithError(err).Error("创建订单失败")
		return nil, errors.New("创建订单失败")
	}

	resp := models.OrderResponse{
		OrderID:        newOrder.OrderID,
		UserID:         newOrder.UserID,
		PurchaseDate:   newOrder.PurchaseDate,
		DurationMonths: newOrder.DurationMonths,
		Amount:         newOrder.Amount,
		PaymentMethod:  newOrder.PaymentMethod,
	}
	log.WithField("orderID", newOrder.OrderID).Info("创建订单成功")
	return &resp, nil
}

func (s *membershipService) GetLatestOrder(userID uint) (*models.OrderResponse, error) {
	log.Info("开始获取最近一条订单")
	log.WithField("userID", userID).Debug("获取最近一条订单请求参数")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.WithError(err).WithField("userID", userID).Warn("用户不存在")
		return nil, errors.New("用户不存在")
	}

	var order models.MembershipOrder
	if err := database.DB.Where("user_id = ?", userID).Order("purchase_date DESC").First(&order).Error; err != nil {
		log.WithError(err).WithField("userID", userID).Warn("用户没有订单记录")
		return nil, errors.New("用户没有订单记录")
	}

	resp := models.OrderResponse{
		OrderID:        order.OrderID,
		UserID:         order.UserID,
		PurchaseDate:   order.PurchaseDate,
		DurationMonths: order.DurationMonths,
		Amount:         order.Amount,
		PaymentMethod:  order.PaymentMethod,
	}
	log.WithField("userID", userID).Info("获取最近一条订单成功")
	return &resp, nil
}

func (s *membershipService) GetRecentOrders(userID uint, n int) ([]models.OrderResponse, error) {
	log.Info("开始获取最近N条订单")
	log.WithFields(log.Fields{
		"userID": userID,
		"n":      n,
	}).Debug("获取最近N条订单请求参数")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.WithError(err).WithField("userID", userID).Warn("用户不存在")
		return nil, errors.New("用户不存在")
	}

	var orders []models.MembershipOrder
	if err := database.DB.Where("user_id = ?", userID).Order("purchase_date DESC").Limit(n).Find(&orders).Error; err != nil {
		log.WithError(err).Error("查询订单失败")
		return nil, errors.New("查询订单失败")
	}

	res := make([]models.OrderResponse, 0, len(orders))
	for _, o := range orders {
		res = append(res, models.OrderResponse{
			OrderID:        o.OrderID,
			UserID:         o.UserID,
			PurchaseDate:   o.PurchaseDate,
			DurationMonths: o.DurationMonths,
			Amount:         o.Amount,
			PaymentMethod:  o.PaymentMethod,
		})
	}
	log.WithFields(log.Fields{
		"userID": userID,
		"n":      n,
	}).Info("获取最近N条订单成功")
	return res, nil
}
