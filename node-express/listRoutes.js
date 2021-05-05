'use strict';
var path = process.cwd();


 
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

    app.route('/images/favicon.ico').get(function (req, res) { res.sendFile(path + '/public/app/images/favicon.ico'); });
    app.route('/images/logo-blue.png').get(function (req, res) { res.sendFile(path + '/public/app/images/logo-blue.png'); });
    app.route('/hist').get(function (req, res) { res.sendFile(path + '/public/hist.html'); });
    app.route('/hist.js').get(function (req, res) { res.sendFile(path + '/public/hist.js'); });
    app.route('/reg').get(function (req, res) { res.sendFile(path + '/public/reg.html'); });
    app.route('/reg.js').get(function (req, res) { res.sendFile(path + '/public/reg.js'); });
    app.route('/dpo').get(function (req, res) { res.sendFile(path + '/public/dpo.html'); });
    app.route('/dpo.js').get(function (req, res) { res.sendFile(path + '/public/dpo.js'); });
    app.route('/opr').get(function (req, res) { res.sendFile(path + '/public/opr.html'); });
    app.route('/opr.js').get(function (req, res) { res.sendFile(path + '/public/opr.js'); });
    app.route('/opa').get(function (req, res) { res.sendFile(path + '/public/opa.html'); });
    app.route('/opa.js').get(function (req, res) { res.sendFile(path + '/public/opa.js'); });
    app.route('/menu').get(function (req, res) { res.sendFile(path + '/public/menu.html'); });
    app.route('/menu.js').get(function (req, res) { res.sendFile(path + '/public/menu.js'); });
}
