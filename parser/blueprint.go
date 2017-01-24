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
	"io/ioutil"
	"unsafe"
)

type BlueprintAPIParser struct {}

func NewBlueprintParser() *BlueprintAPIParser {
	return &BlueprintAPIParser{}
}

func (bp BlueprintAPIParser) Parse(filename string, formatter Formatter) (spec Specification, err error) {
	var raw []byte
	var data map[string]interface{}

	if raw, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	source := C.CString(string(raw))

	defer C.free(unsafe.Pointer(source))

	result := &C.drafter_result{}
	opts := C.drafter_parse_options{}

	if errCode := int(C.drafter_parse_blueprint(source, &result, opts)); errCode != 0 {
		err = fmt.Errorf("Drafter execution failed with code: %d", errCode)
		return
	}

	if err = json.Unmarshal(bp.serialize(result), &data); err == nil {
		spec = formatter.Format(data)
	}

	return
}

func (bp BlueprintAPIParser) serialize(drafterResult *C.drafter_result) ([]byte) {

	opts := C.drafter_serialize_options{sourcemap: false, format: C.DRAFTER_SERIALIZE_JSON}

	serializer := C.drafter_serialize(drafterResult, opts)

	defer C.free(unsafe.Pointer(serializer))

	return []byte(C.GoString(serializer))
}
