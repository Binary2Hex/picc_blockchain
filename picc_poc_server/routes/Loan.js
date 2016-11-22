/**
 * Created by wsh2160132 on 11/20/2016.
 */
var bc_service  = require('../chaincode/bc_service');
module.exports = function (server) {
    server.post('/getFarmById',function (req,res,next) {
        var results = bc_service.query(req,res,next);

        return next;
    });
}