# PICC Blockchain POC

## 快速开始
### Blockchain network
- 你可以使用Bluemix的blockchain service
- 也可以使用local Blockchain network(基于commit level: [3e0e80a898b259fe463295eabff80ee64f20695e](https://github.com/hyperledger-archives/fabric/commit/3e0e80a898b259fe463295eabff80ee64f20695e))
### 如何搭建local Blockchain network
开发过程我们使用local Blockchain network 更加稳定灵活，便于调试。基本步骤如下:
- 在Linux开发机上git clone上面提到的commit level
- 在fabric目录下执行`make images`命令以生成我们需要的*hyperledger/fabric-peer:latest*、*hyperledger/fabric-membersrvc:latest*以及*hyperledger/fabric-baseimage:latest*等docker镜像
- 在本项目的docker-compose目录下执行`docker-compose up`命令，就可以创建一个4节点+membersrvc的网络了

## 应用简介
此POC用来追踪肉牛的生命周期，基本流程如下
1. 进口贸易商(暂不纳入区块链网络)从国外进口肉牛卖给养殖场/农户
2. 养殖场/农户将肉牛在兽医局做防疫
3. 养殖场/农户对肉牛上保险
4. 保险公司根据肉牛防疫情况对保险进行审批
5. 农户向银行申请贷款
6. 银行根据肉牛是否投保进行贷款审批
7. 农户获得贷款，向养殖场购买肉牛
8. 农户利用贷款购买饲料(饲料供应商暂不纳入网络)
  1. 农户饲养的肉牛成熟，卖给屠宰场(暂不纳入网络)
  2. 肉牛发生意外死亡或生病，保险公司介入

### 肉牛属性
| Attribute       | Type                                                                                                  |
| --------------- | ----------------------------------------------------------------------------------------------------- |
| Id              | String, 每一头肉牛有唯一的Id标识,由5个字符组成(字母或数字，大小写敏感)                                        |
| Vaccinated      | String, 用于记录肉牛接种信息,若为空则表示未进行过任何接种                                                   |
| InsuranceID     | String, 用于记录肉牛的投保id,若id为空则代表未投保                                                         |
| Loan            | Int, 记录此头奶牛的贷款额                                                                                |
| LoanID          | Int, 记录此头奶牛的贷款id                                                                                |
| Origin          | String, 记录肉牛产地                                                                                   |
| Trader          | String, 记录进口贸易商名称                                                                              |
| Status          | Int, 用于记录肉牛在整个生命周期中所处的状态                                                                |

## 应用场景

## 权限控制
可以修改docker-compose目录下的membersrvc.yaml中的相关部分，然后重启网络即可
