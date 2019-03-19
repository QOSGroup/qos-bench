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

注：本配置需在创世区块出块前完成


## 运行测试工具
```bash
$ ./qos-bench -v -T 10 -R 10 -home "~/.qoscli" -file "./config.json" localhost:26657
```

## 打印帮助信息
```bash
$ ./qos-bench -h
```



## QOS、Qstar、cassini跨链
对应下述系统参数的压测 TPS 值参考：

|     | qos   | qstar | 中继跨链 |
| --- | :---: | :---: | :---:   |
| 平均TPS | 1000	| 300 | 300 |
| 节点网络配置 | 单节点 | 单节点 | 单qos节点网络，单中继，单qstar节点网络 |
| 节点硬件参数 | 2197.540 MHz 双核CPU * 1 | 2197.306 MHz 双核CPU * 1 | 2197 MHz 双核CPU * 1 |
| 节点带宽 | 1000Mb/s | 1000Mb/s	| 1000Mb/s |
| 测试交易类型 | 单个用户转账 | 单个用户转账 | 单个用户转账 |
| 发币场景 | 网络初始化发币 | CA证书 | CA证书 |

