package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// create a global validator instance // 4.0
var validate = validator.New()

func Validate(i interface{}) error {
	err := validate.Struct(i)
	if err != nil {
		var errMessage []string
		for _, e := range err.(validator.ValidationErrors) {
			errMessage = append(errMessage, fmt.Sprintf("%s is required ", strings.ToLower(e.Field())))
		}

		return fmt.Errorf("%s", strings.Join(errMessage, ","))

	}
	return nil

}
