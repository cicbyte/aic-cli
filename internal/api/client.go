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
	HTTPClient *http.Client
}

type Skill struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Version       string `json:"version"`
	CategoryID    int    `json:"categoryId"`
	CategoryName  string `json:"categoryName"`
	Status        int    `json:"status"`
	IsPublic      bool   `json:"isPublic"`
	IsValid       bool   `json:"isValid"`
	FilePath      string `json:"filePath"`
	DownloadCount int    `json:"downloadCount"`
	StarCount     int    `json:"starCount"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Sort        int    `json:"sort"`
}

type ListResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    ListData `json:"data"`
}

type ListData struct {
	Total int     `json:"total"`
	List  []Skill `json:"skillsList"`
}

type DetailResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    SkillDetail `json:"data"`
}

type SkillDetail struct {
	Skill
	Files []FileNode `json:"files"`
}

type FileNode struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Path     string     `json:"path"`
	Type     string     `json:"type"`
	Content  string     `json:"content,omitempty"`
	Children []FileNode `json:"children,omitempty"`
}

func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = "http://localhost:8000"
	}
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) ListSkills(pageNum, pageSize int, categoryID int) (*ListResponse, error) {
	url := fmt.Sprintf("%s/api/v1/skills/list?pageNum=%d&pageSize=%d", c.BaseURL, pageNum, pageSize)
	if categoryID > 0 {
		url += fmt.Sprintf("&categoryId=%d", categoryID)
	}

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result ListResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

func (c *Client) GetSkillDetail(id int) (*DetailResponse, error) {
	url := fmt.Sprintf("%s/api/v1/skills/detail?id=%d", c.BaseURL, id)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result DetailResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
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

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("写入文件失败: %w", err)
	}

	return outputPath, nil
}

func (c *Client) ListCategories() ([]Category, error) {
	url := fmt.Sprintf("%s/api/v1/categories/list", c.BaseURL)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result struct {
		Code           int         `json:"code"`
		Message        string      `json:"message"`
		Data           interface{} `json:"data"`
		CategoriesList []Category  `json:"categoriesList"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if len(result.CategoriesList) > 0 {
		return result.CategoriesList, nil
	}

	if dataMap, ok := result.Data.(map[string]interface{}); ok {
		if list, ok := dataMap["categoriesList"].([]interface{}); ok {
			data, _ := json.Marshal(list)
			var categories []Category
			json.Unmarshal(data, &categories)
			return categories, nil
		}
	}

	return nil, nil
}

type CreateSkillRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	CategoryID  int      `json:"categoryId"`
	Tags        []string `json:"tags"`
}

type CreateSkillResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		SkillID int `json:"skillId"`
	} `json:"data"`
}

func (c *Client) CreateSkill(req CreateSkillRequest) (*CreateSkillResponse, error) {
	url := fmt.Sprintf("%s/api/v1/skills/create", c.BaseURL)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	resp, err := c.HTTPClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result CreateSkillResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

type SaveFilesRequest struct {
	ID    int        `json:"id"`
	Files []FileNode `json:"files"`
}

type SaveFilesResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (c *Client) SaveFiles(skillID int, files []FileNode) error {
	url := fmt.Sprintf("%s/api/v1/skills/save-files", c.BaseURL)

	req := SaveFilesRequest{
		ID:    skillID,
		Files: files,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %w", err)
	}

	resp, err := c.HTTPClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	var result SaveFilesResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Code != 0 {
		return fmt.Errorf("保存文件失败: %s", result.Message)
	}

	return nil
}

type ImportZipResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		SkillID int    `json:"skillId"`
		Name    string `json:"name"`
	} `json:"data"`
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

	req, err := http.NewRequest("POST", url, &body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result ImportZipResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}
