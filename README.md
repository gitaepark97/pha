## Setup local development

### Install tools[MacOS]

- [Homebrew](https://brew.sh/)

- [Docker](https://docs.docker.com)

  ```bash
  brew install --cask docker
  ```

- [Golang](https://golang.org/)

  ```bash
  brew install go@1.20
  ```

- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

  ```bash
   brew install golang-migrate
   brew install node
  ```

- [DB Docs](https://dbdocs.io/docs)

  ```bash
  npm install -g dbdocs
  dbdocs login
  ```

- [DBML CLI](https://www.dbml.org/cli/#installation)

  ```bash
  npm install -g @dbml/cli
  sql2dbml --version
  ```

- [Sqlc](https://github.com/kyleconroy/sqlc#installation)

  ```bash
  brew install sqlc
  ```

### Documentation

- Generate DB documentation:

  ```bash
  make db_docs
  ```

- Access the DB documentation at [this address](https://dbdocs.io/parkkitae7/pha). Password: `pha_drowssap`

- Access the API documentation at [this address](https://www.notion.so/PHA-API-Doc-3549e5085a4b4d1ea98c59f0e1ea12c7?pvs=4).

- Access the Postman at [this address](https://www.postman.com/gold-rocket-756140/workspace/pha).

- Access the Test case documentation at [this address](https://docs.google.com/spreadsheets/d/1twvQeU4YLSNFk0jrmk1t3T63ch0sEtsTdAXtxpe0pIk/edit#gid=0).

### How to generate code

- Generate DBML file with Schema SQL:

  ```bash
  make db_dbml
  ```

- Generate SQL CRUD with sqlc:

  ```bash
  make sqlc
  ```

- Create a new db migration:
  ```bash
  migrate create -ext sql -dir migration -seq <migration_name>
  ```

### How to run

- Run test:
  ```bash
  make test
  ```

- Run server:
  ```bash
  make server
  ```

- Run docker containers:
  ```bash
  docker-compose up -d
  ```