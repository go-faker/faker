package faker_test

import (
	"testing"

	"github.com/go-faker/faker/v4"
)

func TestFakeDataTaggedArrays(t *testing.T) {
	var value struct {
		Names  [2]string `faker:"oneof: alpha"`
		Counts [3]int    `faker:"oneof: 7"`
		Empty  [0]string `faker:"oneof: alpha"`
		Slice  []string  `faker:"slice_len=2,oneof: alpha"`
	}
	if err := faker.FakeData(&value); err != nil {
		t.Fatal(err)
	}
	if value.Names != [2]string{"alpha", "alpha"} {
		t.Fatalf("names = %v", value.Names)
	}
	if value.Counts != [3]int{7, 7, 7} {
		t.Fatalf("counts = %v", value.Counts)
	}
	if len(value.Slice) != 2 || value.Slice[0] != "alpha" || value.Slice[1] != "alpha" {
		t.Fatalf("slice = %v", value.Slice)
	}
}
