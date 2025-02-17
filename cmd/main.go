package main

import (
	"fmt"
	"os"

	"github.com/NilayYadav/agentos/cmd/cli"
)

func main() {

	err := cli.InitializeCli()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
