package main

import (
	"log"

	"gorgonia.org/gorgonia"
)

func main() {
	g := gorgonia.NewGraph()

	var x, y, z *gorgonia.Node
	var err error

	// define the expression
	x = gorgonia.NewScalar(g, gorgonia.Float64, gorgonia.WithName("x"))
	y = gorgonia.NewScalar(g, gorgonia.Float64, gorgonia.WithName("y"))
	z, err = gorgonia.Add(x, y)
	if err != nil {
		log.Fatal(err)
	}

	// create a VM to run the program on
	machine := gorgonia.NewTapeMachine(g)
	defer machine.Close()

	// set initial values then run
	gorgonia.Let(x, 2.0)
	gorgonia.Let(y, 2.5)

	err = machine.RunAll()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%v\n", z.Value())
}
