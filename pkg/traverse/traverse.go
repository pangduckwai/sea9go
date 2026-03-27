// Package traverse implements traversal of `map[string]any` structures (e.g. from json/yaml).
package traverse

import "fmt"

// Traverse traverse a map[string]interface{} (e.g. from yaml/json)
func Traverse(
	inp map[string]interface{},
	action func([]string, interface{}) (interface{}, error),
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
	action func([]string, interface{}) (interface{}, error),
	keys ...string,
) (
	out map[string]interface{},
	err error,
) {
	keys = append(keys, key)
	out = make(map[string]interface{})
	switch typ := ifc.(type) {
	case []interface{}:
		var itm interface{}
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
	case map[string]interface{}:
		var nxt map[string]interface{}
		nxt, err = Traverse(typ, action, keys...)
		if err != nil {
			break
		}
		out[key] = nxt
	default:
		var act interface{}
		act, err = action(keys, typ)
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
	action func([]string, interface{}) (interface{}, error),
	keys ...string,
) (
	out interface{},
	err error,
) {
	switch typ := ifc.(type) {
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
	case map[string]interface{}:
		var nxt map[string]interface{}
		nxt, err = Traverse(typ, action, keys...)
		if err != nil {
			break
		}
		out = nxt
	default:
		var act interface{}
		act, err = action(keys, typ)
		if err != nil {
			break
		}
		out = act
	}
	return
}
