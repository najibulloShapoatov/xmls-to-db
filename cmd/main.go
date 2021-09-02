package main

import (
	"fmt"
	"os"
	"photouploader/cmd/app"
)

func main() {
	if err := app.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
