/**
 * Created by cocoaWang on 2016/11/8.
 */
// 实现与MySQL交互
var mysql = require('mysql');
var $conf = require('./db');
//var $util = require('../util/util');
var $sql = require('./userSqlMapping');

// 使用连接池，提升性能
var pool  = mysql.createPool( $conf.mysql);

// 向前台返回JSON方法的简单封装
var jsonWrite = function (res, ret) {
    if(typeof ret === 'undefined') {
        res.json({
            code:'1',
            msg: '操作失败'
        });
    } else {
        res.json(ret);
    }
};

module.exports = {
    addUser: function (req, res, next) {
        pool.getConnection(function(err, connection) {
            // 获取前台页面传过来的参数
            var param = req.body;
            console.log(param.userName);
            // 建立连接，向表中插入值
            // 'INSERT INTO user(id, picc_user_name, picc_user_realname,picc_user_secret,picc_user_type ) VALUES(0,?,?,?,?)',
            connection.query($sql.queryByUserName,param.userName.toString(), function(err, result) {
                if(result[0]){
                    res.json({
                        code: '2',
                        msg:'Have already registerd'
                    });
                    connection.release();
                }else{
                    connection.query($sql.insert, [param.userName, param.userRealName ,param.userSecret,param.userType], function(err, result) {
                        console.log(result);
                        // 以json形式，把操作结果返回给前台页面
                        if( result) {
                            res.json({
                                code: '1',
                                msg: 'Register successful'
                            });
                        } else {
                            res.json({
                                code: '0',
                                msg: 'Internal Server Error'
                            });
                        }
                        // 释放连接
                        connection.release();
                    });
                }


            });

        });
    },

    queryByUserName : function (req, res, next) {
        pool.getConnection(function(err, connection) {
            var param = req.body;
            connection.query($sql.queryByUserName, param.userName, function(err, result) {
                console.log("login result is "+result);
                // 用户名密码成功匹配
                if(!result[0]){
                    res.json({
                       code:"0",
                        msg:"Login failed!User does not exist "
                    });
                }
                else if(param.userSecret != result[0].picc_user_secret){
                        res.json({
                            code:"1",
                            msg:"Login failed!Secret is wrong"
                        });
                }else{
                    res.json({
                        code:"2",
                        msg:"login success"

                    });
                }
                connection.release();
            });
        });
    }
};