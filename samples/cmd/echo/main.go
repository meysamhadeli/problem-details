package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/meysamhadeli/problem-details"
	"github.com/meysamhadeli/problem-details/samples/custom-errors"
	custom_problems "github.com/meysamhadeli/problem-details/samples/custom-problems"
	"github.com/pkg/errors"
	"net/http"
)

func main() {
	e := echo.New()

	e.HTTPErrorHandler = EchoErrorHandler

	e.GET("/sample1", sample1)
	e.GET("/sample2", sample2)
	e.GET("/sample3", sample3)

	e.Logger.Fatal(e.Start(":3000"))
}

// handle specific status code to problem details error
func sample1(c echo.Context) error {
	err := errors.New("We have a specific status code error in our endpoint")
	// change status code 'StatusBadGateway' to 'StatusUnauthorized' base on handler config
	return echo.NewHTTPError(http.StatusBadGateway, err)
}

// handle custom type error to problem details error
func sample2(c echo.Context) error {
	err := errors.New("We have a custom type error in our endpoint")
	return custom_errors.BadRequestError{InternalError: err}
}

// handle custom type error to custom problem details error
func sample3(c echo.Context) error {
	err := errors.New("We have a custom error with custom problem details error in our endpoint")
	return custom_errors.ConflictError{InternalError: err}
}

// EchoErrorHandler middleware for handle problem details error on echo
func EchoErrorHandler(error error, c echo.Context) {

	// map custom type error to problem details error
	problem.Map[custom_errors.BadRequestError](func() problem.ProblemDetailErr {
		return problem.New(http.StatusBadRequest, "bad request", error.Error())
	})

	// map custom type error to custom problem details error
	problem.Map[custom_errors.ConflictError](func() problem.ProblemDetailErr {
		return &custom_problems.CustomProblemDetail{
			ProblemDetailErr: problem.New(http.StatusConflict, "conflict", error.Error()),
			AdditionalInfo:   "some additional info...",
			Description:      "some description...",
		}
	})

	// map status code to problem details error
	problem.MapStatus(http.StatusBadGateway, func() problem.ProblemDetailErr {
		return problem.New(http.StatusUnauthorized, "unauthorized", error.Error())
	})

	// resolve problem details error from response in echo
	if !c.Response().Committed {
		if _, err := problem.ResolveProblemDetails(c.Response(), c.Request(), error); err != nil {
			log.Error(err)
		}
	}
}
