package ordered

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

// Traverse traverse a yaml while preserving order.
func Traverse(
	inp []yaml.MapItem,
	action func([]string, interface{}) (interface{}, error),
	keys ...string,
) (
	out []yaml.MapItem,
	err error,
) {
	var nxt []yaml.MapItem
	out = make([]yaml.MapItem, 0)
	for _, itm := range inp {
		nxt, err = _traverse(itm.Key.(string), itm.Value, action, keys...)
		if err != nil {
			break
		}
		out = append(out, nxt...)
	}
	return
}

func _traverse(
	key string,
	ifc interface{},
	action func([]string, interface{}) (interface{}, error),
	keys ...string,
) (
	out []yaml.MapItem,
	err error,
) {
	keys = append(keys, key)
	switch typ := ifc.(type) {
	case []yaml.MapItem:
		var nxt []yaml.MapItem
		nxt, err = Traverse(typ, action, keys...)
		if err != nil {
			break
		}
		out = append(out, yaml.MapItem{Key: key, Value: nxt})
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
		out = append(out, yaml.MapItem{Key: key, Value: nxt})
	default:
		var act interface{}
		act, err = action(keys, typ)
		if err != nil {
			err = fmt.Errorf("[%v]%v", key, err)
			break
		}
		out = append(out, yaml.MapItem{Key: key, Value: act})
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
	case []yaml.MapItem:
		var nxt []yaml.MapItem
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
		act, err = action(keys, typ)
		if err != nil {
			break
		}
		out = act
	}
	return
}
