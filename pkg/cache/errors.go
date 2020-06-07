package cache

import "errors"

var (
	ValueNotFoundErr  = errors.New("Unable to locate value")
	StoreNotActiveErr = errors.New("The Underlying store is not avaialble")
)
