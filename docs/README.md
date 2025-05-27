## Project Structure

- `cmd/` — application entry point  
- `internal/` — core business logic:
  - `logging/` — initialization and configuration of logging
  - `model/` — data structures for use cases, requests, responses, method, and templates
  - `service/` — request handling, routing, and template rendering
    - `serivice/builder` - route, template, response builder
    - `service/constants` - constants for common usage of service
    - `service/matcher` - request matcher
  - `transport/` — HTTP server, handlers, and interfaces
    - `transport/handler` - request handler 
    - `transport/route` - route represents an HTTP route configuration
    -  `transport/server` - server represents an HTTP server that manages multiple routers.
- `vendor/` — external dependencies

## Testing

To run unit tests, use:
```sh
go test ./...
```
Tests check the core functions of the service and its stability.


## Usage Example

Once running, the service listens for HTTP requests, matches them to templates, and returns the corresponding responses.

## Syntax 

### Path Parameters
- `:param_name` - path parameter that can be matched with any value
- `{id:[a-zA-Z0-9-]+}` - path parameter that can be matched with a regular expression

### Allowed values of `MustMethod` field
- `GET` - HTTP method
- `POST` - HTTP method
- `PUT` - HTTP method
- `DELETE` - HTTP method
- `PATCH` - HTTP method

If you do not specify the field, the default value will be method `GET`.


### Placeholder Syntax
- `${...}` - any value
- `${regexp:...}` - value that matches the regular expression, where `...` is custom regexp
- `${req.query:...}` - value from query parameters, where `...` is name of parameter from query  
- `${req.path:...}` - value from path parameters, where `...` is name of parameter from path
- `${req.form:...}` - value from form parameters, where `...` is name of parameter from form
- `${req.headers:...}` - value from headers, where `...` is name of header
- `${req.body:...}` - value from body, where `...` is name of parameter from body

### Requst Matching
- `MustMethod` - method of handled case, is required field
- `MustPathParameters` - path parameters that must be present in the request
- `MustQueryParameters` - query parameters that must be present in the request
- `MustFormParameters` - form parameters that must be present in the request
- `MustHeaders` - headers that must be present in the request
- `MustBody` - body that must be present in the request

### Response Preparation
- `SetStatus` - HTTP status code to return, if you do not specify the field, the default value will be `200`.
- `SetHeaders` - headers to return in the response
- `SetBody` - body to return in the response
- `SetFile` - file to return in the response

### #1 Example Template:

```json
{
    "Path": "/login",
    "Handle": [
        {
            "MatchRequest": {
                "MustMethod": "POST",
                "MustBody": {
                    "username": "test",
                    "password": "password"
                }
            },
            "SetResponse": {
                "SetStatus": 200,
                "SetHeaders": {
                    "SetCookie": "X-Csrf-Token=cookie"
                },
                "SetBody": {
                    "authorized": true
                    "username": ${req.body:username}
                }
            }
        },
        ...
    ]
}
```
### #1 Example response:
```bash
curl -i http://127.0.0.1:5000/login -X POST -d '{"username":"test","password":"password"}'

HTTP/1.1 200 OK
Content-Type: application/json
Setcookie: X-Csrf-Token=cookie
Date: Tue, 27 May 2025 14:50:25 GMT
Content-Length: 19

{"authorized":true, "username":"test"}
```

### #2 Example Template:

```json
{
    "Path": "/",
    "Handle": [
        {}
    ]
}
```
### #2 Example response:
```bash
curl -i http://127.0.0.1:5000/

HTTP/1.1 200 OK
Date: Tue, 27 May 2025 14:54:35 GMT
Content-Length: 0
```
