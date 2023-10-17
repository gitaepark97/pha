package service

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

type CustomErr struct {
	Code int
	Err  error
}

var (
	errDuplicatePhoneNumber   = CustomErr{Code: http.StatusBadRequest, Err: fmt.Errorf("duplicate phone number")}
	errNotFoundUser           = CustomErr{Code: http.StatusNotFound, Err: fmt.Errorf("not found user")}
	errWrongPassword          = CustomErr{Code: http.StatusBadRequest, Err: fmt.Errorf("wrong password")}
	errNotFoundSession        = CustomErr{Code: http.StatusNotFound, Err: fmt.Errorf("not found session")}
	errBlockedSession         = CustomErr{Code: http.StatusUnauthorized, Err: fmt.Errorf("blocked session")}
	errIncorrectSessionUser   = CustomErr{Code: http.StatusUnauthorized, Err: fmt.Errorf("incorrect session user")}
	errMismatchedSessionToken = CustomErr{Code: http.StatusUnauthorized, Err: fmt.Errorf("mismatched session token")}
	errExpiredSession         = CustomErr{Code: http.StatusUnauthorized, Err: fmt.Errorf("expired session")}

	errParseDate        = CustomErr{Code: http.StatusBadRequest, Err: fmt.Errorf("invalid date format")}
	errNotFoundProduct  = CustomErr{Code: http.StatusNotFound, Err: fmt.Errorf("not found product")}
	errForbiddenProduct = CustomErr{Code: http.StatusForbidden, Err: fmt.Errorf("only get your product")}
	errDuplicateBarcode = CustomErr{Code: http.StatusBadRequest, Err: fmt.Errorf("duplicate barcode")}
)

func NewErrInternalServer(err error) CustomErr {
	log.Error().Msg(err.Error())
	return CustomErr{Code: http.StatusInternalServerError, Err: fmt.Errorf("internal server error")}
}

func NewErrBadRequest(err error) CustomErr {
	return CustomErr{Code: http.StatusBadRequest, Err: err}
}
