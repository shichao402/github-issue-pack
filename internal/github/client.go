package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseURL    = "https://api.github.com"
	apiVersion = "2022-11-28"
)

// Client GitHub API 客户端
type Client struct {
	token      string
	httpClient *http.Client
}

// NewClient 创建新的 GitHub 客户端
func NewClient(token string) *Client {
	return &Client{
		token: token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// doRequest 执行 HTTP 请求
func (c *Client) doRequest(method, url string, body interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("序列化请求体失败: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", apiVersion)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API 错误 (%d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// Get 发送 GET 请求
func (c *Client) Get(url string) ([]byte, error) {
	return c.doRequest(http.MethodGet, url, nil)
}

// Post 发送 POST 请求
func (c *Client) Post(url string, body interface{}) ([]byte, error) {
	return c.doRequest(http.MethodPost, url, body)
}

// Patch 发送 PATCH 请求
func (c *Client) Patch(url string, body interface{}) ([]byte, error) {
	return c.doRequest(http.MethodPatch, url, body)
}
