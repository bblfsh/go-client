package tools

import (
	"errors"
	"runtime/debug"
	"sort"
	"sync"
	"unsafe"

	"gopkg.in/bblfsh/sdk.v1/uast"
)

// #cgo CFLAGS: -I/usr/local/include -I/usr/local/include/libxml2 -I/usr/include -I/usr/include/libxml2 -std=c99
// #cgo LDFLAGS: -lxml2
// #include "bindings.h"
import "C"

var findMutex sync.Mutex
var pool cstringPool

func init() {
	C.CreateUast()
}

func nodeToPtr(node *uast.Node) C.uintptr_t {
	return C.uintptr_t(uintptr(unsafe.Pointer(node)))
}

func ptrToNode(ptr C.uintptr_t) *uast.Node {
	return (*uast.Node)(unsafe.Pointer(uintptr(ptr)))
}

// Filter takes a `*uast.Node` and a xpath query and filters the tree,
// returning the list of nodes that satisfy the given query.
// Filter is thread-safe but not current by an internal global lock.
func Filter(node *uast.Node, xpath string) ([]*uast.Node, error) {
	if len(xpath) == 0 {
		return nil, nil
	}

	// Find is not thread-safe bacause of the underlining C API
	findMutex.Lock()
	defer findMutex.Unlock()

	// convert go string to C string
	cquery := pool.getCstring(xpath)

	// Make sure we release the pool of strings
	defer pool.release()

	// stop GC
	gcpercent := debug.SetGCPercent(-1)
	defer debug.SetGCPercent(gcpercent)

	ptr := nodeToPtr(node)
	if !C.Filter(ptr, cquery) {
		return nil, errors.New("error: UastFilter() failed")
	}

	nu := int(C.Size())
	results := make([]*uast.Node, nu)
	for i := 0; i < nu; i++ {
		results[i] = ptrToNode(C.At(C.int(i)))
	}
	return results, nil
}

//export goGetInternalType
func goGetInternalType(ptr C.uintptr_t) *C.char {
	return pool.getCstring(ptrToNode(ptr).InternalType)
}

//export goGetToken
func goGetToken(ptr C.uintptr_t) *C.char {
	return pool.getCstring(ptrToNode(ptr).Token)
}

//export goGetChildrenSize
func goGetChildrenSize(ptr C.uintptr_t) C.int {
	return C.int(len(ptrToNode(ptr).Children))
}

//export goGetChild
func goGetChild(ptr C.uintptr_t, index C.int) C.uintptr_t {
	child := ptrToNode(ptr).Children[int(index)]
	return nodeToPtr(child)
}

//export goGetRolesSize
func goGetRolesSize(ptr C.uintptr_t) C.int {
	return C.int(len(ptrToNode(ptr).Roles))
}

//export goGetRole
func goGetRole(ptr C.uintptr_t, index C.int) C.uint16_t {
	role := ptrToNode(ptr).Roles[int(index)]
	return C.uint16_t(role)
}

//export goGetPropertiesSize
func goGetPropertiesSize(ptr C.uintptr_t) C.int {
	return C.int(len(ptrToNode(ptr).Properties))
}

//export goGetPropertyKey
func goGetPropertyKey(ptr C.uintptr_t, index C.int) *C.char {
	var keys []string
	for k := range ptrToNode(ptr).Properties {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return pool.getCstring(keys[int(index)])
}

//export goGetPropertyValue
func goGetPropertyValue(ptr C.uintptr_t, index C.int) *C.char {
	p := ptrToNode(ptr).Properties
	var keys []string
	for k := range p {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return pool.getCstring(p[keys[int(index)]])
}

//export goHasStartOffset
func goHasStartOffset(ptr C.uintptr_t) C.bool {
	return ptrToNode(ptr).StartPosition != nil
}

//export goGetStartOffset
func goGetStartOffset(ptr C.uintptr_t) C.uint32_t {
	p := ptrToNode(ptr).StartPosition
	if p != nil {
		return C.uint32_t(p.Offset)
	}
	return 0
}

//export goHasStartLine
func goHasStartLine(ptr C.uintptr_t) C.bool {
	return ptrToNode(ptr).StartPosition != nil
}

//export goGetStartLine
func goGetStartLine(ptr C.uintptr_t) C.uint32_t {
	p := ptrToNode(ptr).StartPosition
	if p != nil {
		return C.uint32_t(p.Line)
	}
	return 0
}

//export goHasStartCol
func goHasStartCol(ptr C.uintptr_t) C.bool {
	return ptrToNode(ptr).StartPosition != nil
}

//export goGetStartCol
func goGetStartCol(ptr C.uintptr_t) C.uint32_t {
	p := ptrToNode(ptr).StartPosition
	if p != nil {
		return C.uint32_t(p.Col)
	}
	return 0
}

//export goHasEndOffset
func goHasEndOffset(ptr C.uintptr_t) C.bool {
	return ptrToNode(ptr).EndPosition != nil
}

//export goGetEndOffset
func goGetEndOffset(ptr C.uintptr_t) C.uint32_t {
	p := ptrToNode(ptr).EndPosition
	if p != nil {
		return C.uint32_t(p.Offset)
	}
	return 0
}

//export goHasEndLine
func goHasEndLine(ptr C.uintptr_t) C.bool {
	return ptrToNode(ptr).EndPosition != nil
}

//export goGetEndLine
func goGetEndLine(ptr C.uintptr_t) C.uint32_t {
	p := ptrToNode(ptr).EndPosition
	if p != nil {
		return C.uint32_t(p.Line)
	}
	return 0
}

//export goHasEndCol
func goHasEndCol(ptr C.uintptr_t) C.bool {
	return ptrToNode(ptr).EndPosition != nil
}

//export goGetEndCol
func goGetEndCol(ptr C.uintptr_t) C.uint32_t {
	p := ptrToNode(ptr).EndPosition
	if p != nil {
		return C.uint32_t(p.Col)
	}
	return 0
}
