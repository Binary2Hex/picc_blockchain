/**
 * Created by cocoaWang on 2016/11/18.
 */
// MySQL数据库联接配置
module.exports = {
    mysql: {
        host: '127.0.0.1',
        user: 'root',
        password: 'root',
        database:'picc', // 前面建的user表位于这个数据库中
        port: 3306
    },
    user: {
        insert:'INSERT INTO picc_user(id, picc_user_name, picc_user_realname,picc_user_secret,picc_user_type) VALUES(0,?,?,?,?)',
        update:'update picc_user set picc_user_name=?, picc_user_secret=? where id=?',
        delete: 'delete from picc_user where picc_user_name=?',
        queryById: 'select * from picc_user where id=?',
        queryByUserName: 'select picc_user_secret from picc_user where picc_user_name=?',
        queryAll: 'select * from picc_user'
    }
};

