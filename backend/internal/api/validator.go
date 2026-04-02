package api

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/microcosm-cc/bluemonday"
)

type CustomValidator struct {
	validator *validator.Validate
	sanitizer *bluemonday.Policy
}

func NewCustomValidator() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
		sanitizer: bluemonday.StrictPolicy(),
	}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	cv.sanitizeStruct(reflect.ValueOf(i))
	return cv.validator.Struct(i)
}

func (cv *CustomValidator) sanitizeStruct(v reflect.Value) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		if !field.CanSet() {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			sanitized := cv.sanitizer.Sanitize(field.String())
			field.SetString(sanitized)
		case reflect.Struct:
			cv.sanitizeStruct(field.Addr())
		case reflect.Slice:
			for j := 0; j < field.Len(); j++ {
				elem := field.Index(j)
				if elem.Kind() == reflect.Struct {
					cv.sanitizeStruct(elem.Addr())
				} else if elem.Kind() == reflect.String {
					sanitized := cv.sanitizer.Sanitize(elem.String())
					elem.SetString(sanitized)
				} else if elem.Kind() == reflect.Ptr {
					cv.sanitizeStruct(elem)
				}
			}
		case reflect.Ptr:
			cv.sanitizeStruct(field)
		}
	}
}
