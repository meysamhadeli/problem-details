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

  The data model for problem details is a JSON [RFC7159] object; when
   formatted as a JSON document, it uses the "application/problem+json"
   media type.  Appendix A defines how to express them in an equivalent
   XML format, which uses the "application/problem+xml" media type.

ProblemDetails create a standardized error payload to client, when we have an unhandled error. For implement this approach, we use [RFC 7807](https://datatracker.ietf.org/doc/html/rfc7807) standard to map our error to standard problem details response. It's a JSON or XML format, when formatted as a JSON document, it uses the `"application/problem+json"` media type and for XML format, it uses the "application/problem+xml" media type. This document help us to defines machine-readable details of errors in an HTTP response to avoid the need to define new error response formats for HTTP APIs.

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

## Installation

```bash
go get github.com/meysamhadeli/problem-details
```

There are some samples for using this package [here](./sample/cmd/main.go).
