package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建路由
	r := gin.Default()

	// 绑定路由规则
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})

	// localhost:8080/user
	// localhost:8080/user?name=xxx
	r.GET("/user", func(c *gin.Context) {
		name := c.DefaultQuery("name", "Gopher")
		c.String(http.StatusOK, "Hello "+name)
	})

	// localhost:8080/user/Gopher/reading
	r.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		c.String(http.StatusOK, name+" is "+strings.Trim(action, "/"))
	})

	type LoginReq struct {
		Type         string `form:"type" json:"type"`
		Username     string `form:"username" json:"username"`
		UserPassword string `form:"userpassword" json:"password"`
	}

	r.POST("/form", func(c *gin.Context) {
		var req LoginReq
		// 根据请求的 Content-Type 自动选择解析器
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		respType := req.Type
		if respType == "" {
			respType = "post"
		}
		c.JSON(http.StatusOK, gin.H{
			"username": req.Username,
			"password": req.UserPassword,
			"type":     respType,
		})
	})

	// 限制上传最大尺寸
	r.MaxMultipartMemory = 8 << 20
	r.POST("/upload", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get err %s", err.Error()))
		}
		// 获取所有文件
		files := form.File["files"]
		for _, file := range files {
			// 上传文件至指定目录
			if err := c.SaveUploadedFile(file, fmt.Sprintf("uploads/%s", file.Filename)); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload err %s", err.Error()))
				return
			}
		}
		c.String(http.StatusOK, fmt.Sprintf("upload ok %d files", len(files)))
	})

	// 监听端口
	r.Run()
}
