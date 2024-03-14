# Quake Analytics

This is a simple project, in wip, to read and parse a log game based on rounds and kills. The goal of this project is demonstrate the use of goroutines and some other resources of go.

## How to execute

You need [go 1.22](https://go.dev/dl/) installed in your machine. Just, run the follow command in the root project dir

```shell
# go run cmd/main.go qgame.log if you in root folder
go run cmd/main.go <filepath>
```

## Tests

If you would like to run tests, run the command bellow

```go
go test ./...
```

For get coverage report, run

```go
GOFLAGS="-count=1" go test -cover -v --race ./...
```
