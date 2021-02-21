package controller

import (
	bc "accurchain.com/perishable-food-full/application/blockchain"
	"accurchain.com/perishable-food-full/application/lib"
	"accurchain.com/perishable-food-full/application/repository"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const TemperatureUpdatersInterval = 10 * 1000 // 温度传感器获取间隔 10s

var TemperatureUpdaters []*FbeeTemptureUpdater // 用于上传温度

// 订单请求体
type orderRequest struct {
	CommodityId    string `json:"commodity_id" binding:"required"`   // 商品id
	Id             string `json:"id" binding:"required"`             // 订单id
	DeliverAddress string `json:"deliverAddress" binding:"required"` // 配送地址
	DeliverTime    int64  `json:"deliverTime" binding:"required"`    // 预计送达时间（时间戳）
	Quantity       int64  `json:"quantity" binding:"required"`       // 数量
	BuyerId        string `json:"buyer" binding:"required"`          // 买家
	SellerId       string `json:"seller" binding:"required"`         // 卖家
}

// 创建订单
func CreateOrder(ctx *gin.Context) {
	// 解析请求体
	req := new(orderRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// 格式化时间参数，并打印请求参数
	deliverTime := time.Unix(req.DeliverTime/1000, 0)
	fmt.Println("请求参数：")
	marshal, err := json.Marshal(req)
	fmt.Println(string(marshal))

	// 将请求体参数转化为byte数组，发送给区块链，调用链码的createOrder函数
	resp, err := bc.ChannelExecute("createOrder", [][]byte{
		[]byte(req.CommodityId),
		[]byte(req.Id),
		[]byte(req.DeliverAddress),
		[]byte(deliverTime.Format("2006-01-02T15:04:05Z07:00")),
		[]byte(fmt.Sprintf("%d", req.Quantity)),
		[]byte(req.BuyerId),
		[]byte(req.SellerId),
		[]byte(time.Now().Format("2006-01-02T15:04:05Z07:00")),
	})
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// 可先忽略，完成基本功能后再完善
	// 成功返回后，将orderId与txid存到transactionRecords中
	repository.TransactionRecordList.Push(repository.TransactionRecord{OrderId: req.Id, TxID: resp.TransactionID})
	err = repository.TransactionRecordList.Save() // 写入文件（数据持久化）
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// http返回
	ctx.JSON(http.StatusOK, resp)
}

// 查询订单列表
func OrderList(ctx *gin.Context) {
	// 获取请求的请求参数
	orderId := ctx.Query("orderId")
	var args [][]byte
	if orderId != "" {
		args = append(args, []byte(orderId))
	}

	// 将请求参数发送给区块链，调用链码的queryOrderList
	resp, err := bc.ChannelQuery("queryOrderList", args)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	// 将区块链返回的结果反序列化，存入变量Orders中
	var Orders []lib.Order
	json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &Orders)

	// 可先忽略，完成基本功能后再完善
	// 循环查询该订单号（orderId）对应的txid
	for index, order := range Orders {
		_, tr := repository.TransactionRecordList.FindOneByOrderId(order.Id)
		fmt.Printf("result: %+v\n", tr)
		Orders[index].TransactionId = tr.TxID
	}

	// 将结果返回
	ctx.JSON(http.StatusOK, Orders)
}

// 订单状态枚举键值对
var statusMap = map[string]string{
	"New":        "新建",
	"Processing": "运送中",
	"Done":       "完成",
	"Canceled":   "取消",
}

// 更新订单状态请求体
type updateOrderStatusRequest struct {
	OrderId string `form:"order_id" json:"order_id" binding:"required"`
	Status  string `form:"status" json:"status" binding:"required"`
}

// 更新订单状态
func UpdateOrderStatus(ctx *gin.Context) {
	// 解析请求体
	req := new(updateOrderStatusRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// 检查Status字段是否有错
	if _, ok := statusMap[req.Status]; !ok {
		ctx.String(http.StatusBadRequest, "status字段错误")
		return
	}

	// 打印请求参数
	fmt.Println("请求参数：")
	marshal, err := json.Marshal(req) // 转换成json
	fmt.Println(string(marshal))

	switch req.Status {
	// 运送中时
	case "Processing":
		// 构建温度更新对象
		updater := &FbeeTemptureUpdater{
			OrderId:  req.OrderId,
			Interval: TemperatureUpdatersInterval,
		}
		// 记录当前订单的温度获取线程
		TemperatureUpdaters = append(TemperatureUpdaters, updater)
		// 温度传感器定时获取启动
		updater.Start()

	// 订单完成时
	case "Done":
		for i, updater := range TemperatureUpdaters {
			if updater.OrderId == req.OrderId {
				// 停止定时更新器
				updater.Stop()
				// 从全局更新器数组中删除该更新器
				TemperatureUpdaters = append(TemperatureUpdaters[:i], TemperatureUpdaters[i+1:]...)
			}
		}
	}

	// 调用链码
	resp, err := bc.ChannelExecute("updateOrderStatus", [][]byte{
		[]byte(req.OrderId),
		[]byte(req.Status),
	})
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// 将结果返回
	ctx.JSON(http.StatusOK, resp)
}
