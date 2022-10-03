package main

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	"github.com/meysamhadeli/problem-details"
	"github.com/pkg/errors"
	"net/http"
)

func main() {

	r := gin.Default()

	r.Use(GinErrorHandler())

	r.GET("/sample1", sample1)
	r.GET("/sample2", sample2)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// sample with return specific status code
func sample1(c *gin.Context) {
	err := errors.New("We have a unauthorized error in our endpoint")
	_ = c.AbortWithError(http.StatusUnauthorized, err)
}

// sample with handling unhandled error to customize return status code with problem details
func sample2(c *gin.Context) {
	err := errors.New("We have a custom error in our endpoint")
	_ = c.Error(err)
}

// GinErrorHandler middleware for handle problem details error on gin
func GinErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()

		for _, err := range c.Errors {

			// handle problem details with customize problem map error
			problem.Map(http.StatusInternalServerError, func() *problem.ProblemDetail {
				return &problem.ProblemDetail{
					Type:   "https://httpstatuses.io/400",
					Detail: err.Error(),
					Status: http.StatusBadRequest,
					Title:  "bad-request",
				}
			})

			if err := problem.ResolveProblemDetails(c.Writer, c.Request, err); err != nil {
				log.Error(err)
			}
		}
	}
}
