package customerrors

import "errors"

var ErrUniqueKeyConstrantViolation = errors.New("Unique key violation")
