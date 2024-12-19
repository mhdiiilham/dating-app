# dating-app
An RESTful API for dating app platform (like Tinder and Bumble)

## Setting up migrations
1. Install `golang-migrate` by running the following command:
```shell
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

2. Export an environment variable for convenience:
```shell
export POSTGRESQL_URL='postgres://postgres:password@localhost:5432/example?sslmode=disable'
```

3. Run the migrations based on what you need:

- To run migration up: `make migrate-up`
- To create a new migration file: `make migrate-create`

## Running on Local Machine
1. Create `app.env` based on the `.env.example` file. This file should contain all the necessary environment variables for the application to run properly.

2. Run the unit tests by executing the command `make test`. This will ensure that all the code is functioning as expected before running the server.

3. Spin up the database using Docker based on the `docker-compose.yml` file in the script directory. This file contains all the necessary configurations to run a local instance of the database.
```shell
make dependencies
# or
docker-compose -f script/docker-compose.yml up
```
4. Run the server using the command `make run`. This will start the server and make it available on your specified port on localhost.

<em>It is important to note that, you may need to modify the docker-compose.yml file to match with the database credentials and endpoint you are using. Also, make sure you have docker and docker-compose installed on your machine.</em>

## API Documentation
You can find a Postman exported collection on the project root directory with the name `Dating.postman_collection.json` You can import this collection into Postman to easily test the API endpoints and see examples of the expected request and response formats.
