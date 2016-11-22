/**
 * Created by wsh2160132 on 11/20/2016.
 */
var hfc = require('hfc');
var picc_blockchain = require('../config/blockchain');

// default config
var MEMBERSRVC_ADDRESS   = picc_blockchain.membersrvc.address;
var PEER_ADDRESS = picc_blockchain.peer.address;
var chain = hfc.newChain("mychain");

chain.setKeyValStore( hfc.newFileKeyValStore(picc_blockchain.fileKeyValStore));
chain.setMemberServicesUrl("grpc://"+MEMBERSRVC_ADDRESS);
chain.addPeer("grpc://"+PEER_ADDRESS);

var mode = picc_blockchain.deploy_mode;
console.log("DEPLOY_MODE=" + mode);
if (mode === 'dev') {
    chain.setDevMode(true);
    //Deploy will not take long as the chain should already be running
    chain.setDeployWaitTime(10);
} else {
    chain.setDevMode(false);
    //Deploy will take much longer in network mode
    chain.setDeployWaitTime(120);
}
chain.setInvokeWaitTime(10);

exports.init = function() {
    /*认证三个 client ， gov bank insurance*/
    var clients = picc_blockchain.clients;
    var client  = picc_blockchain.client;
    // clients.forEach(function (client) {
        chain.enroll(client.enrollId,client.enrollSecret,function (err,user){
            if (err) {
                console.log("ERROR: failed to register admin: %s",err);
                process.exit(1);
            }
             deploy(user);
        });
    // });
}

/* bc_service for  deploy chaincodes */
function deploy(user) {
    // Construct the deploy request
    var deployRequest = {
        chaincodeName: picc_blockchain.chaincodes.picc_chain.chaincodeName,
        fcn: "init",
        args: [],
        chaincodePath : picc_blockchain.chaincodes.picc_chain.chaincodePath
    };
    // Issue the deploy request and listen for events
    var tx = user.deploy(deployRequest);
    tx.on('complete', function(results) {
        console.log("deploy complete; results: %j",results);
    });
    tx.on('error', function(error) {
        console.log("Failed to deploy chaincode: request=%j, error=%k",deployRequest,error);
        process.exit(1);
    });
}


/*
*  bc_services for  query
* */
exports.query =  function (req,res,next) {
    /*req 请求中需要包括 secureContext \ chaincodeID \ fcn \ args */
    chain.getUser(req.body.secureContext,function (err,member) {
        if (err) throw Error("Failed to getUser " + users[0].username + ": " + err);
        var queryRequest = {
            chaincodeID: req.chaincodeID,
            fcn: req.body.function,
            args: req.body.args
        };

        var tx = member.query(queryRequest);
        tx.on('complete', function (results) {
            console.log("query completed successfully; results=%j",results);
            return results;
            process.exit(0);
        });
        tx.on('error', function (error) {
            console.log("Failed to query chaincode: request=%j, error=%k",queryRequest,error);
            process.exit(1);
        });
    });
}




