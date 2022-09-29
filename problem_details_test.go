package problem

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func Test_Map_BadRequestErr(t *testing.T) {

	funProblem := func() *ProblemDetail {
		return &ProblemDetail{
			Status:    http.StatusBadRequest,
			Type:      "https://httpstatuses.io/400",
			Detail:    "We have a bad request error",
			Title:     "bad-request",
			Timestamp: time.Date(2022, 9, 12, 8, 0, 0, 0, time.Local),
		}
	}

	Map(http.StatusBadRequest, funProblem)

	assert.Equal(t, http.StatusBadRequest, funProblem().Status)
	assert.Equal(t, "bad-request", funProblem().Title)
	assert.Equal(t, "We have a bad request error", funProblem().Detail)
	assert.Equal(t, "https://httpstatuses.io/400", funProblem().Type)
	assert.Equal(t, time.Date(2022, 9, 12, 8, 0, 0, 0, time.Local), funProblem().Timestamp)
}

func Test_Map_NotFoundErr(t *testing.T) {

	funProblem := func() *ProblemDetail {
		return &ProblemDetail{
			Status:    http.StatusNotFound,
			Type:      "https://httpstatuses.io/404",
			Detail:    "We have a not found error",
			Title:     "not-found",
			Timestamp: time.Date(2022, 9, 12, 8, 0, 0, 0, time.Local),
		}
	}

	Map(http.StatusNotFound, funProblem)

	assert.Equal(t, http.StatusNotFound, funProblem().Status)
	assert.Equal(t, "not-found", funProblem().Title)
	assert.Equal(t, "We have a not found error", funProblem().Detail)
	assert.Equal(t, "https://httpstatuses.io/404", funProblem().Type)
	assert.Equal(t, time.Date(2022, 9, 12, 8, 0, 0, 0, time.Local), funProblem().Timestamp)
}

func Test_Map_InternalServerErr(t *testing.T) {

	funProblem := func() *ProblemDetail {
		return &ProblemDetail{
			Status:    http.StatusInternalServerError,
			Type:      "https://httpstatuses.io/500",
			Detail:    "We have a internal server error",
			Title:     "internal-server-error",
			Timestamp: time.Date(2022, 9, 12, 8, 0, 0, 0, time.Local),
		}
	}

	Map(http.StatusInternalServerError, funProblem)

	assert.Equal(t, http.StatusInternalServerError, funProblem().Status)
	assert.Equal(t, "internal-server-error", funProblem().Title)
	assert.Equal(t, "We have a internal server error", funProblem().Detail)
	assert.Equal(t, "https://httpstatuses.io/500", funProblem().Type)
	assert.Equal(t, time.Date(2022, 9, 12, 8, 0, 0, 0, time.Local), funProblem().Timestamp)
}

func Test_Map_UnauthorizedErr(t *testing.T) {

	funProblem := func() *ProblemDetail {
		return &ProblemDetail{
			Status:    http.StatusUnauthorized,
			Type:      "https://httpstatuses.io/401",
			Detail:    "We have a unauthorized error",
			Title:     "unauthorized",
			Timestamp: time.Date(2022, 9, 12, 8, 0, 0, 0, time.Local),
		}
	}

	Map(http.StatusUnauthorized, funProblem)

	assert.Equal(t, http.StatusUnauthorized, funProblem().Status)
	assert.Equal(t, "unauthorized", funProblem().Title)
	assert.Equal(t, "We have a unauthorized error", funProblem().Detail)
	assert.Equal(t, "https://httpstatuses.io/401", funProblem().Type)
	assert.Equal(t, time.Date(2022, 9, 12, 8, 0, 0, 0, time.Local), funProblem().Timestamp)
}

func Test_Map_ForbiddenErr(t *testing.T) {

	funProblem := func() *ProblemDetail {
		return &ProblemDetail{
			Status:    http.StatusForbidden,
			Type:      "https://httpstatuses.io/403",
			Detail:    "We have a forbidden error",
			Title:     "forbidden",
			Timestamp: time.Date(2022, 9, 12, 8, 0, 0, 0, time.Local),
		}
	}

	Map(http.StatusForbidden, funProblem)

	assert.Equal(t, http.StatusForbidden, funProblem().Status)
	assert.Equal(t, "forbidden", funProblem().Title)
	assert.Equal(t, "We have a forbidden error", funProblem().Detail)
	assert.Equal(t, "https://httpstatuses.io/403", funProblem().Type)
	assert.Equal(t, time.Date(2022, 9, 12, 8, 0, 0, 0, time.Local), funProblem().Timestamp)
}

func Test_Map_UnsupportedMediaTypeErr(t *testing.T) {

	funProblem := func() *ProblemDetail {
		return &ProblemDetail{
			Status:    http.StatusUnsupportedMediaType,
			Type:      "https://httpstatuses.io/415",
			Detail:    "We have a unsupported media type error",
			Title:     "unsupported-media-type",
			Timestamp: time.Date(2022, 9, 12, 8, 0, 0, 0, time.Local),
		}
	}

	Map(http.StatusUnsupportedMediaType, funProblem)

	assert.Equal(t, http.StatusUnsupportedMediaType, funProblem().Status)
	assert.Equal(t, "unsupported-media-type", funProblem().Title)
	assert.Equal(t, "We have a unsupported media type error", funProblem().Detail)
	assert.Equal(t, "https://httpstatuses.io/415", funProblem().Type)
	assert.Equal(t, time.Date(2022, 9, 12, 8, 0, 0, 0, time.Local), funProblem().Timestamp)
}

func Test_Map_BadGatewayErr(t *testing.T) {

	funProblem := func() *ProblemDetail {
		return &ProblemDetail{
			Status:    http.StatusBadGateway,
			Type:      "https://httpstatuses.io/502",
			Detail:    "We have a bad gateway error",
			Title:     "bad-gateway",
			Timestamp: time.Date(2022, 9, 12, 8, 0, 0, 0, time.Local),
		}
	}

	Map(http.StatusBadGateway, funProblem)

	assert.Equal(t, http.StatusBadGateway, funProblem().Status)
	assert.Equal(t, "bad-gateway", funProblem().Title)
	assert.Equal(t, "We have a bad gateway error", funProblem().Detail)
	assert.Equal(t, "https://httpstatuses.io/502", funProblem().Type)
	assert.Equal(t, time.Date(2022, 9, 12, 8, 0, 0, 0, time.Local), funProblem().Timestamp)
}
