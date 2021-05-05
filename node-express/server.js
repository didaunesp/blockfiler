#!/usr/bin/env node

var express = require('express'),
    app = express(),
    //port number
    port = process.env.PORT || 3000,
    bodyParser = require('body-parser');
var cors = require('cors')
app.use(cors())

app.use(bodyParser.urlencoded({ extended: true }));
app.use(bodyParser.json());
app.use(function(req, res, next){
    res.header("Access-Control-Allow-Origin", "*");
    res.header("Access-Control-Allow-Headers", "*");
    next();
});
//importing route
var routes = require('./listRoutes'); 
 //register the route
routes(app);

app.listen(port);

console.log('RESTful API server started on: ' + port);
