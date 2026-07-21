package faker_test

import (
	"fmt"

	"github.com/catdevman/faker"
	"github.com/catdevman/faker/pkg/options"
)

// You can set length for your random strings also set boundary for your integers.
// SomeStructForLanguage ...
type SomeStructForLanguage struct {
	StringENG string `faker:"lang=eng"`
	StringCHI string `faker:"lang=chi"`
	StringRUS string `faker:"lang=rus"`
	StringJPN string `faker:"lang=jpn"`
	StringKOR string `faker:"lang=kor"`
	StringEMJ string `faker:"lang=emj"`
}

func Example_withTagsLang() {
	a := SomeStructForLanguage{}
	_ = faker.FakeData(&a, options.WithRandomStringLength(5))
	fmt.Printf("%+v", a)
	// Result:
	/*
		   {
			   StringENG:VVcaPS
			   StringCHI: 随机字符串
			   StringRUS:ваЩфз
			   StringJPN:びゃほぱヒてふ
			   StringKOR:텻밚쨋큊몉
			   StringEMJ:🐅😄🕢🍪🐡
		   }
	*/
}
