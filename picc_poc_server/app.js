/**
 * Created by cocoaWang on 2016/11/18.
 */
var restify = require('restify');
var config_general = require('./config/general');
var bc_service = require('./chaincode/bc_service');

var server = restify.createServer();
server.use(restify.bodyParser());


require('./routes/login_register')(server);

server.listen(config_general.LISTEN_PORT, function() {
    console.log('%s listening at %s on port %s', server.name, server.url, server.address().port);
});

/*
 *  init 方法，用于启动server 时 enroll 三个 client 到 blockchain中
 * 包括 政府人员、银行贷款人员、保险人员 各一个
 *
 * 部署所有的chaincode
 * */
bc_service.init();
