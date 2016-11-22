/**
 * Created by cocoaWang on 2016/11/18.
 */
module.exports = {
    deploy_mode: '',
    membersrvc: {
        address: 'localhost:50051'
    },
    peer: {
        address: 'localhost:30303'
    },
    /*blockchain的client user*/
    clients: [
        {
            enrollId: 'gov',
            enrollSecret: 'IUZCYDngtwjW',
            affiliation: 'picc_poc'
        },
        {
            enrollId: 'Bank_1',
            enrollSecret: 'mRbbQTpZfVVa',
            affiliation: 'picc_poc'
        },
        {
            enrollId: 'Insurance_1',
            enrollSecret: 'BTaWHtHrCZry',
            affiliation : 'picc_poc'
        }
    ],
    client : {
        enrollId: 'gov',
        enrollSecret: 'IUZCYDngtwjW',
        affiliation: 'picc_poc'
    },
    /*chaincodes to deploy*/
    chaincodes: {
        picc_chain: {
            chaincodeName: 'picc_chaincode',
            //  path 是 chaincode 所在文件夹
            chaincodePath: 'github.com/hyperledger/fabric/picc_poc_chaincodes'
        }
    },

    fileKeyValStore: '/tmp/keyValStore'
};