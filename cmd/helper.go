package cmd

import (
	"os"
	"strings"
)

func splitter(path string) (name, location string) {
	ind := strings.LastIndex(path, string(os.PathSeparator))
	if ind != -1 {
		name = path[ind+1:]
		location = path[:ind]
	} else {
		name = path
		location = "."
	}
	return name, location
}
