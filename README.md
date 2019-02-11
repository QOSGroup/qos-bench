# qos-bench# qos-bench
## install
git clone https://github.com/Shawncles/qos-bench
cd qos-bench
env GO111MODULE=on go build

## run test
./qos-bench -T 10 -r 10 localhost:26657

## qos-bench help
./qos-bench -h
