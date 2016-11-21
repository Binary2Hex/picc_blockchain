/**
 * Created by wsh2160132 on 11/20/2016.
 */
var hfc = require('hfc');
var picc_blockchain = require('../config/blockchain');

exports.init = function() {
    // config 文件中配置 peer 和 membersrvc 的地址
    var MEMBERSRVC_ADDRESS   = picc_blockchain.membersrvc.address;
    var PEER_ADDRESS = picc_blockchain.peer.address;
    var chain = hfc.newChain("mychain");

    chain.setKeyValStore( hfc.newFileKeyValStore(picc_blockchain.fileKeyValStore));
    chain.setMemberServicesUrl("grpc://"+MEMBERSRVC_ADDRESS);
    chain.addPeer("grpc://"+PEER_ADDRESS);

    /*配置是否为 dev mode*/
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

    /*认证三个 client ， gov bank insurance*/
    var clients = picc_blockchain.clients;
    clients.forEach(function (client) {
        chain.enroll(client.enrollId,client.enrollSecret,function (err,user){
            if (err) {
                console.log("ERROR: failed to register admin: %s",err);
                process.exit(1);
            }
        });
    });
}