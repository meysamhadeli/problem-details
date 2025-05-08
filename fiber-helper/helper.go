package fiber_helper

import (
	"github.com/gofiber/fiber/v3"
	"net/http"
	"net/url"
)

type fiberResponseWriter struct {
	ctx     fiber.Ctx
	written bool
	headers http.Header
}

func Response(c fiber.Ctx) *fiberResponseWriter {
	return &fiberResponseWriter{
		ctx:     c,
		headers: make(http.Header),
	}
}

func (f *fiberResponseWriter) Header() http.Header {
	return f.headers
}

func (f *fiberResponseWriter) Write(data []byte) (int, error) {
	f.written = true
	return f.ctx.Response().BodyWriter().Write(data)
}

func (f *fiberResponseWriter) WriteHeader(statusCode int) {
	f.written = true
	f.ctx.Status(statusCode)
}

func Request(c fiber.Ctx) *http.Request {
	fiberURI := c.Request().URI()
	parsedURL, _ := url.Parse(string(fiberURI.FullURI()))

	return &http.Request{
		Method:     c.Method(),
		URL:        parsedURL,
		Header:     make(http.Header),
		RequestURI: string(c.Request().RequestURI()),
	}
}
