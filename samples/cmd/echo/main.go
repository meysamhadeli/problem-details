package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/meysamhadeli/problem-details"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

func main() {
	e := echo.New()

	e.HTTPErrorHandler = EchoErrorHandler

	e.GET("/sample1", sample1)
	e.GET("/sample2", sample2)

	e.Logger.Fatal(e.Start(":3000"))
}

// sample with return specific status code
func sample1(c echo.Context) error {
	err := errors.New("We have a unauthorized error in our endpoint")
	return echo.NewHTTPError(http.StatusUnauthorized, err)
}

// sample with handling unhanded error to customize return status code with problem details
func sample2(c echo.Context) error {
	err := errors.New("We have a custom error in our endpoint")
	return err
}

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

	// resolve problem details error from response in echo or gin or ...
	if !c.Response().Committed {
		if _, err := problem.ResolveProblemDetails(c.Response(), error); err != nil {
			log.Error(err)
		}
	}
}
