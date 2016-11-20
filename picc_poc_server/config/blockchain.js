/**
 * Created by cocoaWang on 2016/11/18.
 */
module.exports = {
    deploy_mode: 'dev',
    membersrvc: {
        address: '172.17.0.3'
    },
    peer: {
        address: '172.17.0.3'
    },
    /*blockchain的client user*/
    clients: [
        {
            enrollId: 'gov',
            enrollSecret: 'IUZCYDngtwjW',
            affiliation: 'government'
        },

        {
            enrollId: 'Bank_1',
            enrollSecret: 'mRbbQTpZfVVa',
            affiliation: 'picc_loan'
        },

        {
            enrollId: 'Insurance_1',
            enrollSecret: 'BTaWHtHrCZry',
            affiliation : 'picc_insurance'
        }
    ],

    fileKeyValStore: '/wsh/Blockchain/Blockchain-Demo/picc_blockchain_poc/picc_poc_server/keyValStore'
};