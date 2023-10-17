package middleware

import (
	"fmt"
	"net/http"

	"github.com/gitaepark/pha/service"
)

var (
	errEmptyAuthorizationHeader   = service.CustomErr{Code: http.StatusUnauthorized, Err: fmt.Errorf("authorization header is not provided")}
	errInvalidAuthorizationHeader = service.CustomErr{Code: http.StatusUnauthorized, Err: fmt.Errorf("invalid authorization header format")}
	errInvalidAuthorizationBearer = service.CustomErr{Code: http.StatusUnauthorized, Err: fmt.Errorf("unsupported authorization type")}
)

func errToken(err error) service.CustomErr {
	return service.CustomErr{Code: http.StatusUnauthorized, Err: err}
}
