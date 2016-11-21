/**
 * Created by cocoaWang on 2016/11/18.
 */
module.exports = {
    deploy_mode: 'dev',
    membersrvc: {
        address: 'localhost:50051'
    },
    peer: {
        address: 'localhost:50051'
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

    fileKeyValStore: '/tmp/keyValStore'
};