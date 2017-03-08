package models

import (
	"fmt"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"log"
)

type APIError struct {
	Status int `json:"status"`

	// Message to show the user.
	Message string `json:"message"`

	// Extra developer info for easier debugging.
	DeveloperInfo string `json:"developerInfo"`

	// Internal error that caused this error. Only for internal debugging use.
	causedBy error
}

func (error *APIError) Error() string {
	if error.Message != "" {
		return fmt.Sprintf("%d: %s", error.Status, error.Message)
	} else {
		return fmt.Sprintf("%d", error.Status)
	}
}

func NewAPIError(status int, message string, causedBy error) *APIError {
	return &APIError{status, message, "", causedBy}
}

func NewHTTP500Error() *APIError {
	return NewAPIError(500, "Internal Server Error", nil)
}

func (error *APIError) RenderAsJson(render render.Render, logger *log.Logger) {
	// All 500 errors are unexpected and should at least be logged.
	if error.Status == 500 {
		message := error.Message
		if error.causedBy != nil {
			message += ": " + error.causedBy.Error()
		}

		logger.Printf("%s", message)

		// Always render "Internal Server Error" for error 500, to make sure not to
		// leak internal error states.
		error.Message = "Internal Server Error"
	}

	render.JSON(error.Status, error)
}

func NewBindingRangeError(field string, min int, max int) binding.Error {
	return binding.Error{
		FieldNames:     []string{field},
		Classification: "ComplaintError",
		Message:        fmt.Sprintf("Name should be between %d and %d characters", min, max),
	}
}
