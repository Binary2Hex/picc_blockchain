## 如何搭建1 peer+ memberservice的环境用于chaincode开发(以chaincode-example02为例)
在docker-compose目录下输入一条命令`docker-compose -f docker-compose-dev.yml`即可

#### 登录admin到网络
```
POST localhost:5000/registrar
{
  "enrollId": "admin",
  "enrollSecret": "Xurw3yU9zI0l"
}
```
#### go build your chaincode
```
cd $GOPATH/src/github.com/chaincode_example02
go build
```
#### 启动chaincode 
```
CORE_CHAINCODE_ID_NAME=mycc CORE_PEER_ADDRESS=0.0.0.0:30303 ./chaincode_example02
```
#### 部署chaincode
```
POST localhost:5000/chaincode
{
  "jsonrpc": "2.0",
  "method": "deploy",
  "params": {
    "type": 1,
    "chaincodeID":{
        "name": "mycc"
    },
    "ctorMsg": {
        "function":"init",
        "args":["a", "100", "b", "200"]
    },
    "secureContext": "admin"
  },
  "id": 1
}
```
#### invoke chaincode
```
POST localhost:5000/chaincode
{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
      "type": 1,
      "chaincodeID":{
          "name":"mycc"
      },
      "ctorMsg": {
         "function":"invoke",
         "args":["a", "b", "10"]
      },
      "secureContext": "admin"
  },
  "id": 3
}
```
#### query chaincode
```
POST localhost:5000/chaincode
{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
      "type": 1,
      "chaincodeID":{
          "name":"mycc"
      },
      "ctorMsg": {
         "function":"query",
         "args":["a"]
      },
      "secureContext": "admin"
  },
  "id": 5
}
```