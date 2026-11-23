package faker

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/go-faker/faker/v4/pkg/options"
)

// chineseMobilePrefixes holds the 3-digit prefixes of Mainland China mobile
// numbers, allocated to China Mobile, China Unicom, China Telecom and China
// Broadnet (including MVNO & IoT segments). Each prefix is followed by 8
// subscriber digits to form an 11-digit number.
var chineseMobilePrefixes = []string{
	"130", "131", "132", "133", "134", "135", "136", "137", "138", "139",
	"145", "146", "147", "148", "149",
	"150", "151", "152", "153", "155", "156", "157", "158", "159",
	"162", "165", "166", "167",
	"170", "171", "172", "173", "174", "175", "176", "177", "178",
	"180", "181", "182", "183", "184", "185", "186", "187", "188", "189",
	"190", "191", "192", "193", "195", "196", "197", "198", "199",
}

// GetPhoner serves as a constructor for Phoner interface
func GetPhoner() Phoner {
	phone := &Phone{}
	return phone
}

// Phoner serves overall tele-phonic contact generator
type Phoner interface {
	PhoneNumber(v reflect.Value) (any, error)
	TollFreePhoneNumber(v reflect.Value) (any, error)
	E164PhoneNumber(v reflect.Value) (any, error)
	ChinesePhoneNumber(v reflect.Value) (any, error)
}

// Phone struct
type Phone struct {
}

// nanpCode returns a three digit code of the North American Numbering Plan, for an area code or an exchange:
// it starts with 2 to 9 and it is not N11, the form of the service codes like 411 and 911
func nanpCode() string {
	for {
		code := strconv.Itoa(2+rand.Intn(8)) + randomStringNumber(2)
		if code[1:] != "11" {
			return code
		}
	}
}

func (p Phone) phonenumber() string {
	return nanpCode() + "-" + nanpCode() + "-" + randomStringNumber(4)
}

// PhoneNumber generates phone numbers of type: "201-886-0269"
func (p Phone) PhoneNumber(v reflect.Value) (any, error) {
	return p.phonenumber(), nil
}

// Phonenumber get fake phone number
func Phonenumber(opts ...options.OptionFunc) string {
	return singleFakeData(PhoneNumber, func() any {
		p := Phone{}
		return p.phonenumber()
	}, opts...).(string)
}

func (p Phone) tollfreephonenumber() string {
	boxDigitsStart := []string{"777", "888"}

	return fmt.Sprintf("(%s) %s-%s", boxDigitsStart[rand.Intn(len(boxDigitsStart))], nanpCode(), randomStringNumber(4))
}

// TollFreePhoneNumber generates phone numbers of type: "(888) 937-7238"
func (p Phone) TollFreePhoneNumber(v reflect.Value) (any, error) {
	return p.tollfreephonenumber(), nil
}

// TollFreePhoneNumber get fake TollFreePhoneNumber
func TollFreePhoneNumber(opts ...options.OptionFunc) string {
	return singleFakeData(TollFreeNumber, func() any {
		p := Phone{}
		return p.tollfreephonenumber()
	}, opts...).(string)
}

func (p Phone) e164PhoneNumber() string {
	boxDigitsStart := []string{"7", "8"}

	return "+" + boxDigitsStart[rand.Intn(len(boxDigitsStart))] + randomStringNumber(10)
}

// E164PhoneNumber generates phone numbers of type: "+27113456789"
func (p Phone) E164PhoneNumber(v reflect.Value) (any, error) {
	return p.e164PhoneNumber(), nil
}

// E164PhoneNumber get fake E164PhoneNumber
func E164PhoneNumber(opts ...options.OptionFunc) string {
	return singleFakeData(E164PhoneNumberTag, func() any {
		p := Phone{}
		return p.e164PhoneNumber()
	}, opts...).(string)
}

func (p Phone) chinesePhoneNumber() string {
	prefix := chineseMobilePrefixes[rand.Intn(len(chineseMobilePrefixes))]
	return prefix + randomStringNumber(8)
}

// ChinesePhoneNumber generates a Mainland-China mobile phone number of the
// form "13812345678" (11 digits, starting with a valid operator prefix)
func (p Phone) ChinesePhoneNumber(v reflect.Value) (any, error) {
	return p.chinesePhoneNumber(), nil
}

// ChinesePhoneNumber get fake Chinese phone number
func ChinesePhoneNumber(opts ...options.OptionFunc) string {
	return singleFakeData(ChinesePhoneNumberTag, func() any {
		p := Phone{}
		return p.chinesePhoneNumber()
	}, opts...).(string)
}
