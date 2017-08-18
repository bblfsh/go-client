package bblfsh

import (
	"errors"
	"reflect"
	"unsafe"
)

// #cgo CFLAGS: -I/usr/local/include -I/usr/local/include/libxml2 -I/usr/include -I/usr/include/libxml2
// #cgo LDFLAGS: -luast -lxml2
// #include "bindings.h"
import "C"

func init() {
	C.create_go_node_api()
}

func Find(node interface{}, xpath string) ([]interface{}, error) {
	cquery := C.CString(xpath)
	defer C.free(unsafe.Pointer(cquery))

	ptr := C.uintptr_t(uintptr(unsafe.Pointer(&node)))
	if C._api_find(ptr, cquery) != 0 {
		return nil, errors.New("error")
	}

	nu := int(C._api_get_nu_results())
	results := make([]interface{}, nu)
	for i := 0; i < nu; i++ {
		results[i] = *(*interface{})(unsafe.Pointer(uintptr(C._api_get_result(C.uint(i)))))
	}
	return results, nil
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

func readString(ptr unsafe.Pointer, attribute string) *C.char {
	return C.CString(readAttribute(ptr, attribute).String())
}

func readLen(ptr unsafe.Pointer, attribute string) C.int {
	return C.int(readAttribute(ptr, attribute).Len())
}

func readIndex(ptr unsafe.Pointer, attribute string, index int) reflect.Value {
	return readAttribute(ptr, attribute).Index(index)
}

//export goGetInternalType
func goGetInternalType(ptr unsafe.Pointer) *C.char {
	return readString(ptr, "InternalType")
}

//export goGetPropertiesSize
func goGetPropertiesSize(ptr unsafe.Pointer) C.int {
	return 0
}

//export goGetToken
func goGetToken(ptr unsafe.Pointer) *C.char {
	return readString(ptr, "Token")
}

//export goGetChildrenSize
func goGetChildrenSize(ptr unsafe.Pointer) C.int {
	return readLen(ptr, "Children")
}

//export goGetChild
func goGetChild(ptr unsafe.Pointer, index C.int) C.uintptr_t {
	a := readIndex(ptr, "Children", int(index)).Elem()
	b := a.Interface()
	return C.uintptr_t(uintptr(unsafe.Pointer(&b)))
}

//export goGetRolesSize
func goGetRolesSize(ptr unsafe.Pointer) C.int {
	return readLen(ptr, "Roles")
	// return 0
}

//export goGetRole
func goGetRole(ptr unsafe.Pointer, index C.int) C.uint16_t {
	// return 0
	return C.uint16_t(readIndex(ptr, "Roles", int(index)).Uint())
}
