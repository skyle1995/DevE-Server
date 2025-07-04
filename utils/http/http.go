package http

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Client HTTP客户端结构体
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Headers    map[string]string
}

// Response HTTP响应结构体
type Response struct {
	StatusCode int
	Body       []byte
	Header     http.Header
	Error      error
}

// NewClient 创建新的HTTP客户端
func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
		Headers: make(map[string]string),
	}
}

// SetHeader 设置HTTP请求头
func (c *Client) SetHeader(key, value string) *Client {
	c.Headers[key] = value
	return c
}

// SetHeaders 批量设置HTTP请求头
func (c *Client) SetHeaders(headers map[string]string) *Client {
	for k, v := range headers {
		c.Headers[k] = v
	}
	return c
}

// SetBasicAuth 设置基本认证
func (c *Client) SetBasicAuth(username, password string) *Client {
	c.Headers["Authorization"] = "Basic " + basicAuth(username, password)
	return c
}

// SetBearerToken 设置Bearer Token认证
func (c *Client) SetBearerToken(token string) *Client {
	c.Headers["Authorization"] = "Bearer " + token
	return c
}

// SetTimeout 设置超时时间
func (c *Client) SetTimeout(timeout time.Duration) *Client {
	c.HTTPClient.Timeout = timeout
	return c
}

// Get 发送GET请求
func (c *Client) Get(path string, params map[string]string) (*Response, error) {
	resourceURL := c.buildURL(path, params)
	req, err := http.NewRequest("GET", resourceURL, nil)
	if err != nil {
		return nil, err
	}

	return c.do(req)
}

// Post 发送POST请求
func (c *Client) Post(path string, body interface{}) (*Response, error) {
	resourceURL := c.buildURL(path, nil)
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", resourceURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return c.do(req)
}

// PostForm 发送表单POST请求
func (c *Client) PostForm(path string, form map[string]string) (*Response, error) {
	resourceURL := c.buildURL(path, nil)
	formValues := url.Values{}
	for k, v := range form {
		formValues.Add(k, v)
	}

	req, err := http.NewRequest("POST", resourceURL, strings.NewReader(formValues.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return c.do(req)
}

// Put 发送PUT请求
func (c *Client) Put(path string, body interface{}) (*Response, error) {
	resourceURL := c.buildURL(path, nil)
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", resourceURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return c.do(req)
}

// Delete 发送DELETE请求
func (c *Client) Delete(path string) (*Response, error) {
	resourceURL := c.buildURL(path, nil)
	req, err := http.NewRequest("DELETE", resourceURL, nil)
	if err != nil {
		return nil, err
	}

	return c.do(req)
}

// Patch 发送PATCH请求
func (c *Client) Patch(path string, body interface{}) (*Response, error) {
	resourceURL := c.buildURL(path, nil)
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", resourceURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return c.do(req)
}

// UploadFile 上传文件
func (c *Client) UploadFile(path string, fieldName, filePath string, params map[string]string) (*Response, error) {
	resourceURL := c.buildURL(path, nil)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", resourceURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	return c.do(req)
}

// Download 下载文件
func (c *Client) Download(path string, destPath string) error {
	resp, err := c.Get(path, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status code: %d", resp.StatusCode)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.Write(resp.Body)
	return err
}

// buildURL 构建完整URL
func (c *Client) buildURL(path string, params map[string]string) string {
	baseURL := strings.TrimRight(c.BaseURL, "/")
	path = strings.TrimLeft(path, "/")
	fullURL := fmt.Sprintf("%s/%s", baseURL, path)

	if len(params) > 0 {
		queryParams := url.Values{}
		for k, v := range params {
			queryParams.Add(k, v)
		}
		fullURL = fmt.Sprintf("%s?%s", fullURL, queryParams.Encode())
	}

	return fullURL
}

// do 执行HTTP请求
func (c *Client) do(req *http.Request) (*Response, error) {
	// 设置请求头
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return &Response{Error: err}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &Response{Error: err}, err
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       body,
		Header:     resp.Header,
		Error:      nil,
	}, nil
}

// basicAuth 生成基本认证字符串
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// GetJSON 解析响应为JSON
func (r *Response) GetJSON(v interface{}) error {
	if r.Error != nil {
		return r.Error
	}

	return json.Unmarshal(r.Body, v)
}

// GetString 获取响应内容为字符串
func (r *Response) GetString() string {
	if r.Error != nil {
		return ""
	}

	return string(r.Body)
}

// IsSuccess 检查响应是否成功
func (r *Response) IsSuccess() bool {
	return r.Error == nil && r.StatusCode >= 200 && r.StatusCode < 300
}

// IsError 检查响应是否为错误
func (r *Response) IsError() bool {
	return r.Error != nil || r.StatusCode >= 400
}

// GetStatusCode 获取HTTP状态码
func (r *Response) GetStatusCode() int {
	return r.StatusCode
}

// GetHeader 获取响应头
func (r *Response) GetHeader(key string) string {
	return r.Header.Get(key)
}

// GetContentType 获取内容类型
func (r *Response) GetContentType() string {
	return r.Header.Get("Content-Type")
}

// GetContentLength 获取内容长度
func (r *Response) GetContentLength() int64 {
	length := r.Header.Get("Content-Length")
	if length == "" {
		return 0
	}

	l, err := strconv.ParseInt(length, 10, 64)
	if err != nil {
		return 0
	}

	return l
}
