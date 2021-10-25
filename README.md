# git-pull-hook
> 用于配置git钩子，用在测试环境时，有push代码时，自动部署测试环境代码


## 编译
> 生成可执行文件名为：git-pull
- 当前机器上编译使用
    > `make`
- 自选：linux 环境
    > `make build-linux`
- 自选：Mac 环境
    > `make build-mac`

## 运行
> `nohup /some_path/git-pull > /dev/null 2>&1 &`
## 配置
```
    项目配置文件为：config.yaml

    配置规则如下：


    # 项目名 :  项目路径
    # 每个项目，配置一行
    如 name1: /path1

    注：编辑配置文件实时生效，无需重启服务
```
## 日志

    运行目录下的 `run.log` 即为运行日志，可查看log信息


## 使用

    运行http端口为 `8081`

    git web hook ，配置链接为 `ip:8081/git_pull?project=name1`

    此处的name1，即为配置文件的中的名字