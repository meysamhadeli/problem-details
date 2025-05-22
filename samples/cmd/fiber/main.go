package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/labstack/gommon/log"
	"github.com/meysamhadeli/problem-details"
	fiber_helper "github.com/meysamhadeli/problem-details/fiber-helper"
	"github.com/meysamhadeli/problem-details/samples/custom-errors"
	custom_problems "github.com/meysamhadeli/problem-details/samples/custom-problems"
	"github.com/pkg/errors"
	"net/http"
)

func main() {
	app := fiber.New()

	// Register error handler middleware
	app.Use(FiberErrorHandler)

	app.Get("/sample1", sample1)
	app.Get("/sample2", sample2)
	app.Get("/sample3", sample3)

	log.Fatal(app.Listen(":3000"))
}

// handle specific status code to problem details error
func sample1(c fiber.Ctx) error {
	err := errors.New("We have a specific status code error in our endpoint")
	// change status code 'StatusBadGateway' to 'StatusUnauthorized' base on handler config
	return fiber.NewError(http.StatusBadGateway, err.Error())
}

// handle custom type error to problem details error
func sample2(c fiber.Ctx) error {
	err := errors.New("We have a custom type error in our endpoint")
	return custom_errors.BadRequestError{InternalError: err}
}

// handle custom type error to custom problem details error
func sample3(c fiber.Ctx) error {
	err := errors.New("We have a custom error with custom problem details error in our endpoint")
	return custom_errors.ConflictError{InternalError: err}
}

// FiberErrorHandler middleware for handling problem details error on Fiber
func FiberErrorHandler(c fiber.Ctx) error {
	err := c.Next()

	if err != nil {
		// map custom type error to problem details error
		problem.Map[custom_errors.BadRequestError](func() problem.ProblemDetailErr {
			return &problem.ProblemDetail{
				Status: http.StatusBadRequest,
				Title:  "bad request",
				Detail: err.Error(),
			}
		})

		// map custom type error to custom problem details error
		problem.Map[custom_errors.ConflictError](func() problem.ProblemDetailErr {
			return &custom_problems.CustomProblemDetail{
				ProblemDetailErr: &problem.ProblemDetail{
					Status: http.StatusConflict,
					Title:  "conflict",
					Detail: err.Error(),
				},
				AdditionalInfo: "some additional info...",
				Description:    "some description...",
			}
		})

		// map status code to problem details error
		problem.MapStatus(http.StatusBadGateway, func() problem.ProblemDetailErr {
			return &problem.ProblemDetail{
				Status: http.StatusUnauthorized,
				Title:  "unauthorized",
				Detail: err.Error(),
			}
		})

		// resolve problem details error
		if _, err := problem.ResolveProblemDetails(fiber_helper.Response(c), fiber_helper.Request(c), err); err != nil {
			log.Error(err)
		}
	}

	return nil
}
