package faker_test

import (
	"fmt"

	"github.com/go-faker/faker/v4"
)

type UserProfile struct {
	FirstName  string `faker:"first_name"`
	LastName   string `faker:"last_name"`
	DomainName string `faker:"domain_name"`

	// Template uses helper 'lower' and 'slug' (slug is lower+dash)
	Username string `faker:"template:{{.FirstName | lower}}.{{.LastName | lower}}"`
	Slug     string `faker:"template:{{.FirstName | slug}}-{{.LastName | slug}}"`
	Email    string `faker:"template:{{.Username}}@{{.DomainName}}"`
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
