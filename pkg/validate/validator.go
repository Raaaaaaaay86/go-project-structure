package validate

type Validator interface {
	Validate() error
}

func Do(validators ...Validator) error {
	for _, validator := range validators {
		err := validator.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
