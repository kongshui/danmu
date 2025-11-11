package service

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// 静态文件下载中间件：为指定路径的请求添加下载头
func StaticDownloadMiddleware(staticPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 仅对静态文件路径的请求生效（例如 "/files/*filepath"）
		if strings.HasPrefix(c.Request.URL.Path, staticPath) {
			// 从 URL 中提取文件名（如 "/files/test.pdf" → "test.pdf"）
			filename := filepath.Base(c.Request.URL.Path)
			if filename == "" {
				filename = "download" // 兜底文件名
			}

			// 编码中文文件名（避免乱码，兼容大部分浏览器）
			encodedFilename := url.QueryEscape(filename)

			// 设置响应头：强制浏览器下载，并用编码后的文件名
			c.Header("Content-Disposition", "attachment; filename*=UTF-8''"+encodedFilename)
			// 可选：设置为二进制流类型（增强兼容性，避免浏览器尝试渲染）
			c.Header("Content-Type", "application/octet-stream")
		}
		// 继续处理请求（交给 StaticFS 读取文件）
		c.Next()
	}
}

// 静态文件处理中间件：目录显示列表，文件触发下载
func StaticFileMiddleware(staticURLPath, localDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 仅处理静态路径下的请求（如 /files/*）
		if !strings.HasPrefix(c.Request.URL.Path, staticURLPath) {
			c.Next()
			return
		}

		// 1. 解析请求路径对应的本地文件路径
		// 从 URL 中提取相对路径（如 /files/docs/test.txt → docs/test.txt）
		relativePath := strings.TrimPrefix(c.Request.URL.Path, staticURLPath)
		// 拼接本地完整路径（确保在 localDir 内，防止路径遍历漏洞）
		localFilePath := filepath.Join(localDir, relativePath)
		// 转换为绝对路径，避免相对路径问题
		absLocalPath, err := filepath.Abs(localFilePath)
		if err != nil {
			c.String(500, "路径解析错误: %v", err)
			c.Abort()
			return
		}

		// 检查路径是否在 localDir 内（安全校验，防止路径遍历）
		localDirAbs, _ := filepath.Abs(localDir)
		if !strings.HasPrefix(absLocalPath, localDirAbs) {
			c.String(403, "禁止访问")
			c.Abort()
			return
		}

		// 2. 判断是文件还是目录
		fileInfo, err := os.Stat(absLocalPath)
		if err != nil {
			// 路径不存在（可能是404），交给 Gin 处理
			c.Next()
			return
		}

		// 3. 若是文件，添加下载头；若是目录，不处理（显示列表）
		if fileInfo.Mode().IsRegular() {
			// 提取文件名（如 test.txt、中文文档.pdf）
			filename := fileInfo.Name()
			// 编码中文文件名（避免乱码）
			encodedFilename := url.QueryEscape(filename)
			// 设置响应头：强制下载，指定文件名
			c.Header("Content-Disposition", "attachment; filename*=UTF-8''"+encodedFilename)
			// 可选：设置为二进制流类型，增强兼容性
			c.Header("Content-Type", "application/octet-stream")
		}

		// 继续处理请求（目录显示列表，文件返回内容）
		c.Next()
	}
}
