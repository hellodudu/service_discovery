#!/bin/bash
echo "停止consul服务..."
docker-compose down

echo "转换txt配置文件为service.json..."
./bin/exporter_mac -import_path=./config/scene/ -export_path=./config/consul/

echo "启动consul容器..."
docker-compose up -d

echo "success..."