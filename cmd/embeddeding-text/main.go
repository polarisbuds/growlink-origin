package main

import (
	"fmt"
	"strings"
)

func main() {
	cuttedStr, _ := strings.CutSuffix("mei.sei@example.com", "@example.com")
	fmt.Println(cuttedStr)

	var lineworksID string
	atMarkIdx := strings.Index("mei.sei@example.com", "@")
	if atMarkIdx >= 0 {
		lineworksID = "mei.sei@example.com"[:atMarkIdx]
	}
	fmt.Println(lineworksID)
}
