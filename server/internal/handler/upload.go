package handler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pxx/pkg/response"

	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil { response.Fail(c, 400, "请选择图片"); return }
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		response.Fail(c, 400, "仅支持jpg/png/webp"); return
	}

	fname := fmt.Sprintf("%d_%06d%s", time.Now().UnixNano(), time.Now().Nanosecond()%1000000, ext)
	dir := "static/uploads"
	os.MkdirAll(dir, 0755)
	fpath := filepath.Join(dir, fname)
	dst, err := os.Create(fpath)
	if err != nil { response.Error(c, err.Error()); return }
	defer dst.Close()
	io.Copy(dst, file)

	response.OK(c, gin.H{"url": "/api/v1/uploads/" + fname})
}
