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

ProblemDetails create a standardized error payload to client, when we have an unhandled error. For implement this approach, we use [RFC 7807 (https://datatracker.ietf.org/doc/html/rfc7807) standard to map our error to standard problem details response. It's a JSON or XML format that help us to defines machine-readable details of errors in an HTTP response to avoid the need to define new error response formats for HTTP APIs.

## Installation

```bash
go get github.com/meysamhadeli/problem-details
```

There are some samples for using this package [here](./sample/cmd/main.go).
