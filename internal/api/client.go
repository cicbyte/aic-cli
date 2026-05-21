package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

type Skill struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Version       string   `json:"version"`
	CategoryID    int      `json:"categoryId"`
	CategoryName  string   `json:"categoryName"`
	Status        int      `json:"status"`
	IsPublic      bool     `json:"isPublic"`
	IsValid       bool     `json:"isValid"`
	FilePath      string   `json:"filePath"`
	License       string   `json:"license"`
	DownloadCount int      `json:"downloadCount"`
	StarCount     int      `json:"starCount"`
	FileSize      int64    `json:"fileSize"`
	CreatedAt     string   `json:"createdAt"`
	UpdatedAt     string   `json:"updatedAt"`
	ValidatedAt   string   `json:"validatedAt"`
	Tags          []string `json:"tags"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Sort        int    `json:"sort"`
}

type FileNode struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Path     string     `json:"path"`
	Type     string     `json:"type"`
	Content  string     `json:"content,omitempty"`
	Children []FileNode `json:"children,omitempty"`
}

type ListResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    ListData `json:"data"`
}

type ListData struct {
	CurrentPage int     `json:"currentPage"`
	Total       int     `json:"total"`
	List        []Skill `json:"skillsList"`
}

type DetailResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Skill  `json:"data"`
}

type FilesResponse struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Data    FilesData  `json:"data"`
}

type FilesData struct {
	Files []FileNode `json:"files"`
}

type CategoriesResponse struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    CategoriesData `json:"data"`
}

type CategoriesData struct {
	CurrentPage     int        `json:"currentPage"`
	Total           int        `json:"total"`
	CategoriesList []Category `json:"categoriesList"`
}

type LoginResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Token string `json:"token"`
	} `json:"data"`
}

type ImportZipResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		SkillID int    `json:"skillId"`
		Name    string `json:"name"`
	} `json:"data"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

func NewClient(baseURL, token string) *Client {
	if baseURL == "" {
		baseURL = "http://localhost:8000"
	}
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) newAuthRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	return req, nil
}

func (c *Client) doRequest(req *http.Request, v interface{}) error {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("认证失败: 请先使用 aic-cli login <token> 登录")
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	return nil
}

func (c *Client) Login(token string) (*LoginResponse, error) {
	url := fmt.Sprintf("%s/api/v1/auth/login", c.BaseURL)

	body, _ := json.Marshal(map[string]string{"token": token})
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	var result LoginResponse
	if err := c.doRequest(req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) ListSkills(pageNum, pageSize, categoryID int, keyword string) (*ListResponse, error) {
	url := fmt.Sprintf("%s/api/v1/skills/list?pageNum=%d&pageSize=%d", c.BaseURL, pageNum, pageSize)
	if categoryID > 0 {
		url += fmt.Sprintf("&categoryId=%d", categoryID)
	}
	if keyword != "" {
		url += fmt.Sprintf("&keyword=%s", keyword)
	}

	req, err := c.newAuthRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var result ListResponse
	if err := c.doRequest(req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetSkillDetail(id int) (*DetailResponse, error) {
	url := fmt.Sprintf("%s/api/v1/skills/detail?id=%d", c.BaseURL, id)

	req, err := c.newAuthRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var result DetailResponse
	if err := c.doRequest(req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetSkillFiles(id int) (*FilesResponse, error) {
	url := fmt.Sprintf("%s/api/v1/skills/files?id=%d", c.BaseURL, id)

	req, err := c.newAuthRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var result FilesResponse
	if err := c.doRequest(req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) DownloadSkill(id int, outputPath string) (string, error) {
	url := fmt.Sprintf("%s/api/v1/skills/download?id=%d", c.BaseURL, id)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("下载失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载失败: HTTP %d", resp.StatusCode)
	}

	contentDisp := resp.Header.Get("Content-Disposition")
	filename := "skill.zip"
	if contentDisp != "" {
		fmt.Sscanf(contentDisp, "attachment; filename=\"%s\"", &filename)
		if len(filename) > 0 && filename[len(filename)-1] == '"' {
			filename = filename[:len(filename)-1]
		}
	}

	if outputPath == "" {
		outputPath = filename
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", fmt.Errorf("写入文件失败: %w", err)
	}

	return outputPath, nil
}

func (c *Client) ListCategories() ([]Category, error) {
	url := fmt.Sprintf("%s/api/v1/categories?all=1", c.BaseURL)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result CategoriesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return result.Data.CategoriesList, nil
}

func (c *Client) ImportZip(zipPath, description string, categoryID int, overwrite bool) (*ImportZipResponse, error) {
	url := fmt.Sprintf("%s/api/v1/skills/import-zip", c.BaseURL)

	file, err := os.Open(zipPath)
	if err != nil {
		return nil, fmt.Errorf("打开ZIP文件失败: %w", err)
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", filepath.Base(zipPath))
	if err != nil {
		return nil, fmt.Errorf("创建表单字段失败: %w", err)
	}

	if _, err = io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("写入文件内容失败: %w", err)
	}

	if description != "" {
		writer.WriteField("description", description)
	}
	if categoryID > 0 {
		writer.WriteField("categoryId", fmt.Sprintf("%d", categoryID))
	}
	if overwrite {
		writer.WriteField("overwrite", "true")
	}

	writer.Close()

	req, err := c.newAuthRequest("POST", url, &body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	var result ImportZipResponse
	if err := c.doRequest(req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) HealthCheck() (*HealthResponse, error) {
	url := fmt.Sprintf("%s/api/v1/health", c.BaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("连接失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result HealthResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}
