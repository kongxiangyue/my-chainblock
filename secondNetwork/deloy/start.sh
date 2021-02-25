#!/bin/bash

SYS_CHANNEL="sys-channel"
CHANNEL_NAME="assetschannel"
echo $SYS_CHANNEL

echo "一 环境清理"
mkdir -p config
mkdir -p crypto-config
mkdir -p channel-artifacts

rm -fr config/*
rm -fr crypto-config/*
rm -fr channel-artifacts/*
echo "清理完毕"


echo "二 生成证书和起始区块信息"
cryptogen generate --config=./crypto-config.yaml
#configtxgen -profile TwoOrgsOrdererGenesis -outputBlock -channelID  ./channel-artifacts/genesis.block
configtxgen -profile OneOrgsOrdererGenesis -outputBlock ./channel-artifacts/genesis.block -channelID $SYS_CHANNEL

echo "三 生成通道配置文件"
mkdir -p channel-artifacts
configtxgen -profile TwoOrgsChannel \
-outputCreateChannelTx ./channel-artifacts/$CHANNEL_NAME.tx \
-channelID $CHANNEL_NAME


echo "四 生成锚节点配置文件"
echo "生成Organization0MSP锚节点配置文件"
configtxgen -profile TwoOrgsChannel \
-outputAnchorPeersUpdate ./channel-artifacts/Organization1MSPanchors.tx \
-channelID $CHANNEL_NAME \
-asOrg Org0MSP


echo "生成Organization1MSP锚节点配置文件"
configtxgen -profile TwoOrgsChannel \
-outputAnchorPeersUpdate ./channel-artifacts/Organization1MSPanchors.tx \
-channelID $CHANNEL_NAME \
-asOrg Org1MSP

echo "生成Organization2MSP锚节点配置文件"
configtxgen -profile TwoOrgsChannel \
-outputAnchorPeersUpdate ./channel-artifacts/Organization1MSPanchors.tx \
-channelID $CHANNEL_NAME \
-asOrg Org2MSP

echo "五 启动docker容器"
docker-compose -f docker-compose-cli.yaml up -d
#docker-compose up -d
sleep 10

echo "打印正在运行的docker容器"
docker ps -a


echo "六 根据通道配置文件生成通道"
docker exec cli peer channel create -o orderer.educhain.accurchain.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/$CHANNEL_NAME.tx


echo "七 将节点加入通道"

echo "将peer0.org1.educhain.accurchain.com 加入通道"
docker exec cli peer channel join -b $CHANNEL_NAME.block

echo "将peer1.org1.educhain.accurchain.com 加入通道"
docker exec -e CORE_PEER_ADDRESS=peer1.org1.educhain.accurchain.com:7051 cli \
peer channel join -b $CHANNEL_NAME.block


echo "将peer0.org0.educhain.accurchain.com 加入通道"
docker exec -e CORE_PEER_ADDRESS=peer0.org0.educhain.accurchain.com:7051 \
-e CORE_PEER_LOCALMSPID=Org0MSP \
-e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/org0.educhain.accurchain.com/users/Admin@org0.educhain.accurchain.com/msp \
cli peer channel join -b $CHANNEL_NAME.block

echo "将peer1.org0.educhain.accurchain.com 加入通道"
docker exec -e CORE_PEER_ADDRESS=peer1.org0.educhain.accurchain.com:7051 \
-e CORE_PEER_LOCALMSPID=Org0MSP \
-e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/org0.educhain.accurchain.com/users/Admin@org0.educhain.accurchain.com/msp \
cli peer channel join -b $CHANNEL_NAME.block


echo "将peer0.org2.educhain.accurchain.com 加入通道"
docker exec -e CORE_PEER_ADDRESS=peer0.org2.educhain.accurchain.com:7051 \
-e CORE_PEER_LOCALMSPID=Org2MSP \
-e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/org2.educhain.accurchain.com/users/Admin@org2.educhain.accurchain.com/msp \
cli peer channel join -b $CHANNEL_NAME.block

echo "将peer1.org2.educhain.accurchain.com 加入通道"
docker exec -e CORE_PEER_ADDRESS=peer1.org2.educhain.accurchain.com:7051 \
-e CORE_PEER_LOCALMSPID=Org2MSP \
-e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/org2.educhain.accurchain.com/users/Admin@org2.educhain.accurchain.com/msp \
cli peer channel join -b $CHANNEL_NAME.block





















