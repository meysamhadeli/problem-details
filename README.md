<div align="center" style="margin-bottom:20px">
  <img src="assets/problem-details.png" alt="problem-details" />
  <h1>ProblemDetails</h1>
  <div align="center">
    <a href="https://github.com/meysamhadeli/problem-details/actions/workflows/ci.yml"><img alt="build-status" src="https://github.com/meysamhadeli/problem-details/actions/workflows/ci.yml/badge.svg?branch=main&style=flat-square"/></a>
    <a><img alt="license" src="https://img.shields.io/badge/go%20version-%3E=1.18-61CFDD.svg?style=flat-square"/></a>
    <a href="https://github.com/meysamhadeli/problem-details/blob/main/LICENSE"><img alt="build-status" src="https://img.shields.io/github/license/meysamhadeli/problem-details?color=%234275f5&style=flat-square"/></a>
    <a href="https://pkg.go.dev/github.com/meysamhadeli/problem-details"><img alt="build-status" src="https://pkg.go.dev/badge/github.com/meysamhadeli/problem-details"/></a>
  </div>
</div>


> ProblemDetails create a standardized error payload to client, when we have an unhandled error. For implement this approach, we use [RFC 7807](https://datatracker.ietf.org/doc/html/rfc7807) standard to map our error to standard problem details response. It's a JSON or XML format, when formatted as a JSON document, it uses the `"application/problem+json"` media type and for XML format, it uses the "application/problem+xml" media type. This document help us to defines machine-readable details of errors in an HTTP response to avoid the need to define new error response formats for HTTP APIs.

Our problem details response body and headers will be look like this:
```go
// Response body

{
"status": 400,                                        // The HTTP status code generated on the problem occurrence
"title": "bad-request",                               // A short human-readable problem summary
"detail": "We have a bad request in our endpoint",    // A human-readable explanation for what exactly happened
"type": "https://httpstatuses.io/400"                 // URI reference to identify the problem type
}
```
```go
// Response headers

 content-type: application/problem+json
 date: Thu,29 Sep 2022 14:07:23 GMT 
```
There are some samples for using this package [here](./sample/cmd/main.go).

## Installation

```bash
go get github.com/meysamhadeli/problem-details
```

#### Creating ProblemDetails Handler
For handeling our error we need to specify a `error handler` on top of Echo, Gin or other framwork:
```go
// ProblemDetailsHandler middleware for handle problem details error on top of echo or gin or ...
func ProblemDetailsHandler(error error, c echo.Context) {
	
	// handle problem details with customize problem map error (optional)
	problem.Map(http.StatusRequestTimeout, func() *problem.ProblemDetail {
		return &problem.ProblemDetail{
			Type:      "https://httpstatuses.io/408",
			Status:    http.StatusRequestTimeout,
			Detail:    error.Error(),
			Title:     "request-timeout",
			Timestamp: time.Now(),
		}
	})

	// resolve problem details error from response in Echo or Gin or ...
	if !c.Response().Committed {
		if _, err := problem.ResolveProblemDetails(c.Response(), error); err != nil {
			c.Logger().Error(err)
		}
	}
}
```

### Built-in function:

For return desired response we can use some built in `handy problem details function` like `BadRequestErr`,... for return our error base on [RFC 7807](https://datatracker.ietf.org/doc/html/rfc7807) standard.

```go
// sample with built in problem details function error
func sample1(c echo.Context) error {

	err := errors.New("We have a bad request in our endpoint")
	return problem.BadRequestErr(err)
}
```
### Creaeting custom error:

For return desired response we more flexibility response we can use function `NewError` for return our error and code base on [RFC 7807](https://datatracker.ietf.org/doc/html/rfc7807) standard.
 
 ```go
 // sample with create custom problem details error
func sample2(c echo.Context) error {

	err := errors.New("We have a request timeout in our endpoint")
	return problem.NewError(http.StatusRequestTimeout, err)
}
 ```
### Handeling unhandel error:

If we return our error directly we handel our response with code [500](https://httpstatuses.io/500) base on [RFC 7807](https://datatracker.ietf.org/doc/html/rfc7807) standard. 

```go
// sample with unhandled server error with problem details
func sample3(c echo.Context) error {

	err := errors.New("We have a unhandled server error in our endpoint")
	return err
}
```

# Support

If you like my work, feel free to:

- ⭐ this repository. And we will be happy together :)

Thanks a bunch for supporting me!

## Contribution

Thanks to all [contributors](https://github.com/meysamhadeli/problem-details/graphs/contributors), you're awesome and this wouldn't be possible without you! The goal is to build a categorized community-driven collection of very well-known resources.
