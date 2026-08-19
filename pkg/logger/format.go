package logger

import (
	"fmt"
	"strings"
)

const INDENT = "  "

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
		for i, f := range typ {
			err = __format(indent, key, i, f, sb)
			if err != nil {
				break
			}
		}
	case map[string]interface{}:
		fmt.Fprintf(sb, "%v%v:\n", indent, key)
		err = format(fmt.Sprintf("%v%v", indent, INDENT), typ, sb)
	case []string:
		fmt.Fprintf(sb, "%v%v: ", indent, key)
		for _, f := range typ {
			fmt.Fprintf(sb, " \"%v\"", f)
		}
		fmt.Fprintf(sb, "\n")
	case string:
		fmt.Fprintf(sb, "%v%v: \"%v\"\n", indent, key, typ)
	default:
		fmt.Fprintf(sb, "%v%v: %v\n", indent, key, typ)
	}
	return
}

func __format(
	indent, key string,
	idx int,
	ifc interface{},
	sb *strings.Builder,
) (err error) {
	switch typ := ifc.(type) {
	case []interface{}:
		for i, f := range typ {
			err = __format(indent, fmt.Sprintf("%v[%v]", key, idx), i, f, sb)
			if err != nil {
				break
			}
		}
	case map[string]interface{}:
		fmt.Fprintf(sb, "%v%v[%v]:\n", indent, key, idx)
		err = format(fmt.Sprintf("%v%v", indent, INDENT), typ, sb)
	case []string:
		fmt.Fprintf(sb, "%v%v: ", indent, key)
		for _, f := range typ {
			fmt.Fprintf(sb, " \"%v\"", f)
		}
		fmt.Fprintf(sb, "\n")
	case string:
		fmt.Fprintf(sb, "%v%v[%v]: \"%v\"\n", indent, key, idx, typ)
	default:
		fmt.Fprintf(sb, "%v%v[%v]: %v\n", indent, key, idx, typ)
	}
	return
}
