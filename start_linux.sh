#!/bin/bash
echo "停止consul服务..."
docker-compose down

echo "启动consul容器..."
docker-compose up -d

echo "等待5秒consul初始化结束"
sleep 5

echo "导入服务配置..."
curl --request PUT --data @StartSceneConfig.txt http://localhost:8500/v1/kv/service_config

echo "success..."