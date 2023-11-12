package validation

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gookit/validate"
)

type (
	bindingWithValidation struct{}
)

func NewBindingWithValidation() *bindingWithValidation {
	newValidators()
	return &bindingWithValidation{}
}

func (b bindingWithValidation) Name() string {
	return "binding_with_validation"
}

func (b bindingWithValidation) Bind(r *http.Request, ptr any) error {
	rules, err := b.extractRules(ptr, "")
	if err != nil {
		return err
	}

	v := b.validate(r, rules)

	if v.Errors.Empty() {
		return v.BindSafeData(ptr)
	}

	b.replaceErrors(v)
	return v.Errors
}

func (b bindingWithValidation) validate(r *http.Request, rules map[string]string) *validate.Validation {
	v := validate.Request(r)

	// immediately checking - if the json is invalid, there will be an error
	if v.IsFail() {
		return v
	}

	v.AddMessages(errorMessages)
	v.StopOnError = false

	for field, rule := range rules {
		v.StringRule(field, rule)
	}

	v.Validate()
	return v
}

func (b bindingWithValidation) extractRules(ptr any, prefix string) (map[string]string, error) {
	ptrType, ok := ptr.(reflect.Type)
	if !ok {
		ptrType = reflect.TypeOf(ptr)
	}

	if ptrType.Kind() == reflect.Ptr {
		ptrType = ptrType.Elem()
	} else if ptrType.Kind() != reflect.Struct {
		return nil, errors.New("invalid type")
	}

	result := make(map[string]string)
	for i := 0; i < ptrType.NumField(); i++ {
		field := ptrType.Field(i)

		key := field.Tag.Get("json")
		if key != "" {
			key = strings.Split(key, ",")[0]
		} else {
			key = field.Name
		}
		key = fmt.Sprintf("%s%s", prefix, key)

		fieldType := field.Type.Kind()
		if fieldType == reflect.Ptr {
			fieldType = field.Type.Elem().Kind()
		}

		if fieldType == reflect.Struct {
			rules, err := b.extractRules(field.Type, fmt.Sprintf("%s.", key))
			if err == nil {
				for k, rule := range rules {
					result[k] = rule
				}
			}
		}

		validateTag := field.Tag.Get("validate")

		var variableType string
		if fieldType == reflect.Bool || fieldType == reflect.String {
			variableType = fieldType.String()
		} else if fieldType == reflect.Float32 || fieldType == reflect.Float64 {
			variableType = "floatType"
		} else if strings.Contains(fieldType.String(), "int") {
			variableType = "intType"
		}

		if variableType != "" {
			if validateTag == "" {
				validateTag = variableType
			} else {
				validateTag = fmt.Sprintf("%s|%s", variableType, validateTag)
			}
		}

		if validateTag == "" {
			continue
		}

		result[key] = validateTag
	}

	return result, nil
}

func (b bindingWithValidation) replaceErrors(v *validate.Validation) {
	// I would like to add a check here via errors.As / errors.Is,
	// but the library’s own type does not support this,
	// since it writes errors as a string

	const (
		jsonInvalidCharacterError = "invalid character"
	)

	if strings.Contains(v.Errors.String(), jsonInvalidCharacterError) {
		v.Errors = validate.Errors{}
		v.AddError("_validate", "_validate", errorMessages["_json"])
	} else if strings.Contains(v.Errors.String(), validate.ErrEmptyData.Error()) {
		v.Errors = validate.Errors{}
		v.Errors.Add("_validate", "_validate", errorMessages["_noData"])
	}
}
