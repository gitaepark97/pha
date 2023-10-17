package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ErrValidate(err validator.ValidationErrors, obj interface{}, tag string) (vErr error) {
	e := reflect.TypeOf(obj).Elem()
	fieldList := getErrFieldList(err)
	tagName := findTagName(e, tag, fieldList)

	switch err[0].ActualTag() {
	case "required":
		vErr = ErrRequired(tagName)
	case "max":
		vErr = ErrMax(tagName, err[0].Param())
	case "phone_number":
		vErr = ErrPhoneNumber(tagName)
	case "product_size":
		vErr = ErrProductSize(tagName)
	case "date":
		vErr = ErrDate(tagName)
	default:
		vErr = err

	}

	return
}

func ErrRequired(field string) error {
	return fmt.Errorf("%s should be required", field)
}

func ErrMax(field string, param string) error {
	return fmt.Errorf("%s's length should be smaller than or equals to %s", field, param)
}

func ErrPhoneNumber(field string) error {
	return fmt.Errorf("%s should be phone number format", field)
}

func ErrProductSize(field string) error {
	return fmt.Errorf("%s should be small or large", field)
}

func ErrDate(field string) error {
	return fmt.Errorf("%s should be 0000-00-00 format", field)
}

func getErrFieldList(err validator.ValidationErrors) []string {
	reg := regexp.MustCompile(`\[[0-9]*\]`)
	return strings.Split(reg.ReplaceAllString(err[0].Namespace(), ""), ".")[1:]
}

func findTagName(t reflect.Type, tag string, fieldList []string) (tagName string) {
	field, _ := t.FieldByName(fieldList[0])
	tagName, _ = field.Tag.Lookup(tag)

	if len(fieldList) == 1 {
		return
	} else {
		field, _ := t.FieldByName(fieldList[0])

		fieldType := field.Type.Elem()
		tagName += "." + findTagName(fieldType, tag, fieldList[1:])

		return
	}
}
