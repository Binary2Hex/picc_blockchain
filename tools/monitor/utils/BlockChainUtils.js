/**
 * Created by zhangfz on 16/6/25.
 */
var logger = require("../utils/logHelper").helper.logger;
var request = require('request-json');
var cccfg = require('./ChainCodeCfg');
var client = request.createClient(cccfg.getPeerApiUrl());
var logger = require("./logHelper").helper.logger;

function bcu() {
    var option = {
        jsonrpc: "2.0",
        method: " ",
        params: {
            type: 1,
            chaincodeID: {
                name: cccfg.getChainCodeName("sycoin")
            },
            ctorMsg: {
                function: " ",
                args: [" "]
            }
        },
        secureContext: 'test_user1',
        id: 3
    }
    this.chaincode = option;

    this.chaincode.invoke = function (func, args, cbfunc) {
        option.method = "invoke";
        option.params.ctorMsg.function = func;
        option.params.ctorMsg.args = args;
        option.params.chaincodeID.name = cccfg.getChainCodeName("sycoin");
        option.id = 3;
        logger.debug("option=" + JSON.stringify(option));
        client.post("chaincode", option, function (err, res, body) {
            cbfunc(err, res, body);
        });
    };
    this.chaincode.query = function (func, args, cbfunc) {
        option.method = "query";
        option.params.ctorMsg.function = func;
        option.params.ctorMsg.args = args;
        option.params.chaincodeID.name = cccfg.getChainCodeName("sycoin");
        option.id = 5;
        logger.debug("option=" + JSON.stringify(option));
        client.post("chaincode", option, function (err, res, body) {
            cbfunc(err, res, body);
        });
    };

    this.startChainCode = function (name) {
        var cmd_sign = "CORE_CHAINCODE_ID_NAME=" + cccfg.getChainCodeName(name) + " CORE_PEER_ADDRESS=0.0.0.0:30303 " + __dirname + "/../chaincode/sycoin/sycoin \n";
        /*var child = require('child_process').exec(cmd_sign, function (err, stdout, stderr) {
         if (err !== null) {
         logger.error("exec error: " + err);
         return;
         }
         });
         //logger.info("child PID="+child.PID);
         */
        var cp = require('child_process');
        var sh = cp.spawn('/bin/sh');
        sh.stdout.on('data', function (d) {
            logger.debug(d.toString());
        });
        sh.on('exit', function () {
            logger.debug("Chain Code Exit!");
        });
        sh.stderr.on('data', function (d) {
            logger.error(d.toString());
        });
        sh.stdin.write(cmd_sign);
        sh.stdin.end();
        logger.info("executed. command = " + cmd_sign);
    }
    this.chaincode.deploy = function (func, args, cbfunc) {
        option.method = "deploy";
        option.params.ctorMsg.function = func;
        option.params.ctorMsg.args = args;
        option.params.chaincodeID.name = cccfg.getChainCodeName("sycoin")
        option.params.chaincodeID.path = "github.com/hyperledger/fabric/chaincode/sycoin"
        option.id = 1;
        logger.debug("option=" + JSON.stringify(option));
        client.post("chaincode", option, function (err, res, body) {
            cbfunc(err, res, body);
        });

    }
    this.chaincode.getChain = function (cbfunc) {
        client.get("chain", function (err, res, body) {
            cbfunc(err, res, body)
        });
    }
    this.chaincode.getBlock = function (blockID, cbfunc) {
        client.get("chain/blocks/" + blockID, function (err, res, body) {
            cbfunc(err, res, body)
        });
    }
};

module.exports = bcu;