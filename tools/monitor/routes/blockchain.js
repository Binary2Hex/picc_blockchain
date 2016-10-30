/**
 * http://usejsdoc.org/
 */
var express = require('express');
var router = express.Router();
var logger = require("../utils/logHelper").helper.logger;
var bcu = require("../utils/BlockChainUtils");
var Obcu = new bcu();

router.get('/blocks', function (req, res) {
    Obcu.chaincode.getChain(function (err, res1, body) {
        if (err) {
            res.send(err);
        } else {
            res.send(body);
        }
    });
});
router.get('/blocks/:id', function (req, res) {
    Obcu.chaincode.getBlock(req.params.id, function (err, res1, body) {
        if (err) {
            res.send(err);
        } else {
            res.send({"block": body});
        }
    });
});

module.exports = router;
