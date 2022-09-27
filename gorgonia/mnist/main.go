package main

import (
	"fmt"
	"log"

	"gorgonia.org/gorgonia/examples/mnist"
	"gorgonia.org/tensor"
)

func main() {
	for _, typ := range []string{"test", "train"} {
		inputs, targets, err := mnist.Load(typ, "./testdata", tensor.Float64)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(typ+" inputs:", inputs.Shape())
		fmt.Println(typ+" data:", targets.Shape())
	}
}
