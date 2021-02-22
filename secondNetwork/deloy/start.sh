#!/bin/bash

SYS_CHANNEL="sys-channel"
CHANNEL_NAME="assetschannel"
echo $SYS_CHANNEL

echo "一、环境清理"
mkdir -p config
mkdir -p crypto-config
mkdir -p channel-artifacts

rm -fr config/*
rm -fr crypto-config/*
rm -fr channel-artifacts/*
echo "清理完毕"


echo "二、生成证书和起始区块信息"
cryptogen generate --config=./crypto-config.yaml
#configtxgen -profile TwoOrgsOrdererGenesis -outputBlock -channelID  ./channel-artifacts/genesis.block
# -channelID $SYS_CHANNEL
configtxgen -profile OneOrgsOrdererGenesis -outputBlock ./channel-artifacts/genesis.block

echo "三 生成通道配置文件"
mkdir -p channel-artifacts
configtxgen -profile TwoOrgsChannel \
-outputCreateChannelTx ./channel-artifacts/channel.tx \
-channelID $CHANNEL_NAME

echo "四 生成锚节点配置文件"
echo "生成Organization0MSP锚节点配置文件"
configtxgen -profile TwoOrgsChannel \
-outputAnchorPeersUpdate ./channel-artifacts/Organization1MSPanchors.tx \
-channelID $CHANNEL_NAME \
-asOrg Organization0MSP


echo "生成Organization1MSP锚节点配置文件"
configtxgen -profile TwoOrgsChannel \
-outputAnchorPeersUpdate ./channel-artifacts/Organization1MSPanchors.tx \
-channelID $CHANNEL_NAME \
-asOrg Organization1MSP

echo "生成Organization2MSP锚节点配置文件"
configtxgen -profile TwoOrgsChannel \
-outputAnchorPeersUpdate ./channel-artifacts/Organization1MSPanchors.tx \
-channelID $CHANNEL_NAME \
-asOrg Organization2MSP

echo "五 启动docker容器"
docker-compose -f docker-compose-cli.yaml up -d




















