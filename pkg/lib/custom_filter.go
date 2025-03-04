package lib

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
)

// FilterItem struct
type FilterItem struct {
	Field     string
	Operator  string
	Value     interface{}
	ValueType string
}

// QueryFilter struct
type QueryFilter struct {
	Item FilterItem
	Type string
}

// CustomFilters func
// => Example
// -> QueryFilters      := [["id","=","6"],["AND"],["status_transaction","=","waiting"],["AND"],["business_id","=","10"]]
// -> QuerySearch       := "value"
// -> columnFilters    := ["trx_id","id"]
//
//gocyclo:ignore
func CustomFilters(QueryFilters, QuerySearch, columnFilter string) (string, []interface{}, string, []interface{}) {
	queryFilters := []string{}
	whereFilters := []interface{}{}
	ResultFilters := ""

	querySearch := []string{}
	whereSearch := []interface{}{}
	ResultSearch := ""

	// search_by _column
	if columnFilter != "" {
		columnFilters := StringToJson(columnFilter)
		if nil != columnFilters {
			if CountLengthIface(columnFilters) > 0 {
				if QuerySearch != "" {
					switch reflect.TypeOf(columnFilters).Kind() {
					case reflect.Slice:
						s := reflect.ValueOf(columnFilters)
						for i := 0; i < s.Len(); i++ {
							value := s.Index(i).Interface().(string)
							querySearch = append(querySearch, ""+value+" LIKE ?")
							whereSearch = append(whereSearch, "%"+QuerySearch+"%")
						}
					}
					if len(querySearch) > 0 {
						ResultSearch = "" + strings.Join(querySearch, " OR ") + ""
					}
				}
			}
		}
	}

	// filters_by_field
	if QueryFilters != "" {
		filters := CreateFilter(QueryFilters)
		if len(filters) == 1 && filters[0].Type == "single" {
			CreateWhereCause(filters[0], &queryFilters, &whereFilters)
		} else if len(filters) >= 1 {
			for b, filter := range filters {
				if filter.Type == "operator" {
					queryFilters = append(queryFilters, filter.Item.Operator)
					continue
				} else if filter.Type == "multiple" {
					if b > 0 && filters[b-1].Type != "operator" {
						queryFilters = append(queryFilters, " OR ")
					}
					CreateWhereCause(filter, &queryFilters, &whereFilters)
				}
			}
		}

		if len(queryFilters) > 0 {
			ResultFilters = "" + strings.Join(queryFilters, " ") + ""
		}
	}

	ResultFilters = strings.ReplaceAll(ResultFilters, "\"", "")

	return ResultFilters, whereFilters, ResultSearch, whereSearch
}

// NormalizeFieldName func
func NormalizeFieldName(field string) string {
	slices := strings.Split(field, "__")
	if len(slices) == 1 {
		return field
	}
	newSlices := []string{}
	if len(slices) > 0 {
		newSlices = append(newSlices, strcase.ToCamel(slices[0]))
		for k, s := range slices {
			if k > 0 {
				newSlices = append(newSlices, s)
			}
		}
	}
	return strings.Join(newSlices, "__")
}

// CreateFilter func
//
//gocyclo:ignore
func CreateFilter(jsonParams string) []QueryFilter {
	var output interface{}
	err := JSONUnmarshal([]byte(jsonParams), &output)
	if nil != err {
		return []QueryFilter{}
	}

	filters := []QueryFilter{}

	iface, ok := output.([]interface{})
	if ok {
		var hasSingle = false
		singleFilter := QueryFilter{
			Item: FilterItem{},
			Type: "single",
		}
		for x, v := range iface {
			item, ok2 := v.([]interface{})
			if ok2 && !hasSingle {
				filter := QueryFilter{
					Type: "multiple",
					Item: FilterItem{},
				}
				for i, a := range item {
					if len(item) == 1 {
						filter.Item.Operator = a.(string)
						filter.Type = "operator"
						continue
					}
					if i == 0 {
						filter.Item.Field = a.(string)
					} else if i == 1 {
						if len(item) == 2 {
							filter.Item.Operator = "="
							SetFilterValue(&filter.Item, a)
						} else {
							filter.Item.Operator = a.(string)
						}
					} else if i == 2 {
						SetFilterValue(&filter.Item, a)
					}
				}
				filter.Item.Field = "\"" + NormalizeFieldName(filter.Item.Field) + "\""
				filter.Item.Operator = strings.ToUpper(filter.Item.Operator)
				filters = append(filters, filter)
			} else {
				hasSingle = true
				if x == 0 {
					fieldName, valid := v.(string)
					if valid {
						singleFilter.Item.Field = "\"" + NormalizeFieldName(fieldName) + "\""
					}
				} else if x == 1 {
					if len(iface) == 2 {
						singleFilter.Item.Operator = "="
						SetFilterValue(&singleFilter.Item, v)
					} else {
						opName, valid := v.(string)
						if valid {
							singleFilter.Item.Operator = strings.ToUpper(opName)
						}
					}
				} else if x == 2 {
					SetFilterValue(&singleFilter.Item, v)
				}
			}
		}
		if hasSingle {
			filters = append(filters, singleFilter)
		}
	}

	return filters
}

// SetFilterValue func
func SetFilterValue(item *FilterItem, a interface{}) {
	stringValue, isString := a.(string)
	boolValue, isBool := a.(bool)
	floatValue, isFloat := a.(float64)
	arrayValue, isArray := a.([]interface{})
	if isString {
		item.Value = stringValue
		item.ValueType = "string"
	} else if isBool {
		item.Value = boolValue
		item.ValueType = "bool"
	} else if isFloat {
		item.Value = floatValue
		item.ValueType = "float64"
	} else if isArray {
		item.Value = arrayValue
		item.ValueType = "array"
	}
}

// CreateWhereCause func
//
//gocyclo:ignore
func CreateWhereCause(filter QueryFilter, queryFilters *[]string, whereParams *[]interface{}) {
	if (filter.Item.Operator == "IS" || filter.Item.Operator == "IS NOT") && filter.Item.Value == nil {
		*queryFilters = append(*queryFilters, fmt.Sprintf("(%s %s NULL)",
			filter.Item.Field,
			filter.Item.Operator,
		))
		return
	}

	if filter.Item.Operator == "LIKE" || filter.Item.Operator == "NOT LIKE" {
		cause := fmt.Sprintf("%s %s ?",
			filter.Item.Field,
			filter.Item.Operator,
		)
		*queryFilters = append(*queryFilters, cause)
		value := ""
		switch filter.Item.ValueType {
		case "string":
			value = "%" + (filter.Item.Value.(string)) + "%"
		case "float64":
			value = "%" + fmt.Sprintf("%v", filter.Item.Value.(float64)) + "%"
		}

		if value != "" {
			value = strings.ReplaceAll(value, " ", "%")
			*whereParams = append(*whereParams, value)
		}

	} else if filter.Item.Operator == "IN" || filter.Item.Operator == "NOT IN" {
		cause := fmt.Sprintf("%s %s ?",
			filter.Item.Field,
			filter.Item.Operator,
		)
		*queryFilters = append(*queryFilters, cause)
		if filter.Item.ValueType == "array" {
			value, ok := filter.Item.Value.([]interface{})
			if ok {
				values := []interface{}{}
				for _, val := range value {
					v, o := val.(string)
					if o {
						values = append(values, strings.ToLower(v))
					} else {
						values = append(values, fmt.Sprintf("%v", v))
					}
				}
				*whereParams = append(*whereParams, values)
			}
		}
	} else if filter.Item.Operator == "BETWEEN" {
		cause := fmt.Sprintf("%s %s ? AND ?",
			filter.Item.Field,
			filter.Item.Operator,
		)
		*queryFilters = append(*queryFilters, cause)
		if filter.Item.ValueType == "array" {
			value, ok := filter.Item.Value.([]interface{})
			if ok {
				values := []interface{}{}
				for _, val := range value {
					v, o := val.(string)
					if o {
						values = append(values, strings.ToLower(v))
					}
				}

				if len(values) == 2 {
					*whereParams = append(*whereParams, values[0], values[1])
				}
			}
		}
	} else {
		cause := fmt.Sprintf("%s %s ?",
			filter.Item.Field,
			filter.Item.Operator,
		)
		*queryFilters = append(*queryFilters, cause)
		*whereParams = append(*whereParams, fmt.Sprintf("%v", filter.Item.Value))
	}
}

// StringToJson func
func StringToJson(value string) interface{} {
	var output interface{}
	JSONUnmarshal([]byte(value), &output)
	return output
}

// CountLengthIface func
func CountLengthIface(data interface{}) int {
	var value int
	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(data)
		value = s.Len()
	}
	return value
}
