package main

import "errors"

var (
	errUnknownDataStock = errors.New("Unknown data stock")
	errInvalidPath      = errors.New("Invalid path")
	errInvalidDataSet   = errors.New("Invalid data set")
	errDataSetExists    = errors.New("Data set already exists")
)
