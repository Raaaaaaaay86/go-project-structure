package validate

type Validator interface {
	Validate() error
}

func Do(validations ...Validator) error {
	for _, validation := range validations {
		err := validation.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
