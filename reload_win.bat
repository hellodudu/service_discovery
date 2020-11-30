@echo off
REM 后续命令使用的是：UTF-8编码
chcp 65001

REM "转换txt配置文件为service.json..."
.\bin\exporter_win.exe -import_path=.\config\scene\ -export_path=.\config\consul\

REM "导入新配置"
curl --request PUT http://127.0.0.1:8500/v1/agent/reload

echo .
echo reload success...
pause