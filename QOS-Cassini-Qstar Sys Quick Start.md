# 部署 QOS，Qstar和 Cassini 跨链网络
## 配置步骤
QOS -> Qstar -> cassini

## 1、QOS

参照 [QOS文档](http://docs.qoschain.info/qos/)

#### 1）启动 QOS

```
## qos 初始化
./qosd init --chain-id test --moniker test

## 为QSC，QCP写入证书
./cmd/qosd config-root-ca --qcp ./ca/root.pub --qsc ./ca/root.pub

## 为 benchmark工具创建账号
./config.sh

## 添加验证人
./qosd add-genesis-validator --name test --owner address1atvp97wwcar9vykxmnqdh75hz0ns7rkt8wctnf --pubkey X4a06QLUtkqb7LKCgV4iuMdqdlsSWQvDBoEVj5P1lXw= --tokens 10
```

#### 2）配置联盟链
```
## 获得 CA 证书之后
./qoscli tx init-qcp --creator bench_test(key name) --qcp.crt ~/ca/qcp.crt(qcp crt file)
./qoscli tx create-qsc --creator bench_test(key name) --qsc.crt ~/ca/qsc.crt(qsc crt file) --accounts bench_test,100000000(init qsc account)
```

使用上述两条命令可获得 bench_test 账户下的指定金额的 QSC 代币，代币名称在 qsc.crt 中指定。

## 2、Qstar
参照 [Qstar文档](http://docs.qoschain.info/qstars/)

其中需补充：

```
## 创建 .qstarcli 目录，需补充如下说明
./qstarscli rest-server ## 该步骤创建 ～/.qstarcli 目录

## Configure ~/.qstarsd/config/qstarsconf.toml 
## //指代 QStarsPrivateKey 部分，需补充说明该值为 qcp.pri 的value
vi ~/.qstarsd/config/qstarsconf.toml ## 二者不可或缺
QStarsPrivateKey = "qcp.pri"
QOSChainName = "qos-test"

## Configure ~/.qstarsd/config/genesis.json 
## //指代 qcps 部分，需补充如下说明
qcps {
	"name": "qstar-name",
	"chain-id": "qstar-chain-id"
	"pub_key": {
              "type": "tendermint/PubKeyEd25519",
              "value": "relay.pub"
    }
}
```
Qstar 的启动参照上述链接的说明文档

## 3、Cassini
参照 [Cassini Quick Start](https://github.com/QOSGroup/cassini/blob/master/docs/quick_start.md)

其中需补充：

#### 1）gnatsd的编译和配置

[gnats文档](https://github.com/nats-io/gnatsd)

```
## 编译 gnats
go get github.com/nats-io/gnatsd
cd $GOPATH/src/github.com/nats-io/gnatsd
go build

## 启动 gnats
nohup ./gnatsd -p 4222 -cluster nats://192.168.1.201:5222 &
```
#### 2）配置 Cassini 的 config 文件

```
vi $GOPATH/src/github.com/cassini/config/config.conf

{
    "nats": "nats://127.0.0.1:4222", ## nats 为Cassini提供消息队列服务，可以是多个 nats 节点，可以和Cassini剥离部署于不同的机器，和Cassini部署于同一个机器可以用localhost，多个IP用‘,’隔开
    "prikey":    "Yv7F3Vt/eJT55TMq9ItybOko6I6TfCGoc0LpqJtbbaLmio/puRig7bAxzhhwqXEEcbP33UN6I3uCjHWs5aDixQ==", ##relay.pri
    "consensus": true,
    "eventWaitMillitime": 2000,
    "useEtcd":true,
    "lock":"etcd://192.168.1.201:2379", ## etcd是cassini内置的分布式锁，对‘交易的处理’进行分布式锁来控制
    "embedEtcd":true,
    "etcd":{
        "name": "dev-cassini", ## any name
        "advertise":"http://192.168.1.201:2379", ## 客户端端口
        "advertisePeer":"http://192.168.1.201:2380", ## 集群通信端口
        "clusterToken":"dev-cassini-cluster", ## any token
        "cluster":"dev-cassini=http://192.168.1.201:2380" ## 所有的 etcd node，“NAME=URI”
    },
    "qscs": [
        {
            "name": "qos-test", ## qos chain id
            "type": "qos",
            "signature": true,
            "nodes": "192.168.1.203:26657" ## qos 如果是个网络，可以写多个IP，单个也可
        },
        {
            "name": "qstars-test", ## qstar chain id
            "type": "qos",
            "nodes": "192.168.1.200:26657" ## qstar ip
        }
    ]
}
```
#### 3）启动 Cassini

```
## 启动 Cassini 服务
#### 启动 gnets（如果没启动，则先启动 qnatsd）
nohup ./gnatsd -p 4222 -cluster nats://192.168.1.201:5222 &

#### 启动 cassini
nohup ./cassini start --config ./config/config.conf &

```

此时，qstar上发起的转账交易，会在qos上被确认和记录，跨链系统搭建完成
