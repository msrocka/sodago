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

// authenticatePath is the prefix of the authentication routes. These routes
// are not protected, so that a client can always request a new token.
const authenticatePath = "/resource/authenticate/"

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
