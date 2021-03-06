package mockery

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const pkg = "test"

func TestGenerator(t *testing.T) {
	parser := NewParser()
	err := parser.Parse(testFile)
	require.NoError(t, err)

	iface, err := parser.Find("Requester")

	gen := NewGenerator(iface, pkg)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `// Requester is an autogenerated mock type for the Requester type
type Requester struct {
	mock.Mock
}

// Get provides a mock function with given fields: path
func (_m *Requester) Get(path string) (string, error) {
	ret := _m.Called(path)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorSingleReturn(t *testing.T) {
	parser := NewParser()
	parser.Parse(testFile2)

	iface, err := parser.Find("Requester2")

	gen := NewGenerator(iface, pkg)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `// Requester2 is an autogenerated mock type for the Requester2 type
type Requester2 struct {
	mock.Mock
}

// Get provides a mock function with given fields: path
func (_m *Requester2) Get(path string) error {
	ret := _m.Called(path)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorNoArguments(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "requester3.go"))

	iface, err := parser.Find("Requester3")

	gen := NewGenerator(iface, pkg)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `// Requester3 is an autogenerated mock type for the Requester3 type
type Requester3 struct {
	mock.Mock
}

// Get provides a mock function with given fields: 
func (_m *Requester3) Get() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorNoNothing(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "requester4.go"))

	iface, err := parser.Find("Requester4")

	gen := NewGenerator(iface, pkg)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `// Requester4 is an autogenerated mock type for the Requester4 type
type Requester4 struct {
	mock.Mock
}

// Get provides a mock function with given fields: 
func (_m *Requester4) Get() {
	_m.Called()
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorUnexported(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "requester_unexported.go"))

	iface, err := parser.Find("requester")

	gen := NewGenerator(iface, pkg)
	gen.ip = true

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `// mockRequester is an autogenerated mock type for the requester type
type mockRequester struct {
	mock.Mock
}

// Get provides a mock function with given fields: 
func (_m *mockRequester) Get() {
	_m.Called()
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorPrologue(t *testing.T) {
	parser := NewParser()
	parser.Parse(testFile)

	iface, err := parser.Find("Requester")
	assert.NoError(t, err)

	gen := NewGenerator(iface, pkg)

	gen.GeneratePrologue("mocks")

	goPath := os.Getenv("GOPATH")
	local, err := filepath.Rel(filepath.Join(goPath, "src"), filepath.Dir(iface.Path))
	assert.NoError(t, err)

	expected := `package mocks

import "` + local + `"
import "github.com/stretchr/testify/mock"

`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorProloguewithImports(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "requester_ns.go"))

	iface, err := parser.Find("RequesterNS")
	assert.NoError(t, err)

	gen := NewGenerator(iface, pkg)

	gen.GeneratePrologue("mocks")

	goPath := os.Getenv("GOPATH")
	local, err := filepath.Rel(filepath.Join(goPath, "src"), filepath.Dir(iface.Path))
	assert.NoError(t, err)

	expected := `package mocks

import "` + local + `"
import "github.com/stretchr/testify/mock"

import "net/http"

`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorPrologueNote(t *testing.T) {
	parser := NewParser()
	parser.Parse(testFile)

	iface, err := parser.Find("Requester")
	assert.NoError(t, err)

	gen := NewGenerator(iface, pkg)

	gen.GeneratePrologueNote("A\\nB")

	expected := `
// A
// B

`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorChecksInterfacesForNilable(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "requester_iface.go"))

	iface, err := parser.Find("RequesterIface")
	assert.NoError(t, err)

	gen := NewGenerator(iface, pkg)
	assert.NoError(t, err)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `// RequesterIface is an autogenerated mock type for the RequesterIface type
type RequesterIface struct {
	mock.Mock
}

// Get provides a mock function with given fields: 
func (_m *RequesterIface) Get() io.Reader {
	ret := _m.Called()

	var r0 io.Reader
	if rf, ok := ret.Get(0).(func() io.Reader); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.Reader)
		}
	}

	return r0
}
`
	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorPointers(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "requester_ptr.go"))

	iface, err := parser.Find("RequesterPtr")
	assert.NoError(t, err)

	gen := NewGenerator(iface, pkg)
	assert.NoError(t, err)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `// RequesterPtr is an autogenerated mock type for the RequesterPtr type
type RequesterPtr struct {
	mock.Mock
}

// Get provides a mock function with given fields: path
func (_m *RequesterPtr) Get(path string) (*string, error) {
	ret := _m.Called(path)

	var r0 *string
	if rf, ok := ret.Get(0).(func(string) *string); ok {
		r0 = rf(path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
`
	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorSlice(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "requester_slice.go"))

	iface, err := parser.Find("RequesterSlice")
	assert.NoError(t, err)

	gen := NewGenerator(iface, pkg)
	assert.NoError(t, err)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `// RequesterSlice is an autogenerated mock type for the RequesterSlice type
type RequesterSlice struct {
	mock.Mock
}

// Get provides a mock function with given fields: path
func (_m *RequesterSlice) Get(path string) ([]string, error) {
	ret := _m.Called(path)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
`
	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorArrayLiteralLen(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "requester_array.go"))

	iface, err := parser.Find("RequesterArray")
	assert.NoError(t, err)

	gen := NewGenerator(iface, pkg)
	assert.NoError(t, err)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `// RequesterArray is an autogenerated mock type for the RequesterArray type
type RequesterArray struct {
	mock.Mock
}

// Get provides a mock function with given fields: path
func (_m *RequesterArray) Get(path string) ([2]string, error) {
	ret := _m.Called(path)

	var r0 [2]string
	if rf, ok := ret.Get(0).(func(string) [2]string); ok {
		r0 = rf(path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([2]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorNamespacedTypes(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "requester_ns.go"))

	iface, err := parser.Find("RequesterNS")
	assert.NoError(t, err)

	gen := NewGenerator(iface, pkg)
	assert.NoError(t, err)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `// RequesterNS is an autogenerated mock type for the RequesterNS type
type RequesterNS struct {
	mock.Mock
}

// Get provides a mock function with given fields: path
func (_m *RequesterNS) Get(path string) (http.Response, error) {
	ret := _m.Called(path)

	var r0 http.Response
	if rf, ok := ret.Get(0).(func(string) http.Response); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Get(0).(http.Response)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorWhereArgumentNameConflictsWithImport(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "requester_arg_same_as_import.go"))

	iface, err := parser.Find("RequesterArgSameAsImport")
	assert.NoError(t, err)

	gen := NewGenerator(iface, pkg)
	assert.NoError(t, err)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `// RequesterArgSameAsImport is an autogenerated mock type for the RequesterArgSameAsImport type
type RequesterArgSameAsImport struct {
	mock.Mock
}

// Get provides a mock function with given fields: _a0
func (_m *RequesterArgSameAsImport) Get(_a0 string) *json.RawMessage {
	ret := _m.Called(_a0)

	var r0 *json.RawMessage
	if rf, ok := ret.Get(0).(func(string) *json.RawMessage); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*json.RawMessage)
		}
	}

	return r0
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorWhereArgumentNameConflictsWithNamedImport(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "requester_arg_same_as_named_import.go"))

	iface, err := parser.Find("RequesterArgSameAsNamedImport")
	assert.NoError(t, err)

	gen := NewGenerator(iface, pkg)
	assert.NoError(t, err)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `// RequesterArgSameAsNamedImport is an autogenerated mock type for the RequesterArgSameAsNamedImport type
type RequesterArgSameAsNamedImport struct {
	mock.Mock
}

// Get provides a mock function with given fields: _a0
func (_m *RequesterArgSameAsNamedImport) Get(_a0 string) *json.RawMessage {
	ret := _m.Called(_a0)

	var r0 *json.RawMessage
	if rf, ok := ret.Get(0).(func(string) *json.RawMessage); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*json.RawMessage)
		}
	}

	return r0
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorWhereArgumentNameConflictsWithPackage(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "requester_arg_same_as_pkg.go"))

	iface, err := parser.Find("RequesterArgSameAsPkg")
	assert.NoError(t, err)

	gen := NewGenerator(iface, pkg)
	assert.NoError(t, err)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `// RequesterArgSameAsPkg is an autogenerated mock type for the RequesterArgSameAsPkg type
type RequesterArgSameAsPkg struct {
	mock.Mock
}

// Get provides a mock function with given fields: _a0
func (_m *RequesterArgSameAsPkg) Get(_a0 string) {
	_m.Called(_a0)
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorHavingNoNamesOnArguments(t *testing.T) {
	parser := NewParser()

	parser.Parse(filepath.Join(fixturePath, "custom_error.go"))

	iface, err := parser.Find("KeyManager")
	assert.NoError(t, err)

	gen := NewGenerator(iface, pkg)
	assert.NoError(t, err)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `// KeyManager is an autogenerated mock type for the KeyManager type
type KeyManager struct {
	mock.Mock
}

// GetKey provides a mock function with given fields: _a0, _a1
func (_m *KeyManager) GetKey(_a0 string, _a1 uint16) ([]byte, *Err) {
	ret := _m.Called(_a0, _a1)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string, uint16) []byte); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 *Err
	if rf, ok := ret.Get(1).(func(string, uint16) *Err); ok {
		r1 = rf(_a0, _a1)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*Err)
		}
	}

	return r0, r1
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorElidedType(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "requester_elided.go"))

	iface, err := parser.Find("RequesterElided")

	gen := NewGenerator(iface, pkg)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `// RequesterElided is an autogenerated mock type for the RequesterElided type
type RequesterElided struct {
	mock.Mock
}

// Get provides a mock function with given fields: path, url
func (_m *RequesterElided) Get(path string, url string) error {
	ret := _m.Called(path, url)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(path, url)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorReturnElidedType(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "requester_ret_elided.go"))

	iface, err := parser.Find("RequesterReturnElided")

	gen := NewGenerator(iface, pkg)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `// RequesterReturnElided is an autogenerated mock type for the RequesterReturnElided type
type RequesterReturnElided struct {
	mock.Mock
}

// Get provides a mock function with given fields: path
func (_m *RequesterReturnElided) Get(path string) (int, int, int, error) {
	ret := _m.Called(path)

	var r0 int
	if rf, ok := ret.Get(0).(func(string) int); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 int
	if rf, ok := ret.Get(1).(func(string) int); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Get(1).(int)
	}

	var r2 int
	if rf, ok := ret.Get(2).(func(string) int); ok {
		r2 = rf(path)
	} else {
		r2 = ret.Get(2).(int)
	}

	var r3 error
	if rf, ok := ret.Get(3).(func(string) error); ok {
		r3 = rf(path)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorVariableArgs(t *testing.T) {

	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "requester_variable.go"))

	iface, err := parser.Find("RequesterVariable")

	gen := NewGenerator(iface, pkg)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `// RequesterVariable is an autogenerated mock type for the RequesterVariable type
type RequesterVariable struct {
	mock.Mock
}

// Get provides a mock function with given fields: values
func (_m *RequesterVariable) Get(values ...string) bool {
	ret := _m.Called(values)

	var r0 bool
	if rf, ok := ret.Get(0).(func(...string) bool); ok {
		r0 = rf(values...)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorFuncType(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "func_type.go"))

	iface, err := parser.Find("Fooer")

	gen := NewGenerator(iface, pkg)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `// Fooer is an autogenerated mock type for the Fooer type
type Fooer struct {
	mock.Mock
}

// Bar provides a mock function with given fields: f
func (_m *Fooer) Bar(f func([]int)) {
	_m.Called(f)
}
// Baz provides a mock function with given fields: path
func (_m *Fooer) Baz(path string) func(string) string {
	ret := _m.Called(path)

	var r0 func(string) string
	if rf, ok := ret.Get(0).(func(string) func(string) string); ok {
		r0 = rf(path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(func(string) string)
		}
	}

	return r0
}
// Foo provides a mock function with given fields: f
func (_m *Fooer) Foo(f func(string) string) error {
	ret := _m.Called(f)

	var r0 error
	if rf, ok := ret.Get(0).(func(func(string) string) error); ok {
		r0 = rf(f)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorChanType(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "async.go"))

	iface, err := parser.Find("AsyncProducer")

	gen := NewGenerator(iface, pkg)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `// AsyncProducer is an autogenerated mock type for the AsyncProducer type
type AsyncProducer struct {
	mock.Mock
}

// Input provides a mock function with given fields: 
func (_m *AsyncProducer) Input() chan<- bool {
	ret := _m.Called()

	var r0 chan<- bool
	if rf, ok := ret.Get(0).(func() chan<- bool); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan<- bool)
		}
	}

	return r0
}
// Output provides a mock function with given fields: 
func (_m *AsyncProducer) Output() <-chan bool {
	ret := _m.Called()

	var r0 <-chan bool
	if rf, ok := ret.Get(0).(func() <-chan bool); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan bool)
		}
	}

	return r0
}
// Whatever provides a mock function with given fields: 
func (_m *AsyncProducer) Whatever() chan bool {
	ret := _m.Called()

	var r0 chan bool
	if rf, ok := ret.Get(0).(func() chan bool); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan bool)
		}
	}

	return r0
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorFromImport(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "io_import.go"))

	iface, err := parser.Find("MyReader")
	require.NoError(t, err)

	gen := NewGenerator(iface, pkg)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `// MyReader is an autogenerated mock type for the MyReader type
type MyReader struct {
	mock.Mock
}

// Read provides a mock function with given fields: p
func (_m *MyReader) Read(p []byte) (int, error) {
	ret := _m.Called(p)

	var r0 int
	if rf, ok := ret.Get(0).(func([]byte) int); ok {
		r0 = rf(p)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorComplexChanFromConsul(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "consul.go"))

	iface, err := parser.Find("ConsulLock")
	require.NoError(t, err)

	gen := NewGenerator(iface, pkg)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `// ConsulLock is an autogenerated mock type for the ConsulLock type
type ConsulLock struct {
	mock.Mock
}

// Lock provides a mock function with given fields: _a0
func (_m *ConsulLock) Lock(_a0 <-chan struct{}) (<-chan struct{}, error) {
	ret := _m.Called(_a0)

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func(<-chan struct{}) <-chan struct{}); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(<-chan struct{}) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
// Unlock provides a mock function with given fields: 
func (_m *ConsulLock) Unlock() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorForEmptyInterface(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "empty_interface.go"))

	iface, err := parser.Find("Blank")

	gen := NewGenerator(iface, pkg)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `// Blank is an autogenerated mock type for the Blank type
type Blank struct {
	mock.Mock
}

// Create provides a mock function with given fields: x
func (_m *Blank) Create(x interface{}) error {
	ret := _m.Called(x)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(x)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorForMapFunc(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "map_func.go"))

	iface, err := parser.Find("MapFunc")

	gen := NewGenerator(iface, pkg)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `// MapFunc is an autogenerated mock type for the MapFunc type
type MapFunc struct {
	mock.Mock
}

// Get provides a mock function with given fields: m
func (_m *MapFunc) Get(m map[string]func(string) string) error {
	ret := _m.Called(m)

	var r0 error
	if rf, ok := ret.Get(0).(func(map[string]func(string) string) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorForMethodUsingInterface(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "mock_method_uses_pkg_iface.go"))

	iface, err := parser.Find("UsesOtherPkgIface")

	gen := NewGenerator(iface, pkg)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `// UsesOtherPkgIface is an autogenerated mock type for the UsesOtherPkgIface type
type UsesOtherPkgIface struct {
	mock.Mock
}

// DoSomethingElse provides a mock function with given fields: obj
func (_m *UsesOtherPkgIface) DoSomethingElse(obj Sibling) {
	_m.Called(obj)
}
`

	assert.Equal(t, expected, gen.buf.String())
}
