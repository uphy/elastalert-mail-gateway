package main

import (
	"fmt"

	"github.com/uphy/elastalert-mail-gateway/cli"
)

func main() {
	if err := cli.Run(); err != nil {
		fmt.Println("failed to run: ", err)
	}
}
