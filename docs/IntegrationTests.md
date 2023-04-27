# Integration Tests

Integration tests are designed to help with testing real API endpoints. Unlike unit tests focused on functionality testing between the API layers etc., integration tests create the API client and call the API the same way the real client would do. The API client for testing purposes is built from OpenAPI specification using the Swagger Codegen tool (https://github.com/clu-ling/docker-swagger-codegen-cli). Integration tests use generated endpoints and models to test predefined flows. The best use case for integration tests is writing the testing flows at the same time with the particular endpoint implementation.

Integration tests should be called manually on a local computer.  Running integration tests on CI is not supported at the moment.

Integration tests are located in the folder `/test/integration`. The file `main.go`  is an executable script. It creates a client and runs the sequence of flows. It logs both requests and responses. 

The subfolder `generate` contains its own copy of the `openapi.yaml` file.  This file is used for client code generation in the `generate/swagger` folder. The generation can be triggered as a part of `make generate` call or manually in `generate.go`. There is intentionally a separate copy of `openapi.yaml` and it is up to the developer to keep it in sync. Otherwise, integration tests would produce possible problems for those who don't want to use them actively.

The integration test flows are independent. One test flow should not affect any other. The requirements for flow initialization will vary from project to project. On IceBreaker each flow starts with wiping the Postgres database and the migration up and populating the DB with initial data.

Part of the test log:
```
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
Flow: Flow2: start
Description: Delete user

WipePostgres: start
WipePostgres: done
MigrateUp: start
...
MigrateUp: done
PopulatePostgres: start
PopulatePostgres: done
-----------------------------
REQUEST:
POST /api/v1/signin/firebase HTTP/1.1
Host: localhost:8080
User-Agent: Swagger-Codegen/1.0.0/go
Content-Length: 0
Accept: application/json; charset=utf-8
Authorization: Bearer 7a3b7a9b-ddae-449b-87e6-7709dc237689
Accept-Encoding: gzip


RESPONSE:
HTTP/1.1 200 OK
Content-Length: 165
Content-Type: application/json; charset=utf-8
Date: Wed, 26 Apr 2023 10:08:10 GMT
Vary: Origin

{"user":{"id":"1649726c-1da7-4cc6-87e6-363b11bf8bac","name":null,"email":null,"instagram":null,"linkedin":null,"imageUrl":null,"finalized":false,"strongAuth":true}}
-----------------------------
REQUEST:
POST /api/v1/signin/firebase HTTP/1.1
Host: localhost:8080
User-Agent: Swagger-Codegen/1.0.0/go
Content-Length: 0
Accept: application/json; charset=utf-8
Authorization: Bearer 94ec0051-98ce-4230-ba73-f1ba41cc18c1
Accept-Encoding: gzip


RESPONSE:
HTTP/1.1 200 OK
Content-Length: 165
Content-Type: application/json; charset=utf-8
Date: Wed, 26 Apr 2023 10:08:10 GMT
Vary: Origin

{"user":{"id":"72fbd4ca-5bf1-45a6-9953-f066dfa9e9c2","name":null,"email":null,"instagram":null,"linkedin":null,"imageUrl":null,"finalized":false,"strongAuth":true}}
-----------------------------
...
Flow: Flow2: success
<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
```

## Mocking the authentication
IceBreaker uses the Firebase authentication. Users can log in using Google or Apple accounts or they can log in as anonymous users. All of this is handled on the Firebase side. For testing purposes, we need to bypass using Firebase. That's why there is `IntegrationMockTokenParser`. The main API can be run in two modes:
 - **Normal mode** - start `/cmd/api/main.go` as usual
 - **Integration tests support mode** - start `/cmd/api/main.go` with `--integration` parameter. The integration test mode for API is expected when you run the integration test client. In this case, the Firebase authentication is replaced with the mocked authentication. The API does not expect the JWT token as the Bearer token. The bearer token is UUID. The integration test client creates the random UUID for each API user. This UUID has the same meaning as the user's FirebaseID. The Postgres connection works slightly differently in this mode (with parameter `default_query_exec_mode=describe_exec`), it is a bit slower, but it handles frequent DB schema changes well.  

## Project Configurations
It is recommended to use integration tests in GoLand IDE. For convenience use this setup - three project configurations.  The IDE allows you to run API in integration test mode in one debugger and at the same time run Integration tests in the second debugger.

### API Server Mode
<p align="center">
    <img src="images/APIServer.png" width="1000">
</p>


### API Server with the Integration Tests support
<p align="center">
    <img src="images/APIServerForIntegration.png" width="1000">
</p>

```
Running in integration tests mode
{"level":"info","time":"2023-04-26T10:07:53Z","caller":"api/main.go:164","msg":"starting application","version":"0.0.0","addr":":8080"}
{"level":"info","time":"2023-04-26T10:07:53Z","msg":"server started","caller":"httpx.Server"}
```

### Integration tests
<p align="center">
    <img src="images/IntegrationTests.png" width="1000">
</p>
