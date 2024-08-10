package customerrors

import "errors"

var ErrUniqueKeyConstrantViolation = errors.New("unique key violation")
