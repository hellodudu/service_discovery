@echo off
REM 后续命令使用的是：UTF-8编码
chcp 65001
echo .
echo 停止consul服务...
docker-compose down

echo .
echo 启动consul容器...
docker-compose up -d

echo .
echo 等待5秒consul初始化结束...
timeout /t 5

echo .
echo 导入服务配置...
curl.exe --request PUT --data "@StartSceneConfig.txt" http://localhost:8500/v1/kv/service_config

echo .
echo success...
pause