package custom_problems

import "github.com/meysamhadeli/problem-details"

type CustomProblemDetail struct {
	problem.ProblemDetailErr
	Description    string `json:"description,omitempty"`
	AdditionalInfo string `json:"additionalInfo,omitempty"`
}
