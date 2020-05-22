# Competency - Elasticsearch - Level 1

Elasticsearch is a powerful tool which is essential for all search engineers in today's climate. Many online retailers and search providers have turned to ES as the core of their search platform. While sometimes treated as a NoSQL database, Elasticsearch is better described as a document store and search enginge, as there are a number of mechanisms and measures present in, say, an SQL RDB that ES is lazy about in the interest of performance. At it's core, Elasticsearch runs Lucene to track and manage data in the form of documents; the rest of ES might be regarded in a simplified sense as a Java wrapper which co-ordinates distributed instances of that data, accessible via a well documented, RESTful API.  

Truly, there is much, much more to Elasticsearch in terms of features and applications than could be written here. For more information on Elasticsearch, I'd start with their site - elastic.co - which is rich in learning resources.  

## How do you prove it?

- You can describe what a node is, as well as the roles it might serve in a cluster.

- You can describe an index, both in terms of configuration (types, mappings, etc.) and composition (shards, documents).

- You know how to interact with a cluster via the CLI. You are able to perform status checks as well as search queries, with a comprehensive understanding of the responses from both.

- You can describe quorum based decision making, where it is used in ES, and the ubiquitous arithmetic formula associated with it.

- You can fill in the blank: "You know, for _____."

- You can describe why ES is so commonly used in e-commerce search.

## How do you improve it?

- Play with the API in a sandbox and RTFM!

- Read any of the numerous blogs and articles which have been published surrounding the fundamentals of ES.

- See Elasticsearch Level 2 (TBD)