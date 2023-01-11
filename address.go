package faker

import (
	_ "embed"
	"encoding/json"
	"reflect"

	"github.com/go-faker/faker/v4/pkg/options"
)

var (
	//go:embed misc/addresses-us-1000.min.json
	addressesUSBytes []byte
	addressesUS      []RealAddress
)

func init() {
	data := struct {
		Addresses []RealAddress `json:"addresses"`
	}{}
	if err := json.Unmarshal(addressesUSBytes, &data); err != nil {
		panic(err)
	}
	addressesUS = data.Addresses
}

// GetAddress returns a new Addresser interface of Address
func GetAddress() Addresser {
	address := &Address{}
	return address
}

// Addresser is logical layer for Address
type Addresser interface {
	Latitude(v reflect.Value) (interface{}, error)
	Longitude(v reflect.Value) (interface{}, error)
}

// Address struct
type Address struct{}

func (i Address) latitude() float32 {
	return (rand.Float32() * 180) - 90
}

// Latitude sets latitude of the address
func (i Address) Latitude(v reflect.Value) (interface{}, error) {
	kind := v.Kind()
	val := i.latitude()
	if kind == reflect.Float32 {
		return val, nil
	}
	return float64(val), nil
}

func (i Address) longitude() float32 {
	return (rand.Float32() * 360) - 180
}

// Longitude sets longitude of the address
func (i Address) Longitude(v reflect.Value) (interface{}, error) {
	kind := v.Kind()
	val := i.longitude()
	if kind == reflect.Float32 {
		return val, nil
	}
	return float64(val), nil
}

// Longitude get fake longitude randomly
func Longitude(opts ...options.OptionFunc) float64 {
	return singleFakeData(LONGITUDE, func() interface{} {
		address := Address{}
		return float64(address.longitude())
	}, opts...).(float64)
}

// Latitude get fake latitude randomly
func Latitude(opts ...options.OptionFunc) float64 {
	return singleFakeData(LATITUDE, func() interface{} {
		address := Address{}
		return float64(address.latitude())
	}, opts...).(float64)
}

type RealCoordinates struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type RealAddress struct {
	Address1    string          `json:"address1"`
	Address2    string          `json:"address2"`
	City        string          `json:"city"`
	State       string          `json:"state"`
	PostalCode  string          `json:"postalCode"`
	Coordinates RealCoordinates `json:"coordinates"`
}

// GetRealAddress get real world address randomly
func GetRealAddress(opts ...options.OptionFunc) RealAddress {
	return singleFakeData(RealAddressTag, func() interface{} {
		return addressesUS[rand.Intn(len(addressesUS))]
	}, opts...).(RealAddress)
}
