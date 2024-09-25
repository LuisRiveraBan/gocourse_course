package course

import (
	"errors"
	"fmt"
)

var ErrNameRequired = errors.New("first name is required")
var ErrEndDateRequired = errors.New("first End Date is required")
var ErrStartDateRequired = errors.New("first Start Date is required")
var ErrEndLesserStart = errors.New("end date mustn't be lesser than start date")
var ErrInvalidStartDate = errors.New("invalid start date")
var ErrInvalidEndDate = errors.New("invalid end date")

type ErrNotFound struct {
	Coursed string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("Course %s doesnt't exist", e.Coursed)
}
