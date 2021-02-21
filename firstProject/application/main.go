package main

import (
	"net/http"

	"firstProject/application/farbicSDK"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化SDK
	farbicSDK.Init()
	// 设置路由
	router := getRouter()
	// 开启web服务
	err := router.Run()
	if err != nil {
		panic(err.Error())
	}
}

func getRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/query", query)
	router.GET("/invoke", invoke)
	return router
}

func query(context *gin.Context) {
	fcn := context.Query("fcn")
	arg := context.Query("arg")
	respByte, err := farbicSDK.ChannelQuery(fcn, [][]byte{
		[]byte(arg),
	})
	if err != nil {
		context.String(http.StatusInternalServerError, "调用SDK出错, 错误信息： %v", err.Error())
		return
	}
	context.JSON(http.StatusOK, string(respByte))
}

func invoke(context *gin.Context) {
	fcn := context.Query("fcn")
	arg1 := context.Query("arg1")
	arg2 := context.Query("arg2")
	arg3 := context.Query("arg3")
	respByte, err := farbicSDK.ChannelExecute(fcn, [][]byte{
		[]byte(arg1),
		[]byte(arg2),
		[]byte(arg3),
	})
	if err != nil {
		context.String(http.StatusInternalServerError, "调用SDK出错, 错误信息： %v", err.Error())
		return
	}
	context.JSON(http.StatusOK, string(respByte))
}
