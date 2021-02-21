package main

import (
	"net/http"

	"./blockchain"
	"./controller"
	"./fbeecloud"
	"./repository"// by kong 更正了一下路径
	"github.com/gin-gonic/gin"
)

const (
	TransactionRecordFileName = "transactionRecord.json" // 存储 [{orderId，txid}] 数组的json文件名
)

// 设置路由
func setupRouter() *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)			// 开启生产模式（关闭DEBUG模式）
	router := gin.Default()

	// 加载中间件
	router.Use(controller.CORS)         // 先加载跨域模块（便于开发）
	router.Use(controller.FixVueRouter) // 修复前端内部路由问题

	// 定义路由（当地址匹配时，调用相应函数）
	router.POST("/createCommodity", controller.CreateCommodity)
	router.POST("/createOrder", controller.CreateOrder)
	router.GET("/commodityList", controller.CommodityList)
	router.GET("/orderList", controller.OrderList)
	router.POST("/accountList", controller.AccountList)
	router.POST("/updateOrderTempture", controller.UpdateOrderTempture)
	router.POST("/updateOrderStatus", controller.UpdateOrderStatus)
	router.GET("/queryBlockByTxID", controller.QueryBlockByTxID)
	router.GET("/queryTransaction", controller.QueryTransaction)

	// 静态文件路由
	router.StaticFS("/web/", http.Dir("./public/"))

	// 返回路由对象
	return router
}

func main() {
	// 初始化fabric sdk
	blockchain.Init()

	// 加载存储[{orderId，txid}]的json文件
	repository.TransactionRecordList.FilePath = TransactionRecordFileName
	_ = repository.TransactionRecordList.LoadTransactionRecords()

	// 初始化fbeeCloud
	var err error
	controller.Fbee, err = fbeecloud.InitFbeeCloud()
	if err != nil {
		panic(err)
	}

	// 加载路由
	router := setupRouter()

	// 启动http服务器
	router.Run() // listen and serve on 0.0.0.0:8080
}
