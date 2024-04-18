package main

import (
	"fmt"
	"os"
)

func main() {
	dest := "./backend/dat.go"
	f, _ := os.Create(dest)
	defer f.Close()

	body := fmt.Sprintf(
		"package backend\n"+
			"\n"+
			"const (\n"+
			"\tAppClientID     = \"%v\"\n"+
			"\tAppClientSecret = \"%v\"\n"+
			")\n",
		os.Getenv("AppClientID"), os.Getenv("AppClientSecret"),
	)
	f.WriteString(body)
}
