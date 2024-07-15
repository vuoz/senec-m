package db

type UserInputError struct {
	reason string
}

func (err *UserInputError) Error() string {
	return err.reason

}
func NewUserInputError(reason string) error {
	return &UserInputError{
		reason: reason,
	}

}
