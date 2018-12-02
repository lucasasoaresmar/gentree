# Gentree

API to manage a genealogical tree

## Getting Started

You can clone this project and deploy it using docker. 

## Prerequisites

* [Docker](https://www.docker.com/get-started) - Get started in docker

## Install

After cloning the project, create a .env file in the root directory of the project and define SERVER and DATABASE as bellow.
```
SERVER=mongo:27017
DATABASE=persons_db
```

If you change the server, remember to change it in the docker-compose.yaml file.
```
...
  mongo:
    image: 'mongo'
    restart: 'always'
    ports:
      - '<SERVER>:27017'
```

In the root directory, build and start the image.
```
$ docker-compose up -d --build
```
That's it!

## Usage

### Properties

Field | Description
------|------------
**id** | The item's unique id.
name | Person's name.
parents | Person's parents.
order | Order of a person in genealogical tree starting from 1 (last descendants).

### Endpoits

`/persons`
* GET - Get everyone
* DELETE - Delete all Persons
* POST - Add a Person

`/persons/{id}`
* GET - Get a Person by its id
* PUT - Update a Person
* PATCH - Update a Person specific field
* DELETE - Delete a Person by its id

`/persons/{id}/parents`
* GET - Get a Person's parents by id

`/persons/{id}/children`
* GET - Get a Person's children by id

`/persons/{parent_id}/isparentof/{child_id}`
* PATCH - Relate a parent to a child

`/persons/{id}/tree`
* GET - Get a Person's genealogical tree by its id

Example of `/persons/5c00aa4b62496c0007eb7f45/tree` response:
```
[
  {
    "id": "5c00aa4b62496c0007eb7f45",
    "name": "Jr",
    "order": 3,
    "parents": [
      "5c00c7204bb9b20007f9c79d",
      "5c00d04e01a9860006fe02e9"
    ]
  },
  {
    "id": "5c00c7204bb9b20007f9c79d",
    "name": "Maria",
    "order": 4
  },
  {
    "id": "5c00d04e01a9860006fe02e9",
    "name": "Mario",
    "order": 4
  }
]
```

### Adding and relating Persons

You can add as many Persons as you want with nothing but a name in the request and relate than with the `/persons/{parent_id}/isparentof/{child_id}` endpoint.

## Built With

* [Mongodb](https://www.mongodb.com/) - The most popular database for modern apps
* [Golang](https://golang.org/) - Backend programming language
* [Docker](https://www.docker.com/) - Containers platform

## Versioning

This project is still in development