package options

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"sync/atomic"
	"time"
	"unsafe"

	fakerErrors "github.com/go-faker/faker/v4/pkg/errors"
	"github.com/go-faker/faker/v4/pkg/interfaces"
)

var (
	// global settings, read/write must be concurrent safe
	generateUniqueValues atomic.Value
	ignoreInterface      atomic.Value
	randomStringLen      int32 = 25
	lang                 unsafe.Pointer
	randomMaxSize        int32 = 100
	randomMinSize        int32
	iBoundary            unsafe.Pointer
)

func init() {
	generateUniqueValues.Store(false)
	ignoreInterface.Store(false)
	lang = unsafe.Pointer(&interfaces.LangENG)
	iBoundary = unsafe.Pointer(&interfaces.DefaultIntBoundary)
}

// Options represent all available option for faker.
type Options struct {
	// IgnoreFields used for ignoring a field when generating the fake data
	IgnoreFields map[string]struct{}
	// FieldProviders used for storing the custom provider function
	FieldProviders map[string]interfaces.CustomProviderFunction
	// StructTypeProviders used for storing the struct type of custom provider function
	StructTypeProviders map[reflect.Type]interfaces.CustomProviderFunction
	// MaxDepthOption used for configuring the max depth of nested identical structs for faker
	MaxDepthOption *MaxDepthOption
	// MaxFieldDepthOption used for configuring the max depth of fields that are filled for a struct
	MaxFieldDepthOption int
	// IgnoreInterface used for ignoring any interface field
	IgnoreInterface bool
	// StringLanguage used for setting the language for any string in faker
	StringLanguage *interfaces.LangRuneBoundary
	// GenerateUniqueValues to ensure the generated data is unique
	GenerateUniqueValues bool
	// RandomStringLength to ensure the generated string is expected as we want
	RandomStringLength int
	// RandomMaxSliceSize used for setting the maximum of slice size, or map size that will be generated
	RandomMaxSliceSize int
	// RandomMinSliceSize used for setting the minimum of slize, array, map size that will be generated
	RandomMinSliceSize int
	// MaxGenerateStringRetries set how much tries for generating random string
	MaxGenerateStringRetries int
	// SetSliceMapNilIfLenZero allows to set nil for the slice and maps, if size is 0.
	SetSliceMapNilIfLenZero bool
	// SetSliceMapRandomToZero sets random integer generation to zero for slice and maps
	SetSliceMapRandomToZero bool
	// RandomIntegerBoundary sets boundary random integer value generation. Boundaries can not exceed integer(4 byte...)
	RandomIntegerBoundary *interfaces.RandomIntegerBoundary
	// RandomFloatBoundary sets the boundary for random float value generation. Boundaries should comply with float values constraints (IEEE 754)
	RandomFloatBoundary *interfaces.RandomFloatBoundary
	// SetTagName sets the tag name that should be used
	TagName string
	// CustomDomain is used for specifying a custom domain when generating email
	CustomDomain *string
	// OnlyZeroFields skips any field that already holds a non-zero value, leaving it unchanged
	OnlyZeroFields bool
	// RandomNestedMaxSliceSize controls the max size for slices/maps that are nested inside
	// another slice or map. When set, this prevents exponential memory growth when generating
	// large outer slices containing structs with nested slice/map fields.
	// If unset (0), RandomMaxSliceSize applies at all depths (original behavior).
	RandomNestedMaxSliceSize int
	// RandomNestedMinSliceSize controls the min size for slices/maps nested inside another
	// slice or map. Pair with RandomNestedMaxSliceSize.
	RandomNestedMinSliceSize int
	// sliceDepth tracks the current nesting depth of slice/map generation.
	// Access via IsNested() and Nested(); do not read or write this field directly.
	sliceDepth int
}

// MaxDepthOption used for configuring the max depth of nested struct for faker
type MaxDepthOption struct {
	typeSeen          map[reflect.Type]int
	recursionMaxDepth int
}

func (o *MaxDepthOption) RememberType(t reflect.Type) {
	o.typeSeen[t]++
}

func (o *MaxDepthOption) ForgetType(t reflect.Type) {
	o.typeSeen[t]--
}

func (o *MaxDepthOption) RecursionOutOfLimit(t reflect.Type) bool {
	return o.typeSeen[t] > o.recursionMaxDepth
}

// IsNested reports whether faker is currently generating elements inside an outer
// slice, array, or map. Used by randomSliceAndMapSize to pick nested sizes.
func (o Options) IsNested() bool {
	return o.sliceDepth > 0
}

// Nested returns a copy of o with the nesting depth incremented by one.
// Call this before recursing into the elements of a slice, array, or map.
//
// The value receiver is intentional: Go copies o on entry, so incrementing
// sliceDepth and returning that copy leaves the caller's Options unchanged.
// A pointer receiver would mutate the original and corrupt depth tracking
// across sibling iterations.
func (o Options) Nested() Options {
	o.sliceDepth++
	return o
}

// BuildOptions build all option functions into one option
func BuildOptions(optFuncs []OptionFunc) *Options {
	ops := DefaultOption()

	for _, optFunc := range optFuncs {
		optFunc(ops)
	}

	return ops
}

// DefaultOption build the default option
func DefaultOption() *Options {
	ops := &Options{}
	ops.StructTypeProviders = make(map[reflect.Type]interfaces.CustomProviderFunction)
	ops.StructTypeProviders[reflect.TypeFor[time.Time]()] = func() (any, error) {
		return time.Now().Add(time.Duration(rand.Int63())), nil
	}
	ops.MaxDepthOption = &MaxDepthOption{
		typeSeen:          make(map[reflect.Type]int, 1),
		recursionMaxDepth: 1,
	}
	ops.MaxFieldDepthOption = -1
	ops.GenerateUniqueValues = generateUniqueValues.Load().(bool)
	ops.IgnoreInterface = ignoreInterface.Load().(bool)
	ops.StringLanguage = (*interfaces.LangRuneBoundary)(atomic.LoadPointer(&lang))
	ops.RandomStringLength = int(atomic.LoadInt32(&randomStringLen))
	ops.RandomMaxSliceSize = int(atomic.LoadInt32(&randomMaxSize))
	ops.RandomMinSliceSize = int(atomic.LoadInt32(&randomMinSize))
	ops.MaxGenerateStringRetries = 1000000 //default
	ops.RandomIntegerBoundary = (*interfaces.RandomIntegerBoundary)(atomic.LoadPointer(&iBoundary))
	ops.RandomFloatBoundary = &interfaces.DefaultFloatBoundary
	ops.TagName = "faker"
	return ops
}

// OptionFunc define the options contract
type OptionFunc func(oo *Options)

// WithFieldsToIgnore used for ignoring a field when generating the fake data
func WithFieldsToIgnore(fieldNames ...string) OptionFunc {
	return func(oo *Options) {
		if oo.IgnoreFields == nil {
			oo.IgnoreFields = make(map[string]struct{}, len(fieldNames))
		}
		for _, f := range fieldNames {
			oo.IgnoreFields[f] = struct{}{}
		}
	}
}

// WithCustomDomain is used to set a custom domain for generating fake email
func WithCustomDomain(domain string) OptionFunc {
	return func(oo *Options) {
		oo.CustomDomain = &domain
	}
}

// WithCustomFieldProvider used for storing the custom provider function
func WithCustomFieldProvider(fieldName string, provider interfaces.CustomProviderFunction) OptionFunc {
	return func(oo *Options) {
		if oo.FieldProviders == nil {
			oo.FieldProviders = make(map[string]interfaces.CustomProviderFunction, 1)
		}
		oo.FieldProviders[fieldName] = provider
	}
}

// WithRecursionMaxDepth used for configuring the max depth of recursion struct for faker
func WithRecursionMaxDepth(depth uint) OptionFunc {
	return func(oo *Options) {
		if oo.MaxDepthOption == nil {
			oo.MaxDepthOption = &MaxDepthOption{
				recursionMaxDepth: 1, // default
				typeSeen:          make(map[reflect.Type]int, 1),
			}
		}
		oo.MaxDepthOption.recursionMaxDepth = int(depth)
	}
}

// WithMaxFieldDepthOption used for configuring the max depth of fields that are filled for a struct
func WithMaxFieldDepthOption(depth int) OptionFunc {
	return func(oo *Options) {
		oo.MaxFieldDepthOption = depth
	}
}

// WithStructTypeProviders used for configuring the custom provider of struct type
func WithStructTypeProviders(t any, provider interfaces.CustomProviderFunction) OptionFunc {
	if reflect.TypeOf(t).Kind() != reflect.Struct {
		panic(fakerErrors.ErrOnlyStructTypeSupported)
	}
	return func(oo *Options) {
		oo.StructTypeProviders[reflect.TypeOf(t)] = provider
	}
}

// WithIgnoreInterface allows to set a flag to ignore found interface{}s.
func WithIgnoreInterface(value bool) OptionFunc {
	return func(oo *Options) {
		oo.IgnoreInterface = value
	}
}

// WithStringLanguage sets language of random string generation (LangENG, LangCHI, LangRUS, LangJPN, LangKOR, EmotEMJ)
func WithStringLanguage(l interfaces.LangRuneBoundary) OptionFunc {
	return func(oo *Options) {
		oo.StringLanguage = &l
	}
}

// WithGenerateUniqueValues allows to set the single fake data generator functions to generate unique data.
func WithGenerateUniqueValues(unique bool) OptionFunc {
	return func(oo *Options) {
		oo.GenerateUniqueValues = unique
	}
}

// WithRandomStringLength sets a length for random string generation
func WithRandomStringLength(size uint) OptionFunc {
	return func(oo *Options) {
		oo.RandomStringLength = int(size)
	}
}

// WithRandomMapAndSliceMaxSize sets the max size for maps and slices for random generation.
func WithRandomMapAndSliceMaxSize(size uint) OptionFunc {
	if size < 1 {
		err := fmt.Errorf(fakerErrors.ErrSmallerThanOne, size)
		panic(err)
	}
	return func(oo *Options) {
		oo.RandomMaxSliceSize = int(size)
	}
}

// WithRandomMapAndSliceMinSize sets the min size for maps and slices for random generation.
func WithRandomMapAndSliceMinSize(size uint) OptionFunc {
	return func(oo *Options) {
		oo.RandomMinSliceSize = int(size)
	}
}

// WithNestedRandomMapAndSliceSize sets the min and max size for slices and maps that are
// generated as fields inside the elements of an outer slice or map. Use this together with
// WithRandomMapAndSliceMaxSize to avoid exponential memory growth when the generated struct
// contains nested slice/map fields.
//
// Example: generate 1000 users where each user's nested slices stay small:
//
//	faker.FakeData(&users,
//	    options.WithRandomMapAndSliceMaxSize(1000),
//	    options.WithNestedRandomMapAndSliceSize(1, 5),
//	)
func WithNestedRandomMapAndSliceSize(minSize, maxSize uint) OptionFunc {
	if maxSize < 1 {
		panic(fmt.Errorf(fakerErrors.ErrSmallerThanOne, maxSize))
	}
	if minSize > maxSize {
		panic(errors.New(fakerErrors.ErrStartValueBiggerThanEnd))
	}
	return func(oo *Options) {
		oo.RandomNestedMinSliceSize = int(minSize)
		oo.RandomNestedMaxSliceSize = int(maxSize)
	}
}

// WithMaxGenerateStringRetries set how much tries for generating random string
func WithMaxGenerateStringRetries(retries uint) OptionFunc {
	return func(oo *Options) {
		oo.MaxGenerateStringRetries = int(retries)
	}
}

// WithNilIfLenIsZero allows to set nil for the slice and maps, if size is 0.
func WithNilIfLenIsZero(setNil bool) OptionFunc {
	return func(oo *Options) {
		oo.SetSliceMapNilIfLenZero = setNil
	}
}

// WithSliceMapRandomToZero Sets random integer generation to zero for slice and maps
func WithSliceMapRandomToZero(setNumberToZero bool) OptionFunc {
	return func(oo *Options) {
		oo.SetSliceMapRandomToZero = setNumberToZero
	}
}

// WithRandomIntegerBoundaries sets boundary random integer value generation. Boundaries can not exceed integer(4 byte...)
func WithRandomIntegerBoundaries(boundary interfaces.RandomIntegerBoundary) OptionFunc {
	if boundary.Start > boundary.End {
		err := errors.New(fakerErrors.ErrStartValueBiggerThanEnd)
		panic(err)
	}
	return func(oo *Options) {
		oo.RandomIntegerBoundary = &boundary
	}
}

// WithRandomFloatBoundaries sets the boundary for random float value generation. Boundaries should comply with float values constraints (IEEE 754)
func WithRandomFloatBoundaries(boundary interfaces.RandomFloatBoundary) OptionFunc {
	if boundary.Start > boundary.End {
		err := errors.New(fakerErrors.ErrStartValueBiggerThanEnd)
		panic(err)
	}
	return func(oo *Options) {
		oo.RandomFloatBoundary = &boundary
	}
}

// WithOnlyZeroFields makes faker skip any field that already holds a non-zero value.
// Useful for partially pre-initializing a struct before passing it to faker.
func WithOnlyZeroFields() OptionFunc {
	return func(oo *Options) {
		oo.OnlyZeroFields = true
	}
}

// WithTagName sets the tag name to use. Default tag name is 'faker'.
func WithTagName(tagName string) OptionFunc {
	if tagName == "" {
		err := errors.New(fakerErrors.ErrFieldTagIdentifierInvalid)
		panic(err)
	}
	return func(oo *Options) {
		oo.TagName = tagName
	}
}

// SetGenerateUniqueValues allows to set the single fake data generator functions to generate unique data.
func SetGenerateUniqueValues(unique bool) {
	generateUniqueValues.Store(unique)
}

// SetIgnoreInterface allows to set a flag to ignore found interface{}s.
func SetIgnoreInterface(ignore bool) {
	ignoreInterface.Store(ignore)
}

// SetRandomStringLength sets a length for random string generation
func SetRandomStringLength(size int) error {
	if size < 0 {
		return fmt.Errorf(fakerErrors.ErrSmallerThanZero, size)
	}
	atomic.StoreInt32(&randomStringLen, int32(size))
	return nil
}

// SetStringLang sets language of random string generation (LangENG, LangCHI, LangRUS, LangJPN, LangKOR, EmotEMJ)
func SetStringLang(l interfaces.LangRuneBoundary) {
	atomic.StorePointer(&lang, unsafe.Pointer(&l))
}

// SetRandomMapAndSliceSize sets the size for maps and slices for random generation.
// deprecates, currently left for old version usage
func SetRandomMapAndSliceSize(size int) error {
	return SetRandomMapAndSliceMaxSize(size)
}

// SetRandomMapAndSliceMaxSize sets the max size for maps and slices for random generation.
func SetRandomMapAndSliceMaxSize(size int) error {
	if size < 1 {
		return fmt.Errorf(fakerErrors.ErrSmallerThanOne, size)
	}
	atomic.StoreInt32(&randomMaxSize, int32(size))
	return nil
}

// SetRandomMapAndSliceMinSize sets the min size for maps and slices for random generation.
func SetRandomMapAndSliceMinSize(size int) error {
	if size < 0 {
		return fmt.Errorf(fakerErrors.ErrSmallerThanZero, size)
	}
	atomic.StoreInt32(&randomMinSize, int32(size))
	return nil
}

// SetRandomNumberBoundaries sets boundary for random number generation
func SetRandomNumberBoundaries(start, end int) error {
	if start > end {
		return errors.New(fakerErrors.ErrStartValueBiggerThanEnd)
	}
	ptr := &interfaces.RandomIntegerBoundary{Start: start, End: end}
	atomic.StorePointer(&iBoundary, unsafe.Pointer(ptr))
	return nil
}
