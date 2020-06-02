package web

import (
	"encoding/json"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strings"
)

//
// Ctx - Context
//
type Ctx struct {
	*gin.Context
	Router *Router
}

//
// BindAndValidate - Validate Json Body
//
func (ctx *Ctx) BindAndValidate(obj interface{}) *ErrorResponse {
	var errs ErrorResponse
	var addError = func(key string, value interface{}) {
		errs.AddErrors(FieldError{Field: key, Message: value})
	}
	if err := ctx.ShouldBindJSON(obj); err != nil {
		switch err.(type) {
		case validator.ValidationErrors:
			valErr := err.(validator.ValidationErrors)
			for _, fe := range valErr {
				errs.AddErrors(FromFE(ctx, fe, "json", obj, ""))
			}
		case *json.UnmarshalTypeError:
			typeErr := err.(*json.UnmarshalTypeError)
			addError("_json", fmt.Sprintf("%v contain invalid data type", typeErr.Field))
		default:
			addError("_syntax", "Syntax ErrorMessage: "+err.(error).Error())
		}
	}
	if len(errs.Errors) < 1 {
		return nil
	}
	return &errs
}

//
// BindAndValidateQuery - Validate query string from struct
//
func (ctx *Ctx) BindAndValidateQuery(obj interface{}) *ErrorResponse {
	var errs ErrorResponse
	var addError = func(key string, value interface{}) {
		errs.AddErrors(FieldError{Field: fmt.Sprintf("query_%s", key), Message: value})
	}
	if err := ctx.ShouldBindQuery(obj); err != nil {
		switch err.(type) {
		case validator.ValidationErrors:
			valErr := err.(validator.ValidationErrors)
			for _, fe := range valErr {
				errs.AddErrors(FromFE(ctx, fe, "form", obj, "query"))
			}
		default:
			addError("_syntax", "Syntax ErrorMessage: "+err.(error).Error())
		}
	}
	if len(errs.Errors) < 1 {
		return nil
	}
	return &errs
}

//
// ErrorWithStatus - Prepare error response
//
func (ctx *Ctx) ErrorWithStatus(status int, err string, errs ...error) {
	res := ErrorResponse{
		Status:       status,
		ErrorMessage: err,
		Errors:       errs,
	}
	ctx.Error(&res)
}

//
// Send JSON ErrorMessage structure
//
func (ctx *Ctx) Error(err error) {
	var e *ErrorResponse
	switch (err).(type) {
	case IErrorResponse:
		e = (err).(*ErrorResponse)
	default:
		e = BadRequestError()
		e.ErrorMessage = err.Error()
	}
	e.Path = ctx.Request.RequestURI
	e = e.Build()
	ctx.JSON(e.Status, e)
	ctx.Abort()
}

//
// RequiredParam - required parameters
//
func (ctx Ctx) RequiredParam(paramName string) (*string, error) {
	if id := ctx.Param("id"); len(strings.TrimSpace(id)) == 0 {
		return nil, fmt.Errorf("please provide %s url parameter", paramName)
	} else {
		return &id, nil
	}
}

//
// GetJwtClaims - Get JWT claims from the context
//
func (ctx Ctx) GetJwtClaims() jwt.MapClaims {
	return jwt.ExtractClaims(ctx.Context)
}

//
// GetCurrentUsername - Get the current logged in user.
//
func (ctx Ctx) GetCurrentUsername() string {
	claims := ctx.GetJwtClaims()
	//TODO: "auth.identity" should come from configuration
	if username, ok := claims["auth.identity"].(string); ok {
		return username
	}
	return ""
}
