package controller

import (
	bc "accurchain.com/perishable-food-full/application/blockchain"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"net/http"
)

//
// 区块链数据查询
//

type queryBlockByTxIDRequest struct {
	TxID fab.TransactionID `json:"tx_id" form:"tx_id" binding:"required"`
}

// 依据txid查询所在区块
func QueryBlockByTxID(ctx *gin.Context) {
	// 解析请求体
	req := new(queryBlockByTxIDRequest)
	if err := ctx.ShouldBindQuery(req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// 创建上下文，表明查询哪个通道，自己属于哪个组织，用户是谁
	bcCtx := bc.SDK.ChannelContext(bc.ChannelName, fabsdk.WithOrg(bc.Org), fabsdk.WithUser(bc.User))

	// 获取一个账本客户端对象
	cli, err := ledger.New(bcCtx)
	if err != nil {
		panic(err)
	}

	// 表明该客户端要查询的目标节点
	_, err = cli.QueryInfo(ledger.WithTargetEndpoints("peer0.org1.educhain.accurchain.com"))
	if err != nil {
		panic(err)
	}

	// 使用txid查询所在的block
	block, err := cli.QueryBlockByTxID(req.TxID)

	// 将结果返回
	ctx.JSON(http.StatusOK, block)
}

type queryTransactionRequest struct {
	TxID fab.TransactionID `json:"tx_id" form:"tx_id" binding:"required"`
}

// 根据txid查询该笔交易的详细信息
func QueryTransaction(ctx *gin.Context) {
	// 解析请求体
	req := new(queryTransactionRequest)
	if err := ctx.ShouldBindQuery(req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// 创建上下文，表明查询哪个通道，自己属于哪个组织，用户是谁
	bcCtx := bc.SDK.ChannelContext(bc.ChannelName, fabsdk.WithOrg(bc.Org), fabsdk.WithUser(bc.User))

	// 获取一个账本客户端对象
	cli, err := ledger.New(bcCtx)
	if err != nil {
		panic(err)
	}

	// 表明该客户端要查询的目标节点
	_, err = cli.QueryInfo(ledger.WithTargetEndpoints("peer0.org1.educhain.accurchain.com"))
	if err != nil {
		panic(err)
	}

	// 根据txid查询该笔交易的详细信息
	block, err := cli.QueryTransaction(req.TxID)

	// 将结果返回
	ctx.JSON(http.StatusOK, block)
}
