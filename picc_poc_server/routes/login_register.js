var express = require('express');
var router = express.Router();
var userDao = require('../localdb/userDao');


router.get("/login", function (req, res) {
  res.render('login', { title: "用户登录"});
});

router.get("/register", function (req, res) {
  res.render('register', { title: "用户注册"});
});

router.get("/logout", function (req, res) {
  req.session.user = null;
  res.redirect('/');
});

router.get("/home", function (req, res) {
  res.render('home', {
    title: 'Home'
    //username: req.session.username.toString()
  });
});

// 注册用户
router.post('/register', function(req, res, next) {
  console.log('registering a new user ');
  userDao.addUser(req, res, next);
});


// 登录用户
router.post('/login', function(req, res, next) {
  console.log('user login ');
  userDao.queryByUserName(req, res, next);
});


module.exports = router;
