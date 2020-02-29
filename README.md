Trion
=====

Trion is meant to be headless CMS with access to content available via API.

- [API concept](#api-concept)
	- [node.js + express.js](#nodejs--expressjs)
	- [REST vs GraphQL](#rest-vs-graphql)
	- [Content model and planned API endpoints](#content-model-and-planned-api-endpoints)
		- [Projects](#projects)
		- [Basic content container: Trion](#basic-content-container-trion)
		- [Defining Trions](#defining-trions)
		- [Trion entries](#trion-entries)
- [Trivia](#trivia)
	- [Where did the name Trion come from?](#where-did-the-name-trion-come-from)
	- [Fun fact](#fun-fact)


## API concept

### node.js + express.js

Due to initial assumptions for the project:

 - free or cheap hosting in the cloud available
 - performance for large amount of requests
 - popular & free to use framework in order to enable quick contributions and cooperation with multiple types of databases

 an educated-guess decision was to use express.js + node.js combo as a API backend due to it's reported performance, popularity and thus availability of learning materials, ease of development and available free hosting options.

### REST vs GraphQL

Assuming a lot of content will be served in the end as static pages and multiple requests and responses may be simply repeating, utilizing network caching capabilities seem to be more valuable than preventing over- and underfetching. This translates to REST API, at least for now. Over time, if internal optimization/caching mechanisms are built, a switch towards GraphQL may be reconsidered.

### Content model and planned API endpoints

![content hierarchy concept](https://raw.githubusercontent.com/MichalRybinski/Trion/master/documentation/res/TrionCMS.png "Content hierarchy concept")

#### Projects

The system is meant to support content CRUD for multiple projects within one deployment. It serves purpose to be extensible for much larger tool than it is now w/o large architectural change.

Planned Endpoints:
```
GET     /projects            # get a list of projects
POST    /projects            # create a project
DELETE  /projects/:project   # delete a project
GET     /:project            # retrieve specific project info
PUT     /:project            # update some of specific project data
```

#### Basic content container: Trion

A basic content container is named here ['Trion'](#where-did-the-name-trion-come-from). It is meant as a container for single instance of specific content model (e.g. instance of article, post or comment is considered to be trion).

Trion doesn't have pre-defined schema.
Trions are defined and used in a context of [Project](#projects).

#### Defining Trions

Defining a Trion (specific content container model) is done via means of JSON Schema (TODO: draftv4 or v6?).

**Consideration**: JSON schema parsing and validation may still cost significant time and negatively impact API performance.

*Assessment*: JSON Schema definition of Trion can be easily re-used for API-level validation of payload when manipulating Trion entries and, though complex, gives a great flexibility in content modeling. Validity and stability of content entries delivered and read seems to be prevailing over potential performance loss, especially given the fact that the latter may be mitigated to some extent by resource scaling or/and code optimizations.

*TODOS*: Deliver some pre-defined JSON Schemas for more complex Trions, like Rich Text Format, which can be easily reused as a part of the project.

**Consideration**: Updating or deleting Trion definition while some Trion entries are already existing may cause all kinds of validation related problems, when fetching/updating these Trion entries.

*Assessment*: This scnearios may occur through a lifetime of a Project, therefore any handling should focus on
- not deleting/loosing accidentally data (existing Trion entries)
- allowing to fetch that data via API
- signalling to API caller an irregular situation, so the API caller may act upon it (e.g. ensure displaying the content in 'old format' somehow)

*TODO*: Trion definitions should keep revision history and Trion entries hold a reference to these revisions as well as Trion definitions. 
In definition update/delete scenario: For fetching, a response should contain some kind of warning code or message, upon which API caller may act. For updating, current revision definition should be applied or an error if the relevant Trion definition is not present (resulting in a specific endpoint being not available in routing).

Planned Endpoints:
```
GET     /:project/trions        # list trions defined so far in :project
POST    /:project/trions        # create new Trion definition. JSON Schema in payload.
GET     /:project/trions/:trion # Retrieve specific Trion definition (JSON Schema)
PUT     /:project/trions/:trion # update specific Trion definition. JSON Schema in payload.
                                # TODO: assess handling Trion entries created with outdated schema.
DELETE  /:project/trions/:trion # delete specific Trion definition. 
                                # TODO: assess handling Trion entries for deleted schema.
```

#### Trion entries

A single Trion entry is a record, row in table, a file with content.
Example:
        
        /:project = blog
        /:project/:trion = content model definition of a post for this blog
        /:project/:trion/:id = a specific post in format defined by :trion


Planned Endpoints:
```
GET     /:project/:trion        # list all :trion entries
POST    /:project/:trion        # create new :trion entry.
                                # Validated against current :trion definition.
GET     /:project/:trion/:id    # Get specific :trion entry
PUT     /:project/:trion/:id    # Update specific :trion entry.
                                # Validated against current :trion definition.
DELETE  /:project/:trion/:id    # Delete specific :trion entry
```

## Trivia

### Where did the name Trion come from?

The term 'trion' is taken in the meaning described in a book by [Stanisław Lem](https://en.wikipedia.org/wiki/Stanis%C5%82aw_Lem), titled 'The Magellanic Cloud' (1955). The single trion is a crystallic content container, which can store not only text and images, which can be easily displayed upon request, but also voice and music and even scents. A collection of trions ('Trion Library') was a unimagineably huge database, to which all members of society had direct and instant access and could search for and consume the content. This is the basic function of Internet nowadays for many of it's users, though we're still waiting for the format for the scent ;-)

Since the trigger for this project was to deliver brand new, shining CMS adapted to specific and constantly shifting needs of science fiction/fiction/fantasy community site, the project and [basic content container](#basic-content-container-trion) has been named after the 'trion'.

### Fun fact

In his non-fiction work of philosophy, 'Dialogs' (1957, English translation: Frank Prengel) mr Lem foresaw the creation of global network. He lived long enough to see that prediction to come to live and commented on it:

> Until I used the Internet I wasn't aware, that there are so many idiots in the world




Markdown ToC thanks to: https://magnetikonline.github.io/markdown-toc-generate/ 