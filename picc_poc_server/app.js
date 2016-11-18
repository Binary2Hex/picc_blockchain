/**
 * Created by cocoaWang on 2016/11/18.
 */
var restify = require('restify');
var hfc = require('hfc');
var login_register = require('./routes/login_register');
var picc_blockchain = require('./config/picc_blockchain');
var server = restify.createServer();
server.use(restify.bodyParser());
/*
*  init 方法，用于启动server 时 enroll 三个 client 到 blockchain中
 * 包括 政府人员、银行贷款人员、保险人员 各一个
 *
 * 部署所有的chaincode
* */
init();

function init() {
    // config 文件中配置 peer 和 membersrvc 的地址
    var MEMBERSRVC_ADDRESS   = picc_blockchain.membersrvc.address;
    var PEER_ADDRESS = picc_blockchain.peer.address;
    var chain = hfc.newChain("mychain");

    chain.setKeyValStore( hfc.newFileKeyValStore('/tmp/keyValStore') );
    chain.setMemberServicesUrl("grpc://"+MEMBERSRVC_ADDRESS);
    chain.addPeer("grpc://"+PEER_ADDRESS);

    /*配置是否为 dev mode*/
    var mode = picc_blockchain.peer.peer_mode;
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
        chain.enroll(client.enrollId,client.enrollSecret);
    });

}


/*
*  用于 保险、银行、政府人员的登录、注册
*
* */
server.post('/login',function login(req,res,next) {
    console.log("login processing");
    login_register.queryByUserName(req,res,next);
    return next();

});

server.post('register',function register(req,res,next) {
    console.log("register processing");
    login_register.addUser(req,res,next);
    return next();
});

server.listen(3900, function() {
    console.log('%s listening at %s', server.name, server.url);
});

