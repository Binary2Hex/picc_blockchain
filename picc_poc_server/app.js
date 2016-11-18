/**
 * Created by cocoaWang on 2016/11/18.
 */
var restify = require('restify');
var login_register = require('./routes/login_register');
var picc_blockchain = require('./config/picc_blockchain');
var server = restify.createServer();
server.use(restify.bodyParser());
/*
*  init 方法，用于启动server 时 enroll 三个 client 到 blockchain中
 * 包括 政府人员、银行贷款人员、保险人员 各一个
* */
init();

function init() {

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

