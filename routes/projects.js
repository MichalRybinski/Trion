//fake projects db
var projects = [
    { name: 'fantastyka', id: '1', owner: 'community' },
    { name: 'second', id: '2', owner: 'corporate' }
  ];

var express = require('express'),
    router = express.Router();

router
    // Add a binding to handle '/projects'
    .get('/', function(req,res){
      res.send(projects)
    })

module.exports = router;