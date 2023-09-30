package go_jeans

import (
	"errors"
	"sync"
)

var count uint32

var lock sync.Mutex

var ErrOfBytesToBaseType_float error = errors.New("float err: of BytesToBaseType float bounds out of max or min value")
var ErrOfBytesToBaseType_String error = errors.New("string err: of BytesToBaseType resolution length is greater than the remaining length")
var ErrOfBytesToBaseType_SliceBytes error = errors.New("slice byte err: of BytesToBaseType  resolution length is greater than the remaining length")
