package services

import (
	"backend/models"
	"errors"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
)

type WebService struct{}

func NewWebService() *WebService { return &WebService{} }

// 抓取并抽取 <title> 与正文纯文本（去 script/style）
func (s *WebService) FetchAndParse(url string) (title, content string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", "", errors.New("bad status: " + resp.Status)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(b)))
	if err != nil {
		return "", "", err
	}
	title = strings.TrimSpace(doc.Find("title").First().Text())
	// 去掉 script/style
	doc.Find("script,style,noscript").Each(func(i int, s *goquery.Selection) { s.Remove() })
	content = strings.TrimSpace(doc.Find("body").Text())
	// 归一化空白
	content = strings.Join(strings.Fields(content), " ")
	return title, content, nil
}

// 简单分块：按字符数量切
func (s *WebService) ChunkContent(content string, maxChars int) ([]models.ContentChunk, error) {
	if maxChars <= 0 {
		maxChars = 2000
	}
	var res []models.ContentChunk
	runes := []rune(content)
	for i := 0; i < len(runes); i += maxChars {
		end := i + maxChars
		if end > len(runes) {
			end = len(runes)
		}
		res = append(res, models.ContentChunk{Content: string(runes[i:end])})
	}
	return res, nil
}

// 简短摘录（搜索结果摘要）
func Snippet(text, q string, size int) string {
	if size <= 0 {
		size = 240
	}
	text = strings.TrimSpace(strings.ReplaceAll(text, "\n", " "))
	pos := strings.Index(strings.ToLower(text), strings.ToLower(q))
	if pos < 0 {
		if utf8.RuneCountInString(text) <= size {
			return text
		}
		return string([]rune(text)[:size])
	}
	start := pos - size/2
	if start < 0 {
		start = 0
	}
	end := start + size
	if end > len([]rune(text)) {
		end = len([]rune(text))
	}
	return string([]rune(text)[start:end])
}
