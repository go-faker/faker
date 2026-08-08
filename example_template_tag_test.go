package faker_test

import (
	"fmt"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
)

// UserProfile derives Username, Slug and Email from the fields faker generates.
// Note that Email is declared before the Username it reads: template fields are
// evaluated in dependency order, not declaration order.
type UserProfile struct {
	FirstName  string `faker:"first_name"`
	LastName   string `faker:"last_name"`
	DomainName string `faker:"domain_name"`

	Email    string `faker:"template:{{.Username}}@{{.DomainName}}"`
	Username string `faker:"template:{{.FirstName | lower}}.{{.LastName | lower}}"`
	Slug     string `faker:"template:{{.FirstName | slug}}-{{.LastName | slug}}"`
}

func Example_templateTag() {
	// The inputs are pinned so the derived values are reproducible; in real use
	// you would let faker generate them.
	var v UserProfile
	err := faker.FakeData(&v,
		options.WithCustomFieldProvider("FirstName", func() (any, error) { return "John", nil }),
		options.WithCustomFieldProvider("LastName", func() (any, error) { return "Doe", nil }),
		options.WithCustomFieldProvider("DomainName", func() (any, error) { return "example.org", nil }),
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("Username:", v.Username)
	fmt.Println("Slug:", v.Slug)
	fmt.Println("Email:", v.Email)
	// Output:
	// Username: john.doe
	// Slug: john-doe
	// Email: john.doe@example.org
}
