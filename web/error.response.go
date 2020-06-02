package web

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/sayuthisobri/goutils/list"
	"github.com/sayuthisobri/goutils/text"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"time"
)

//
// Field related error
//
type FieldError struct {
	Field   string      `json:"field,omitempty"`
	Message interface{} `json:"message"`
}

//
// Error - error
//
func (e FieldError) Error() string {
	switch e.Message.(type) {
	case string:
		return e.Message.(string)
	}
	return ""
}

//
// Construct error from validation result
//
func FromFE(ctx *Ctx, fe validator.FieldError, preferedTag string, obj interface{}, fieldPrefix string) FieldError {
	e := FieldError{}
	field := fe.Field()
	var jsonField string
	if field, ok := reflect.TypeOf(obj).Elem().FieldByName(field); ok {
		jsonField = field.Tag.Get(preferedTag)
		if jsons := list.Filter(strings.Split(jsonField, ","), func(i interface{}) bool {
			return sort.SearchStrings([]string{"omitempty", "-"}, i.(string)) > -1
		}).([]string); len(jsons) > 0 {
			jsonField = jsons[0]
		}
	}
	if len(jsonField) > 0 {
		field = jsonField
	} else {
		field = text.NewString(field).ToSnakeCase()
	}
	if fieldPrefix != "" {
		field = fmt.Sprintf("%s_%s", fieldPrefix, field)
	}
	e.Field = field
	e.Message = fe.Translate(ctx.Router.DI.Get("translator").(ut.Translator))
	return e
}

//
// IErrorResponse - Interface for error response
//
type IErrorResponse interface {
	error
	Build() *ErrorResponse
}

//
// ErrorResponse - error response
//
type ErrorResponse struct {
	Timestamp    string  `json:"timestamp"`
	Status       int     `json:"status"`
	ErrorMessage string  `json:"error,omitempty"`
	Errors       []error `json:"errors"`
	Path         string  `json:"path"`
}

//
// ER - error response
//
func ER(statusCode int, err string, path string, errs []error) *ErrorResponse {
	return (&ErrorResponse{
		Status:       statusCode,
		ErrorMessage: err,
		Errors:       errs,
		Path:         path,
	}).Build()
}

//
// BadRequestError - bad request error
//
func BadRequestError(errs ...error) *ErrorResponse {
	return (&ErrorResponse{Errors: append([]error{}, errs...)}).Build()
}

//
// AddErrors - add errors
//
func (r *ErrorResponse) AddErrors(errs ...error) {
	r.Errors = append(r.Errors, errs...)
}

//
// Error - error
//
func (r ErrorResponse) Error() string {
	return r.ErrorMessage
}

//
// Build - build
//
func (r ErrorResponse) Build() *ErrorResponse {
	r.Timestamp = time.Now().UTC().Format(time.RFC3339Nano)
	for _, err := range r.Errors {
		switch err.(type) {
		case FieldError:
			break
		case error:
			err = FieldError{
				Message: err.Error(),
			}

		}
	}
	if r.Status == 0 {
		r.Status = http.StatusBadRequest
	}
	if len(r.ErrorMessage) == 0 {
		r.ErrorMessage = http.StatusText(r.Status)
	}
	return &r
}
