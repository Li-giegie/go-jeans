package go_jeans

import (
	"errors"
	"sync"
)

var count uint32

var lock sync.Mutex

var ErrOfBytesToBaseType_float error = errors.New("float err: of BytesToBaseType float bounds out of max or min value")
