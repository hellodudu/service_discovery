@echo off
REM 后续命令使用的是：UTF-8编码
chcp 65001
echo .
echo 停止consul服务...
docker-compose down

REM "转换txt配置文件为service.json..."
.\bin\exporter_win.exe -import_path=.\config\scene\ -export_path=.\config\consul\

echo .
echo 启动consul容器...
docker-compose up -d

echo .
echo success...
pause