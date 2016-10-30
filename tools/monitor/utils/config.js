/**
 * config设定
 * 
 */
function resetConfig() {
	global.SYSTEM = require('../config/config.json').SYSTEM;
	delete require.cache[require.resolve('../config/config.json')];
}

exports.resetConfig = resetConfig;
