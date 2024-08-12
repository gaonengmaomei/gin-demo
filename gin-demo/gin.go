package main

import (
	"gin-demo/startupInit"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 创建一个默认的Gin引擎
	server := gin.Default()

	// 2. 定义路由：当使用 GET 方法访问路径为 /hello 时，执行回调函数
	server.GET("/hello", func(c *gin.Context) {
		// 3. 在回调函数中，返回一个字符串 "hello, go" 并设置HTTP状态码为200 OK
		c.String(http.StatusOK, "hello, go")
	})

	startupInit.Load()

	// 4. 启动服务，监听在 0.0.0.0:8080 上
	server.Run(":8080") // 如果不指定端口号，默认为8080

}
