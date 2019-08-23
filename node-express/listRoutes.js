'use strict';
module.exports = function(app) {
    //path to the file with the function to call
    var createTx = require('../scripts/createTx');
    var queryTx = require('../scripts/query');

    // Route the webservices
    app.route('/createTx').post(createTx.createTx);
    app.route('/queryTx').post(queryTx.queryTx);
}
