package shop

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

func GetProductDetailFailure(errs error, printer *message.Printer) validation.Errors {
	errMap := errs.(validation.Errors)

	var errorString string
	for key, err := range errMap {
		ozzoError := err.(validation.Error)
		code := ozzoError.Code()
		params := ozzoError.Params()

		switch code {
		case validation.ErrRequired.Code():
			errorString = printer.Sprintf("cannot be blank")
		case validation.ErrLengthOutOfRange.Code():
			errorString = printer.Sprintf("the length must be between %v and %v", number.Decimal(params["min"]), number.Decimal(params["max"]))
		case validation.ErrMinGreaterEqualThanRequired.Code():
			errorString = printer.Sprintf("must be no less than %v", number.Decimal(params["threshold"]))
		default:
			errorString = printer.Sprintf("unknown validation error")
		}

		newError := ozzoError.SetMessage(errorString)
		errMap[key] = newError
	}
	return errMap
}
