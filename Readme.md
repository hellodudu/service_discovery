# 服务发现中间件consul

## 目录结构

- 依赖组件:
    * `docker` 
    * `docker-compose` 
    * `curl`

- 全服统一配置文件:
    * `StartSceneConfig.txt`

- 执行脚本开启后可访问`http://localhost:8500`进行操作

- 连接服务配置文件:
    * `config/service.json`

## 新增服务流程
1. 更改全服统一配置文件`StartSceneConfig.txt`
2. `config/service.json`新增`services`相关配置
3. `config/service.json`新增`watches`相关配置
> 举例如果要新增一台`Gate3`Server，先将`StartSceneConfig.txt`表全服更新，然后在`config/service.json`中添加新的`services`和`watches`，`service`中的`check http`必须保证是通的，`watches`要注意`key`和`service`类型都需要添加，否则新加入的服务无法正确添加到ServiceMesh中，无法被其他服务访问到。