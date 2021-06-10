package main

const (
	contactPath      = "contacts"
	flowPath         = "flows"
	flowPropertyPath = "flowproperties"
	methodPath       = "lciamethods"
	processPath      = "processes"
	sourcePath       = "sources"
	unitGroupPath    = "unitgroups"
)

func isValidPath(path string) bool {
	switch path {
	case
		processPath,
		flowPath,
		flowPropertyPath,
		unitGroupPath,
		contactPath,
		sourcePath,
		methodPath:
		return true
	}
	return false
}
