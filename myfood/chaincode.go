/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// PerishableFood example
type PerishableFood struct {
}

// 账户
type Account struct {
	Id      string  `json:"id"`      //账号ID
	Name    string  `json:"name"`    //账号名
	Balance float64 `json:"balance"` //余额
}

// 商品
type Commodity struct {
	Name            string  `json:"name"`            // 商品名
	Id              string  `json:"id"`              // 商品ID
	Origin          string  `json:"origin"`          //产地
	LowTemperature  float64 `json:"lowTemperature"`  //最低温
	HighTemperature float64 `json:"highTemperature"` //最高温
	Price           float64 `json:"price"`           //单价
	OwnerId         string  `json:"owner"`           //所有者
}

//订单
type Order struct {
	Commodity            *Commodity     `json:"commodity"`            //商品
	DeliverAddress       string         `json:"deliverAddress"`       //配送地址
	Id                   string         `json:"id"`                   // 订单ID
	OrderTime            time.Time      `json:"orderTime"`            //  下单时间
	DeliverTime          time.Time      `json:"deliverTime"`          //配送时间
	Quantity             int64          `json:"quantity"`             //数量
	Status               string         `json:"status"`               //订单状态
	TemperatureVariation []*Temperature `json:"temperatureVariation"` //温度变化
	BuyerId              string         `json:"buyer"`                //买家
	SellerId             string         `json:"seller"`               //卖家
}

//温度
type Temperature struct {
	Temperature float64   `json:"temperature"` //温度
	RecordTime  time.Time `json:"record_time"` //记录时间
}

//订单状态
type Status struct {
	New        string //新建
	Processing string //运送中
	Done       string //完成
	Canceled   string //取消
}

//状态枚举
func newStatus() *Status {
	return &Status{
		New:        "新建",
		Processing: "运送中",
		Done:       "完成",
		Canceled:   "取消",
	}
}

var enumStatus = newStatus()

//枚举键值对
var statusMap = map[string]string{
	"New":        enumStatus.New,
	"Processing": enumStatus.Processing,
	"Done":       enumStatus.Done,
	"Canceled":   enumStatus.Canceled,
}

//链码初始化
func (t *PerishableFood) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("链码初始化")
	//初始化默认数据
	var names = [3]string{"国光", "红星", "红富士"}
	var ids = [3]string{
		"88efd7ea-bec6-4994-8ed1-f3f7b6f8cac7",
		"36bf5c7f-4cf7-4926-b0f6-0c5c18515752",
		"d9ce807b-e308-11e8-a47c-3e1591a6f5bb"}
	var accountsName = [3]string{"供货商", "物流商", "买家"}
	var accountList []string
	//初始化账号数据，为“供应商”，“物流商”，“买家”账号初始化账号
	for i, val := range accountsName {
		account := &Account{// 将各用户序列化成json by kong
			Name:    val,
			Id:      strconv.Itoa(i + 1),
			Balance: 1000,
		}
		// 序列化对象
		bytes, err := json.Marshal(account)
		if err != nil {
			return shim.Error(fmt.Sprintf("marshal account error %s", err))
		}

		var key string
		// 以"account" + id 为key by kong
		if val, err := stub.CreateCompositeKey("account", []string{account.Id}); err != nil {
			return shim.Error(fmt.Sprintf("create key error %s", err))
		} else {
			key = val
		}

		if err := stub.PutState(key, bytes); err != nil {
			return shim.Error(fmt.Sprintf("put account error %s", err))
		}

		// 辅助变量，用于记住当前有多少个用户 by kong
		accountList = append(accountList, account.Id)
	}

	//初始化商品数据，"国光", "红星", "红富士" 3种商品
	for i, val := range names {
		price := 6.00 + float64(i)// by kong这里以给一个不同的价格
				commodity := &Commodity{
				Name:            val,
				Id:              ids[i],
			Origin:          "中国",
			LowTemperature:  -2,
			HighTemperature: 0,
			Price:           price, //单价
			OwnerId:         accountList[i],
		}

		// 序列化对象

		commodityBytes, err := json.Marshal(commodity)
		if err != nil {
			return shim.Error(fmt.Sprintf("marshal commodity error %s", err))
		}

		fmt.Printf( "-by kong json string %s\n", commodityBytes)



		var key string// by kong 以"commodity"+id为key
		if val, err := stub.CreateCompositeKey("commodity", []string{commodity.Id}); err != nil {
			return shim.Error(fmt.Sprintf("create key error %s", err))
		} else {
			key = val
			fmt.Printf( "-by kong CreateCompositeKey %s\n", key)

		}



		if err := stub.PutState(key, commodityBytes); err != nil {
			return shim.Error(fmt.Sprintf("put commodity error %s", err))
		}
	}

	fmt.Print( "-by kong Init success\n")
	return shim.Success(nil)
}

//实现Invoke接口调用智能合约，例子中的所有智能合约都会集中在这个接口实现
func (t *PerishableFood) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	funcName, args := stub.GetFunctionAndParameters()

	switch funcName {
	//创建商品
	case "createCommodity":
		return createCommodity(stub, args)
	//创建订单
	case "createOrder":
		return createOrder(stub, args)
	//查询商品列表
	case "queryCommodityList":
		return queryCommodityList(stub, args)
	//查询订单列表
	case "queryOrderList":
		return queryOrderList(stub, args)
	//查询账户
	case "queryAccount":
		return queryAccount(stub, args)
	//查询订单温度
	case "updateOrderTempture":
		return updateOrderTempture(stub, args)
	//更新订单状态
	case "updateOrderStatus":
		return updateOrderStatus(stub, args)
	default:
		return shim.Error(fmt.Sprintf("unsupported function: %s", funcName))
	}
}

//新建商品
func createCommodity(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 检查参数的个数
	if len(args) != 7 {
		return shim.Error("not enough args")
	}

	// 验证参数的正确性
	name := args[0]
	id := args[1]
	origin := args[2]
	lowTemperature := args[3]
	highTemperature := args[4]
	price   := args[5]
	ownerId := args[6]

	if name == "" || id == "" || origin == "" || lowTemperature == "" || highTemperature == "" || price == "" || ownerId == "" {
		return shim.Error("invalid args")
	}

	//创建主键
	var key string // by kong 这里就可以跟前init那里可以要对应起来
	if val, err := stub.CreateCompositeKey("commodity", []string{id}); err != nil {
		return shim.Error(fmt.Sprintf("create key error %s", err))
	} else {
		key = val
		fmt.Printf( "-by kong createCommodity:key %s\n", key)
	}

	// 验证数据是否存在 应该存在 or 不应该存在
	if commodityBytes, err := stub.GetState(key); err == nil && len(commodityBytes) != 0 {
		return shim.Error("commodity already exist")
	}

	//数据格式转换
	var formattedLowTemp float64
	if val, err := strconv.ParseFloat(lowTemperature, 64); err != nil {
		return shim.Error("format low temperature error")
	} else {
		formattedLowTemp = val
	}

	//数据格式转换
	var formattedHighTemp float64
	if val, err := strconv.ParseFloat(highTemperature, 64); err != nil {
		return shim.Error("format high temperature error")
	} else {
		formattedHighTemp = val
	}

	//数据格式转换
	var formattedPrice float64
	if val, err := strconv.ParseFloat(price, 64); err != nil {
		return shim.Error("format price error")
	} else {
		formattedPrice = val
	}

	// 写入状态
	commodity := &Commodity{
		Name:            name,
		Id:              id,
		Origin:          origin,
		LowTemperature:  formattedLowTemp,
		HighTemperature: formattedHighTemp,
		Price:           formattedPrice, //单价
		OwnerId:         ownerId,
	}

	// 序列化对象
	commodityBytes, err := json.Marshal(commodity)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal commodity error %s", err))
	}

	fmt.Printf( "-by kong createCommodity json string %s\n", commodityBytes)

	//写入区块链账本
	if err := stub.PutState(key, commodityBytes); err != nil {
		return shim.Error(fmt.Sprintf("put commodity error %s", err))
	}

	// 成功返回
	return shim.Success(commodityBytes)// 2021的样题要求返回200 by kong
}

//新建订单
func createOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 检查参数的个数
	if len(args) != 8 {
		return shim.Error("not enough args")
	}

	// 验证参数的正确性 by kong 为什么参数顺序是这样子的——需要跟application的代码联动起来看
	commodity_id := args[0]
	id := args[1]
	deliverAddress := args[2]
	deliverTime := args[3]
	quantity := args[4]
	buyerId := args[5]
	sellerId := args[6]
	orderTime := args[7]

	if commodity_id == "" || id == "" || deliverAddress == "" || deliverTime == "" || quantity == "" || buyerId == "" || sellerId == "" || orderTime == "" {
		return shim.Error("invalid args")
	}

	commodity := new(Commodity)
	// 验证数据是否存在 应该存在 or 不应该存在
	result, err := stub.GetStateByPartialCompositeKey("commodity", []string{commodity_id})
	if err != nil {
		return shim.Error(fmt.Sprintf("Get commodity error %s", err))
	}
	defer result.Close()
	for result.HasNext() {
		val, err := result.Next()
		if err != nil {
			return shim.Error(fmt.Sprintf("Get commodity error %s", err))
		}

		if err := json.Unmarshal(val.GetValue(), commodity); err != nil {
			return shim.Error(fmt.Sprintf("Commodity failed to convert from bytes, error %s", err))
		}
	}

	//数据格式转换
	var formattedTime time.Time
	if val, err := time.Parse("2006-01-02T15:04:05Z07:00", deliverTime); err != nil {
		return shim.Error(fmt.Sprintf("formate deliverTime error: %s", err))
	} else {
		formattedTime = val
	}

	//数据格式转换
	var formattedQuantity int64
	if val, err := strconv.ParseInt(quantity, 10, 64); err != nil {
		return shim.Error(fmt.Sprintf("formate quantity error: %s", err))
	} else {
		formattedQuantity = val
	}

	//数据格式转换
	var formattedOrderTime time.Time
	if val, err := time.Parse("2006-01-02T15:04:05Z07:00", orderTime); err != nil {
		return shim.Error(fmt.Sprintf("formate deliverTime error: %s", err))
	} else {
		formattedOrderTime = val
	}

	// 写入状态
	order := &Order{
		Commodity:            commodity,
		Id:                   id,
		OrderTime:            formattedOrderTime,
		DeliverAddress:       deliverAddress,
		DeliverTime:          formattedTime,
		Quantity:             formattedQuantity,
		Status:               enumStatus.New,
		TemperatureVariation: make([]*Temperature, 0),
		BuyerId:              buyerId,
		SellerId:             sellerId,
	}

	// 序列化对象
	orderBytes, err := json.Marshal(order)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal order error %s", err))
	}

	//创建主键
	var key string
	if val, err := stub.CreateCompositeKey("order", []string{id}); err != nil {
		return shim.Error(fmt.Sprintf("create key error %s", err))
	} else {
		key = val
	}

	//写入区块链账本
	if err := stub.PutState(key, orderBytes); err != nil {
		return shim.Error(fmt.Sprintf("put order error %s", err))
	}

	// 成功返回
	return shim.Success(nil)
}

//查询商品列表
func queryCommodityList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 检查参数的个数
	if len(args) != 0 {
		return shim.Error("no args required.")
	}

	//通过主键从区块链查找相关的数据
	keys := make([]string, 0)
	result, err := stub.GetStateByPartialCompositeKey("commodity", keys)
	if err != nil {
		return shim.Error(fmt.Sprintf("query commodity error: %s", err))
	}
	defer result.Close()

	//检查返回的数据是否为空，不为空则遍历数据，否则返回空数组
	commoditylist := make([]*Commodity, 0)
	for result.HasNext() {
		val, err := result.Next()
		if err != nil {
			return shim.Error(fmt.Sprintf("query commodity error: %s", err))
		}

		commodity := new(Commodity)//by kong val.GetValue() 返回是json字符串
		if err := json.Unmarshal(val.GetValue(), commodity); err != nil {
			return shim.Error(fmt.Sprintf("unmarshal error: %s", err))
		}

		commoditylist = append(commoditylist, commodity)
	}

	//序列化数据
	bytes, err := json.Marshal(commoditylist)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal error: %s", err))
	}

	fmt.Printf( "-by kong CreateCompositeKey %s\n", bytes)
	return shim.Success(bytes)
}

//查询订单列表
func queryOrderList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 检查参数的个数
	if len(args) > 1 {// by kong 只会有一个参数，联动application代码
		return shim.Error("no args required.")
	}

	keys := make([]string, 0)

	if len(args) == 1 {
		orderId := args[0]
		if orderId == "" {
			return shim.Error("invalid args")
		}
		keys = append(keys, orderId)
	}

	//通过主键从区块链查找相关的数据
	result, err := stub.GetStateByPartialCompositeKey("order", keys)
	if err != nil {
		return shim.Error(fmt.Sprintf("query order error: %s", err))
	}
	defer result.Close()

	//检查返回的数据是否为空，不为空则遍历数据，否则返回空数组
	orders := make([]*Order, 0)
	for result.HasNext() {

		val, err := result.Next()
		if err != nil {
			return shim.Error(fmt.Sprintf("query orders error: %s", err))
		}

		order := new(Order)
		if err := json.Unmarshal(val.GetValue(), order); err != nil {
			return shim.Error(fmt.Sprintf("unmarshal error: %s", err))
		}

		orders = append(orders, order)

	}

	//序列化数据
	bytes, err := json.Marshal(orders)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal error: %s", err))
	}

	return shim.Success(bytes)
}

//查询账号
func queryAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 检查参数的个数
	if len(args) != 1 {
		return shim.Error("not enough args.")
	}

	// 验证参数的正确性
	account_id := args[0]

	if account_id == "" {
		return shim.Error("invalid args")
	}

	keys := make([]string, 0)
	if account_id != "all" {
		keys = append(keys, account_id)
	}

	//通过主键从区块链查找相关的数据
	result, err := stub.GetStateByPartialCompositeKey("account", keys)
	if err != nil {
		return shim.Error(fmt.Sprintf("query account error: %s", err))
	}
	defer result.Close()

	//检查返回的数据是否为空，不为空则遍历数据，否则返回空数组
	accounts := make([]*Account, 0)
	for result.HasNext() {
		val, err := result.Next()
		if err != nil {
			return shim.Error(fmt.Sprintf("query accounts error: %s", err))
		}

		account := new(Account)
		if err := json.Unmarshal(val.GetValue(), account); err != nil {
			return shim.Error(fmt.Sprintf("unmarshal error: %s", err))
		}

		accounts = append(accounts, account)
	}

	//序列化数据
	bytes, err := json.Marshal(accounts)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal error: %s", err))
	}

	return shim.Success(bytes)
}

//更新订单温度列表
func updateOrderTempture(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 检查参数的个数
	if len(args) != 3 {
		return shim.Error("not enough args.")
	}

	// 验证参数的正确性
	order_id := args[0]
	temperature := args[1]
	recordTime := args[2]

	if order_id == "" || temperature == "" {
		return shim.Error("invalid args")
	}

	//数据格式转换
	var formattedTemperature float64
	if val, err := strconv.ParseFloat(temperature, 64); err != nil {
		return shim.Error(fmt.Sprintf("formate temperature error: %s", err))
	} else {
		formattedTemperature = val
	}

	//数据格式转换
	var formattedTime time.Time
	if val, err := time.Parse("2006-01-02T15:04:05Z07:00", recordTime); err != nil {
		return shim.Error(fmt.Sprintf("formate deliverTime error: %s", err))
	} else {
		formattedTime = val
	}

	//通过主键从区块链查找相关的数据
	keys := make([]string, 0)
	keys = append(keys, order_id)
	result, err := stub.GetStateByPartialCompositeKey("order", keys)
	if err != nil {
		return shim.Error(fmt.Sprintf("query order error: %s", err))
	}
	defer result.Close()

	//检查返回的数据是否为空，不为空则遍历数据，否则返回空数组
	order := new(Order)
	for result.HasNext() {
		val, err := result.Next()
		if err != nil {
			return shim.Error(fmt.Sprintf("query orders error: %s", err))
		}

		if err := json.Unmarshal(val.GetValue(), order); err != nil {
			return shim.Error(fmt.Sprintf("unmarshal error: %s", err))
		}

		temperature := new(Temperature)
		temperature.Temperature = formattedTemperature
		temperature.RecordTime = formattedTime

		order.TemperatureVariation = append(order.TemperatureVariation, temperature)
	}

	// 序列化对象
	orderBytes, err := json.Marshal(order)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal order error %s", err))
	}

	//创建主键
	var key string
	if val, err := stub.CreateCompositeKey("order", []string{order_id}); err != nil {
		return shim.Error(fmt.Sprintf("create key error %s", err))
	} else {
		key = val
	}

	//写入区块链账本
	if err := stub.PutState(key, orderBytes); err != nil {
		return shim.Error(fmt.Sprintf("put commodity error %s", err))
	}

	return shim.Success(nil)
}

//更新订单状态
func updateOrderStatus(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 检查参数的个数
	if len(args) != 2 {
		return shim.Error("not enough args.")
	}

	// 验证参数的正确性
	order_id := args[0]
	status := args[1]

	if order_id == "" || status == "" {
		return shim.Error("invalid args")
	}

	// 通过主键从区块链查找相关的数据
	keys := make([]string, 0)
	keys = append(keys, order_id)
	result, err := stub.GetStateByPartialCompositeKey("order", keys)
	if err != nil {
		return shim.Error(fmt.Sprintf("query order error: %s", err))
	}
	defer result.Close()

	//检查返回的数据是否为空，不为空则遍历数据，否则返回空数组
	order := new(Order)
	for result.HasNext() {
		val, err := result.Next()
		if err != nil {
			return shim.Error(fmt.Sprintf("query orders error: %s", err))
		}

		if err := json.Unmarshal(val.GetValue(), order); err != nil {
			return shim.Error(fmt.Sprintf("unmarshal error: %s", err))
		}

		order.Status = statusMap[status]
	}

	//订单完成状态的处理逻辑
	/**
	  若温度超出约定范围，将按如下公式进行计算。
	  没有超出范围将正常付款。
	  扣款计算公式：扣款 = 最低温度偏差值 * 0.1 * 货物数量 + 最高温度偏差值 * 0.2 * 货物数量
	*/
	if status == "Done" {
		//转账的逻辑
		var totalPrice float64
		totalPrice = float64(order.Quantity) * order.Commodity.Price
		// 统计温差的变化
		lossFee := 0.0

		accounts := make([]*Account, 0)
		//获取买家账号
		buyerResult, buyErr := getStateByPartialCompositeKey(stub, order.BuyerId)
		if buyErr != nil {
			return shim.Error(fmt.Sprintf("query account error: %s", buyErr))
		}
		defer buyerResult.Close()

		for buyerResult.HasNext() {
			val, err := buyerResult.Next()
			if err != nil {
				return shim.Error(fmt.Sprintf("query accounts error: %s", err))
			}

			account := new(Account)
			if err := json.Unmarshal(val.GetValue(), account); err != nil {
				return shim.Error(fmt.Sprintf("unmarshal error: %s", err))
			}

			accounts = append(accounts, account)
		}

		//获取卖家账号
		sellerResult, sellerErr := getStateByPartialCompositeKey(stub, order.SellerId)
		if sellerErr != nil {
			return shim.Error(fmt.Sprintf("query account error: %s", sellerErr))
		}
		defer sellerResult.Close()

		for sellerResult.HasNext() {
			val, err := sellerResult.Next()
			if err != nil {
				return shim.Error(fmt.Sprintf("query accounts error: %s", err))
			}

			account := new(Account)
			if err := json.Unmarshal(val.GetValue(), account); err != nil {
				return shim.Error(fmt.Sprintf("unmarshal error: %s", err))
			}

			accounts = append(accounts, account)
		}

		//获取账号中的余额
		var buyerAcc *Account
		var sellerAcc *Account
		for _, val := range accounts {
			if val.Id == order.BuyerId {
				buyerAcc = val
			} else if val.Id == order.SellerId {
				sellerAcc = val
			}
		}

		now := time.Now()
		var formattedTime = order.DeliverTime
		sumD := now.Sub(formattedTime)
		// 超时买家则不用付钱
		if sumD.Hours()/24 > 0.0 {
			totalPrice = 0
		} else {
			if len(order.TemperatureVariation) != 0 {
				temperatureArr := order.TemperatureVariation
				lowestTemp := temperatureArr[0].Temperature
				highestTemp := temperatureArr[0].Temperature

				for _, val := range temperatureArr {
					if val.Temperature <= lowestTemp {
						lowestTemp = val.Temperature
					}
					if val.Temperature >= highestTemp {
						highestTemp = val.Temperature
					}
				}

				//计算最低温度偏差值
				lowTempDif := 0.0
				if lowestTemp < order.Commodity.LowTemperature {
					lowTempDif = order.Commodity.LowTemperature - lowestTemp
				}
				//计算最高温度偏差值
				highTempDif := 0.0
				if highestTemp > order.Commodity.HighTemperature {
					highTempDif = highestTemp - order.Commodity.HighTemperature
				}

				lossFee = lowTempDif*0.1*float64(order.Quantity) + highTempDif*0.2*float64(order.Quantity)
			}

		}
		//转账
		totalPrice -= lossFee
		buyerAcc.Balance -= totalPrice
		sellerAcc.Balance += totalPrice

		// 序列化对象
		buyerBytes, sellerErr := json.Marshal(buyerAcc)
		if sellerErr != nil {
			return shim.Error(fmt.Sprintf("marshal buyer error %s", sellerErr))
		}

		var buyerKey string
		if val, err := stub.CreateCompositeKey("account", []string{order.BuyerId}); err != nil {
			return shim.Error(fmt.Sprintf("create key error %s", err))
		} else {
			buyerKey = val
		}

		if err := stub.PutState(buyerKey, buyerBytes); err != nil {
			return shim.Error(fmt.Sprintf("put buyer account error %s", err))
		}

		// 序列化对象
		sellerBytes, sellerErr := json.Marshal(sellerAcc)
		if sellerErr != nil {
			return shim.Error(fmt.Sprintf("marshal seller error %s", sellerErr))
		}

		var sellerKey string
		if val, err := stub.CreateCompositeKey("account", []string{order.SellerId}); err != nil {
			return shim.Error(fmt.Sprintf("create key error %s", err))
		} else {
			sellerKey = val
		}

		if err := stub.PutState(sellerKey, sellerBytes); err != nil {
			return shim.Error(fmt.Sprintf("put seller account error %s", err))
		}
	}

	// 序列化对象
	orderBytes, err := json.Marshal(order)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal order error %s", err))
	}

	//构建主键
	var key string
	if val, err := stub.CreateCompositeKey("order", []string{order_id}); err != nil {
		return shim.Error(fmt.Sprintf("create key error %s", err))
	} else {
		key = val
	}

	//写入区块链账本
	if err := stub.PutState(key, orderBytes); err != nil {
		return shim.Error(fmt.Sprintf("put commodity error %s", err))
	}

	return shim.Success(nil)
}

func getStateByPartialCompositeKey(stub shim.ChaincodeStubInterface, key string) (shim.StateQueryIteratorInterface, error) {
	keys := make([]string, 0)
	keys = append(keys, key)
	result, err := stub.GetStateByPartialCompositeKey("account", keys)
	return result, err
}

func main() {
	err := shim.Start(new(PerishableFood))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
