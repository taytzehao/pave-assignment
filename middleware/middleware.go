package middleware

import (
	"encore.dev/middleware"
	"encore.dev/beta/errs"
    "github.com/go-playground/validator/v10"

)

//encore:middleware global target=all
func ValidationMiddleware(req middleware.Request, next middleware.Next) middleware.Response {
    
    payload := req.Data().Payload

    // Create a new validator instance
    validate := validator.New()

    // Validate the payload using the validator
    if payload == nil {
        return next(req)
    }
    if err := validate.Struct(payload); err != nil {
        // If the validation fails, return an InvalidArgument error.
        err = errs.WrapCode(err, errs.InvalidArgument, "validation failed")
        return middleware.Response{Err: err}
    }

    return next(req)
}