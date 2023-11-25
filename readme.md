# Polls app

This project is a demonstration of best practices while building software.

From product side application is a simple app that supports creating account, logging in with existing accounts, logging in with google, browsing polls with pagination, creating polls, and voting on polls.

It utilizes some of the popular technologies like golang REST api, echo, golangci-lint, postgre, atlas db migrations, terraform, docker, helm, make, openapi, typescript, nextjs, react, redux, tailwind, google login.

## Local running

To run this project locally:

1. Rename [/.envrc.example](./envrc.example) to .`/.envrc`
2. Rename [./front/.envrc.example](./front/envrc.example) to .`./front/.envrc`
3. Run `direnv allow`
4. Run it using `docker compose up`
5. Apply db schema migration using atlas `make db_migrate`
    - If you don't have atlas installed you can install it by running `make dep`
6. visit `http://localhost:1314`

## Cloud running

To run this project in cloud:

1. Install yq, terraform, helm, atlas, aws cli
2. run `make infra_deploy`
    - it will provision your eks cluster, vpc, and ecr's, and create unique secrets
    - it will push front image and api image to appropriate ecrs
    - it will use helm to deploy 3 services (front, db/postgre, api) and a single ingress
    - it will connect to db service in order to migrate database schemas

## Directories explanations

This project is organized in the manner of monorepo
You can find below explanations and overview of project directory structure.

1. [./cmd](./cmd)
    - there is a single golang command and it is to start a server
2. [./database](./database)
    - contains atlas.hcl database schema
3. [./front](./front)
   contains frontend nextjs project in the manner of monorepo
    1. [./front/src/components](./front/src/components)
        - contains all the react components that don't utilize state
    2. [./front/src/containers](./front/src/containers)
        - contains all the react components that are using state
    3. [./front/src/lib](./front/src/lib)
        - contains frontend helpers
    4. [./front/src/pages](./front/src/pages)
        - contains nextjs pages
    5. [./front/src/store](./front/src/store)
        - contains redux store and typings support
        1. [./front/src/store/reducers](./front/src/store/reducers)
            - contains implementation reducers
    6. [./front/src/style](./front/src/style) - contains simple global styles
4. [./infrastructure](./infrastucture)
   contains infrastucture files needed to deploy to cloud
    1. [./infrastructure/terraform](./infrastructure/terraform)
        - contains terraform project used to provision aws cloud infrastructure
    2. [./infrastructure/helm](./infrastructure/helm)
        - contains helm templates used to deploy services on kubernetes cluster
5. [./pkg](./pkg)
   contains golang rest api app
    1. [./pkg/app](./pkg/app)
        - used to perform app initalization
    2. [./pkg/config](./pkg/config)
        - utlizes viper to perform config initalization
    3. [./pkg/domain](./pkg/domain)
       domain bound code
        1. [./pkg/domain/db](./pkg/domain/db)
            - postgre database access layer, containing all the queries
        2. [./pkg/domain/model](./pkg/domain/model)
            - domain bound business logic
        3. [./pkg/domain/util](./pkg/domain/util)
            - utilities used for api
    4. [./pkg/server/http](./pkg/server)
       contain http server logic
        1. [./pkg/server/http/api](./pkg/server/http/api)
            - contains all api controllers
            1. [./pkg/server/http/middleware](./pkg/server/http/middleware)
                - contains api middleware such as error handler or auth logic
        2. [./pkg/server/http/router](./pkg/server/http/router)
            - contains route mapping
        3. [./pkg/server/http/schema](./pkg/server/http/schema)
            - contains public schemas for api input and outputs as well as possible public errors

### API

Api is build using common web library called [echo](https://echo.labstack.com/).

It contains unit test which provide a decent coverage and enable efficient code scaling.

It follows common best practice of abstracting [database layer](./pkg/domain/db) allowing underlying db to be easily replaced.

It also exposes [public schemas and errors](./pkg/server/http/schema) allowing other software to reuse it.

[API configuration](./pkg/config) is done trough most common library called viper.

It is documented trough [openapi (ex Swagger)](./openapi.yaml)

Code is structured in a way that enables adding different communication mechanisms such as grpc.

### Front

[Front](./front) is built using [next.js](https://nextjs.org/) and statically compiled.

It uses redux for keeping the state, as well as typescript in order to make the code more scalable.

Some containers and pages use local state to simplify the logic.

Styles are done trough extensible tailwind tailwind classes.

### Infrastructure

Infrastructure is setup using terraform in order to make it repeatable, collaborative and extensible.

1. It can easily support multiple different environments such as development, staging, production

Helm is used to produce kubernetes manifest files encouraging code reuse.

1. Postgre `db` service is meant to scale vertically while `front` and `api` services can scale horizontally (that's why there are 2 node groups)
2. `db` pod uses ebs in order to make database persistent trough pod upgrades and restarts
3. Eks cluster is setup in a way that ingress provisions it's own public elb allowing multiple different ingresses to exist

Docker image tagging uses md5 hash of commit instead of `latest` to increase transparency of what code is running in cluster.

### CI/CD

The reason why I haven't setup ci/cd for the project is that it would require always running cloud infrastructure (which would cost a bit of $$).

Since this is a demo project I've decided not to follow trunk based development, nor git-flow.

So performing testing and linting doesn't make sense.

However....

API code quality is ensured through running `make lint`.

API code corectness is ensured through running `make test`.

I've tried to follow git semantic commits while developing this project.

### Outro

Thank you for having a look at this project, I hope that you've enjoyed it.
