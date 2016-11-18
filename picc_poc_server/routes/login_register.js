/**
 * Created by cocoaWang on 2016/11/18.
 */

var mysql = require('mysql');
var $confdb = require('../config/serverdb');
var pool = mysql.createPool( $confdb.mysql);
var user = $confdb.user;


/*functions for sql query insert delete and update*/
module.exports={
    addUser : function(req,res,next) {
        pool.getConnection(function(err,connection) {
            // 获取前台页面传过来的参数
            var param = req.body;
            console.log(param.userName);
            // 建立连接，向表中插入值
            connection.query(user.queryByUserName,param.userName.toString(), function(err, result) {
                if(result[0]){
                    res.json({
                        code: '2',
                        msg:'Have already registerd'
                    });
                    connection.release();
                }else{
                    connection.query(user.insert, [param.userName, param.userRealName ,param.userSecret,param.userType], function(err, result) {
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

    queryByUserName : function(req,res,next) {
        pool.getConnection(function(err, connection) {
            var param = req.body;
            console.log("username is %s",param.userName);
            connection.query(user.queryByUserName, param.userName, function(err, result) {
                console.log("login result is "+ result);
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
}




