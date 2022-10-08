package main

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	"github.com/meysamhadeli/problem-details"
	custom_errors "github.com/meysamhadeli/problem-details/samples/custom-errors"
	custom_problems "github.com/meysamhadeli/problem-details/samples/custom-problems"
	"github.com/pkg/errors"
	"net/http"
)

func main() {

	r := gin.Default()

	r.Use(GinErrorHandler())

	r.GET("/sample1", sample1)
	r.GET("/sample2", sample2)
	r.GET("/sample3", sample3)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// handle specific status code to problem details error
func sample1(c *gin.Context) {
	err := errors.New("We have a specific status code error in our endpoint")
	// change status code 'StatusBadGateway' to 'StatusUnauthorized' base on handler config
	_ = c.AbortWithError(http.StatusBadGateway, err)
}

// handle custom type error to problem details error
func sample2(c *gin.Context) {

	err := errors.New("We have a custom type error in our endpoint")
	customBadRequestError := custom_errors.BadRequestError{InternalError: err}
	_ = c.Error(customBadRequestError)
}

// handle custom type error to custom problem details error
func sample3(c *gin.Context) {
	err := errors.New("We have a custom error with custom problem details error in our endpoint")
	customConflictError := custom_errors.ConflictError{InternalError: err}
	_ = c.Error(customConflictError)
}

// GinErrorHandler middleware for handle problem details error on gin
func GinErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()

		for _, err := range c.Errors {

			// map custom type error to problem details error
			problem.Map[custom_errors.BadRequestError](func() problem.ProblemDetailErr {
				return problem.New(http.StatusBadRequest, "bad request", err.Error())
			})

			// map custom type error to custom problem details error
			problem.Map[custom_errors.ConflictError](func() problem.ProblemDetailErr {
				return &custom_problems.CustomProblemDetail{
					ProblemDetailErr: problem.New(http.StatusConflict, "conflict", err.Error()),
					AdditionalInfo:   "some additional info...",
					Description:      "some description...",
				}
			})

			// map status code to problem details error
			problem.MapStatus(http.StatusBadGateway, func() problem.ProblemDetailErr {
				return problem.New(http.StatusUnauthorized, "unauthorized", err.Error())
			})

			if _, err := problem.ResolveProblemDetails(c.Writer, c.Request, err); err != nil {
				log.Error(err)
			}
		}
	}
}
