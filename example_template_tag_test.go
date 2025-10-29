package faker_test

import (
	"fmt"

	"github.com/go-faker/faker/v4"
)

type Address struct {
	Street string `faker:"street_name"`
	City   string `faker:"city"`
}

type UserProfile struct {
	FirstName string `faker:"first_name"`
	LastName  string `faker:"last_name"`
	Domain    string `faker:"domain"`

	// Template uses helper 'lower' and 'slug' (slug is lower+dash)
	Username string `faker:"template:{{.FirstName | lower}}.{{.LastName | lower}}"`
	Slug     string `faker:"template:{{.FirstName | slug}}-{{.LastName | slug}}"`
	Email    string `faker:"template:{{.Username}}@{{.Domain}}"`
	// (only template-related fields kept for this example)
}

func Example_templateTag() {
	var v UserProfile
	if err := faker.FakeData(&v); err != nil {
		fmt.Println("error:", err)
		return
	}
	// output the filled struct; values are random so we don't assert exact output
	fmt.Printf("%+v\n", v)
}
