package bblfsh

import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/bblfsh/sdk/uast"
)

// #cgo CFLAGS: -I/usr/local/include -I/usr/local/include/libxml2 -I/usr/include -I/usr/include/libxml2
// #cgo LDFLAGS: -luast -lxml2
// #include "bindings.h"
import "C"

func init() {
	C.create_go_node_api()
}

func Find(node *uast.Node, xpath string) ([]*uast.Node, error) {
	cquery := C.CString(xpath)
	defer C.free(unsafe.Pointer(cquery))

	ptr := C.uintptr_t(uintptr(unsafe.Pointer(node)))
	if C._api_find(ptr, cquery) != 0 {
		return nil, errors.New("error")
	}

	nu := int(C._api_get_nu_results())
	results := make([]*uast.Node, nu)
	for i := 0; i < nu; i++ {
		results[i] = (*uast.Node)(unsafe.Pointer(uintptr(C._api_get_result(C.uint(i)))))
	}
	return results, nil
}

func readNode(ptr unsafe.Pointer) *uast.Node {
	return (*uast.Node)(ptr)
}

func readAttribute(ptr unsafe.Pointer, attribute string) reflect.Value {
	obj := *((*interface{})(ptr))
	value := reflect.ValueOf(obj)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	value = value.FieldByName(attribute)
	return value
}

//export goGetInternalType
func goGetInternalType(ptr unsafe.Pointer) *C.char {
	return C.CString(readNode(ptr).InternalType)
}

//export goGetPropertiesSize
func goGetPropertiesSize(ptr unsafe.Pointer) C.int {
	return 0
}

//export goGetToken
func goGetToken(ptr unsafe.Pointer) *C.char {
	return C.CString(readNode(ptr).Token)
}

//export goGetChildrenSize
func goGetChildrenSize(ptr unsafe.Pointer) C.int {
	return C.int(len(readNode(ptr).Children))
}

//export goGetChild
func goGetChild(ptr unsafe.Pointer, index C.int) C.uintptr_t {
	child := readNode(ptr).Children[int(index)]
	return C.uintptr_t(uintptr(unsafe.Pointer(child)))
}

//export goGetRolesSize
func goGetRolesSize(ptr unsafe.Pointer) C.int {
	return C.int(len(readNode(ptr).Roles))
}

//export goGetRole
func goGetRole(ptr unsafe.Pointer, index C.int) C.uint16_t {
	role := readNode(ptr).Roles[int(index)]
	return C.uint16_t(role)
}
