const express = require('express')
const app = express()
const port = 3000
var projects = require('./routes/api/projects')
var project = require('./nodejs_experiment/routes/api/projects/project')


app.use('/projects',projects)
//app.use('/:project',project)
app.get('/', (req, res) => res.send('Hello 11World!'))

app.listen(port, () => console.log(`Example app listening on port ${port}!`))

//hints on injecting additional params: https://stackoverflow.com/a/30234851
//var rootRouter = express.Router()
/*rootRouter.use('/:project', function (req, res, next) {
    req.project = req.params.project
    next();        
},
project)
*/
/*
rootRouter.use('/:project', function(req,res) {
    res.send(req.params)
}
)
*/
//rootRouter
//  .get('/', (req, res) => res.send('Hello 11World!'))
//app.use('/', rootRouter)