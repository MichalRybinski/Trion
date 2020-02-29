const express = require('express')
const app = express()
const port = 3000
var projects = require('./routes/projects')

app.use('/projects',projects)

app.get('/', (req, res) => res.send('Hello 11World!'))

app.listen(port, () => console.log(`Example app listening on port ${port}!`))