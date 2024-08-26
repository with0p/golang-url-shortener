package customerrors

import "errors"

var ErrUniqueKeyConstrantViolation = errors.New("unique key violation")
var ErrRecordHasBeenDeleted = errors.New("record has been deleted")
