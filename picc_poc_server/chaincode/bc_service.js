/**
 * Created by wsh2160132 on 11/20/2016.
 */
var hfc = require('hfc');
var picc_blockchain = require('../config/blockchain');


// config 文件中配置 peer 和 membersrvc 的地址
var MEMBERSRVC_ADDRESS   = picc_blockchain.membersrvc.address;
var PEER_ADDRESS = picc_blockchain.peer.address;
var chain = hfc.newChain("mychain");

chain.setKeyValStore( hfc.newFileKeyValStore(picc_blockchain.fileKeyValStore));

// Set the URL to membership services and to the peer
console.log("member services address ="+MEMBERSRVC_ADDRESS);
console.log("peer address ="+PEER_ADDRESS);

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


// Enroll a user.
function enroll(clients) {
    console.log("enrolling user admin ...");

   // Enroll "admin" which is preregistered in the membersrvc.yaml
    chain.enroll("admin", "Xurw3yU9zI0l", function(err, admin) {
        if (err) {
            console.log("ERROR: failed to register admin: %s",err);
            process.exit(1);
        }
        // Set this user as the chain's registrar which is authorized to register other users.
        chain.setRegistrar(admin);

        clients.forEach(function (client) {

            var registrationRequest = {
                enrollmentID: client.enrollId,
                account:client.affiliation,
                affiliation: "000001"
            };

            chain.registerAndEnroll(registrationRequest, function(error, user) {
                if (error) throw Error(" Failed to register and enroll " + client.enrollId + ": " + error);
                console.log("Enrolled %s successfully\n", client.enrollId);
                // deploy(user);
            });
        });

    });
}

// Deploy chaincode
function deploy(user) {
    console.log("deploying chaincode; please wait ...");
    // Construct the deploy request
    var deployRequest = {
        chaincodeName: process.env.CORE_CHAINCODE_ID_NAME,
        fcn: "init",
        args: ["a", "100", "b", "200"]
    };
    // where is the chain code, ignored in dev mode
    deployRequest.chaincodePath = "github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02";

    // Issue the deploy request and listen for events
    var tx = user.deploy(deployRequest);
    tx.on('complete', function(results) {
        // Deploy request completed successfully
        console.log("deploy complete; results: %j",results);
        // Set the testChaincodeID for subsequent tests
        chaincodeID = results.chaincodeID;
        invoke(user);
    });
    tx.on('error', function(error) {
        console.log("Failed to deploy chaincode: request=%j, error=%k",deployRequest,error);
        process.exit(1);
    });
}

exports.init = function() {

    /*认证三个 client ， gov bank insurance*/
    var clients = picc_blockchain.clients;
    enroll(clients);

}