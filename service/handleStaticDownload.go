package service

import (
	"net/url"
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
