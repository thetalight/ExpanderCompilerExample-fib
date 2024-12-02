package main

import (
	"math/big"
	"os"
	"time"
	"fmt"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"

	"github.com/PolyhedraZK/ExpanderCompilerCollection/ecgo"
	"github.com/PolyhedraZK/ExpanderCompilerCollection/test"
)

type Circuit struct {
	A, B   frontend.Variable
	Result frontend.Variable `gnark:",public"`
}

func (circuit *Circuit) Define(api frontend.API) error {
	a := circuit.A
	b := circuit.B

	// next = a^2 + b^2
	for i := 2; i <= 10000; i++ {
		aPower := api.Mul(a,a);
		bPower := api.Mul(b,b);
		next := api.Add(aPower, bPower)
		a, b = b, next
	}
	api.AssertIsEqual(circuit.Result, b)
	return nil
}

func main() {
	proveGnark()
//   expander()
}

func proveGnark() {
	result, _ := new(big.Int).SetString("20158110143150017413199636048408137057957725491036497926379928306369476614404", 10)
	assignment := &Circuit{A: 0, B: 1, Result: result}

	r1cs, _ := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &Circuit{})

	println("Nb Constraints: ", r1cs.GetNbConstraints())

	witness, _ := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	publicWitness, _ := witness.Public()
	pk, vk, _ := groth16.Setup(r1cs)
	start := time.Now()
	proof, _ := groth16.Prove(r1cs, pk, witness)
	end := time.Now()
	fmt.Printf("time: %v\n",  end.Sub(start))
	err := groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		panic(err)
	}
}

func expander() {
	circuit, err := ecgo.Compile(ecc.BN254.ScalarField(), &Circuit{})
	if err != nil {
		panic(err)
	}

	c := circuit.GetLayeredCircuit()

	result, _ := new(big.Int).SetString("20158110143150017413199636048408137057957725491036497926379928306369476614404", 10)
	assignment := &Circuit{A: 0, B: 1, Result: result}

	inputSolver := circuit.GetInputSolver()
	witness, err := inputSolver.SolveInputAuto(assignment)
	if err != nil {
		panic(err)
	}

	if !test.CheckCircuit(c, witness) {
		panic("verification failed")
	}

	os.WriteFile("./fib/circuit.txt", c.Serialize(), 0o644)
	os.WriteFile("./fib/witness.txt", witness.Serialize(), 0o644)
}
