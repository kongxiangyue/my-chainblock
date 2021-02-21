#!/bin/bash

# 本脚本从头构建一个区块链网络
# 请确保cryptogen 和 configtxgen 这两个可执行文件已经被正确安装
# 创建一个通道 assetschannel

# 遇到错误直接退出程序
set -e

# 一、环境清理
echo "一、环境清理"
mkdir -p config
mkdir -p crypto-config
rm -fr config/*
rm -fr crypto-config/*
echo "清理完毕"

# 二、生成证书和起始区块信息
echo "二、生成证书和起始区块信息"
cryptogen generate --config=./crypto-config.yaml
configtxgen -profile OneOrgOrdererGenesis -outputBlock ./config/genesis.block

# 三、启动区块链网络
echo "三、区块链 ： 启动"
docker-compose up -d        # 按照docker-compose.yaml的配置启动区块链网络并在后台运行
echo "正在等待节点的启动完成，等待10秒"
sleep 10                    # 启动整个区块链网络需要一点时间，所以此处等待10s，让区块链网络完全启动

# 四、生成通道(这个动作会创建一个创世交易，也是该通道的创世交易)
# channel === 通道
echo "四、生成通道的TX文件(这个动作会创建一个创世交易，也是该通道的创世交易)"
configtxgen -profile TwoOrgChannel -outputCreateChannelTx ./config/assetschannel.tx -channelID assetschannel

# 五、在区块链上按照刚刚生成的TX文件去创建通道
# 该操作和上面操作不一样的是，这个操作会写入区块链
echo "五、在区块链上按照刚刚生成的TX文件去创建通道"
docker exec cli peer channel create -o orderer.educhain.accurchain.com:7050 -c assetschannel -f /etc/hyperledger/config/assetschannel.tx

# 六、让节点去加入到通道
echo "六、让节点去加入到通道"
docker exec cli peer channel join -b assetschannel.block
docker exec -e CORE_PEER_ADDRESS=peer1.org1.educhain.accurchain.com:7051 cli peer channel join -b assetschannel.block

docker exec -e CORE_PEER_ADDRESS=peer0.org0.educhain.accurchain.com:7051 -e CORE_PEER_LOCALMSPID=Org0MSP -e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/org0.educhain.accurchain.com/users/Admin@org0.educhain.accurchain.com/msp cli peer channel join -b assetschannel.block
docker exec -e CORE_PEER_ADDRESS=peer1.org0.educhain.accurchain.com:7051 -e CORE_PEER_LOCALMSPID=Org0MSP -e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/org0.educhain.accurchain.com/users/Admin@org0.educhain.accurchain.com/msp cli peer channel join -b assetschannel.block

docker exec -e CORE_PEER_ADDRESS=peer0.org2.educhain.accurchain.com:7051 -e CORE_PEER_LOCALMSPID=Org2MSP -e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/org2.educhain.accurchain.com/users/Admin@org2.educhain.accurchain.com/msp cli peer channel join -b assetschannel.block
docker exec -e CORE_PEER_ADDRESS=peer1.org2.educhain.accurchain.com:7051 -e CORE_PEER_LOCALMSPID=Org2MSP -e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/org2.educhain.accurchain.com/users/Admin@org2.educhain.accurchain.com/msp cli peer channel join -b assetschannel.block

# 七、链码安装
# -n 是链码的名字，可以自己随便设置
# -v 就是版本号，就是composer的bna版本
# -p 是目录，目录是基于cli这个docker里面的$GOPATH相对的
# 此处安装的是示例链码，后续课程会自己编写
echo "七、链码安装"
#docker exec cli peer chaincode install -n assets -v 1.0.0 -l golang -p github.com/chaincode/sacc
docker exec cli peer chaincode install -n food -v 1.0.0 -l golang -p github.com/chaincode/perishable-food-full
docker exec -e CORE_PEER_ADDRESS=peer1.org1.educhain.accurchain.com:7051 cli peer chaincode install -n food -v 1.0.0 -l golang -p github.com/chaincode/perishable-food-full

docker exec -e CORE_PEER_ADDRESS=peer0.org0.educhain.accurchain.com:7051 -e CORE_PEER_LOCALMSPID=Org0MSP -e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/org0.educhain.accurchain.com/users/Admin@org0.educhain.accurchain.com/msp cli peer chaincode install -n food -v 1.0.0 -l golang -p github.com/chaincode/perishable-food-full
docker exec -e CORE_PEER_ADDRESS=peer1.org0.educhain.accurchain.com:7051 -e CORE_PEER_LOCALMSPID=Org0MSP -e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/org0.educhain.accurchain.com/users/Admin@org0.educhain.accurchain.com/msp cli peer chaincode install -n food -v 1.0.0 -l golang -p github.com/chaincode/perishable-food-full

docker exec -e CORE_PEER_ADDRESS=peer0.org2.educhain.accurchain.com:7051 -e CORE_PEER_LOCALMSPID=Org2MSP -e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/org2.educhain.accurchain.com/users/Admin@org2.educhain.accurchain.com/msp cli peer chaincode install -n food -v 1.0.0 -l golang -p github.com/chaincode/perishable-food-full
docker exec -e CORE_PEER_ADDRESS=peer1.org2.educhain.accurchain.com:7051 -e CORE_PEER_LOCALMSPID=Org2MSP -e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/org2.educhain.accurchain.com/users/Admin@org2.educhain.accurchain.com/msp cli peer chaincode install -n food -v 1.0.0 -l golang -p github.com/chaincode/perishable-food-full



#八、实例化链码
#-n 对应前文安装链码的名字 其实就是composer network start bna名字
#-v 为版本号，相当于composer network start bna名字@版本号
#-C 是通道，在fabric的世界，一个通道就是一条不同的链，composer并没有很多提现这点，composer提现channel也就在于多组织时候的数据隔离和沟通使用
#-c 为传参，传入init参数
echo "八、实例化链码"
#docker exec cli peer chaincode instantiate -o orderer.educhain.accurchain.com:7050 -C assetschannel -n assets -l golang -v 1.0.0 -c '{"Args":["a","10"]}'
docker exec cli peer chaincode instantiate -o orderer.educhain.accurchain.com:7050 -C assetschannel -n food -l golang -v 1.0.0 -c '{"Args":["init"]}' -P 'AND("Org0MSP.member","Org1MSP.member","Org2MSP.member")'
#sleep 2
#请注意，安装链码是文件的复制，其实不等于我们电脑的安装，实例化才是真正的安装

# 进行链码交互，验证链码是否正确安装及区块链网络能否正常工作
#echo "修改 a-->20"
#docker exec cli peer chaincode invoke -C assetschannel -n assets -c '{"Args":["set", "a", "20"]}'
#sleep 2

# 进行链码查询，查询刚刚执行的智能合约是否将数据写入区块链
#echo "查询: "
#docker exec cli peer chaincode query -C assetschannel -n assets -c '{"Args":["query", "a"]}'
