package main

const (
	processPath      = "processes"
	flowPath         = "flows"
	flowPropertyPath = "flowproperties"
	unitGroupPath    = "unitgroups"
	contactPath      = "contacts"
	sourcePath       = "sources"
	methodPath       = "lciamethods"
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
