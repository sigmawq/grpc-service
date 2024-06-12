# Short overview
A three component service that fetches data from a file, stores it in elastic search and exposes it though an API using GraphQL.
`worker` is responsible for fetching data and sending it to `service` over GRPC.
`service` is responsible for handling data update requests from `worker` and data fetch requests from `api`. 
`api` provides visualization for this data using GraphQL.

# Run using Docker
To build and deploy every service except `worker` run:
```
docker compose up --build
```
Wait until every service is initialized (elasticsearch can take some time)
Once every service from the list above finished initializing the following to build and run `worker`.
```
docker compose -f worker-compose.yaml --build
```
It will do it's job and quickly exit.

You can now go to API and explore 2 possible requests: 
```
retrieve(search, from, size)
```
```
aggregate()
```

## Default configuration
Default configuration is ready for testing, you don't need to change anything if you don't want to. API port is `8080`, elastic search sits on `9200`

## Configuration
Environment parameters for `api`:
| Name        | Description | 
| ------------- |:-------------:|
| PORT      | API will listen on the following port
| SERVICE_HOST | `service` host string

Environment parameters for `service`:
| Name        | Description | 
| ------------- |:-------------:|
| PORT      | Service will listen on the following port |
| ELASTIC_HOST      | `elasticsearch` host string


Environment parameters for `worker`:
| Name        | Description | 
| ------------- |:-------------:|
| DATA_SOURCE      | Service will listen on the following port |
| ELASTIC_HOST      | `elasticsearch` host string 

if a configuration value is not specified in the environment then a default value is going to be used.
Also for ease of use there are things that are not configurable from env such as `worker`'s buffer size and maximum number of objects accumulated before sending a store request to `service`

## Important details on `worker` component
Worker is buffered to not cause memory problems when consuming huge JSON files. Default read buffer size is 1MB. Worker sends a batch to `service` once it's internal buffer is full, internal buffer has a default size of 10000 objects.

## Important details on `service` component
Correct initialization is required the first time to aggregate and do full text search:
1. We are defining a mapping on `subcategory` field to do full text search. (elasticsearch doesn't allow full text search on text fields by default)
2. We are defining an analyzer to correctly handle lowercase and diacritics when performing search.

## Important details on GRPC and GraphQL
If you are going to change GraphQL or GRPC schemas then the corresponding code generation commands will need to be run
GRPC:
```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative grpc/grpc.proto
```

GraphQL:
```
cd api
go run github.com/99designs/gqlgen generate
```
