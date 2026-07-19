package faker

import (
	"fmt"
	"reflect"

	"github.com/go-faker/faker/v4/pkg/options"
)

var bloodTypes = []string{"O", "A", "B", "AB"}
var bloodRhFactors = []string{"+", "-"}

func GetBlood(opts ...options.OptionFunc) Blooder {
	opt := options.BuildOptions(opts)
	db := &Blood{
		fakerOption: *opt,
	}
	return db
}

type Blooder interface {
	BloodType(v reflect.Value) (any, error)
	BloodRHFactor(v reflect.Value) (any, error)
	BloodGroup(v reflect.Value) (any, error)
}

// Internet struct
type Blood struct {
	fakerOption options.Options
}

func (b Blood) bloodType() string {
	return randomElementFromSliceString(bloodTypes)
}

func (b Blood) BloodType(v reflect.Value) (any, error) {
	return b.bloodType(), nil
}

func (b Blood) bloodRhFactor() string {
	return randomElementFromSliceString(bloodRhFactors)
}

func (b Blood) BloodRHFactor(v reflect.Value) (any, error) {
	return b.bloodRhFactor(), nil
}

func (b Blood) bloodGroup() string {
	return fmt.Sprintf("%s%s", b.bloodType(), b.bloodRhFactor())
}

func (b Blood) BloodGroup(v reflect.Value) (any, error) {
	return b.bloodGroup(), nil
}
