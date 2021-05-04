package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Order struct {
	Id          string `json"id"`
	Status      string `json"status"`
	CommodityId string `json"commodityid"`
	BuyerId     string `json"buyerid"`
	SellerId    string `json"sellerid"`
}
//
//订单状态
type Status struct {
	New               string //新建
	Confirm           string //确认
	Apply_loan        string //申请贷款
	Accept_loan       string //接受贷款
	Confirm_payment   string //确认货款
	Deliver           string //发货
	Receive           string //确认收货
	Repayment         string //还款
	Repayment_confirm string //确认还款
}

//状态枚举
func newStatus() *Status {
	return &Status{
		New               :"新建",
		Confirm           :"确认",
		Apply_loan        :"申请贷款",
		Accept_loan       :"接受贷款",
		Confirm_payment   :"确认货款",
		Deliver           :"发货",
		Receive           :"确认收货",
		Repayment         :"还款",
		Repayment_confirm :"确认还款",
	}
}

var enumStatus = newStatus()

//以map来实现 枚举体
var statusMap = map[string]string{
	"New"               : enumStatus.New              ,
	"Confirm"           : enumStatus.Confirm          ,
	"Apply_loan"        : enumStatus.Apply_loan       ,
	"Accept_loan"       : enumStatus.Accept_loan      ,
	"Confirm_payment"   : enumStatus.Confirm_payment  ,
	"Deliver"           : enumStatus.Deliver          ,
	"Receive"           : enumStatus.Receive          ,
	"Repayment"         : enumStatus.Repayment        ,
	"Repayment_confirm" : enumStatus.Repayment_confirm,
}


// OrderChaincode example simple Chaincode implementation
type OrderChaincode struct {
}

func (t *OrderChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("order Init")


	fmt.Println("order Init createOrder 001")
	//-一般加一些初始数据(空的数据也无所谓),可以直接调createOrder
	res := t.createOrder(stub,
		[]string{
		"001",
		statusMap["New"],
		"commodity001",
		"buyer001",
		"seller001",
		})

	if res.Status == shim.ERROR {
		return res;
	}

	res = t.createOrder(stub,
		[]string{
			"002",
			statusMap["New"],
			"commodity002",
			"buyer002",
			"seller002",
		})

	if res.Status == shim.ERROR {
		return res;
	}

	fmt.Println("order Init createOrder 002")


	return shim.Success(nil)
}

func (t *OrderChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {


	funcName, args := stub.GetFunctionAndParameters()
	fmt.Println("order Invoke ", funcName)

	switch funcName {
	case "createOrder":
		return t.createOrder(stub, args)
	case "updateOrderStatus":
		return t.updateOrderStatus(stub, args)
	case "queryOrderList":
		return t.queryOrderList(stub, args)

	default:
		return pb.Response{
			Status: shim.ERRORTHRESHOLD,
			Message: "not such function",
		}
	}
}


func (t *OrderChaincode) createOrder(stub shim.ChaincodeStubInterface, args [] string) pb.Response {
	/*
		type Order struct {
			Id          string `json"id"`
			Status      string `json"status"`
			CommodityId string `json"commodityid"`
			BuyerId     string `json"buyerid"`
			SellerId    string `json"sellerid"`
		}
	*/

	//1-检查参数个数
	//2-将参数结构化,然后变成json字符串
	//3-生成key (使用CreateCompositeKey)
	//4-生成kv对(使用PutState)

	if len(args) != 5 {
		return pb.Response{
			Status: shim.ERRORTHRESHOLD,
			Message: "count of args not correct",
		}
	}

	newOrder := &Order{
		Id:args[0],
		Status:statusMap[args[1]],
		CommodityId:args[2],
		BuyerId:args[3],
		SellerId:args[4],
	}

	newOrderJson,err := json.Marshal(newOrder)
	if err != nil {
		return shim.Error("createOrder json wrong")
	}
	fmt.Println("createOrder json:", string(newOrderJson))


	key, err := stub.CreateCompositeKey("order", []string{args[0]})
	if err != nil {
		return shim.Error("createOrder json wrong")
	}
	fmt.Println("createOrder key:", key)

	stub.PutState(key, newOrderJson)


	return shim.Success(nil)
}


func (t *OrderChaincode) updateOrderStatus(stub shim.ChaincodeStubInterface, args [] string) pb.Response {
	//1-检查参数,id status
	//2-生成key,查询(有或没有判断)
	//3-GetState (要判断是否成功)
	//4-将val(是个string)转换成Order结构体
	//5-改Order中的status(严谨一些的话可以验证结构体的id与args[0]是否一致)
	//6-将order转换成json 字符串
	//7-PutState
	if len(args) != 2 {
		return pb.Response{
			Status: shim.ERRORTHRESHOLD,
			Message: "count of args not correct",
		}
	}

	key, err := stub.CreateCompositeKey("order", []string{args[0]})
	if err != nil {
		return shim.Error("updateOrderStatus json wrong")
	}

	json_str, err := stub.GetState(key)
	if err != nil {
		return shim.Error("updateOrderStatus GetState wrong")
	}

	order := &Order{}
	err = json.Unmarshal(json_str, order)
	if err != nil {
		return shim.Error("updateOrderStatus Unmarshal json wrong")
	}

	// order.Id != args[0]

	order.Status = statusMap[args[1]]

	fmt.Println("new status", order.Id, order.Status)

	newOrderJson,err := json.Marshal(order)
	if err != nil {
		return shim.Error("updateOrderStatus json wrong")
	}
	fmt.Println("updateOrderStatus json:", string(newOrderJson))


	stub.PutState(key, newOrderJson)


	return shim.Success(newOrderJson)
}

func (t *OrderChaincode) queryOrderList(stub shim.ChaincodeStubInterface, args [] string) pb.Response {

	//1-拿到所有key (GetStateByPartialCompositeKey)
	//2-循环将keys对应value拿出来 (GetState)
	//	3-合成json(反序列化,append)
	//	4-将Order数线序列化后返回
	fmt.Println("queryOrderList")

	keys := make([]string, 0)
	res, err := stub.GetStateByPartialCompositeKey("order", keys)
	if err != nil {
		return shim.Error("queryOrderList GetStateByPartialCompositeKey wrong")
	}
	defer res.Close()


	orders := make([]*Order, 0)
	for res.HasNext() {
		val, err := res.Next()
		if err != nil {
			return shim.Error("queryOrderList res.Next() wrong")
		}

		tmpOrder := &Order{}
		err = json.Unmarshal(val.GetValue(), tmpOrder)
		if err != nil {
			return shim.Error("queryOrderList res.Next() Unmarshal wrong")
		}

		orders = append(orders, tmpOrder)
	}

	retjson, err := json.Marshal(orders)

	fmt.Println("queryOrderList res json:", string(retjson))



	return shim.Success(retjson)
}



func main() {
	err := shim.Start(new(OrderChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}