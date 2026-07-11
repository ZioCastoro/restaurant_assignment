package validators

import (
	"fmt"
	"math"
	"os"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

var (
	validate     *validator.Validate
	onceValidate sync.Once
)

func getValidatorInstance() *validator.Validate {
	onceValidate.Do(func() {
		validate = validator.New()
		err := validate.RegisterValidation("notblank", validators.NotBlank)
		if err != nil {
			fmt.Fprintf(os.Stderr, "getValidatorInstance() error: %v\n", err)

			os.Exit(1)
		}

		err = validate.RegisterValidation("cent", validateCent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "getValidatorInstance() error: %v\n", err)

			os.Exit(1)
		}
	})

	return validate
}

func ValidateStruct(s any) error {
	if err := getValidatorInstance().Struct(s); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}

func ValidateSlice(s any) error {
	if err := getValidatorInstance().Var(s, "required,min=1,dive"); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}

func validateCent(fl validator.FieldLevel) bool {
	n := fl.Field().Float()
	scaled := math.Round(n * 100)

	return math.Abs(n*100-scaled) < 1e-9
}
