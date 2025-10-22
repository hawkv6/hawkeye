package version

import "fmt"

const MAJOR uint = 1
const MINOR uint = 0
const PATCH uint = 1

func GetVersion() string {
	return fmt.Sprintf("%d.%d.%d", MAJOR, MINOR, PATCH)
}
