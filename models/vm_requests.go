package models

func (req *VMsRequest) Validate() error {
	return nil
}

func (request *VMByCidRequest) Validate() error {
	var validationError ValidationError

	if request.Cid == 0 {
		validationError = validationError.Append(ErrInvalidField{"cid"})
	}

	if !validationError.Empty() {
		return validationError
	}

	return nil
}
