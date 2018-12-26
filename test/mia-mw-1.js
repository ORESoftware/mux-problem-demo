#!/usr/bin/env node
'use strict';


// start server on port 3000 using instructions in readme

const assert = require('assert');
const http = require('http');


const req = http.get({

    hostname: 'localhost',
    port: 3000,
    path: '/api/v1/dogs'

}, res => {


    let json = '';

    res.on('data', d => {
        json += String(d);
    });

    res.once('end', () => {

        let parsed = null;

        try{
            parsed = JSON.parse(json);
        }
        catch(err){
            console.error('coult not parse json:', json);
            console.error(err.message);
            return;
        }


        assert.strictEqual(parsed.Message,"this should get trapped by error handler 3")

    })


});


req.end();