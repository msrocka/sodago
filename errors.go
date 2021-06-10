package main

import "errors"

var (
	errUnknownDataStock = errors.New("unknown data stock")
	errInvalidPath      = errors.New("invalid path")
	errInvalidDataSet   = errors.New("invalid data set")
	errDataSetExists    = errors.New("data set already exists")
	errDataSetNotExists = errors.New("data set does not exist")
)
