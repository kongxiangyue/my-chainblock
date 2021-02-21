package lib

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"time"
)

// 商品
type Commodity struct {
	Name            string  `json:"name"` // 商品名
	Id              string  `json:"id"`
	Origin          string  `json:"origin"`          //产地
	LowTemperature  float64 `json:"lowTemperature"`  //最低温
	HighTemperature float64 `json:"highTemperature"` //最高温
	Price           float64 `json:"price"`           //单价
	OwnerId         string  `json:"owner"`           //所有者
}

// 温度
type Temperature struct {
	Temperature float64   `json:"temperature"`
	RecordTime  time.Time `json:"record_time"`
}

// 订单
type Order struct {
	Commodity            *Commodity        `json:"commodity"`      //商品
	DeliverAddress       string            `json:"deliverAddress"` //配送地址
	Id                   string            `json:"id"`
	OrderTime            time.Time         `json:"orderTime"`
	DeliverTime          time.Time         `json:"deliverTime"`          //配送时间
	Quantity             int64             `json:"quantity"`             //数量
	Status               string            `json:"status"`               //订单状态
	TemperatureVariation []*Temperature    `json:"temperatureVariation"` //温度变化
	BuyerId              string            `json:"buyer"`                //买家
	SellerId             string            `json:"seller"`               //卖家
	TransactionId        fab.TransactionID `json:"transaction_id"`
}
