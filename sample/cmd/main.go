package main

import (
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/problem"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

func main() {
	e := echo.New()

	e.HTTPErrorHandler = ProblemDetailsHandler

	e.GET("/sample1", sample1)
	e.GET("/sample2", sample2)
	e.GET("/sample3", sample3)

	e.Logger.Fatal(e.Start(":5000"))
}

// sample with built in problem details function error
func sample1(c echo.Context) error {

	err := errors.New("We have a bad request in our endpoint")
	return problem.BadRequestErr(err)
}

// sample with create custom problem details error
func sample2(c echo.Context) error {

	err := errors.New("We have a request timeout in our endpoint")
	return problem.NewError(http.StatusRequestTimeout, err)
}

// sample with unhandled server error with problem details
func sample3(c echo.Context) error {

	err := errors.New("We have a unhandled server error in our endpoint")
	return err
}

// ProblemDetailsHandler middleware for handle problem details error on top of echo or gin or ...
func ProblemDetailsHandler(error error, c echo.Context) {

	// handle problem details with customize problem map error (optional)
	problem.Map(http.StatusBadRequest, func() *problem.ProblemDetail {
		return &problem.ProblemDetail{
			Type:      "https://httpstatuses.io/400",
			Detail:    error.Error(),
			Title:     "bad-request",
			Timestamp: time.Now(),
		}
	})

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

	// resolve problem details error from response in echo or gin or ...
	if !c.Response().Committed {
		if _, err := problem.ResolveProblemDetails(c.Response(), error); err != nil {
			c.Logger().Error(err)
		}
	}
}
