package errors

import (
	. "errors"
)

var ValueParseError = New("ValueParseError")
var PriorityQueueEmptyError = New("PriorityQueueEmptyError")
