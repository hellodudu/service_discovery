@echo off
REM 后续命令使用的是：UTF-8编码
chcp 65001
echo .
echo 停止consul服务...
docker-compose down

echo .
echo 启动consul容器...
docker-compose up -d

@REM echo .
@REM echo 等待5秒consul初始化结束...
@REM timeout /t 5

@REM echo .
@REM echo 导入服务配置...
@REM curl.exe --request PUT --data "@StartSceneConfig.txt" http://localhost:8500/v1/kv/service_config

echo .
echo success...
pause