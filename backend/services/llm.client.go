package services

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"strings"
)

type LLMMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type llmOnceBody struct {
	Messages []LLMMessage `json:"messages"`
	Model    string       `json:"model,omitempty"`
}

type llmOnceResp struct {
	Content string `json:"content"`
}

type LLMClient struct {
	base   string
	client *http.Client
}

func NewLLMClient() *LLMClient {
	base := os.Getenv("LLM_BASE")
	if base == "" {
		base = "http://127.0.0.1:7207"
	}
	return &LLMClient{
		base: base,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (c *LLMClient) ChatOnce(ctx context.Context, messages []LLMMessage, model string) (string, error) {
	body := llmOnceBody{Messages: messages}
	if model != "" {
		body.Model = model
	}
	bs, _ := json.Marshal(body)
	req, _ := http.NewRequestWithContext(ctx, "POST", c.base+"/api/chat", bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		raw, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("LLM /api/chat failed: %s", string(raw))
	}
	var out llmOnceResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.Content, nil
}

// ChatStream 以行分隔 JSON 事件的方式透传，并把最终文本返回给回调保存
func (c *LLMClient) ChatStream(ctx context.Context, w io.Writer, messages []LLMMessage, model string, onDone func(fullText string)) error {
	body := llmOnceBody{Messages: messages}
	if model != "" {
		body.Model = model
	}
	bs, _ := json.Marshal(body)
	req, _ := http.NewRequestWithContext(ctx, "POST", c.base+"/api/chat/stream", bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		raw, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("LLM /api/chat/stream failed: %s", string(raw))
	}
	sc := bufio.NewScanner(resp.Body)
	full := strings.Builder{}
	for sc.Scan() {
		line := sc.Bytes()
		w.Write(line)
		w.Write([]byte("\n"))
		// 累积 delta
		var evt map[string]any
		if err := json.Unmarshal(line, &evt); err == nil {
			if d, ok := evt["delta"].(string); ok {
				full.WriteString(d)
			}
		}
	}
	if sc.Err() != nil {
		// 流中断就不保存
		return sc.Err()
	}
	if onDone != nil {
		onDone(full.String())
	}
	return nil
}
