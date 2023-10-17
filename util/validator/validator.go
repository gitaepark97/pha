package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

const (
	PHONE_NUMBER_REGEX = `^010([0-9]{4})([0-9]{4})$`

	DATE_REGEX = `^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])$`

	Small = "small"
	Large = "large"
)

type Validate = validator.Validate
type ValidationErrors = validator.ValidationErrors

// validator 휴대폰 번호 양식(01000000000) 검증 함수
var ValidatePhoneNumber validator.Func = func(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(string); ok {
		return validateRegex(PHONE_NUMBER_REGEX, value)
	}
	return true
}

// validator 날짜 양식(0000-00-00) 검증 함수
var ValidateDate validator.Func = func(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(string); ok {
		return validateRegex(DATE_REGEX, value)
	}
	return true
}

// validator 상품 사이즈 종류(small, large) 검증 함수
var ValidateProductSize validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if value, ok := fieldLevel.Field().Interface().(string); ok {
		return IsSupportedProductSize(value)
	}

	return false
}

// 정규식 검증 함수
func validateRegex(regex, value string) bool {
	reg := regexp.MustCompile(regex)

	return reg.MatchString(value)
}

// 상품 사이즈 종류(small, large) 검증 함수
func IsSupportedProductSize(size string) bool {
	switch size {
	case Small, Large:
		return true
	}

	return false
}
