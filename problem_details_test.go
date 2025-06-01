package problem

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v3"
	"github.com/labstack/echo/v4"
	custom_errors "github.com/meysamhadeli/problem-details/samples/custom-errors"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMap_CustomType_Echo(t *testing.T) {

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "http://echo_endpoint1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := echo_endpoint1(c)

	Map[custom_errors.BadRequestError](func() ProblemDetailErr {
		return &ProblemDetail{
			Status: http.StatusBadRequest,
			Title:  "bad-request",
			Detail: err.Error(),
		}
	})

	p, _ := ResolveProblemDetails(c.Response(), c.Request(), err)

	assert.Equal(t, http.StatusBadRequest, c.Response().Status)

	assert.Equal(t, err.Error(), p.GetDetails())
	assert.Equal(t, "bad-request", p.GetTitle())
	assert.Equal(t, "https://httpstatuses.io/400", p.GetType())
	assert.Equal(t, http.StatusBadRequest, p.GetStatus())
}

func TestMap_Custom_Problem_Err_Echo(t *testing.T) {

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "http://echo_endpoint4", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := echo_endpoint4(c)

	Map[custom_errors.ConflictError](func() ProblemDetailErr {
		return &CustomProblemDetailTest{
			ProblemDetailErr: &ProblemDetail{
				Status: http.StatusConflict,
				Title:  "conflict",
				Detail: err.Error(),
			},
			AdditionalInfo: "some additional info...",
			Description:    "some description...",
		}
	})

	p, _ := ResolveProblemDetails(c.Response(), c.Request(), err)
	cp := p.(*CustomProblemDetailTest)

	assert.Equal(t, http.StatusConflict, c.Response().Status)

	assert.Equal(t, err.Error(), cp.GetDetails())
	assert.Equal(t, "conflict", cp.GetTitle())
	assert.Equal(t, "https://httpstatuses.io/409", cp.GetType())
	assert.Equal(t, http.StatusConflict, cp.GetStatus())
	assert.Equal(t, "some description...", cp.Description)
	assert.Equal(t, "some additional info...", cp.AdditionalInfo)
}

func TestMap_Status_Echo(t *testing.T) {

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "http://echo_endpoint2", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := echo_endpoint2(c)

	MapStatus(http.StatusBadGateway, func() ProblemDetailErr {
		return &ProblemDetail{
			Status: http.StatusUnauthorized,
			Title:  "unauthorized",
			Detail: err.Error(),
		}
	})

	p, _ := ResolveProblemDetails(c.Response(), c.Request(), err)

	assert.Equal(t, http.StatusUnauthorized, c.Response().Status)

	assert.Equal(t, "We have a specific status code error in our endpoint", p.GetDetails())
	assert.Equal(t, "unauthorized", p.GetTitle())
	assert.Equal(t, "https://httpstatuses.io/401", p.GetType())
	assert.Equal(t, http.StatusUnauthorized, p.GetStatus())
}

func TestMap_Status2_Echo(t *testing.T) {

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "http://echo_endpoint2", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := echo_endpoint2b(c)

	MapStatus(http.StatusBadGateway, func() ProblemDetailErr {
		return &ProblemDetail{
			Status: http.StatusUnauthorized,
			Title:  "unauthorized",
			Detail: err.Error(),
		}
	})

	p, _ := ResolveProblemDetails(c.Response(), c.Request(), err)

	assert.Equal(t, c.Response().Status, http.StatusUnauthorized)
	assert.Equal(t, "We have a specific status code error in our endpoint", p.GetDetails())
	assert.Equal(t, "unauthorized", p.GetTitle())
	assert.Equal(t, "https://httpstatuses.io/401", p.GetType())
	assert.Equal(t, http.StatusUnauthorized, p.GetStatus())
}

func TestMap_Status_Details_Echo(t *testing.T) {

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "http://echo_endpoint2", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := echo.HTTPError{Code: http.StatusNotFound, Message: "Entity not found. Please contact admin.", Internal: errors.New("Error with additional analysis/audit logging information")}

	MapStatus(http.StatusBadGateway, func() ProblemDetailErr {
		return &ProblemDetail{
			Status: http.StatusUnauthorized,
			Title:  "unauthorized",
			Detail: err.Error(),
		}
	})

	p, _ := ResolveProblemDetails(c.Response(), c.Request(), &err)

	assert.Equal(t, http.StatusNotFound, c.Response().Status)
	assert.Equal(t, "Entity not found. Please contact admin.", p.GetDetails())
	assert.Contains(t, p.GetStackTrace(), "Error with additional analysis/audit logging information")
	assert.Equal(t, "Not Found", p.GetTitle())
	assert.Equal(t, "https://httpstatuses.io/404", p.GetType())
	assert.Equal(t, http.StatusNotFound, p.GetStatus())
}

func TestMap_Unhandled_Err_Echo(t *testing.T) {

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "http://echo_endpoint3", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := echo_endpoint3(c)

	p, _ := ResolveProblemDetails(c.Response(), c.Request(), err)

	assert.Equal(t, http.StatusInternalServerError, c.Response().Status)

	assert.Equal(t, err.Error(), p.GetDetails())
	assert.Equal(t, "Internal Server Error", p.GetTitle())
	assert.Equal(t, "https://httpstatuses.io/500", p.GetType())
	assert.Equal(t, http.StatusInternalServerError, p.GetStatus())
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

		Map[custom_errors.BadRequestError](func() ProblemDetailErr {
			return &ProblemDetail{
				Status: http.StatusBadRequest,
				Title:  "bad-request",
				Detail: err.Error(),
			}
		})

		p, _ := ResolveProblemDetails(w, req, err)

		assert.Equal(t, http.StatusBadRequest, p.GetStatus())
		assert.Equal(t, err.Error(), p.GetDetails())
		assert.Equal(t, "bad-request", p.GetTitle())
		assert.Equal(t, "https://httpstatuses.io/400", p.GetType())
	}
}

func TestMap_Custom_Problem_Err_Gin(t *testing.T) {

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	r := gin.Default()

	r.GET("/gin_endpoint4", func(ctx *gin.Context) {
		err := errors.New("We have a custom error with custom problem details error in our endpoint")
		customConflictError := custom_errors.ConflictError{InternalError: err}
		_ = c.Error(customConflictError)
	})

	req, _ := http.NewRequest(http.MethodGet, "/gin_endpoint4", nil)
	r.ServeHTTP(w, req)

	for _, err := range c.Errors {

		Map[custom_errors.ConflictError](func() ProblemDetailErr {
			return &CustomProblemDetailTest{
				ProblemDetailErr: &ProblemDetail{
					Status: http.StatusConflict,
					Title:  "conflict",
					Detail: err.Error(),
				},
				AdditionalInfo: "some additional info...",
				Description:    "some description...",
			}
		})

		p, _ := ResolveProblemDetails(w, req, err)
		cp := p.(*CustomProblemDetailTest)

		assert.Equal(t, http.StatusConflict, cp.GetStatus())
		assert.Equal(t, err.Error(), cp.GetDetails())
		assert.Equal(t, "conflict", cp.GetTitle())
		assert.Equal(t, "https://httpstatuses.io/409", cp.GetType())
		assert.Equal(t, "some description...", cp.Description)
		assert.Equal(t, "some additional info...", cp.AdditionalInfo)
	}
}

func TestMap_Status_Gin(t *testing.T) {

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	r := gin.Default()

	r.GET("/gin_endpoint2", func(ctx *gin.Context) {
		err := errors.New("We have a specific status code error in our endpoint")
		_ = c.AbortWithError(http.StatusBadGateway, err)
	})

	req, _ := http.NewRequest(http.MethodGet, "/gin_endpoint2", nil)
	r.ServeHTTP(w, req)

	for _, err := range c.Errors {

		MapStatus(http.StatusBadGateway, func() ProblemDetailErr {
			return &ProblemDetail{
				Status: http.StatusUnauthorized,
				Title:  "unauthorized",
				Detail: err.Error(),
			}
		})

		p, _ := ResolveProblemDetails(w, req, err)

		assert.Equal(t, http.StatusUnauthorized, p.GetStatus())
		assert.Equal(t, err.Error(), p.GetDetails())
		assert.Equal(t, "unauthorized", p.GetTitle())
		assert.Equal(t, "https://httpstatuses.io/401", p.GetType())
	}
}

func TestMap_Unhandled_Err_Gin(t *testing.T) {

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	r := gin.Default()

	r.GET("/gin_endpoint3", func(ctx *gin.Context) {
		err := errors.New("We have a unhandeled error in our endpoint")
		_ = c.Error(err)
	})

	req, _ := http.NewRequest(http.MethodGet, "/gin_endpoint3", nil)
	r.ServeHTTP(w, req)

	for _, err := range c.Errors {

		p, _ := ResolveProblemDetails(w, req, err)

		assert.Equal(t, http.StatusInternalServerError, p.GetStatus())
		assert.Equal(t, err.Error(), p.GetDetails())
		assert.Equal(t, "Internal Server Error", p.GetTitle())
		assert.Equal(t, "https://httpstatuses.io/500", p.GetType())
	}
}

func TestMap_CustomType_Fiber(t *testing.T) {
	app := fiber.New()

	app.Get("/fiber_endpoint1", func(c fiber.Ctx) error {
		return fiber_endpoint1(c)
	})

	// Create fasthttp request context
	fctx := &fasthttp.RequestCtx{}
	fctx.Request.SetRequestURI("/fiber_endpoint1")
	fctx.Request.Header.SetMethod(http.MethodGet)

	// Create Fiber context
	ctx := app.AcquireCtx(fctx)
	defer app.ReleaseCtx(ctx)

	// Execute the handler
	handlerErr := fiber_endpoint1(ctx)

	Map[custom_errors.BadRequestError](func() ProblemDetailErr {
		return &ProblemDetail{
			Status: http.StatusBadRequest,
			Title:  "bad-request",
			Detail: handlerErr.Error(),
		}
	})

	p, _ := ResolveProblemDetails(Response(ctx), Request(ctx), handlerErr)

	assert.Equal(t, http.StatusBadRequest, ctx.Response().StatusCode())
	assert.Equal(t, handlerErr.Error(), p.GetDetails())
	assert.Equal(t, "bad-request", p.GetTitle())
	assert.Equal(t, "https://httpstatuses.io/400", p.GetType())
	assert.Equal(t, http.StatusBadRequest, p.GetStatus())
}

func TestMap_Custom_Problem_Err_Fiber(t *testing.T) {
	app := fiber.New()

	app.Get("/fiber_endpoint4", func(c fiber.Ctx) error {
		return fiber_endpoint4(c)
	})

	fctx := &fasthttp.RequestCtx{}
	fctx.Request.SetRequestURI("/fiber_endpoint4")
	fctx.Request.Header.SetMethod(http.MethodGet)

	ctx := app.AcquireCtx(fctx)
	defer app.ReleaseCtx(ctx)

	handlerErr := fiber_endpoint4(ctx)

	Map[custom_errors.ConflictError](func() ProblemDetailErr {
		return &CustomProblemDetailTest{
			ProblemDetailErr: &ProblemDetail{
				Status: http.StatusConflict,
				Title:  "conflict",
				Detail: handlerErr.Error(),
			},
			AdditionalInfo: "some additional info...",
			Description:    "some description...",
		}
	})

	p, _ := ResolveProblemDetails(Response(ctx), Request(ctx), handlerErr)
	cp := p.(*CustomProblemDetailTest)

	assert.Equal(t, http.StatusConflict, ctx.Response().StatusCode())
	assert.Equal(t, handlerErr.Error(), cp.GetDetails())
	assert.Equal(t, "conflict", cp.GetTitle())
	assert.Equal(t, "https://httpstatuses.io/409", cp.GetType())
	assert.Equal(t, http.StatusConflict, cp.GetStatus())
	assert.Equal(t, "some description...", cp.Description)
	assert.Equal(t, "some additional info...", cp.AdditionalInfo)
}

func TestMap_Status_Fiber(t *testing.T) {
	app := fiber.New()

	app.Get("/fiber_endpoint2", func(c fiber.Ctx) error {
		return fiber_endpoint2(c)
	})

	fctx := &fasthttp.RequestCtx{}
	fctx.Request.SetRequestURI("/fiber_endpoint2")
	fctx.Request.Header.SetMethod(http.MethodGet)

	ctx := app.AcquireCtx(fctx)
	defer app.ReleaseCtx(ctx)

	handlerErr := fiber_endpoint2(ctx)

	MapStatus(http.StatusBadGateway, func() ProblemDetailErr {
		return &ProblemDetail{
			Status: http.StatusUnauthorized,
			Title:  "unauthorized",
			Detail: handlerErr.Error(),
		}
	})

	p, _ := ResolveProblemDetails(Response(ctx), Request(ctx), handlerErr)

	assert.Equal(t, http.StatusUnauthorized, ctx.Response().StatusCode())
	assert.Equal(t, handlerErr.Error(), p.GetDetails())
	assert.Equal(t, "unauthorized", p.GetTitle())
	assert.Equal(t, "https://httpstatuses.io/401", p.GetType())
	assert.Equal(t, http.StatusUnauthorized, p.GetStatus())
}

func TestMap_Unhandled_Err_Fiber(t *testing.T) {
	app := fiber.New()

	app.Get("/fiber_endpoint3", func(c fiber.Ctx) error {
		return fiber_endpoint3(c)
	})

	fctx := &fasthttp.RequestCtx{}
	fctx.Request.SetRequestURI("/fiber_endpoint3")
	fctx.Request.Header.SetMethod(http.MethodGet)

	ctx := app.AcquireCtx(fctx)
	defer app.ReleaseCtx(ctx)

	handlerErr := fiber_endpoint3(ctx)

	p, _ := ResolveProblemDetails(Response(ctx), Request(ctx), handlerErr)

	assert.Equal(t, http.StatusInternalServerError, ctx.Response().StatusCode())
	assert.Equal(t, handlerErr.Error(), p.GetDetails())
	assert.Equal(t, "Internal Server Error", p.GetTitle())
	assert.Equal(t, "https://httpstatuses.io/500", p.GetType())
	assert.Equal(t, http.StatusInternalServerError, p.GetStatus())
}

func echo_endpoint1(c echo.Context) error {
	err := errors.New("We have a custom type error in our endpoint")
	return custom_errors.BadRequestError{InternalError: err}
}

func echo_endpoint2(c echo.Context) error {
	err := errors.New("We have a specific status code error in our endpoint")
	return echo.NewHTTPError(http.StatusBadGateway, err)
}

func echo_endpoint2b(c echo.Context) error {
	return echo.NewHTTPError(http.StatusBadGateway, "We have a specific status code error in our endpoint")
}

func echo_endpoint3(c echo.Context) error {
	err := errors.New("We have a unhandeled error in our endpoint")
	return err
}

func echo_endpoint4(c echo.Context) error {
	err := errors.New("We have a custom error with custom problem details error in our endpoint")
	return custom_errors.ConflictError{InternalError: err}
}

func fiber_endpoint1(c fiber.Ctx) error {
	err := errors.New("We have a custom type error in our endpoint")
	return custom_errors.BadRequestError{InternalError: err}
}

func fiber_endpoint2(c fiber.Ctx) error {
	err := errors.New("We have a specific status code error in our endpoint")
	return fiber.NewError(http.StatusBadGateway, err.Error())
}

func fiber_endpoint3(c fiber.Ctx) error {
	err := errors.New("We have an unhandled error in our endpoint")
	return err
}

func fiber_endpoint4(c fiber.Ctx) error {
	err := errors.New("We have a custom error with custom problem details error in our endpoint")
	return custom_errors.ConflictError{InternalError: err}
}

type CustomProblemDetailTest struct {
	ProblemDetailErr
	Description    string `json:"description,omitempty"`
	AdditionalInfo string `json:"additionalInfo,omitempty"`
}
