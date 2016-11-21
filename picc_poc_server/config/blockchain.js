/**
 * Created by cocoaWang on 2016/11/18.
 */
module.exports = {
    deploy_mode: 'dev',
    membersrvc: {
        address: '9.12.22.11:10000'
    },
    peer: {
        address: '9.12.22.11:20001'
    },
    /*blockchain的client user*/
    clients: [
        {
            enrollId: 'gov',
            enrollSecret: 'IUZCYDngtwjW',
            affiliation: 'institution_a'
        },

        {
            enrollId: 'Bank_1',
            enrollSecret: 'mRbbQTpZfVVa',
            affiliation: 'institution_a'
        },

        {
            enrollId: 'Insurance_1',
            enrollSecret: 'BTaWHtHrCZry',
            affiliation : 'institution_a'
        }
    ],

    fileKeyValStore: '/wsh/Blockchain/Blockchain-Demo/picc_blockchain_poc/picc_poc_server/keyValStore'
};