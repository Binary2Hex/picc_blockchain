#Chaincode接口描述(API描述)

此文档用于描述chaincode/picc/beef_cattles.go的接口

##Contents

* [Deploy](#deploy)
	* [监管机构部署合约](#监管机构部署合约)

* [Invoke](#invoke)
	* [创建一批肉牛](#创建一批肉牛)
	* [肉牛从进口贸易商移交给养殖场](#肉牛从进口贸易商移交给养殖场)
	* [肉牛从养殖场到屠宰场](#肉牛从养殖场到屠宰场)
	* [防疫局对肉牛进行防疫](#防疫局对肉牛进行防疫)
	* [养殖场向保险公司申请保险](#养殖场向保险公司申请保险)
	* [保险公司对保险申请进行审批](#保险公司对保险申请进行审批)
	* [养殖场申请贷款](#养殖场申请贷款)
	* [贷款机构对贷款申请进行审批](#贷款机构对贷款申请进行审批)
	* [养殖场取消贷款或还清贷款](#养殖场取消贷款或还清贷款)
	* [贷款机构确认贷款还清或取消](#贷款机构确认贷款还清或取消)
	
* [Query](#query)
	* [查询网络中的所有肉牛信息](#查询网络中的所有肉牛信息)
	* [查询自己所拥有的肉牛信息](#查询自己所拥有的肉牛信息)
	* [查询某一肉牛信息](#查询某一肉牛信息)

	
##Deploy
###监管机构部署合约
#####Chaincode Spec: 

	{
		"jsonrpc": "2.0",
		"method": "deploy",
		"params": {
			"type": 1,
			"chaincodeID": {
				"path": “github.com/hyperledger/fabric/chaincode/picc”
			},
			"ctorMsg": {
				"function": "init",
				"args": []
			},
			"secureContext": "gov"
		},
		"id": 1
	}


##Invoke

###创建一批肉牛
监管方或者保险公司在网络中创建一批肉牛

参数: 肉牛id， 进口贸易商
#####Chaincode Spec: 
	{
		"jsonrpc": "2.0",
		"method": "deploy",
		"params": {
			"type": 1,
			"chaincodeID": {
				"name": <chaincode_name>
			},			
			"ctorMsg": {
				"function": "createCattle",
				"args": [
				    <id>, <trader>
				]
			},
			"secureContext": "gov"
		},
		"id": 1
	}
	
###肉牛从进口贸易商移交给养殖场
进口贸易商将自己所拥有的肉牛卖给养殖场。只有肉牛的拥有者以及监管机构有权限做此操作

参数: 肉牛id, 贸易商, 养殖场
#####Chaincode Spec: 
			"ctorMsg": {
				"function": "traderToFarm",
				"args": [
				    <id>, <trader>, <farm>
				]
			}
	
###肉牛从养殖场到屠宰场
养殖场的肉牛送到屠宰场. 只有肉牛的拥有者以及监管机构有权限做此操作

参数: 肉牛id, 养殖场, 屠宰场
#####Chaincode Spec: 
			"ctorMsg": {
				"function": "farmToAbattoir",
				"args": [
				    <id>, <farm>, <abattoir>
				]
			}

###防疫局对肉牛进行防疫
防疫局进行此操作

参数: 肉牛id
####Chaincode Spec
			"ctorMsg": {
				"function": "vaccinate",
				"args": [
				    <id>
				]
			}

###养殖场向保险公司申请保险
养殖场针对自己所拥有的肉牛发起

参数: 肉牛id, 保险公司名称
####Chaincode Spec
			"ctorMsg": {
				"function": "applyForInsurance",
				"args": [
				    <id>, <insuranceCorp>
				]
			}

###保险公司对保险申请进行审批
保险公司操作...

保险通过　  参数: 肉牛id, 保险号
保险未通过  参数: 肉牛id, "_REJECTED" + dateAndTime
####Chaincode Spec
			"ctorMsg": {
				"function": "insure",
				"args": [
				    <id>, <InsuranceID>
				]
			}

###养殖场申请贷款
养殖场发起操作...

参数: 肉牛id, 贷款机构名称
####Chaincode Spec
			"ctorMsg": {
				"function": "applyForLoan",
				"args": [
				    <id>, <LoanCorp>
				]
			}

###贷款机构对贷款申请进行审批

参数: 肉牛id, 贷款数额
####Chaincode Spec
			"ctorMsg": {
				"function": "loan",
				"args": [
				    <id>, <loan>
				]
			}


##Query

###查询网络中的所有肉牛信息
只有监管机构有权限查询所有的肉牛信息
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
				"function": "getAllCattles",
				"args": []
			},
			"secureContext": <role>
		},
		"id": 1
	}

###查询自己所拥有的肉牛信息
任何人都有权限查看自己所拥有的肉牛信息及状态
#####Chaincode Spec: 
			"ctorMsg": {
				"function": "getAllMyCattles",
				"args": []
			}
			
###查询某一肉牛信息
肉牛拥有者，拥有者制定的保险公司以及贷款机构都可以查看此肉牛的信息
#####Chaincode Spec: 
			"ctorMsg": {
				"function": "getCattleByID",
				"args": [
				    <id>
				]
			}

