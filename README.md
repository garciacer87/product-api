# Product API

This is a basic API to manage CRUD operations over products. <br/>
The service was made in [Go](https://go.dev/) and the database was implemented in [Postgresql](https://www.postgresql.org/).

## Setting up dev environment
### Tooling

In order to test this app, you must install the following dependencies:
* [go](https://go.dev/)
* [migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
* [docker](https://www.docker.com/)
  
For code analysis:
* revive: go install github.com/mgechev/revive@latest
* staticcheck: go install honnef.co/go/tools/cmd/staticcheck@latest

For swagger doc generation:
* swaggo: go install github.com/swaggo/swag/cmd/swag@latest

<br/>

### Environment Variables
The application requires the following environment variables:
* **PORT:** server port. e.g.: 8080
* **DATABASE_URI:** e.g.: postgres://productapi:password@localhost:5432/productapi

<br/>

## Build and run

### 1. Locally using Go
To build and run this app locally using native Go, you need a running postgresql database in your localhost. To achieve that, go to the root of the app and run the following command from available from the makefile:

```console
make create-postgres
```

It will download the official postgresql docker image based in alpine.Then, it will create and run the docker container. In a couple of seconds, there will be available a postgresql database running in a docker container with the following params:
* Database name: productapi
* User: productapi
* password: password

Now, to create the structure of the database, run the following command:

```console
make migrate-up
```

Finally, run the following commands to build and run the API (Remember to set the environment variables mentioned before):

```console
make go-build
make run
```

<br/>

### 2. Using docker-compose 
With this approach, you will use the docker-compose file located in the root. This will create a container with the database ready to use, and also, it will run the API from a docker container. So, prior to run the docker-compose, execute the following to build a docker image of the API: 

```console
make docker-build
```

The above command will build a new image with the tag: `product-api:1.0.0`

Then, you can execute the docker-compose on the root of the app:

```console
docker-compose up
```
<br/>

## Endpoints
You can generate swagger documentation from the API to check the information of the endpoints available. <br/>
First, run the following:

```console
make swagger
```

* It will create `swagger.json` and `swagger.yml` in the `/api` folder
* It will create the `docs.go` in the `/docs` folder

Then, go to http://localhost:8080/swagger/index.html to see the Swagger UI.

<br/>

## Unit testing

To locally run the unit tests, first execute the following to create the portable postgresql database:

```console
make create-postgresql
make create-test-unit-db
```

Then, to run the unit tests:

```console
make test
```

