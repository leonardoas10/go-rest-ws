### REST & WS with golang

## Description

App for practice Rest and WS with golang, using libraries as:

-   jwt
-   mux
-   websocket
-   pq
-   cors
-   crypto
-   ksuid

## Start App

1. Install [ Docker Engine ](https://docs.docker.com/engine/install/) :fire:
2. Go to database folder and run `docker build . -t go-rest-ws-db`
3. Then run => `docker run -p 5432:54321 go-rest-ws-db up`
4. Go to root and execute `go mod download`
5. Execute `go run main.go`
6. [ App ](http://localhost:5050/)

## Compile server with docker

7. Go to root and run `docker build . -t go-rest-ws-app`
8. Then run => `docker run -p 5432:54321 go-rest-ws-app up`
9. [ App ](http://localhost:5050/)
