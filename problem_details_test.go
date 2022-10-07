package problem

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	custom_errors "github.com/meysamhadeli/problem-details/samples/custom-errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_BadRequest_Err(t *testing.T) {
	badRequestErr := New(http.StatusBadRequest, "bad-request", "We have a bad request error")

	assert.Equal(t, "We have a bad request error", badRequestErr.GetDetails())
	assert.Equal(t, "bad-request", badRequestErr.GetTitle())
	assert.Equal(t, "https://httpstatuses.io/400", badRequestErr.GetType())
	assert.Equal(t, http.StatusBadRequest, badRequestErr.GetStatus())
}

func TestMap_CustomType_Echo(t *testing.T) {

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "http://echo_endpoint1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := echo_endpoint1(c)

	var problemErr ProblemDetailErr

	Map[custom_errors.BadRequestError](func() ProblemDetailErr {
		problemErr = New(http.StatusBadRequest, "bad-request", err.Error())
		return problemErr
	})

	_ = ResolveProblemDetails(c.Response(), c.Request(), err)

	assert.Equal(t, c.Response().Status, http.StatusBadRequest)
	assert.Equal(t, err.Error(), problemErr.GetDetails())
	assert.Equal(t, "bad-request", problemErr.GetTitle())
	assert.Equal(t, "https://httpstatuses.io/400", problemErr.GetType())
	assert.Equal(t, http.StatusBadRequest, problemErr.GetStatus())
}

func TestMap_Status_Echo(t *testing.T) {

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "http://echo_endpoint2", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := echo_endpoint2(c)

	var problemErr ProblemDetailErr

	// map status code to problem details error
	MapStatus(http.StatusBadGateway, func() ProblemDetailErr {
		problemErr = New(http.StatusUnauthorized, "unauthorized", err.Error())
		return problemErr
	})

	_ = ResolveProblemDetails(c.Response(), c.Request(), err)

	assert.Equal(t, c.Response().Status, http.StatusUnauthorized)
	assert.Equal(t, err.(*echo.HTTPError).Message.(error).Error(), problemErr.GetDetails())
	assert.Equal(t, "unauthorized", problemErr.GetTitle())
	assert.Equal(t, "https://httpstatuses.io/401", problemErr.GetType())
	assert.Equal(t, http.StatusUnauthorized, problemErr.GetStatus())
}

func TestMap_CustomType_Gin(t *testing.T) {

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	r := gin.Default()

	r.GET("/gin_endpoint1", func(ctx *gin.Context) {
		err := errors.New("We have a custom type error in our endpoint")
		customBadRequestError := custom_errors.BadRequestError{InternalError: err}
		_ = c.Error(customBadRequestError)
	})

	req, _ := http.NewRequest(http.MethodGet, "/gin_endpoint1", nil)
	r.ServeHTTP(w, req)

	for _, err := range c.Errors {

		var problemErr ProblemDetailErr

		Map[custom_errors.BadRequestError](func() ProblemDetailErr {
			problemErr = New(http.StatusBadRequest, "bad-request", err.Error())
			return problemErr
		})

		_ = ResolveProblemDetails(w, req, err)

		assert.Equal(t, http.StatusBadRequest, problemErr.GetStatus())
		assert.Equal(t, err.Error(), problemErr.GetDetails())
		assert.Equal(t, "bad-request", problemErr.GetTitle())
		assert.Equal(t, "https://httpstatuses.io/400", problemErr.GetType())
	}
}

func TestMap_Status_Gin(t *testing.T) {

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	r := gin.Default()

	r.GET("/gin_endpoint2", func(ctx *gin.Context) {
		err := errors.New("We have a specific status code error in our endpoint")
		// change status code 'StatusBadGateway' to 'StatusUnauthorized' base on handler config
		_ = c.AbortWithError(http.StatusBadGateway, err)
	})

	req, _ := http.NewRequest(http.MethodGet, "/gin_endpoint2", nil)
	r.ServeHTTP(w, req)

	for _, err := range c.Errors {

		var problemErr ProblemDetailErr

		// map status code to problem details error
		MapStatus(http.StatusBadGateway, func() ProblemDetailErr {
			problemErr = New(http.StatusUnauthorized, "unauthorized", err.Error())
			return problemErr
		})

		_ = ResolveProblemDetails(w, req, err)

		assert.Equal(t, http.StatusUnauthorized, problemErr.GetStatus())
		assert.Equal(t, err.Error(), problemErr.GetDetails())
		assert.Equal(t, "unauthorized", problemErr.GetTitle())
		assert.Equal(t, "https://httpstatuses.io/401", problemErr.GetType())
	}
}

func echo_endpoint1(c echo.Context) error {
	err := errors.New("We have a custom type error in our endpoint")
	return custom_errors.BadRequestError{InternalError: err}
}

func echo_endpoint2(c echo.Context) error {
	err := errors.New("We have a specific status code error in our endpoint")
	// change status code 'StatusBadGateway' to 'StatusUnauthorized' base on handler config
	return echo.NewHTTPError(http.StatusBadGateway, err)
}
