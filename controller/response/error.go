package response

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gitaepark/pha/service"
	"github.com/gitaepark/pha/util/validator"
)

var ErrParseString = fmt.Errorf("params invalid syntax")

func NewErrResponse(ctx *gin.Context, cErr service.CustomErr) {
	ctx.AbortWithStatusJSON(cErr.Code, gin.H{"meta": gin.H{"code": cErr.Code, "message": cErr.Err.Error()}, "data": nil})
}

func NewErrBindingResponse(ctx *gin.Context, err error, obj interface{}, tag string) {

	if _, ok := err.(*strconv.NumError); ok {
		err = ErrParseString
	}
	if tErr, ok := err.(*json.UnmarshalTypeError); ok {
		err = ErrType(tErr.Field, tErr.Type.Name())
	}
	if vErrs, ok := err.(validator.ValidationErrors); ok {
		err = validator.ErrValidate(vErrs, obj, tag)
	}

	NewErrResponse(ctx, service.NewErrBadRequest(err))
}

func ErrType(field string, fieldType string) error {
	return fmt.Errorf("%s should be %s type", field, fieldType)
}
