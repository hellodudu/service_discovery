# 服务发现中间件consul

## 目录结构

- 依赖组件:
    * `docker` 
    * `docker-compose` 
    * `curl`

- 全服统一配置文件:
    * `config/scene/StartSceneConfig.txt`
    * `config/scene/StartProcessConfig.txt`
    * `config/scene/StartMachineConfig.txt`

- 执行脚本开启后可访问`http://localhost:8500`进行操作

- consul服务配置文件:
    * `config/consul/service.json`

## 新增服务流程
1. 更改`StartSceneConfig.txt`，`StartProcessConfig.txt`，`StartMachineConig.txt`后，更新到线上服务器的`config/scene`目录中
2. 执行不同系统对应的`start`脚本，脚本会将`txt`配置文件合并转换为`consul`可读取的配置文件`service.json`，并且开启`consul`的`docker`镜像，服务发现功能将可用
3. 如果要热更配置的话，先执行第一步，然后执行系统对应的`reload`脚本
> 注意：在本机测试时如果服务的`InnerIP`配置为`127.0.0.1`，则需要将`config/consul/service.json`文件中所有`127.0.0.1`的ip地址修改为本地局域网地址例如`192.168.1.137`，因为在`docker`镜像中解析`localhost`地址解析不到宿主机中，所以本地测试时最好将`InnerIP`配置为本机局域网ip地址