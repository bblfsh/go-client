package bblfsh

import (
	"errors"
	"runtime/debug"
	"sync"
	"unsafe"

	"github.com/bblfsh/sdk/uast"
)

// #cgo CFLAGS: -I/usr/local/include -I/usr/local/include/libxml2 -I/usr/include -I/usr/include/libxml2
// #cgo LDFLAGS: -luast -lxml2
// #include "bindings.h"
import "C"

var findMutex sync.Mutex
var pool cstringPool

func init() {
	C.create_go_node_api()
}

func nodeToPtr(node *uast.Node) C.uintptr_t {
	return C.uintptr_t(uintptr(unsafe.Pointer(node)))
}

func ptrToNode(ptr C.uintptr_t) *uast.Node {
	return (*uast.Node)(unsafe.Pointer(uintptr(ptr)))
}

// Find takes a `*uast.Node` and a xpath query and filters the tree,
// returning the list of nodes that satisfy the given query.
// Find is thread-safe but not current by an internal global lock.
func Find(node *uast.Node, xpath string) ([]*uast.Node, error) {
	// Find is not thread-safe bacuase of the underlining C API
	findMutex.Lock()
	defer findMutex.Unlock()

	// convert xpath string to a NULL-terminated c string
	cquery := pool.getCPtr(xpath)

	// Make sure we release the pool of strings
	defer pool.release()

	// stop GC
	gcpercent := debug.SetGCPercent(-1)
	defer debug.SetGCPercent(gcpercent)

	ptr := nodeToPtr(node)
	if C._api_find(ptr, cquery) != 0 {
		return nil, errors.New("error: node_api_find() failed")
	}

	nu := int(C._api_get_nu_results())
	results := make([]*uast.Node, nu)
	for i := 0; i < nu; i++ {
		results[i] = ptrToNode(C._api_get_result(C.uint(i)))
	}
	return results, nil
}

//export goGetInternalType
func goGetInternalType(ptr C.uintptr_t) *C.char {
	return pool.getCPtr(ptrToNode(ptr).InternalType)
}

//export goGetToken
func goGetToken(ptr C.uintptr_t) *C.char {
	return pool.getCPtr(ptrToNode(ptr).Token)
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
