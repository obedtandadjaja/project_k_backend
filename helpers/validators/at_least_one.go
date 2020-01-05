package validators

import (
	"fmt"

	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type AtLeastOne struct {
	Name       string
	Validators []validate.Validator
	Message    string
}

func (v *AtLeastOne) IsValid(errors *validate.Errors) {
	for _, validator := range v.Validators {
		newErrors := validate.NewErrors()
		if validator.IsValid(newErrors); !newErrors.HasAny() {
			return
		}
	}

	errors.Add(validators.GenerateKey(v.Name), fmt.Sprintf("None of the validators pass"))
}
