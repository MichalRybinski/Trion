//fake system db
var sysdb = {
    //fake projects db
    projects: [
        { id: '1', name: 'fantastyka', owner: 'community' },
        { id: '2', name: 'second', owner: 'corporate' }
    ],
    schema_revision: "1"
};

module.exports.sysdb = sysdb