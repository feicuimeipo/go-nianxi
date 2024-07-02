const auth = require("./auth")
const EventEmitter = require('events').EventEmitter;
const life = new EventEmitter();
module.exports = (req, res, next) => {



  next();
};

