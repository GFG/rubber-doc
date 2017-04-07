package walker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObjectWalker_NewWalker(t *testing.T) {
	wlkr := NewObjectWalker(struct{}{})

	assert.IsType(t, ObjectWalker{}, wlkr)
}

func TestObjectWalker_Children(t *testing.T) {
	data := []string{
		"string value 1",
		"string value 2",
	}

	wlkr := NewObjectWalker(data)

	children, _ := wlkr.Children()

	assert.Len(t, children, 2)
}

func TestObjectWalker_ChildrenMap(t *testing.T) {
	data := map[string]interface{}{
		"prop1": "string value 1",
		"prop2": "string value 2",
	}

	wlkr := NewObjectWalker(data)

	children, _ := wlkr.ChildrenMap()

	assert.Len(t, children, 2)
}

func TestObjectWalker_Path(t *testing.T) {
	data := map[string]interface{}{
		"prop1": map[string]interface{}{ "prop2": "string value" },
	}

	wlkr := NewObjectWalker(data)

	children, _ := wlkr.Children()

	assert.Equal(t, "string value", children[0].Path("prop2").String())
}

func TestObjectWalker_Exists(t *testing.T) {
	data := map[string]interface{}{
		"prop1": map[string]interface{}{ "prop2": "string value" },
	}

	wlkr := NewObjectWalker(data)

	children, _ := wlkr.Children()

	assert.True(t, children[0].Exists("prop2"))
	assert.False(t, children[0].Exists("doesntExist"))
}

func TestObjectWalker_Index(t *testing.T) {
	data := []string{
		"string value 1",
		"string value 2",
	}

	wlkr := NewObjectWalker(data)

	child := wlkr.Index(1)

	assert.Equal(t, "string value 2", child.String())
}