package walker

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type ObjectWalker struct {
	object interface{}
}

func NewObjectWalker(data interface{}) ObjectWalker {
	return ObjectWalker{object: data}
}

func (b *ObjectWalker) Object() interface{} {
	return b.object
}

func (b *ObjectWalker) Value() reflect.Value {
	return reflect.ValueOf(b.object)
}

func (b *ObjectWalker) String() string {
	v := b.Value()

	if v.IsValid() {
		switch v.Kind() {
		case reflect.String:
			return v.String()
		case reflect.Float64:
			return strconv.FormatFloat(v.Float(), 'E', 0, 64)
		}
	}

	return ""
}

func (b *ObjectWalker) Path(key string) *ObjectWalker {
	return b.search(b.hierarchy(key)...)
}

func (b *ObjectWalker) Exists(key string) bool {
	return b.Path(key).Value().IsNil()
}

func (b *ObjectWalker) Index(i int) *ObjectWalker {
	val := b.Value()

	if val.Kind() == reflect.Slice {
		if i >= val.Len() {
			return &ObjectWalker{nil}
		}

		return &ObjectWalker{val.Index(i).Interface()}
	}

	return &ObjectWalker{nil}
}

func (b *ObjectWalker) Children() ([]*ObjectWalker, error) {
	val := b.Value()

	switch val.Kind() {
	case reflect.Slice:
		children := make([]*ObjectWalker, val.Len())
		for i := 0; i < val.Len(); i++ {
			children[i] = &ObjectWalker{val.Index(i).Interface()}
		}

		return children, nil
	case reflect.Map:
		children := []*ObjectWalker{}

		for _, key := range val.MapKeys() {
			children = append(children, &ObjectWalker{val.MapIndex(key).Interface()})
		}

		return children, nil
	}

	return nil, errors.New("is not object or array")
}

func (b *ObjectWalker) ChildrenMap() (map[string]*ObjectWalker, error) {
	val := b.Value()

	if val.Kind() == reflect.Map {
		children := map[string]*ObjectWalker{}

		for _, key := range val.MapKeys() {
			children[key.String()] = &ObjectWalker{val.MapIndex(key).Interface()}
		}

		return children, nil
	}

	return nil, errors.New("is not an object")
}

func (b *ObjectWalker) hierarchy(key string) []string {
	return strings.Split(key, ".")
}

func (b *ObjectWalker) search(keys ...string) *ObjectWalker {
	var object interface{}

	object = b.object

	for i := 0; i < len(keys); i++ {
		if bmap, ok := object.(map[string]interface{}); ok {
			object = bmap[keys[i]]
		} else if barr, ok := object.([]interface{}); ok {
			tmpArr := []interface{}{}

			for i := range barr {
				tmpObj := &ObjectWalker{barr[i]}
				val := tmpObj.search(keys[i:]...).object
				if val != nil {
					tmpArr = append(tmpArr, val)
				}
			}

			if len(tmpArr) == 0 {
				return &ObjectWalker{nil}
			}

			return &ObjectWalker{tmpArr}
		} else {
			return &ObjectWalker{nil}
		}
	}

	return &ObjectWalker{object}
}
