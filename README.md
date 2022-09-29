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
There are some samples for using this package on top of Echo [here](./sample/cmd/echo/main.go) and for Gin [here](./sample/cmd/gin/main.go).

## Installation

```bash
go get github.com/meysamhadeli/problem-details
```

### Creating EchoErrorHandler
For handling our error we need to specify an `Error Handler` on top of `Echo` framework:
```go
// EchoErrorHandler middleware for handle problem details error on echo
func EchoErrorHandler(error error, c echo.Context) {

	// handle problem details with customize problem map error
	problem.Map(http.StatusInternalServerError, func() *problem.ProblemDetail {
		return &problem.ProblemDetail{
			Type:      "https://httpstatuses.io/400",
			Detail:    error.Error(),
			Status:    http.StatusBadRequest,
			Title:     "bad-request",
			Timestamp: time.Now(),
		}
	})

	// resolve problem details error from response in echo
	if !c.Response().Committed {
		if _, err := problem.ResolveProblemDetails(c.Response(), error); err != nil {
			log.Error(err)
		}
	}
}
```

### Creaeting specific status code error for echo:

In this sample we get error response with specific code.
 
 ```go
// sample with return specific status code
func sample1(c echo.Context) error {
	err := errors.New("We have a unauthorized error in our endpoint")
	return echo.NewHTTPError(http.StatusUnauthorized, err)
}
 ```
### Handeling unhandled error for echo:

If we don't have specific status code by default our status code is `500` and we can write a `config option` for problem details in our `ErrorHandler` and override a new staus code and additinal info for our error.

```go
// sample with handling unhanded error to customize return status code with problem details
func sample2(c echo.Context) error {
	err := errors.New("We have a custom error in our endpoint")
	return err
}
```

### Creating GinErrorHandler
For handling our error we need to specify an `Error Handler` on top of `Gin` framework:
```go
// GinErrorHandler middleware for handle problem details error on gin
func GinErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()

		for _, err := range c.Errors {

			// handle problem details with customize problem map error
			problem.Map(http.StatusInternalServerError, func() *problem.ProblemDetail {
				return &problem.ProblemDetail{
					Type:      "https://httpstatuses.io/400",
					Detail:    err.Error(),
					Status:    http.StatusBadRequest,
					Title:     "bad-request",
					Timestamp: time.Now(),
				}
			})

			if _, err := problem.ResolveProblemDetails(c.Writer, err); err != nil {
				log.Error(err)
			}
		}
	}
}
```

### Creaeting specific status code error for gin:

In this sample we get error response with specific code.
 
 ```go
// sample with return specific status code
func sample1(c *gin.Context) {
	err := errors.New("We have a unauthorized error in our endpoint")
	_ = c.AbortWithError(http.StatusUnauthorized, err)
}
 ```
### Handeling unhandled error for gin:

If we don't have specific status code by default our status code is `500` and we can write a `config option` for problem details in our `ErrorHandler` and override a new staus code and additinal info for our error.

```go
// sample with handling unhandled error to customize return status code with problem details
func sample2(c *gin.Context) {
	err := errors.New("We have a custom error in our endpoint")
	_ = c.Error(err)
}
```


# Support

If you like my work, feel free to:

- ⭐ this repository. And we will be happy together :)

Thanks a bunch for supporting me!

## Contribution

Thanks to all [contributors](https://github.com/meysamhadeli/problem-details/graphs/contributors), you're awesome and this wouldn't be possible without you! The goal is to build a categorized community-driven collection of very well-known resources.
