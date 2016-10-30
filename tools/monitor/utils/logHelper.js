var helper = {};
var log4js = require('log4js');
var fs = require("fs");
var path = require("path");

// 加载配置文件
var objConfig = JSON.parse(fs.readFileSync("./config/logconfig.json", "utf8"));

// 目录创建完毕，才加载配置，不然会出异常
log4js.configure(objConfig);

helper.logger = log4js.getLogger('normal');

process.on('uncaughtException', function (err) {
	//打印出错误
	helper.logger.error(err);
	//打印出错误的调用栈方便调试
	helper.logger.error(err.stack);
});
// 配合express用的方法
exports.use = function(app) {
	// 页面请求日志, level用auto时,默认级别是WARN
	app.use(log4js.connectLogger(helper.logger, {
		level : 'debug',
		format : ':method :url'
	}));
};

exports.helper = helper;