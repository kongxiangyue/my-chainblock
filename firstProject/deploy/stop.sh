#!/bin/bash

# 停止相关容器
docker-compose -f docker-compose-cli.yaml down --volumes --remove-orphans
# 删除相关容器
docker rm -f $(docker ps -a | awk '($2 ~ /dev-node.*/) {print $1}')
# 删除相关镜像
docker rmi -f $(docker images | awk '($1 ~ /dev-node.*/) {print $3}')

# 清除证书与各种配置文件
rm -rf channel-artifacts/*
rm -rf crypto-config/*