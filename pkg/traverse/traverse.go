package traverse

import "fmt"

// Traverse traverse a map[string]interface{} (e.g. from yaml/json)
func Traverse(
	inp map[string]interface{},
	action func(interface{}) (interface{}, error),
) (
	out map[string]interface{},
	err error,
) {
	var nxt map[string]interface{}
	out = make(map[string]interface{})
	for k, v := range inp {
		nxt, err = _traverse(k, v, action)
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
	action func(interface{}) (interface{}, error),
) (
	out map[string]interface{},
	err error,
) {
	out = make(map[string]interface{})
	switch typ := ifc.(type) {
	case []interface{}:
		var itm interface{}
		nxt := make([]interface{}, len(typ))
		for i, f := range typ {
			itm, err = __traverse(f, action)
			if err != nil {
				err = fmt.Errorf("[%v][%v]%v", key, i, err)
				break
			}
			nxt[i] = itm
		}
		out[key] = nxt
	case map[string]interface{}:
		var nxt map[string]interface{}
		nxt, err = Traverse(typ, action)
		if err != nil {
			break
		}
		out[key] = nxt
	default:
		if key == "header" {
			fmt.Printf("TEMP!!!!!!!!! 'default'\n")
		}
		var act interface{}
		act, err = action(typ)
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
	action func(interface{}) (interface{}, error),
) (
	out interface{},
	err error,
) {
	switch typ := ifc.(type) {
	case []interface{}:
		var itm interface{}
		nxt := make([]interface{}, len(typ))
		for i, f := range typ {
			itm, err = __traverse(f, action)
			if err != nil {
				err = fmt.Errorf("[%v]%v", i, err)
				break
			}
			nxt[i] = itm
		}
		out = nxt
	case map[string]interface{}:
		var nxt map[string]interface{}
		nxt, err = Traverse(typ, action)
		if err != nil {
			break
		}
		out = nxt
	default:
		var act interface{}
		act, err = action(typ)
		if err != nil {
			break
		}
		out = act
	}
	return
}
