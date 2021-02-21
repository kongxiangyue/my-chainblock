package main

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestAddCommodity(t *testing.T) {
	stub := shim.NewMockStub("test", new(Smart))
	stub.MockInit("1", [][]byte{})
	resp := stub.MockInvoke("2", [][]byte{
		[]byte("addCommodity"),
		[]byte("火龙果"),
		[]byte("1000"),
	})
	if resp.Status == shim.ERRORTHRESHOLD {
		t.Log(resp.Message)
		t.Fail()
	}
	fmt.Println(string(resp.Payload))
}
