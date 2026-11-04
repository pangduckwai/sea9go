// Package ordered implements traversal of `[]yaml.MapItem` structures from yaml while keeping items order of the file.
package ordered

import (
	"fmt"

	"github.com/pangduckwai/sea9go/pkg/traverse"
	"gopkg.in/yaml.v2"
)

// Traverse traverse a yaml while preserving order.
func Traverse(
	inp []yaml.MapItem,
	action func(
		traverse.NodeType, // node type
		[]string, // keys
		interface{}, // value
	) (
		interface{}, // processed value
		error,
	),
	keys ...string,
) (
	out []yaml.MapItem,
	err error,
) {
	var nxt []yaml.MapItem
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
	action func(
		traverse.NodeType,
		[]string,
		interface{},
	) (
		interface{},
		error,
	),
	keys ...string,
) (
	out []yaml.MapItem,
	err error,
) {
	keys = append(keys, key)
	switch typ := ifc.(type) {
	case []yaml.MapItem:
		var trm interface{}
		trm, err = action(traverse.OBJT, keys, ifc)
		if err != nil && err == traverse.ShouldStop && trm != nil {
			out = append(out, yaml.MapItem{Key: key, Value: trm})
			err = nil
			break
		}
		var nxt []yaml.MapItem
		nxt, err = Traverse(typ, action, keys...)
		if err != nil {
			break
		}
		out = append(out, yaml.MapItem{Key: key, Value: nxt})
	case []interface{}:
		var itm interface{}
		itm, err = action(traverse.LIST, keys, ifc)
		if err != nil && err == traverse.ShouldStop && itm != nil {
			out = append(out, yaml.MapItem{Key: key, Value: itm})
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
		out = append(out, yaml.MapItem{Key: key, Value: nxt})
	default:
		var act interface{}
		act, err = action(traverse.LEAF, keys, typ)
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
	action func(
		traverse.NodeType,
		[]string,
		interface{},
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
	case []yaml.MapItem:
		var trm interface{}
		trm, err = action(traverse.OBJT, keys, ifc)
		if err != nil && err == traverse.ShouldStop && trm != nil {
			out = trm
			err = nil
			break
		}
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
		act, err = action(traverse.LEAF, keys, typ)
		if err != nil {
			break
		}
		out = act
	}
	return
}
