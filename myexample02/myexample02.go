/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"

	"encoding/json"
)

type Mystruct struct {
	Key string     `json:key`
	Value string   `json:value`
}


// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {

	fmt.Println("myex02 Init")
	_, args := stub.GetFunctionAndParameters()


	var key,val string

	key, err := stub.CreateCompositeKey("person", []string{args[0]})
	if err != nil {
		return shim.Error("CreateCompositeKey error")
	}

	val = args[1]


	fmt.Printf("%s-%s\n", key, val)
	stub.PutState(key, []byte(val))


	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("myex02 Invoke")
	funname, args := stub.GetFunctionAndParameters()
	if funname == "create" {
		return t.create(stub, args)
	} else if funname == "query" {
		return t.query(stub, args)
	} else if funname == "queryall" {
		return t.queryall(stub, args)
	}else if funname == "update" {
		return t.create(stub, args)
	}


	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}


func (t *SimpleChaincode) create(stub shim.ChaincodeStubInterface, args []string) pb.Response {


	var key,val string

	key, err := stub.CreateCompositeKey("person", []string{args[0]})
	if err != nil {
		return shim.Error("CreateCompositeKey error")
	}
	val = args[1]

	fmt.Println("before Put val:", val,"\n")


	ret := fmt.Sprintf("%s-%s\n", key, val)
	fmt.Println(ret, "\n")
	stub.PutState(key, []byte(val))


	return shim.Success([]byte(ret))
}


func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var key string
	var err error

	key, err = stub.CreateCompositeKey("person", []string{args[0]})
	if err != nil {
		return shim.Error("CreateCompositeKey error")
	}

	val, err := stub.GetState(key)
	if err != nil {
		return shim.Error("query GetState error")
	}

	fmt.Println("query:value:", val, "\n")
	// {"key":"A","value":"100"}

	mystruct := &Mystruct{
		Key: args[0],
		Value: string(val[:]),
	}

	jsonbytes,err :=  json.Marshal(mystruct)
	if err != nil {
		return shim.Error("query json error")
	}

	fmt.Printf("query json %s \n", jsonbytes)

	return shim.Success(jsonbytes)
}


func (t *SimpleChaincode) queryall(stub shim.ChaincodeStubInterface, args []string) pb.Response {


	keys := make([]string, 0)

	result, err := stub.GetStateByPartialCompositeKey("person", keys)
	if err != nil {
		return shim.Error("queryall result error")
	}
	defer result.Close()
	fmt.Print("queryall keys len",len(keys),"\n")

	mystructlist := make([]*Mystruct, 0)
	for result.HasNext() {
		val,err := result.Next()
		if err != nil {
			shim.Error("queryall result.Next error")
		}

		mystruct := &Mystruct{
			Key: val.GetKey(),
			Value: string(val.GetValue()[:]),
		}

		mystructlist = append(mystructlist, mystruct)
	}

	jsonbytes,err :=  json.Marshal(mystructlist)
	if err != nil {
		return shim.Error("queryall json error")
	}

	fmt.Printf("query json %s \n", jsonbytes)

	return shim.Success(jsonbytes)
}

func (t *SimpleChaincode) update(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}



func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}