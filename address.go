package faker

import (
	_ "embed"
	"encoding/json"
	"reflect"

	"github.com/cj/faker/pkg/options"
)

var (
	//go:embed misc/addresses-us-1000.min.json
	addressesUSBytes []byte
	addressesUS      []RealAddress

	//go:embed misc/country_info.json
	countriesBytes []byte
	countries      []CountryInfo
)

func init() {
	data := struct {
		Addresses []RealAddress `json:"addresses"`
		Countries []CountryInfo `json:"countries"`
	}{}
	if err := json.Unmarshal(addressesUSBytes, &data); err != nil {
		panic(err)
	}
	addressesUS = data.Addresses

	if err := json.Unmarshal(countriesBytes, &data); err != nil {
		panic(err)
	}
	countries = data.Countries
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
	RealWorld(v reflect.Value) (interface{}, error)
	CountryInfo(v reflect.Value) (interface{}, error)
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

func (i Address) realWorld() RealAddress {
	return addressesUS[rand.Intn(len(addressesUS))]
}

func (i Address) countryInfo() CountryInfo {
	return countries[rand.Intn(len(countries))]
}

func (i Address) CountryInfo(_ reflect.Value) (interface{}, error) {
	return i.countryInfo(), nil
}

// RealWorld sets real world address
func (i Address) RealWorld(_ reflect.Value) (interface{}, error) {
	return i.realWorld(), nil
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
	Address     string          `json:"address1"`
	City        string          `json:"city"`
	State       string          `json:"state"`
	PostalCode  string          `json:"postalCode"`
	Coordinates RealCoordinates `json:"coordinates"`
}

// GetRealAddress get real world address randomly
func GetRealAddress(opts ...options.OptionFunc) RealAddress {
	return singleFakeData(RealAddressTag, func() interface{} {
		address := Address{}
		return address.realWorld()
	}, opts...).(RealAddress)
}

type CountryInfo struct {
	Abbr       string `json:"abbr"`
	Name       string `json:"name"`
	Capital    string `json:"capital"`
	Population string `json:"population"`
	Continent  string `json:"continent"`
}

func GetCountryInfo(opts ...options.OptionFunc) CountryInfo {
	return singleFakeData(CountryInfoTag, func() interface{} {
		address := Address{}
		return address.countryInfo()
	}, opts...).(CountryInfo)
}
