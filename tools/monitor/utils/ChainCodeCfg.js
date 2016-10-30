const low = require('lowdb');
const storage = require('lowdb/file-sync');
function getPeerApiUrl() {
    var chainCodeCfghandel = low('./config/chaincodeCfg.json', {storage: storage});
    peer = chainCodeCfghandel("peers").value()[0]
    return "http://" + peer.host + ":" + peer.APIPort + "/"
}

function getChainCodeName(inputName) {
    if (SYSTEM.DevMode) {
        return inputName;
    } else {
        var chainCodeCfghandel = low('./config/chaincodeCfg.json', {storage: storage});
        var res = chainCodeCfghandel("ChaincodeID").chain().find({name: inputName}).value()
        if (res) {
            return res.hash
        } else {
            return inputName;
        }
    }
}
function setChainCodeName(inputName, inputhash) {

    var chainCodeCfghandel = low('./config/chaincodeCfg.json', {storage: storage});
    var res = chainCodeCfghandel("ChaincodeID").chain().find({name: inputName}).value();
    if (res) {
        chainCodeCfghandel("ChaincodeID").chain().find({
            name : inputName
        }).assign({
            name : inputName,
            hash : inputhash
        }).value();
    } else {
        chainCodeCfghandel("ChaincodeID").push({
            name : inputName,
            hash : inputhash
        });
    }

}

exports.getPeerApiUrl = getPeerApiUrl
exports.getChainCodeName = getChainCodeName
exports.setChainCodeName = setChainCodeName