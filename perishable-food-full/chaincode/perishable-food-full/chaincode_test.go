/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// 检查初始化
func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

// 检查商品获取
func checkCommodityResult(stub *shim.MockStub, t *testing.T, key string, value string) {
	keys := make([]string, 0)
	keys = append(keys, key)
	result, err := stub.GetStateByPartialCompositeKey("commodity", keys)
	if err != nil {
		fmt.Println("State ", key, "failed to get value")
		t.FailNow()
	}
	defer result.Close()
	for result.HasNext() {
		val, err := result.Next()
		if err != nil {
			t.FailNow()
		}

		commodity := new(Commodity)
		if err := json.Unmarshal(val.GetValue(), commodity); err != nil {
			fmt.Println("Commodity ", key, "failed to convert from bytes")
			t.FailNow()
		}

		if commodity.Name != value {
			fmt.Println("Commodity name", key, "was not", value, "as expected")
			t.FailNow()
		}
	}
}

// 检查账户获取
func checkAccountResult(stub *shim.MockStub, t *testing.T, key string, value string) {
	keys := make([]string, 0)
	keys = append(keys, key)
	result, err := stub.GetStateByPartialCompositeKey("account", keys)
	if err != nil {
		fmt.Println("State ", key, "failed to get value")
		t.FailNow()
	}
	defer result.Close()
	for result.HasNext() {
		val, err := result.Next()
		if err != nil {
			t.FailNow()
		}

		account := new(Account)
		if err := json.Unmarshal(val.GetValue(), account); err != nil {
			fmt.Println("Account ", key, "failed to convert from bytes")
			t.FailNow()
		}

		if account.Name != value {
			fmt.Println("Account name", key, "was not", value, "as expected")
			t.FailNow()
		}
	}
}

// 检查是否写入订单
func checkOrderResult(stub *shim.MockStub, t *testing.T, key string, value string) {
	keys := make([]string, 0)
	keys = append(keys, key)
	result, err := stub.GetStateByPartialCompositeKey("order", keys)
	if err != nil {
		fmt.Println("State ", key, "failed to get value")
		t.FailNow()
	}
	defer result.Close()
	for result.HasNext() {
		val, err := result.Next()
		if err != nil {
			t.FailNow()
		}

		order := new(Order)
		if err := json.Unmarshal(val.GetValue(), order); err != nil {
			fmt.Println("Commodity ", key, "failed to convert from bytes")
			t.FailNow()
		}

		if order.Id != value {
			fmt.Println("Commodity name", key, "was not", value, "as expected")
			t.FailNow()
		}
	}
}

func checkCommodityQuery(t *testing.T, stub *shim.MockStub, name string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("queryCommodityList")})
	if res.Status != shim.OK {
		fmt.Println("Query failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query failed to get value")
		t.FailNow()
	}

	var list []Commodity
	if err := json.Unmarshal(res.Payload, &list); err != nil {
		fmt.Println("Commodity list failed to convert from bytes")
		t.FailNow()
	}

	if list[0].Name != name {
		fmt.Println("Query value", name, "was not as expected")
		t.FailNow()
	}
}

func checkOrderQuery(t *testing.T, stub *shim.MockStub, id string, attr string, expected string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("queryOrderList")})
	if res.Status != shim.OK {
		fmt.Println("Query failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query failed to get value")
		t.FailNow()
	}

	var list []Order
	if err := json.Unmarshal(res.Payload, &list); err != nil {
		fmt.Println("Order list failed to convert from bytes")
		t.FailNow()
	}

	// fmt.Println(list)
	// fmt.Println(list[0].TemperatureVariation[0].RecordTime)//温度对象为空，会报错

	if list[0].Id == id {
		// fmt.Println("Query value", id, "was not as expected")
		// t.FailNow()
		if attr != "" && expected != "" {
			if attr == "TemperatureVariation" {
				var formattedTempture float64
				if val, err := strconv.ParseFloat(expected, 64); err != nil {
					fmt.Printf("Format %s value error", attr)
					t.FailNow()
				} else {
					formattedTempture = val
				}
				if list[0].TemperatureVariation[0].Temperature != formattedTempture {
					fmt.Printf("Query attr %s was not as %s", attr, expected)
					t.FailNow()
				}
			} else if attr == "Status" {
				if list[0].Status != expected {
					fmt.Printf("Query attr %s was not as %s", attr, expected)
					t.FailNow()
				}
			}
		}
	}
}

func checkAccountQuery(t *testing.T, stub *shim.MockStub, id string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("queryAccount"), []byte(id)})
	if res.Status != shim.OK {
		fmt.Println("Query failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query failed to get value")
		t.FailNow()
	}

	var list []Account
	if err := json.Unmarshal(res.Payload, &list); err != nil {
		fmt.Println("Account list failed to convert from bytes")
		t.FailNow()
	}

	if id == "all" {
		if list[0].Id != "1" {
			fmt.Println("Query value", id, "was not as expected")
			t.FailNow()
		}
	} else {
		if list[0].Id != id {
			fmt.Println("Query value", id, "was not as expected")
			t.FailNow()
		}
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
}

//测试初始化和获取账户/商品
func TestPerishableFood_Init(t *testing.T) {
	scc := new(PerishableFood)
	stub := shim.NewMockStub("ex01", scc)

	// checkInit(t, stub, [][]byte{[]byte("init")})
	checkInit(t, stub, nil)

	var names = [3]string{"国光", "红星", "红富士"}
	var ids = [3]string{
		"88efd7ea-bec6-4994-8ed1-f3f7b6f8cac7",
		"36bf5c7f-4cf7-4926-b0f6-0c5c18515752",
		"d9ce807b-e308-11e8-a47c-3e1591a6f5bb"}
	for i, val := range names {
		checkCommodityResult(stub, t, ids[i], val)
	}

	var accountsName = [3]string{"供货商", "物流商", "买家"}
	for i, val := range accountsName {
		checkAccountResult(stub, t, strconv.Itoa(i+1), val)
	}

}

//测试创建商品
func TestCreateCommodity(t *testing.T) {
	scc := new(PerishableFood)
	stub := shim.NewMockStub("ex01", scc)

	// checkInit(t, stub, [][]byte{[]byte("init")})
	checkInit(t, stub, nil)

	checkInvoke(t, stub, [][]byte{[]byte("createCommodity"), []byte("蛇果"), []byte("004"), []byte("中国"),
		[]byte("-2"), []byte("0"), []byte("8"), []byte("1")})

	checkCommodityResult(stub, t, "004", "蛇果")
}

func TestCreateOrder(t *testing.T) {
	scc := new(PerishableFood)
	stub := shim.NewMockStub("ex01", scc)

	// checkInit(t, stub, [][]byte{[]byte("init")})
	checkInit(t, stub, nil)

	checkInvoke(t, stub, [][]byte{[]byte("createOrder"), []byte("88efd7ea-bec6-4994-8ed1-f3f7b6f8cac7"), []byte("001"), []byte("广州市"),
		[]byte("2018-12-07T18:00:00+08:00"), []byte("5"), []byte("3"), []byte("1"), []byte("2018-12-02T18:00:00+08:00")})

	checkOrderResult(stub, t, "001", "001")
}

func TestQueryCommodity(t *testing.T) {
	scc := new(PerishableFood)
	stub := shim.NewMockStub("ex01", scc)

	// checkInit(t, stub, [][]byte{[]byte("init")})
	checkInit(t, stub, nil)

	checkInvoke(t, stub, [][]byte{[]byte("createCommodity"), []byte("蛇果"), []byte("004"), []byte("中国"),
		[]byte("-2"), []byte("0"), []byte("8"), []byte("1")})

	checkCommodityQuery(t, stub, "蛇果")
}

func TestQueryOrder(t *testing.T) {
	scc := new(PerishableFood)
	stub := shim.NewMockStub("ex01", scc)

	// checkInit(t, stub, [][]byte{[]byte("init")})
	checkInit(t, stub, nil)

	checkInvoke(t, stub, [][]byte{[]byte("createOrder"), []byte("88efd7ea-bec6-4994-8ed1-f3f7b6f8cac7"), []byte("001"), []byte("广州市"),
		[]byte("2018-12-02T18:00:00+08:00"), []byte("5"), []byte("3"), []byte("1"), []byte("2018-12-02T18:00:00+08:00")})

	checkOrderQuery(t, stub, "001", "", "")
}

func TestQueryAccount(t *testing.T) {
	scc := new(PerishableFood)
	stub := shim.NewMockStub("ex01", scc)

	// checkInit(t, stub, [][]byte{[]byte("init")})
	checkInit(t, stub, nil)

	checkAccountQuery(t, stub, "1")
}

func TestUpdateOrderTempture(t *testing.T) {
	scc := new(PerishableFood)
	stub := shim.NewMockStub("ex01", scc)

	// checkInit(t, stub, [][]byte{[]byte("init")})
	checkInit(t, stub, nil)

	checkInvoke(t, stub, [][]byte{[]byte("createOrder"), []byte("88efd7ea-bec6-4994-8ed1-f3f7b6f8cac7"), []byte("001"), []byte("广州市"),
		[]byte("2018-12-02T18:00:00+08:00"), []byte("5"), []byte("3"), []byte("1"), []byte("2019-04-11T18:00:00+08:00")})

	checkInvoke(t, stub, [][]byte{[]byte("updateOrderTempture"), []byte("001"), []byte("10"), []byte("2018-12-02T18:00:00+08:00")})

	checkOrderQuery(t, stub, "001", "TemperatureVariation", "10")
}

func TestUpdateOrderStatus(t *testing.T) {
	scc := new(PerishableFood)
	stub := shim.NewMockStub("ex01", scc)

	// checkInit(t, stub, [][]byte{[]byte("init")})
	checkInit(t, stub, nil)

	checkCommodityQuery(t, stub, "红星")

	checkInvoke(t, stub, [][]byte{[]byte("createOrder"), []byte("88efd7ea-bec6-4994-8ed1-f3f7b6f8cac7"), []byte("001"), []byte("广州市"),
		[]byte("2018-12-02T18:00:00+08:00"), []byte("5"), []byte("3"), []byte("1"), []byte("2018-12-02T18:00:00+08:00")})

	// checkInvoke(t, stub, [][]byte{[]byte("updateOrderStatus"), []byte("001"), []byte("Processing")})
	// checkOrderQuery(t, stub, "001", "Status", "运送中")

	// checkInvoke(t, stub, [][]byte{[]byte("updateOrderTempture"), []byte("001"), []byte("1"), []byte("2018-12-01T18:00:00+08:00")})
	// checkInvoke(t, stub, [][]byte{[]byte("updateOrderTempture"), []byte("001"), []byte("2"), []byte("2018-12-01T18:00:00+08:00")})
	// checkInvoke(t, stub, [][]byte{[]byte("updateOrderTempture"), []byte("001"), []byte("-3"), []byte("2018-12-01T18:00:00+08:00")})
	checkInvoke(t, stub, [][]byte{[]byte("updateOrderStatus"), []byte("001"), []byte("Done")})
	checkOrderQuery(t, stub, "001", "Status", "完成")

	checkAccountResult(stub, t, strconv.Itoa(1), "供货商")
}
