package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type Smart struct {
}

// 订单
type Order struct {
	OrderID         string
	Price           int64
	Commodity       *Commodity
	Count           int64
	Buyer           string
	BuyTime         time.Time
	DeliveryAddress string
	DeliveryTime    time.Time
	Status          string
	Logistics       string
}

type Commodity struct {
	CommodityID string
	Name        string
	Owner       string
	price       int
}

func (s *Smart) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (s *Smart) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	fcn, args := stub.GetFunctionAndParameters()

	if fcn == "addCommodity" {
		return addCommodity(stub, args)
	} else if fcn == "queryCommodity" {
		return queryCommodity(stub, args)
	} else if fcn == "addOrder" {
		return addOrder(stub, args)
	} else if fcn == "queryOrder" {
		return queryOrder(stub, args)
	}
	return shim.Success(nil)
}

func queryOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}

func addOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}

func queryCommodity(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}

func addCommodity(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// 1.验证
	if len(args) != 2 {
		return pb.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "参数个数不对",
		}
	}
	if args[0] == "" || args[1] == "" {
		return pb.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "无效参数",
		}
	}

	// 2. 修改数据
	//  数据类型的转换
	price, err := strconv.Atoi(args[1])
	if err != nil {
		return pb.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: fmt.Sprintf("价格 数据类型 错误：%v", err.Error()),
		}
	}

	commodity := Commodity{
		CommodityID: strconv.FormatInt(time.Now().Unix(), 10),
		Name:        args[0],
		Owner:       "1", // 1 供货商 2 物流商 3 购买者
		price:       price,
	}

	// 3. 保存数据
	key, err := stub.CreateCompositeKey("commodity", []string{commodity.CommodityID})
	if err != nil {
		return pb.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: fmt.Sprintf("CreateCompositeKey 错误：%v", err.Error()),
		}
	}

	commodityByte, err := json.Marshal(commodity)
	if err != nil {
		return pb.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: fmt.Sprintf("commodity Marshal 错误：%v", err.Error()),
		}
	}

	err = stub.PutState(key, commodityByte)
	if err != nil {
		return pb.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: fmt.Sprintf("putState 错误：%v", err.Error()),
		}
	}

	return shim.Success([]byte(commodity.CommodityID))
}

func main() {
	err := shim.Start(new(Smart))
	if err != nil {
		panic(err.Error())
	}
}
