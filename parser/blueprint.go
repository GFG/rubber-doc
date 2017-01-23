package parser

/*
#cgo CFLAGS: -I"${SRCDIR}/../ext/drafter/src/" -I"${SRCDIR}/../ext/drafter/ext/snowcrash/src/"
#cgo darwin LDFLAGS: -L"${SRCDIR}/../ext/drafter/build/out/Release/" -ldrafter -lsos -lsnowcrash -lmarkdownparser -lsundown  -lc++
#cgo linux LDFLAGS: -L"${SRCDIR}/../ext/drafter/build/out/Release/" -ldrafter -lsos -lsnowcrash -lmarkdownparser -lsundown  -lstdc++
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include "drafter.h"
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"unsafe"
)

type Blueprint struct{}

func (e Blueprint) Parse(r io.Reader) (map[string]interface{}, error) {
	b, err := ioutil.ReadAll(r)
	var objmap map[string]interface{}

	if err != nil {
		return nil, err
	}

	source := C.CString(string(b))
	defer C.free(unsafe.Pointer(source))

	result := &C.drafter_result{}
	opts := C.drafter_parse_options{}

	errCode := int(C.drafter_parse_blueprint(source, &result, opts))
	if errCode != 0 {
		return nil, fmt.Errorf("Drafter execution failed with code: %d", errCode)
	}

	err = json.Unmarshal(e.serialize(result), &objmap)

	if err != nil {
		return nil, err
	}

	return objmap, nil
}

func (e Blueprint) serialize(r *C.drafter_result) []byte {
	options := C.drafter_serialize_options{sourcemap: false, format: C.DRAFTER_SERIALIZE_JSON}
	cResult := C.drafter_serialize(r, options)
	defer C.free(unsafe.Pointer(cResult))

	results := C.GoString(cResult)

	return []byte(results)
}

// PrintRecursiveMap fixme remove this method as it's here for testing purpose only
func (e Blueprint) PrintRecursiveMap(mapval map[string]interface{}) {
	for key, value := range mapval {
		if str, ok := value.(string); ok {
			fmt.Println("Key:", key, "Value:", str)
		} else if strmap, ok := value.(map[string]interface{}); ok {
			fmt.Println("Key:", key)
			e.PrintRecursiveMap(strmap)
		} else if strmap, ok := value.([]interface{}); ok {
			fmt.Println("Key:", key)
			e.printRecursiveArray(strmap)
		}
	}
}

// printRecursiveArray fixme remove this method as it's here for testing purpose only
func (e Blueprint) printRecursiveArray(arrval []interface{}) {
	for key, value := range arrval {
		if str, ok := value.(string); ok {
			fmt.Println("Key:", key, "Value:", str)
		} else if strmap, ok := value.(map[string]interface{}); ok {
			fmt.Println("Key:", key)
			e.PrintRecursiveMap(strmap)
		} else if strmap, ok := value.([]interface{}); ok {
			fmt.Println("Key:", key)
			e.printRecursiveArray(strmap)
		}
	}
}
