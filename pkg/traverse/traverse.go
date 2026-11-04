// Package traverse implements traversal of `map[string]any` structures (e.g. from json/yaml).
package traverse

import (
	"errors"
	"fmt"
)

var ShouldStop error = errors.New("traversal of the current branch should stop")

type NodeType int

const (
	LEAF NodeType = iota // 0
	LIST
	OBJT
)

// ParseKey parse key array (`keys`) for array index
func ParseKey(key string) (idx int) {
	ln := len(key)
	if ln < 3 {
		return -3
	}
	if key[0] != '[' || key[ln-1] != ']' || key[1] == '-' {
		return -2
	}
	for i := 1; i < ln-1; i++ {
		if key[i] < '0' || key[i] > '9' {
			return -1
		}
		idx = idx*10 + int(key[i]-'0')
	}
	return idx
}

// Traverse traverse a map[string]interface{} (e.g. from yaml/json)
func Traverse(
	inp map[string]interface{},
	action func(
		NodeType, // node type
		[]string, // keys
		interface{}, // value
	) (
		interface{},
		error,
	),
	keys ...string,
) (
	out map[string]interface{},
	err error,
) {
	var nxt map[string]interface{}
	out = make(map[string]interface{})
	for k, v := range inp {
		nxt, err = _traverse(k, v, action, keys...)
		if err != nil {
			break
		}
		out[k] = nxt
	}
	return
}

func _traverse(
	key string,
	ifc interface{},
	action func(
		NodeType, // node type
		[]string, // keys
		interface{}, // value
	) (
		interface{},
		error,
	),
	keys ...string,
) (
	out map[string]interface{},
	err error,
) {
	keys = append(keys, key)
	out = make(map[string]interface{})
	switch typ := ifc.(type) {
	case map[string]interface{}:
		var trm interface{}
		trm, err = action(OBJT, keys, ifc)
		if err != nil && err == ShouldStop && trm != nil {
			out[key] = trm
			err = nil
			break
		}
		var nxt map[string]interface{}
		nxt, err = Traverse(typ, action, keys...)
		if err != nil {
			break
		}
		out[key] = nxt
	case []interface{}:
		var itm interface{}
		itm, err = action(LIST, keys, ifc)
		if err != nil && err == ShouldStop && itm != nil {
			out[key] = itm
			err = nil
			break
		}
		nxt := make([]interface{}, len(typ))
		for i, f := range typ {
			itm, err = __traverse(f, action, append(keys, fmt.Sprintf("[%v]", i))...)
			if err != nil {
				err = fmt.Errorf("[%v][%v]%v", key, i, err)
				break
			}
			nxt[i] = itm
		}
		out[key] = nxt
	default:
		var act interface{}
		act, err = action(LEAF, keys, typ)
		if err != nil {
			err = fmt.Errorf("[%v]%v", key, err)
			break
		}
		out[key] = act
	}
	return
}

func __traverse(
	ifc interface{},
	action func(
		NodeType, // node type
		[]string, // keys
		interface{}, // value
	) (
		interface{},
		error,
	),
	keys ...string,
) (
	out interface{},
	err error,
) {
	switch typ := ifc.(type) {
	case map[string]interface{}:
		var trm interface{}
		trm, err = action(OBJT, keys, ifc)
		if err != nil && err == ShouldStop && trm != nil {
			out = trm
			err = nil
			break
		}
		var nxt map[string]interface{}
		nxt, err = Traverse(typ, action, keys...)
		if err != nil {
			break
		}
		out = nxt
	case []interface{}:
		var itm interface{}
		nxt := make([]interface{}, len(typ))
		for i, f := range typ {
			itm, err = __traverse(f, action, append(keys, fmt.Sprintf("[%v]", i))...)
			if err != nil {
				err = fmt.Errorf("[%v]%v", i, err)
				break
			}
			nxt[i] = itm
		}
		out = nxt
	default:
		var act interface{}
		act, err = action(LEAF, keys, typ)
		if err != nil {
			break
		}
		out = act
	}
	return
}
