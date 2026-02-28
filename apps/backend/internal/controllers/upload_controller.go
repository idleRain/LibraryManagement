package controllers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/library-system/backend/internal/utils"
)

type UploadController struct {
	uploadDir string
}

func NewUploadController(uploadDir string) *UploadController {
	// 确保上传目录存在
	os.MkdirAll(uploadDir, 0755)
	return &UploadController{uploadDir: uploadDir}
}

// UploadImage 上传图片
func (uc *UploadController) UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请选择要上传的文件")
		return
	}
	defer file.Close()

	// 检查文件大小（最大 5MB）
	if header.Size > 5*1024*1024 {
		utils.Error(c, 400, "文件大小不能超过 5MB")
		return
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !allowedExts[ext] {
		utils.Error(c, 400, "只支持 jpg、jpeg、png、gif、webp 格式的图片")
		return
	}

	// 生成文件名
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// 按日期分目录存储
	dateDir := time.Now().Format("2006/01/02")
	fullDir := filepath.Join(uc.uploadDir, dateDir)
	os.MkdirAll(fullDir, 0755)

	// 创建文件
	filePath := filepath.Join(fullDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		utils.ServerError(c, "创建文件失败")
		return
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, file); err != nil {
		utils.ServerError(c, "保存文件失败")
		return
	}

	// 返回文件 URL
	fileURL := fmt.Sprintf("/uploads/%s/%s", dateDir, filename)

	utils.Success(c, gin.H{
		"url":      fileURL,
		"filename": header.Filename,
		"size":     header.Size,
	})
}

// UploadFile 上传文件（通用）
func (uc *UploadController) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请选择要上传的文件")
		return
	}
	defer file.Close()

	// 检查文件大小（最大 10MB）
	if header.Size > 10*1024*1024 {
		utils.Error(c, 400, "文件大小不能超过 10MB")
		return
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
		".pdf":  true,
		".doc":  true,
		".docx": true,
		".xls":  true,
		".xlsx": true,
		".txt":  true,
	}

	if !allowedExts[ext] {
		utils.Error(c, 400, "不支持的文件格式")
		return
	}

	// 生成文件名
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// 按日期分目录存储
	dateDir := time.Now().Format("2006/01/02")
	fullDir := filepath.Join(uc.uploadDir, "files", dateDir)
	os.MkdirAll(fullDir, 0755)

	// 创建文件
	filePath := filepath.Join(fullDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		utils.ServerError(c, "创建文件失败")
		return
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, file); err != nil {
		utils.ServerError(c, "保存文件失败")
		return
	}

	// 返回文件 URL
	fileURL := fmt.Sprintf("/uploads/files/%s/%s", dateDir, filename)

	utils.Success(c, gin.H{
		"url":      fileURL,
		"filename": header.Filename,
		"size":     header.Size,
	})
}

// DeleteFile 删除文件
func (uc *UploadController) DeleteFile(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	// 安全检查：确保 URL 以 /uploads/ 开头
	if !strings.HasPrefix(req.URL, "/uploads/") {
		utils.Error(c, 400, "无效的文件路径")
		return
	}

	// 构建文件路径
	filePath := filepath.Join(uc.uploadDir, strings.TrimPrefix(req.URL, "/uploads/"))

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utils.NotFound(c, "文件不存在")
		return
	}

	// 删除文件
	if err := os.Remove(filePath); err != nil {
		utils.ServerError(c, "删除文件失败")
		return
	}

	utils.SuccessWithMsg(c, "删除成功", nil)
}
