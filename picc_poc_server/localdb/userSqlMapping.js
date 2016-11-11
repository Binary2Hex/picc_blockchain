/**
 * Created by cocoaWang on 2016/11/8.
 */
// dao/userSqlMapping.js
// CRUD SQL语句
var user = {
    insert:'INSERT INTO picc_user(id, picc_user_name, picc_user_realname,picc_user_secret,picc_user_type) VALUES(0,?,?,?,?)',
    update:'update picc_user set picc_user_name=?, picc_user_secret=? where id=?',
    delete: 'delete from picc_user where picc_user_name=?',
    queryById: 'select * from picc_user where id=?',
    queryByUserName: 'select picc_user_secret from picc_user where picc_user_name=?',
    queryAll: 'select * from picc_user'
};

module.exports = user;