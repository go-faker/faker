package faker_test

import (
	"fmt"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
)

// SomeStructWithSliceLength ...
type SomeStructWithSliceLength struct {
	EmptyList       []string `faker:"slice_len=0"`
	FixedStringList []string `faker:"slice_len=2"`
	FixedIntList    []int64  `faker:"slice_len=4"`
	RandomIntList   []int64
	FixedURLList    []string `faker:"slice_len=3,url"` // any built-in tag fills each element
	RandomUUIDList  []string `faker:"uuid_digit"`      // random length, every element a uuid_digit
}

// You can set the size for your random slices.
func Example_withTagsSliceLength() {
	a := SomeStructWithSliceLength{}
	_ = faker.FakeData(&a, options.WithRandomMapAndSliceMaxSize(20)) // If no slice_len is set, this sets the max of the random size
	fmt.Printf("%+v", a)
	// Result:
	/*
		   {
		       EmptyList:[]
		       FixedStringList:[
			           geHYIpEoQhQdijFooVEAOyvtTwJOofbQPJdbHvEEdjueZaKIgI
			           WVJBBtmrrVccyIydAiLSkMwWbFzFMEotEXsyUXqcmBTVORlkJK
			   ]
			   FixedIntList:[10,25,60,15]
		       RandomIntList:[5,16,134,6235,123,53,123] //Random len() with a max of 20
		       FixedURLList:[https://xkcqnze.net/ http://qvrmzpo.com/cPmbsAB.html https://gwtdkla.org/]
		       RandomUUIDList:[3a6b9e0f1c2d4e5fa6b7c8d9e0f1a2b3 9f8e7d6c5b4a4c3d8e2f1a0b9c8d7e6f]
		   }
	*/
}

// You can set the size for your random slices.
func Test_withTagsSliceLength(t *testing.T) {
	a := SomeStructWithSliceLength{}
	err := faker.FakeData(&a, options.WithRandomMapAndSliceMaxSize(20)) // If no slice_len is set, this sets the max of the random size
	if err != nil {
		t.Errorf("want %v, got %v", nil, err)
	}
	if len(a.EmptyList) != 0 {
		t.Errorf("want %v, got %v", 0, len(a.EmptyList))
	}
	if len(a.FixedStringList) > 2 {
		t.Errorf("want <=%v, got %v", 2, len(a.FixedStringList))
	}
	if len(a.FixedIntList) > 4 {
		t.Errorf("want <=%v, got %v", 4, len(a.FixedIntList))
	}
	if len(a.FixedURLList) != 3 {
		t.Errorf("want %v, got %v", 3, len(a.FixedURLList))
	}
}
