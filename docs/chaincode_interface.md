#Chaincode接口描述(API描述)

此文档用于描述目录[chaincode/picc/src/main](../chaincode/picc/src/main)下的chaincode的接口

##Contents

* [Deploy](#deploy)
	* [监管机构部署合约](#监管机构部署合约)

* [Invoke](#invoke)

	
* [Query](#query)
	* [查询某一养殖场的信息](#查询某一养殖场的信息)
	* [查询网络中所有养殖场的数量](#查询网络中所有养殖场的数量)
	* [查询某一养殖场的所有肉牛信息](#查询某一养殖场的所有肉牛信息)
	* [查询某一养殖场的所有保险信息](#查询某一养殖场的所有保险信息)
    * [查询某一养殖场的所有贷款信息](#查询某一养殖场的所有贷款信息)
    * [查询某一城市的所有养殖场ID](#查询某一城市的所有养殖场id)
    * [查询某一省的所有养殖场ID](#查询某一省的所有养殖场id)
    * [通过养殖场ID和肉牛耳标查询肉牛信息](#通过养殖场id和肉牛耳标查询肉牛信息)
    * [通过养殖场名字查询养殖场ID](#通过养殖场名字查询养殖场id)
    * [查询放贷员的放贷记录ID](#查询放贷员的放贷记录id)

	
##Deploy
###监管机构部署合约
#####Chaincode Spec: 

	{
		"jsonrpc": "2.0",
		"method": "deploy",
		"params": {
			"type": 1,
			"chaincodeID": {
				"path": “path/to/chaincode”
			},
			"ctorMsg": {
				"function": "init",
				"args": []
			},
			"secureContext": "admin"
		},
		"id": 1
	}


##Invoke


##Query

###查询某一养殖场的信息
输入参数`farmID`, 比如`1234567`或者`1234568`

返回符合查询条件的`Farm`对象字符串, 参考[farm.proto](../chaincode/picc/src/main/farm.proto)
#####Chaincode Spec: 
	{
		"jsonrpc": "2.0",
		"method": "query",
		"params": {
			"type": 1,
			"chaincodeID": {
				"name":  <chaincode_name>
			},
			"ctorMsg": {
				"function": "getFarmById",
				"args": ["farmID"]
			},
			"secureContext": <role>
		},
		"id": 1
	}

###查询网络中所有养殖场的数量
输入参数: 无

返回参数: blockchain网络中记录的养殖场总数
#####Chaincode Spec: 
    {
      "jsonrpc": "2.0",
      "method": "query",
      "params": {
        "type": 1,
        "chaincodeID":{
            "name": <chaincode_name>
        },
        "ctorMsg": {
            "function":"getFarmAmount",
            "args":[]
        },
        "secureContext": "admin"
      },
      "id": 1
    }
			
###查询某一养殖场的所有肉牛信息
输入参数: 养殖场id `farmID`, 比如`1234567`

返回参数: 包含Beef对象(参考[beef.proto](../chaincode/picc/src/main/beef.proto))的json数组
#####Chaincode Spec: 
    {
      "jsonrpc": "2.0",
      "method": "query",
      "params": {
        "type": 1,
        "chaincodeID":{
            "name": <chaincode_name>
        },
        "ctorMsg": {
            "function":"getAllBeevesByFarm",
            "args":[farmID]
        },
        "secureContext": "admin"
      },
      "id": 1
    }
    
###查询某一养殖场的所有保险信息
输入参数: 养殖场id `farmID`, 比如`1234567`

返回参数: 包含Insurance对象(参考[insurance.proto](../chaincode/picc/src/main/insurance.proto))的json数组
#####Chaincode Spec: 
    {
      "jsonrpc": "2.0",
      "method": "query",
      "params": {
        "type": 1,
        "chaincodeID":{
            "name": <chaincode_name>
        },
        "ctorMsg": {
            "function":"getAllInsurancesByFarm",
            "args":[farmID]
        },
        "secureContext": "admin"
      },
      "id": 1
    }
    
###查询某一养殖场的所有贷款信息
输入参数: 养殖场id `farmID`, 比如`1234567`

返回参数: 包含Loan对象(参考[loan.proto](../chaincode/picc/src/main/loan.proto))的json数组
#####Chaincode Spec: 
    {
      "jsonrpc": "2.0",
      "method": "query",
      "params": {
        "type": 1,
        "chaincodeID":{
            "name": <chaincode_name>
        },
        "ctorMsg": {
            "function":"getAllLoansByFarm",
            "args":[farmID]
        },
        "secureContext": "admin"
      },
      "id": 1
    }
    
###查询某一城市的所有养殖场ID
输入参数: 省份名称province,比如`HEBEI`; 城市名称city,比如`CHENGDE`

返回参数: 包含符合查询条件的所有养殖场id的json数组，比如`["1234567", "1234568"]`
#####Chaincode Spec: 
    {
      "jsonrpc": "2.0",
      "method": "query",
      "params": {
        "type": 1,
        "chaincodeID":{
            "name": <chaincode_name>
        },
        "ctorMsg": {
            "function":"getAllFarmIdsByCity",
            "args":[province,city]
        },
        "secureContext": "admin"
      },
      "id": 1
    }
  
###查询某一省的所有养殖场ID
输入参数: 省份名称province,比如`HEBEI`

返回参数: 包含符合查询条件的所有养殖场id的json数组，比如`["1234567", "1234568"]`
#####Chaincode Spec: 
    {
      "jsonrpc": "2.0",
      "method": "query",
      "params": {
        "type": 1,
        "chaincodeID":{
            "name": <chaincode_name>
        },
        "ctorMsg": {
            "function":"getAllFarmIdsByCity",
            "args":[province]
        },
        "secureContext": "admin"
      },
      "id": 1
    }    


###通过养殖场ID和肉牛耳标查询肉牛信息
输入参数: 养殖场id `farmID`, 比如`1234567`;肉牛耳标`earLabel`，比如`Z5TC923U81`

返回参数: Beef对象(参考[beef.proto](../chaincode/picc/src/main/beef.proto))的json形式
#####Chaincode Spec: 
    {
      "jsonrpc": "2.0",
      "method": "query",
      "params": {
        "type": 1,
        "chaincodeID":{
            "name": <chaincode_name>
        },
        "ctorMsg": {
            "function":"getBeefByFarmAndLabel",
            "args":[farmID, earLabel]
        },
        "secureContext": "admin"
      },
      "id": 1
    }
    
###通过养殖场名字查询养殖场ID
输入参数: 省份名称`province`，比如`HEBEI`;城市名称`city`,比如`CHENGDE`;养殖场名称包含的字符`partialName`,比如`承德第一`

返回参数: 包含符合查询条件的所有养殖场id的json数组，比如`["1234567", "1234568"]`
#####Chaincode Spec: 
    {
      "jsonrpc": "2.0",
      "method": "query",
      "params": {
        "type": 1,
        "chaincodeID":{
            "name": <chaincode_name>
        },
        "ctorMsg": {
            "function":"getAllFarmIdsByName",
            "args":[province, city, partialName]
        },
        "secureContext": "admin"
      },
      "id": 1
    }

###查询放贷员的放贷记录ID
输入参数: 放贷员id `lenderID`, 比如`BETA1290`

返回参数: 包含符合查询条件的所有贷款记录ID, 比如`[I2KS12SJJS]`
#####Chaincode Spec: 
    {
      "jsonrpc": "2.0",
      "method": "query",
      "params": {
        "type": 1,
        "chaincodeID":{
            "name": <chaincode_name>
        },
        "ctorMsg": {
            "function":"getAllLoanIdByLender",
            "args":[lenderID]
        },
        "secureContext": "admin"
      },
      "id": 1
    }
