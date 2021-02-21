package controller

import (
	bc "../blockchain" // by kong 更正了路径
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 商品请求体
type commodityRequest struct {
	Name            string  `json:"name" form:"name" binding:"required"`                       // 商品名
	Id              string  `json:"id" form:"id" binding:"required"`                           // id
	Origin          string  `json:"origin" form:"origin" binding:"required"`                   // 产地
	LowTemperature  string  `json:"lowTemperature" form:"lowTemperature" binding:"required"`   // 最低温
	HighTemperature string  `json:"highTemperature" form:"highTemperature" binding:"required"` // 最高温
	Price           float64 `json:"price" form:"price" binding:"required"`                     // 单价
	OwnerId         string  `json:"owner" form:"owner" binding:"required"`                     // 所有者
}

// 创建商品
func CreateCommodity(ctx *gin.Context) {
	// 解析请求体
	req := new(commodityRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// 打印请求体
	fmt.Println("请求参数：")
	marshal, _ := json.Marshal(req)
	fmt.Println(string(marshal))

	// 将请求体参数转化为byte数组，发送给区块链，调用链码的createCommodity函数
	resp, err := bc.ChannelExecute("createCommodity", [][]byte{
		[]byte(req.Name),
		[]byte(req.Id),
		[]byte(req.Origin),
		[]byte(req.LowTemperature),
		[]byte(req.HighTemperature),
		[]byte(fmt.Sprintf("%v", req.Price)),
		[]byte(req.OwnerId),
	})
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// http返回
	ctx.JSON(http.StatusOK, resp)
}

// 查询商品列表
func CommodityList(ctx *gin.Context) {
	// 向区块链发起query，调用链码的queryCommodityList函数
	resp, err := bc.ChannelQuery("queryCommodityList", [][]byte{})
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// 反序列化json
	var data []map[string]interface{}
	json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data)

	// 将结果返回
	ctx.JSON(http.StatusOK, data)
}
