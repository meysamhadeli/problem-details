package problem

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

// ProblemDetail error struct
type ProblemDetail struct {
	Status     int       `json:"status,omitempty"`
	Title      string    `json:"title,omitempty"`
	Detail     string    `json:"detail,omitempty"`
	Type       string    `json:"type,omitempty"`
	Timestamp  time.Time `json:"timestamp,omitempty"`
	StackTrace string    `json:"stackTrace,omitempty"`
}

var mappers = map[int]func() *ProblemDetail{}

// WriteTo writes the JSON Problem to an HTTP Response Writer
func (p *ProblemDetail) writeTo(w http.ResponseWriter) (int, error) {
	p.writeHeaderTo(w)
	return w.Write(p.json())
}

// Map map error to problem details error
func Map(statusCode int, funcProblem func() *ProblemDetail) {
	mappers[statusCode] = funcProblem
}

// ResolveProblemDetails retrieve and resolve error with format problem details error
func ResolveProblemDetails(w http.ResponseWriter, err error) (int, error) {

	var statusCode int = http.StatusInternalServerError

	var echoError *echo.HTTPError

	var ginError *gin.Error

	if errors.As(err, &echoError) {
		if err.(*echo.HTTPError).Code != http.StatusOK {
			statusCode = err.(*echo.HTTPError).Code
		}
		err = err.(*echo.HTTPError).Message.(error)
	} else if errors.As(err, &ginError) {
		var rw = w.(gin.ResponseWriter)
		if rw.Status() != http.StatusOK {
			statusCode = rw.Status()
		}
		err = err.(*gin.Error)
	}

	problem := mappers[statusCode]

	if problem != nil {
		problem := problem()

		validationProblems(problem, err, statusCode)

		val, err := problem.writeTo(w)

		if err != nil {
			return 0, err
		}

		return val, err
	}

	defaultProblem := func() *ProblemDetail {
		return &ProblemDetail{
			Type:      getDefaultType(statusCode),
			Status:    statusCode,
			Detail:    err.Error(),
			Timestamp: time.Now(),
			Title:     http.StatusText(statusCode),
		}
	}

	val, err := defaultProblem().writeTo(w)

	if err != nil {
		return 0, err
	}

	return val, nil
}

func validationProblems(problem *ProblemDetail, err error, statusCode int) {
	problem.Detail = err.Error()

	if problem.Status == 0 {
		problem.Status = statusCode
	}
	if problem.Timestamp.IsZero() {
		problem.Timestamp = time.Now()
	}
	if problem.Type == "" {
		problem.Type = getDefaultType(problem.Status)
	}
	if problem.Title == "" {
		problem.Title = http.StatusText(problem.Status)
	}
}

func (p *ProblemDetail) writeHeaderTo(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/problem+json")

	w.WriteHeader(p.Status)
}

func (p *ProblemDetail) json() []byte {
	res, _ := json.Marshal(&p)
	return res
}

func getDefaultType(statusCode int) string {
	return fmt.Sprintf("https://httpstatuses.io/%d", statusCode)
}
