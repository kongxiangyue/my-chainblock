package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"strings"
)

var (
	GOPATH        = os.Getenv("GOPATH") // 从环境变量中读取GOPATH
	IndexHtmlPath = path.Join(GOPATH,
		"src/accurchain.com/perishable-food-full/application/public/index.html") // 拼接index.html的路径
)

// 允许跨域
func CORS(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "*")

	if c.Request.Method == "OPTIONS" {
		c.String(http.StatusOK, "")
	}

	// 调用下个中间件
	c.Next()
}

// 修复前端内部路由的问题
func FixVueRouter(c *gin.Context) {
	// 调用下个中间件
	c.Next()

	// 请求返回前进行判断，如果找不到文件，且访问的路径是前端的路径，则将index.html返回
	if strings.Contains(c.Request.RequestURI, "/web/") && c.Writer.Status() == 404 {
		// 返回index.html
		c.File(IndexHtmlPath)
	}
}
