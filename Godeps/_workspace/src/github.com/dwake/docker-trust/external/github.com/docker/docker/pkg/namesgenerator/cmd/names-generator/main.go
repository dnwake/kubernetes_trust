package main

import (
	"fmt"

	"github.com/dwake/docker-trust/external/github.com/docker/docker/pkg/namesgenerator"
)

func main() {
	fmt.Println(namesgenerator.GetRandomName(0))
}
