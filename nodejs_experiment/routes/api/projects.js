var express = require('express'),
// you need to set mergeParams: true on the router,
// if you want to access params from the parent router
    router = express.Router({mergeParams: true}),
    project = require('./projects/project')
    db = require('./dummy-data/sysdb'),
    sysdb = db.sysdb;

router.route('/')               // Add a binding to handle '/projects'
  .get(function(req,res){
    res.send(sysdb.projects)
  })

router.use('/:project', project)

module.exports = router;