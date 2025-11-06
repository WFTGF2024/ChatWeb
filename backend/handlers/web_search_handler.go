package handlers

import (
	"backend/database"
	"backend/models"
	"backend/services"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type searchReq struct {
	Q     string   `json:"q"`        // 关键词（可选）
	URLs  []string `json:"urls"`     // 直接抓取这些 URL（可选）
	Fetch bool     `json:"fetch"`    // 是否抓取外网（默认 true）
	TopK  int      `json:"top_k"`    // 关键词检索返回上限
}

type pageResp struct {
	ID      uint   `json:"id"`
	URL     string `json:"url"`
	Title   string `json:"title"`
	Content string `json:"content,omitempty"`
	Snippet string `json:"snippet,omitempty"`
	Score   int    `json:"score,omitempty"`
}

func currentUserID(c *gin.Context) (uint, error) {
    // 依次尝试常见键
    keys := []string{"user_id", "uid", "sub"}
    for _, k := range keys {
        if v, ok := c.Get(k); ok {
            if u := toUint(v); u > 0 { return u, nil }
        }
    }
    if m, ok := c.Get("claims"); ok {
        switch mm := m.(type) {
        case map[string]any:
            for _, k := range []string{"user_id","uid","sub"} {
                if v, ok := mm[k]; ok {
                    if u := toUint(v); u > 0 { return u, nil }
                }
            }
        }
    }
    return 0, errors.New("missing user")
}
func toUint(v any) uint {
    switch t := v.(type) {
    case uint: return t
    case int: return uint(t)
    case int64: return uint(t)
    case float64: return uint(t)
    default: return 0
    }
}

// POST /web/search
// 1) 若传 URLs：抓取并入库（仅当前用户可见）
// 2) 若传 Q：在“当前用户的库”里按标题/正文 LIKE 检索
func HandleWebSearch(c *gin.Context) {
	var req searchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad payload"})
		return
	}
	if req.TopK <= 0 {
		req.TopK = 10
	}
	uid, err := currentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	ws := services.NewWebService()
	db := database.DB

	var results []pageResp

	// 抓取并 upsert URL（仅作用于当前用户空间）
	for _, u := range req.URLs {
		u = strings.TrimSpace(u)
		if u == "" {
			continue
		}
		title, content, err := ws.FetchAndParse(u)
		if err != nil {
			continue
		}
		var page models.WebPage
		err = db.Where("user_id=? AND url=?", uid, u).First(&page).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			page = models.WebPage{UserID: uid, URL: u, Title: title, Content: content}
			if err := db.Create(&page).Error; err != nil {
				continue
			}
		} else if err == nil {
			page.Title = title
			page.Content = content
			_ = db.Save(&page).Error
		}
		results = append(results, pageResp{
			ID: page.ID, URL: page.URL, Title: page.Title,
			Snippet: services.Snippet(page.Content, req.Q, 240),
			Score:   100,
		})
	}

	// 关键词检索：仅在“当前用户”的库里搜
	if strings.TrimSpace(req.Q) != "" {
		var pages []models.WebPage
		q := "%" + strings.TrimSpace(req.Q) + "%"
		db.Where("user_id=? AND (title LIKE ? OR content LIKE ?)", uid, q, q).
			Limit(req.TopK).
			Find(&pages)
		for _, p := range pages {
			results = append(results, pageResp{
				ID: p.ID, URL: p.URL, Title: p.Title,
				Snippet: services.Snippet(p.Content, req.Q, 240),
				Score:   80,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}

// GET /web/page/:id  （只看自己的）
func HandleGetPage(c *gin.Context) {
	uid, err := currentUserID(c)
	if err != nil { c.JSON(http.StatusUnauthorized, gin.H{"error":"unauthorized"}); return }
	id, _ := strconv.Atoi(c.Param("id"))
	var p models.WebPage
	if err := database.DB.Where("user_id=? AND id=?", uid, id).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error":"not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

/************** CRUD：只在用户自己空间 **************/

// GET /web/items?offset=&limit=&q=
func ListPages(c *gin.Context) {
	uid, err := currentUserID(c)
	if err != nil { c.JSON(http.StatusUnauthorized, gin.H{"error":"unauthorized"}); return }
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	q := strings.TrimSpace(c.Query("q"))

	var pages []models.WebPage
	db := database.DB
	tx := db.Where("user_id=?", uid)
	if q != "" {
		like := "%" + q + "%"
		tx = tx.Where("title LIKE ? OR content LIKE ?", like, like)
	}
	if err := tx.Order("updated_at DESC").Limit(limit).Offset(offset).Find(&pages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"db error"}); return
	}
	c.JSON(http.StatusOK, gin.H{"items": pages})
}

type upsertReq struct {
	URL     string `json:"url"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Fetch   bool   `json:"fetch"` // 如果只给 URL 可选抓取
}

// POST /web/items   （创建/抓取一个页面）
func CreatePage(c *gin.Context) {
	uid, err := currentUserID(c)
	if err != nil { c.JSON(http.StatusUnauthorized, gin.H{"error":"unauthorized"}); return }

	var req upsertReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":"bad payload"}); return
	}

	title, content := req.Title, req.Content
	if req.Fetch && req.URL != "" && content == "" {
		ws := services.NewWebService()
		t, cnt, err := ws.FetchAndParse(req.URL)
		if err == nil {
			if title == "" { title = t }
			content = cnt
		}
	}

	// upsert (user_id, url)
	var p models.WebPage
	err = database.DB.Where("user_id=? AND url=?", uid, req.URL).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		p = models.WebPage{UserID: uid, URL: req.URL, Title: title, Content: content}
		if err := database.DB.Create(&p).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"db error"}); return
		}
	} else if err == nil {
		if title != "" { p.Title = title }
		if content != "" { p.Content = content }
		_ = database.DB.Save(&p).Error
	}
	c.JSON(http.StatusOK, p)
}

// PUT /web/items/:id
func UpdatePage(c *gin.Context) {
	uid, err := currentUserID(c)
	if err != nil { c.JSON(http.StatusUnauthorized, gin.H{"error":"unauthorized"}); return }
	id, _ := strconv.Atoi(c.Param("id"))
	var req upsertReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":"bad payload"}); return
	}
	var p models.WebPage
	if err := database.DB.Where("user_id=? AND id=?", uid, id).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error":"not found"}); return
	}
	if strings.TrimSpace(req.Title) != "" { p.Title = req.Title }
	if strings.TrimSpace(req.Content) != "" { p.Content = req.Content }
	if err := database.DB.Save(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"db error"}); return
	}
	c.JSON(http.StatusOK, p)
}

// DELETE /web/items/:id
func DeletePage(c *gin.Context) {
	uid, err := currentUserID(c)
	if err != nil { c.JSON(http.StatusUnauthorized, gin.H{"error":"unauthorized"}); return }
	id, _ := strconv.Atoi(c.Param("id"))
	if err := database.DB.Where("user_id=? AND id=?", uid, id).Delete(&models.WebPage{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"db error"}); return
	}
	// 顺手删掉块
	_ = database.DB.Where("user_id=? AND web_page_id=?", uid, id).Delete(&models.ContentChunk{}).Error
	c.JSON(http.StatusOK, gin.H{"success": true})
}
