package utils

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"math"
	"os"
	"strings"
)

var Template = "import 'package:freezed_annotation/freezed_annotation.dart';\nimport 'package:flutter/foundation.dart';\n%[4]s\n\npart '%[1]s.freezed.dart';\npart '%[1]s.g.dart';\n\n@freezed\nclass %[2]s with _$%[2]s {\n  const factory %[2]s({ %[3]s }) = _%[2]s;\n  \n  factory %[2]s.fromJson(Map<String, dynamic> json) => _$%[2]sFromJson(json);\n}"

var BuiltinTypes = StrSlice{
	"int",
	"num",
	"double",
	"String",
	"bool",
	"dynamic",
	"int?",
	"num?",
	"double?",
	"String?",
	"bool?",
	"dynamic?",
	"List<String>",
	"List<int>",
	"List<num>",
	"List<double>",
	"List<bool>",
	"List",
	"List<String>?",
	"List<int>?",
	"List<num>?",
	"List<double>?",
	"List<bool>?",
	"List?",
	"Map<String, dynamic>",
	"Map<int, dynamic>",
	"Map<num, dynamic>",
	"Map<double, dynamic>",
	"Map<bool, dynamic>",
	"Map<dynamic, dynamic>",
	"Map<String, dynamic>?",
	"Map<int, dynamic>?",
	"Map<num, dynamic>?",
	"Map<double, dynamic>?",
	"Map<bool, dynamic>?",
	"Map<dynamic, dynamic>?",
}

func IsBuiltInType(typeStr string) bool {
	return BuiltinTypes.Has(strings.TrimSpace(typeStr))
}

func GetType(value interface{}, importsSlice *[]string, tag string) string {
	switch value.(type) {
	case int, int16, int32, int64, uint, uint16, uint32, uint64:
		return "int"
	case float64, float32:
		value := value.(float64)
		if math.Round(value) == value {
			return "int"
		}
		return "double"
	case bool:
		return "bool"
	case nil:
		return "String?"
	case map[string]interface{}:
		return "Map<String, dynamic>"
	case []interface{}:
		return "List"
	case string:
		value := value.(string)
		if strings.HasPrefix(value, fmt.Sprintf("%s[]", tag)) {
			value = strings.TrimPrefix(value, fmt.Sprintf("%s[]", tag))
			if IsBuiltInType(value) {
				return fmt.Sprintf("List<%s>", value)
			}

			*importsSlice = append(*importsSlice, value)
			return fmt.Sprintf("List<%s>", strcase.ToCamel(value))
		}
		if strings.HasPrefix(value, tag) {
			value = strings.TrimPrefix(value, tag)
			if IsBuiltInType(value) {
				return value
			}

			*importsSlice = append(*importsSlice, value)
			return strcase.ToCamel(value)
		}

		if strings.HasPrefix(value, "@") {
			return value
		}

		return "String"
	}

	return "dynamic"
}

func JsonToDart(json map[string]interface{}, tag string, name string) string {
	var attrSlice []string
	var importsSlice []string

	for key, value := range json {
		typeStr := GetType(value, &importsSlice, tag)
		if strings.HasSuffix(typeStr, "?") {
			if strcase.ToLowerCamel(key) != key {
				attrSlice = append(attrSlice, fmt.Sprintf("@JsonKey(name: \"%s\") %s %s", key, typeStr, strcase.ToLowerCamel(key)))
			} else {
				attrSlice = append(attrSlice, fmt.Sprintf("%s %s", typeStr, key))
			}
		} else {
			if strcase.ToLowerCamel(key) != key {
				attrSlice = append(attrSlice, fmt.Sprintf("@JsonKey(name: \"%s\") required %s %s", key, typeStr, strcase.ToLowerCamel(key)))
			} else {
				attrSlice = append(attrSlice, fmt.Sprintf("required %s %s", typeStr, key))
			}
		}
	}

	attrString := strings.Join(attrSlice[:], ", ")
	importsStr := ""
	for _, value := range importsSlice {
		importsStr += fmt.Sprintf("\nimport '%s.dart';", strcase.ToKebab(value))
	}

	return fmt.Sprintf(Template, strcase.ToKebab(name), strcase.ToCamel(name), attrString, importsStr)
}

func EnsureDir(dirName string) error {
	if _, err := os.Stat(dirName); err != nil {
		err := os.MkdirAll(dirName, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
