#!/bin/bash
echo "转换txt配置文件为service.json..."
./bin/exporter_linux -import_path=./config/scene/ -export_path=./config/consul/

echo "导入新配置"
curl --request PUT http://127.0.0.1:8500/v1/agent/reload

echo "success..."