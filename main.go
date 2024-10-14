package main

import (
	"math/big"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"

	"github.com/PolyhedraZK/ExpanderCompilerCollection"
	"github.com/PolyhedraZK/ExpanderCompilerCollection/test"
)

type Circuit struct {
	A, B   frontend.Variable `gnark:",public"`
	Result frontend.Variable `gnark:",public"`
}

func (circuit *Circuit) Define(api frontend.API) error {
	a := circuit.A
	b := circuit.B

	for i := 2; i <= 100; i++ {
		next := api.Add(a, b)
		a, b = b, next
	}
	api.AssertIsEqual(circuit.Result, b)
	return nil
}

func main() {
	circuit, err := ExpanderCompilerCollection.Compile(ecc.BN254.ScalarField(), &Circuit{})
	if err != nil {
		panic(err)
	}

	c := circuit.GetLayeredCircuit()

	result, _ := new(big.Int).SetString("354224848179261915075", 10)
	assignment := &Circuit{A: 0, B: 1, Result: result}

	inputSolver := circuit.GetInputSolver()
	witness, err := inputSolver.SolveInputAuto(assignment)
	if err != nil {
		panic(err)
	}

	if !test.CheckCircuit(c, witness) {
		panic("verification failed")
	}

	os.WriteFile("circuit.txt", c.Serialize(), 0o644)
	os.WriteFile("witness.txt", witness.Serialize(), 0o644)
}
