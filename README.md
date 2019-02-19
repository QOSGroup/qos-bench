# QOS benchmark tool

## 依赖
本工具运行在QOS节点上，安装步骤请参照：

[QOS节点安装文档](https://github.com/QOSGroup/qos/blob/master/docs/install/installation.md
)

## Benchmark 工具安装
```bash
$ git clone https://github.com/QOSGroup/qos-bench
$ cd qos-bench
$ env GO111MODULE=on go build
```

## 配置
因为QOS公链采用账户机制，所以开始测试前需准备测试账号，并保证测试账号中有充足的余额

使用该脚本来创建账户和配置账户资产
```bash
$ ./config.sh
```

该脚本在 QOS 网络中配置账户和资产，并在当前目录生成 `config.json` 文件，文件中记录账户名、地址、密码。


## 运行测试工具
```bash
$ ./qos-bench -v -T 10 -R 10 -home "~/.qoscli" -file "./config.json" localhost:26657
```

## 打印帮助信息
```bash
$ ./qos-bench -h
```
