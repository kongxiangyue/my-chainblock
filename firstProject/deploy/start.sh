#!/bin/bash

# 遇到错误停止程序
set -e

# 清理旧环境
./stop.sh

echo "第一步：生成配置文件"
cryptogen generate --config=./crypto-config.yaml

echo "第二步：生成起始区块"
mkdir -p channel-artifacts
configtxgen -profile TwoOrgsOrdererGenesis \
-channelID syschannel -outputBlock \
./channel-artifacts/genesis.block

echo "第三步：生成通道配置文件"
configtxgen -profile TwoOrgsChannel \
-outputCreateChannelTx ./channel-artifacts/channel.tx \
-channelID mychannel

echo "第四步：生成锚节点配置文件"
echo "生成Organization1MSP锚节点配置文件"
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate \
./channel-artifacts/Organization1MSPanchors.tx \
-channelID mychannel -asOrg Organization1MSP

echo "生成Organization2MSP锚节点配置文件"
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate \
./channel-artifacts/Organization2MSPanchors.tx \
-channelID mychannel -asOrg Organization2MSP

echo "生成Organization3MSP锚节点配置文件"
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate \
./channel-artifacts/Organization3MSPanchors.tx \
-channelID mychannel -asOrg Organization3MSP

echo "第五步：启动docker容器"
docker-compose -f docker-compose-cli.yaml up -d

echo "打印正在运行的docker容器"
docker ps -a

echo "第六步： 根据通道配置文件生成通道"
docker exec cli peer channel create -o orderer.test.com:7050 \
-c mychannel -f ./channel-artifacts/channel.tx

echo "第七步： 将节点加入通道"
echo "将node1.Organization1.test.com 加入通道"
docker exec cli peer channel join -b mychannel.block

echo "将node2.Organization1.test.com 加入通道"
docker exec -e CORE_PEER_ADDRESS=node2.Organization1.test.com:7051 cli \
peer channel join -b mychannel.block

echo "将node1.Organization2.test.com 加入通道"
docker exec -e CORE_PEER_ADDRESS=node1.Organization2.test.com:7051 \
-e CORE_PEER_LOCALMSPID=Organization2MSP \
-e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/Organization2.test.com/users/Admin@Organization2.test.com/msp \
cli peer channel join -b mychannel.block

echo "将node2.Organization2.test.com 加入通道"
docker exec -e CORE_PEER_ADDRESS=node2.Organization2.test.com:7051 \
-e CORE_PEER_LOCALMSPID=Organization2MSP \
-e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/Organization2.test.com/users/Admin@Organization2.test.com/msp \
cli peer channel join -b mychannel.block

echo "将node1.Organization3.test.com 加入通道"
docker exec -e CORE_PEER_ADDRESS=node1.Organization3.test.com:7051 \
-e CORE_PEER_LOCALMSPID=Organization3MSP \
-e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/Organization3.test.com/users/Admin@Organization3.test.com/msp \
cli peer channel join -b mychannel.block

echo "将node2.Organization3.test.com 加入通道"
docker exec -e CORE_PEER_ADDRESS=node2.Organization3.test.com:7051 \
-e CORE_PEER_LOCALMSPID=Organization3MSP \
-e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/Organization3.test.com/users/Admin@Organization3.test.com/msp \
cli peer channel join -b mychannel.block

echo "第八步： 根据锚节点配置文件更新锚节点"
echo "更新Organization1锚节点"
docker exec cli peer channel update -o orderer.test.com:7050 -c mychannel \
-f ./channel-artifacts/Organization1MSPanchors.tx

echo "更新Organization2锚节点"
docker exec -e CORE_PEER_ADDRESS=node1.Organization2.test.com:7051 \
-e CORE_PEER_LOCALMSPID=Organization2MSP \
-e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/Organization2.test.com/users/Admin@Organization2.test.com/msp \
cli peer channel update -o orderer.test.com:7050 -c mychannel \
-f ./channel-artifacts/Organization2MSPanchors.tx

echo "更新Organization3锚节点"
docker exec -e CORE_PEER_ADDRESS=node1.Organization3.test.com:7051 \
-e CORE_PEER_LOCALMSPID=Organization3MSP \
-e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/Organization3.test.com/users/Admin@Organization3.test.com/msp \
cli peer channel update -o orderer.test.com:7050 -c mychannel \
-f ./channel-artifacts/Organization3MSPanchors.tx

echo "第九步： 安装链码"
echo "将链码安装到node1.Organization1"
docker exec cli peer chaincode install -n mycc -v 1.0 -l golang \
-p "github.com/chaincode/chaincode_example02/go/"

echo "将链码安装到node1.Organization2"
docker exec -e CORE_PEER_ADDRESS=node1.Organization2.test.com:7051 \
-e CORE_PEER_LOCALMSPID=Organization2MSP \
-e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/Organization2.test.com/users/Admin@Organization2.test.com/msp \
cli peer chaincode install -n mycc -v 1.0 -l golang \
-p "github.com/chaincode/chaincode_example02/go/"

echo "将链码安装到node1.Organization3"
docker exec -e CORE_PEER_ADDRESS=node1.Organization3.test.com:7051 \
-e CORE_PEER_LOCALMSPID=Organization3MSP \
-e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/Organization3.test.com/users/Admin@Organization3.test.com/msp \
cli peer chaincode install -n mycc -v 1.0 -l golang \
-p "github.com/chaincode/chaincode_example02/go/"

echo "第十步： 实例化链码"
docker exec cli peer chaincode instantiate -o orderer.test.com:7050 \
-C mychannel -n mycc -l golang -v 1.0 -c '{"Args":["init","a","100","b","200"]}' \
-P "AND ('Organization1MSP.peer','Organization2MSP.peer','Organization3MSP.peer')"