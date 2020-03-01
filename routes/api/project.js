var express = require('express'),
    router = express.Router({mergeParams: true}),
    db = require('../../../dummy-data/sysdb');
    sysdb = db.sysdb;
    const IndexNotFound = -1

router.route('/') // Add a binding to handle '/:project'
    .get(function(req,res){    
      var projectIndex = IndexNotFound
      if ( (projectIndex = findProjectIndex(req.params.project)) > IndexNotFound ) {
        res.status(200)
          .send(sysdb.projects[projectIndex])
      }
      else
        res.status(404)
        .send('Not found');
    })
  
function findProjectIndex(projName) {
  var res = IndexNotFound; 
  for (let i = 0; i < sysdb.projects.length; i++) {  //loop over dummy data instead of query to DB
    if (projName == sysdb.projects[i].name) {
      res = i
      break
    }
  }
  return res
}
module.exports = router;