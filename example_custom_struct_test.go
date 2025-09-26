package faker_test

import (
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
)

type RedefinedTime time.Time

func Example_structTypeProviders() {
	var s struct {
		T    *RedefinedTime
		RP   *big.Rat
		R    big.Rat
		UUID faker.UUID
	}

	_ = faker.FakeData(&s, options.WithStructTypeProviders(RedefinedTime{}, func() (interface{}, error) {
		return time.Now(), nil
	}), options.WithStructTypeProviders(big.Rat{}, func() (interface{}, error) {
		return *big.NewRat(rand.Int63(), rand.Int63n(1<<63-1)+1), nil
	}))
	fmt.Printf("time: %v\n", time.Time(*s.T))
	fmt.Printf("pointer rational number: %s\n", s.RP.FloatString(12))
	fmt.Printf("rational number: %s\n", s.R.FloatString(12))
	/*
		Result:
		{
			time: 2025-09-26 19:32:30.091437 +0200 CEST m=+0.003062168
			pointer rational number: 3.781240379647
			rational number: 0.081930416292
		}
	*/

}
