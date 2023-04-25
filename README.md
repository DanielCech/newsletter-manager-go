# STRV Go backend template

[![Latest release][release-img]][release]
[![codecov][codecov-img]][codecov]
[![tests][tests-img]][tests]
[![build][build-img]][build]
[![linter][linter-img]][linter]
[![vulnerabilities][vuln-scan-img]][vuln-scan]

[//]: # ([![Go Reference][goreference-img]][goreference])
[//]: # ([![Go Report Card][goreportcard-img]][goreportcard])
[//]: # ([![License: MIT][license-img]][license])

Go API (microservice) template repository.

# Standard Go Project Layout
This template was inspired by https://github.com/golang-standards/project-layout.

## Overview
This is a basic layout for Go application projects. It's **`not an official standard defined by the core Go dev team`**; however, it is a set of common
historical and emerging project layout patterns in the Go ecosystem. Some of these patterns are more popular than others. It also has a number of small
enhancements along with several supporting directories common to any large enough real world application.

> If you are trying to learn Go or if you are building a PoC or a simple project for yourself this project layout is an overkill.
Start with something really simple instead (a single `main.go` file and `go.mod` is more than enough). As your project grows keep in mind
that it'll be important to make sure your code is well structured otherwise you'll end up with a messy code with lots of hidden
dependencies and global state. When you have more people working on the project you'll need even more structure. That's when it's important
to introduce a common way to manage packages/libraries. When you have an open source project or when you know other projects import the
code from your project repository that's when it's important to have private (aka `internal`) packages and code. Clone the repository, keep
what you need and delete everything else! Just because it's there it doesn't mean you have to use it all. None of these patterns are used
in every single project. Even the `vendor` pattern is not universal.

If you need help with naming, formatting and style start by running [`gofmt`](https://golang.org/cmd/gofmt/) and [`golint`](https://github.com/golang/lint).
Also, make sure to read these [code style guidelines and recommendations](https://github.com/strvcom/backend-docs/tree/master/guidelines).

## Notes
This is an example repository. This is by no means a structure that is set in stone, feel free to customize it to your project's needs.

## Configuration
### Example
```yaml
port: 8080
database:
  secret:
    # Uncomment and add path to file with content like secret manager: {"username":"postgres","password":"postgres","engine":"postgres","host":"localhost","port":"5432","dbname":"db_name"}
    # path: db_config.json
    secret_arn: arn:aws:secretsmanager:us-east-1:1234567890ab:secret:alpha/1234567890ab
hash_pepper: 8O0bIesBzlyL1SMO40zpczY5WpoP9BXzwI2s6K3qfWfnXUd3t690Isk4XW8Dpunk
auth_secret: PwA5tHpOl5nXM0y93DdeskRfNj5MBV97pr9ztWMp6T7SInccCLU3ZOuserF1rOck
session:
  access_token_expiration: 1h
  refresh_token_expiration: 30d
metrics:
  port: 9178
  namespace: strv
  subsystem: template
log_level: debug
cors:
  allowed_origins: ["*"]
  allowed_methods: ["GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"]
  allowed_headers: ["Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Origin"]
  allowed_credentials: true
  max_age: 300
# Uncomment to use LocalStack instead of AWS
# aws_endpoint_url: http://localhost:4566
```
`hash_pepper` and `auth_secret` should be generated using cryptographically secure PRNG. Example:
```shell
$ openssl rand -base64 64
```

The prepared file `config.yaml.template` should be used for generating `config.yaml` using the `tea` utility.
Provided config example here in `README.md` could be helpful for a deployment configuration.

## Initializing the Project
The template contains two transport layers you can choose from. It is very uncommon that you would want to use both - `REST` and `GraphQL`.
That's why it is recommended to initialize the project via Makefile. `make init-api CMD={rest|graphql}` command will remove respective files.
After this, you can build the project using `make build`.

## Running the Project
You can run this project via Makefile. `make start-local` will start the docker-compose and create necessary tables in DB.

## Choose between REST vs GraphQL
Both APIs are better in something different. When one API is better in some criteria it will have a bold background. But bear in mind, it
is very difficult to say which API is better in general. The table 
None of these APIs are the silver bullet, and both make sense for their respective use cases.

| CRITERIA                                                 | GraphQL                                                                                                                                                  | REST                                                                                                                                                                                                            |
|----------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| STATUS CODE                                              | Everything is provided with HTTP status code 200 (SUCCESS) with a message in the response body, even in error cases.                                     | `Well defined status codes for each kind of response.`                                                                                                                                                          |
| LEGACY SYSTEMS                                           | Younger technology, not adapted by all of the systems yet.                                                                                               | `More traditional API, used on a large scale in legacy systems, older but far from dead.`                                                                                                                       |
| MICROSERVICES                                            | `Federation saves the day, front-end does not have to know about underlying architecture.`                                                               | By default, it is needed to call each service separately, or use some API gateway.                                                                                                                              |
| SCHEMA-CODE CORRESPONDENCE                               | `Schema first, code generated from the schema. Only one endpoint to care about. Generators are often well maintained for various programming languages.` | Multiple endpoints, need to carefully design the paths. Generators are poorly supported or not existing for some programming languages. It is often needed to maintain the schema-code correspondence manually. |
| FRONT-END AND BACK-END COMMUNICATION                     | `Everything that can be published is published and the front-end then fetches whatever they need. Back-end does not need to sync changes very often.`    | Sometimes needed to discuss endpoints between back-end and front-end. There is a chance that REST will require more time for communication.                                                                     |
| REQUEST TYPES                                            | Everything is a HTTP method POST request.                                                                                                                | `Can organize requests by HTTP method types, each method type has a convention for what it should do.`                                                                                                          |
| FETCHING COMPLEX DATA THAT SPREAD OVER MULTIPLE SERVICES | `One query if federation is used. Sometimes dataloaders need to be used to reduce database load and to solve N+1 problem.`                               | Without a well designed gateway, it is needed to call every endpoint separately. This leads to increased network traffic which is longer and more expensive in comparison with database load.                   |
| UPLOAD FILES                                             | Possible, but ugly as hell. You definitely won't to do it.                                                                                               | `Widely used, simple.`                                                                                                                                                                                          |
| PERFORMANCE                                              | `Heavier on service and database layers of the system.`                                                                                                  | Heavier on network and infrastructure.                                                                                                                                                                          |

### Points to aid API choice decision
- Migrating between these two APIs is usually not worth it, since one is not that overwhelmingly better than the other. None of these is missing crucial features.
- When you are extending a system that uses only one of the APIs, then stick to the API already in use. Migrating to another API is not worth all the time and effort and trouble if the project is already in use. Another reason why to stick with the API already in use is that forcefully bending something new into an existing solution will bring more problems than benefits.
- When starting a resource-driven project, go for REST. Downloads and uploads are easier to implement and manage in REST.
- When starting a new MVP that is not heavily resource oriented (file uploads/downloads), the Graphql should be preferred as long as the tooling for Graphql is well maintained for your programming language.
- When existing tooling for Graphql for your programming language doesn't look promising, then choose REST.
- You can choose an API based on what front-end is more familiar with, in the case you are still unsure which one to pick ;)

## GraphQL
Graphql transport layer is based on `github.com/99designs/gqlgen` since it is the most widely used and maintained one.
To automatically install the generator with `go mod tidy` , there should be a file with imported `gqlgen` package (see `api/graphql/dependency.go`).
This is because the dependency is not used in code, but rather as a package for `go run` (in `go:generate`). The build tag can be arbitrary, it is a simple workaround to keep `go.mod` consistent.

### Structure
Most of GraphQL files in `api/graphql` are generated by `gqlgen`. A brief explanation for each of the generated file:
- `gqlgen.yaml` is config file for `gqlgen` generator, also paths to custom models can be specified here in the autobind section.
- `server.go` is file for running the graphql server quickly, contains minimal configs, can be deleted if we intend to add the graphql endpoint to an existing server. We will not be using this in our projects.
- `schema.graphqls` is source of truth for the graphql. Here are defined models, queries and mutations. This source of truth applies both for FE and BE.
- `resolver.go` is place where dependencies such as service can be plugged into the resolver.
- `schema.resolvers.go` contains resolver functions, generated from schema.graphqls. These functions are meant to be calling service layer in most cases. This file is partially generated, one should proceed with care. It is entirely possible that some functions will be left unimplemented (like querying a password, which makes no sense from a security perspective).
- `model/models_gen.go` contains generated models from schema in schema.graphqls. If existing models are specified in the autobind section in gqlgen.yaml, then they are not generated and existing ones are used.
- `generated.go` contains thousands of lines of generated code. It consists of all the bindings for graphql and Go. Ideally one would not want to even open this file ;)
After schema modification, don't forget to run `go generate ./...` so resolvers and models are updated.

### Validation and custom models
Validation tags can be still used as we are used to. Custom structures can be autobinded in `gqlgen.yaml`, and these models 
are omitted from generation by the generator itself. Our custom models are then binded by generated code in `generated.go`.

Custom non-struct types are defined as scalars in the schema - they should be included in autobind. Our custom struct 
models can contain more fields than are specified in the schema, but if a field is present in the schema and not present in 
our custom struct model, then you have to implement a way to obtain that value, otherwise it will not work. What will 
happen is that the generator will create a resolver function for that missing field, containing panic(unimplemented) by default.

There are more options where to do validation of the query or mutation:
- Validation of request body would be right before calling the service layer in `schema.resolvers.go`, technically the same way as we are used to in REST API.
- Define custom non-struct fields as scalars in the schema and do validation in `UnmarshalGQL`, so there would be less of the 
logic in partially generated `schema.resolvers.go` which should be as lightweight as possible. For a reference check `api/graphql/graph/model/user.go`, types password and email.

### 3rd party models
You will often face a situation where you will need to put an external model such as uuid.UUID into the graphql 
models and tell the schema to use this specific model. This is done in config and in the package with graphql models.

Add an entry for the external model in the models section in `gqlgen.yaml`. Example:
```yaml
UUID:
  model:
    - strv-template-backend-go-api/api/graphql/graph/model.UUID
```

Add serialization functions in the models package which is linked in the autobind section of the config above:
```go
func MarshalUUID(u uuid.UUID) graphql.Marshaler {
	return graphql.MarshalString(u.String())
}

func UnmarshalUUID(v any) (uuid.UUID, error) {
	switch v := v.(type) {
	case string:
		u, err := uuid.Parse(v)
		if err != nil {
			return uuid.UUID{}, fmt.Errorf("invalid UUID: %w", err)
		}
		return u, nil
	default:
		return uuid.UUID{}, fmt.Errorf("%T must be a string", v)
	}
}
```

There must be functions with the prefix `Unmarshal` and `Marshal`, strictly followed by the name of the model. The name must 
be the same as it is in the schema. This is because of the generator, which will take care of linking everything together. 
So the name of the model in the schema here in the code above is `UUID`. The model name in the Go code does not have to 
be in a specific format, but as always, it should be named somehow consistently for clarity.

### Routing
Graphql usually utilizes only one endpoint. It is often handled by the router from `github.com/go-chi`. 
This means that even in graphql we can use the router we are used to from REST and write middlewares
in the typical HTTP way. Check injection of Graphql handler into chi router in `api/graphql/graph/handler.go`.

### Authentication with middlewares
Graphql offers a mechanism called directives. They take form of an annotation in the graphql 
schema and are preceded by `@`. Example:
```graphql
directive @auth(roles: [Role!]!) on OBJECT | FIELD_DEFINITION
```

`OBJECT` and `FIELD_DEFINITION` are places where this directive can be used. The `roles` is an argument that 
takes an array of roles. This will cause the generator to generate following function in generated go code:
```go
type DirectiveRoot struct {
  Auth func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []model.Role) (res interface{}, err error)
}
```

All of the directives defined inside the schema `schema.graphqls` will be generated in `DirectiveRoot` in `generated.go`. 
Then we need to pass in our implementation for that directive - check `api/graphql/graph/directive.go`.

Example use of directives in `schema.graphqls`:
```graphql
destroySession(input: DestroySessionInput): DestroySessionResponse! @auth(roles: [USER, ADMIN])
```

This states that before the destroy session `Auth` directive will be invoked, with argument `[USER, ADMIN]`, behavior 
depends on implementation. After we have prepared everything regarding the directives, we can move to authentication 
logic which is implemented in an `Auth` directive in `directive.go`. Functionality of this directive expects the user 
role to be present in a context. It then checks if the user has sufficient role - that means role specified 
in arguments here `@auth(roles: [USER, ADMIN])`.

### Data Loaders
Graphql saves a lot of traffic, by fetching all the required data in one API call, but at the cost of many accesses to
the database. In case of listing nested objects, there would be a database call for each object in the list (also known as N+1 problem).
This puts the database under unnecessary load. There is a way however, how to fetch the nested objects in one database call.
It is done via data loaders. We use them here as a secondary database layer (`domain/user/postgres/dataloader`), for each API call,
loaders are initialized in the middleware in `api/graphql/middleware/dataloader.go`. Inside the dataloader implementation 
are queries, which select N objects for N keys at once. We use dataloaders from `github.com/graph-gophers/dataloader`.
Useful short reading with examples at `https://gqlgen.com/reference/dataloaders/`.

[build]: https://github.com/strvcom/strv-template-backend-go-api/actions/workflows/build.yaml
[build-img]: https://github.com/strvcom/strv-template-backend-go-api/actions/workflows/build.yaml/badge.svg
[codecov]: https://codecov.io/gh/strvcom/strv-template-backend-go-api
[codecov-img]: https://codecov.io/gh/strvcom/strv-template-backend-go-api/branch/master/graph/badge.svg?token=V10H1DOLYK
[linter]: https://github.com/strvcom/strv-template-backend-go-api/actions/workflows/lint.yaml
[linter-img]: https://github.com/strvcom/strv-template-backend-go-api/actions/workflows/lint.yaml/badge.svg
[release]: https://github.com/strvcom/strv-template-backend-go-api/releases
[release-img]: https://img.shields.io/github/v/release/strvcom/strv-template-backend-go-api?display_name=tag&sort=semver
[tests]: https://github.com/strvcom/strv-template-backend-go-api/actions/workflows/tests.yaml
[tests-img]: https://github.com/strvcom/strv-template-backend-go-api/actions/workflows/tests.yaml/badge.svg
[vuln-scan]: https://github.com/strvcom/strv-template-backend-go-api/actions/workflows/vuln-scan.yaml
[vuln-scan-img]: https://github.com/strvcom/strv-template-backend-go-api/actions/workflows/vuln-scan.yaml/badge.svg

[//]: # ([goreference]: https://pkg.go.dev/go.strv.io/<open-source-lib-name>)
[//]: # ([goreference-img]: https://pkg.go.dev/badge/go.strv.io/<open-source-lib-name>.svg)
[//]: # ([goreportcard]: https://goreportcard.com/report/go.strv.io/<open-source-lib-name>)
[//]: # ([goreportcard-img]: https://goreportcard.com/badge/go.strv.io/<open-source-lib-name>)
[//]: # ([license]: https://opensource.org/licenses/MIT)
[//]: # ([license-img]: https://img.shields.io/github/license/strvcom/strv-template-backend-go-api)
