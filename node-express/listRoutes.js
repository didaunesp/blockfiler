'use strict';
module.exports = function (app) {
    //path to the file with the function to call
    var createUser = require('../scripts/createUser');
    var queryDPO = require('../scripts/queryGetDPOList');
    var queryInfo = require('../scripts/queryGetInfo');
    var queryHistory = require('../scripts/queryGetHistory');
    var getDpoKeys = require('../scripts/getDpoKeys');

    // Route the webservices
    app.route('/createUser').post(createUser.createUser);
    app.route('/queryDPO').post(queryDPO.queryDPO);
    app.route('/queryInfo').post(queryInfo.queryInfo);
    app.route('/queryHistory').post(queryHistory.queryHistory);
    app.route('/getDpoKeys').post(getDpoKeys.getDpoKeys);
}
