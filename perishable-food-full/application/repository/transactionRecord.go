package repository

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"io/ioutil"
)

// transactionRecords中单个元素的结构体 {订单id，txid}
type TransactionRecord struct {
	OrderId string            `json:"order_id" bson:"order_id"` // orderId 订单Id
	TxID    fab.TransactionID `json:"tx_id" bson:"tx_id"`       // txid 区块链的transactionId
}

// 存储 [{orderId，txid}] 数组，用于持久化存储订单id与txid的对应关系
// 生成订单时会由前端传入一个订单号，将该交易发送到区块链后会返回一个区块链的txid
// 在订单列表中需要展示，从区块链查询的数据只有orderId,所以需要找到它对应的txid
type TransactionRecords struct {
	TxRecord []TransactionRecord
	FilePath string
}

// 将一个TransactionRecord加进数组
func (trs *TransactionRecords) Push(tr TransactionRecord) {
	trs.TxRecord = append(trs.TxRecord, tr)
}

// 将本结构体内的TxRecord数组由内存写到磁盘，默认当前目录生成一个json文件
func (trs *TransactionRecords) Save() (e error) {
	marshal, e := json.Marshal(trs.TxRecord) // 序列化为json
	if e != nil {
		return e
	}

	// 写文件
	e = ioutil.WriteFile(trs.FilePath, marshal, 0644)
	if e != nil {
		return e
	}

	return nil
}

// 用orderId来查找对应的txid
func (trs *TransactionRecords) FindOneByOrderId(orderId string) (index int, tr TransactionRecord) {
	for index, item := range trs.TxRecord {
		if item.OrderId == orderId {
			return index, item
		}
	}
	return -1, TransactionRecord{}
}

// 根据传入的orderid，将其从数组中删除。
// 若成功删除，则返回该结构体；若不存在，则返回空。
func (trs *TransactionRecords) DeleteOne(orderId string) (tr TransactionRecord) {
	index, tr := trs.FindOneByOrderId(orderId)
	if index != -1 {
		trs.TxRecord = append(trs.TxRecord[:index], trs.TxRecord[index+1:]...)
	}
	return tr
}

// 实现String方法，用于打印日志
func (trs *TransactionRecords) String() (s string) {
	return fmt.Sprintf("%+v", trs.TxRecord)
}

var TransactionRecordList TransactionRecords // 订单id与交易id对应关系

// 加载存储orderId与txid对应关系的JSON
func (trs *TransactionRecords) LoadTransactionRecords() error {
	// 读文件
	file, e := ioutil.ReadFile(trs.FilePath)
	if e != nil {
		return e
	}

	// 反序列化JSON
	json.Unmarshal(file, &trs.TxRecord)
	return nil
}
