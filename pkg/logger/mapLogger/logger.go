package maplogger

import (
	"fmt"
	"strings"
)

const IDENT = "  " // TODO TEMP

func Format(inp map[string]interface{}) (string, error) {
	var sb strings.Builder
	err := format("", inp, &sb)
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

func format(
	indent string,
	inp map[string]interface{},
	sb *strings.Builder,
) (err error) {
	for k, v := range inp {
		fmt.Printf("TEMP!!! format 0 \"%v\"\n", k)
		err = _format(indent, k, v, sb)
		if err != nil {
			break
		}
	}
	return
}

func _format(
	indent, key string,
	ifc interface{},
	sb *strings.Builder,
) (err error) {
	switch typ := ifc.(type) {
	case []interface{}:
		fmt.Printf("TEMP!!! _format 1 \"%v[%v]\"\n", key, len(typ))
		for i, f := range typ {
			err = _format(fmt.Sprintf("%v%v", indent, IDENT), fmt.Sprintf("[%v]", i), f, sb)
			if err != nil {
				break
			}
		}
	case map[string]interface{}:
		fmt.Printf("TEMP!!! _format 2 \"%v{%v}\"\n", key, len(typ))
		err = format(fmt.Sprintf("%v%v", indent, IDENT), typ, sb)
	case string:
		fmt.Printf("TEMP!!! _format 3 \"%v\": \"%v\"\n", key, typ)
		fmt.Fprintf(sb, "%v\"%v\": \"%v\"\n", indent, key, typ)
	default:
		fmt.Printf("TEMP!!! _format 4 \"%v\": %v\n", key, typ)
		fmt.Fprintf(sb, "%v\"%v\": %v\n", indent, key, typ)
	}
	return
}
