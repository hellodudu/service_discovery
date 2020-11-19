#!/bin/bash
echo "停止consul服务..."
docker-compose down

echo "启动consul容器..."
docker-compose up -d

echo "导入服务配置..."
curl --request PUT --data @StartSceneConfig.txt http://localhost:8500/v1/kv/service_config

echo "success..."