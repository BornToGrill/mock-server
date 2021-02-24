# mock server


## Summary
mock-server provides a mock http server where response variables can be defined both at startup and in URL parameters.
Configurable parameters are:
- status code
- content type
- response time (delay)
- body

## Install

`go >= 1.16`
```bash
go install github.com/borntogrill/mock-server@latest
```
`older`
```bash
go get -v github.com/borntogrill/mock-server
```

## Usage
To run the mock server simply run the command
### Default
```bash
mock-server
```
This will start the server with default values.

### Startup arguments
You can list available arguments using
```bash
mock-server --help
```
Using arguments you can overwrite default response parameters, for example.
```bash
mock-server \
    --status 500 \
    --contentType 'application/json' \
    --body '{ "error": "Something went wrong" }'
```

### URL parameters
You can also override default values using URL parameters.
```bash
curl localhost:8080/?status=201&delay=5000
```
In this case your http client can decide what values it wants. In this case status code `201` is returned and the response time is `~5000 miliseconds`.

## Override precedence
Response values are calculated using the following precedence (descending).
- URL parameter (?status=)
- Stdin (only for body)
- Startup flag (--status)
- Application defaults


## Note
If you get the error `command not found: mock-server` you need to add your `$GOPATH/bin` directory to your `PATH`, see [here](https://golang.org/doc/gopath_code#GOPATH).