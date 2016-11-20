/**
 * Created by wsh2160132 on 11/20/2016.
 */
var login_register = require('../db/dbservice');
/*
 *  用于 保险、银行、政府人员的登录、注册
 *
 * */
module.exports = function (server) {
    server.post('/login', function login(req, res, next) {
        console.log("login processing");
        login_register.queryByUserName(req, res, next);
        return next();

    });

    server.post('/register', function register(req, res, next) {
        console.log("register processing");
        login_register.addUser(req, res, next);
        return next();
    });
};