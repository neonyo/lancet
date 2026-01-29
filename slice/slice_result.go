package slice

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/neonyo/lancet/v2/enum"
)

// Pluck returns a slice of values from each item in the slice.
func Pluck[S any, T comparable](items []S, fieldName string) (result []T) {
	result = make([]T, 0)

	for _, item := range items {
		// 使用反射获取字段值
		val := reflect.ValueOf(item)
		if !val.IsValid() {
			continue
		}

		// 安全地解引用指针
		val = safeDeref(val)
		if !val.IsValid() {
			continue
		}

		if val.Kind() == reflect.Ptr {
			if val.IsNil() {
				continue
			}
			val = val.Elem()
		}

		field := val.FieldByName(fieldName)
		if !field.IsValid() {
			continue
		}

		// 类型断言
		fieldValue, ok := field.Interface().(T)
		if !ok {
			continue
		}

		result = append(result, fieldValue)
	}

	return result
}

type FieldSelector struct {
	FieldName string
	Result    any // 指向切片的指针
}

// PluckBySelectors returns a slice of values from each item in the slice.
func PluckBySelectors[S any](items []S, selectors []FieldSelector) {
	if len(items) == 0 {
		return
	}

	// 第一个元素用于确定字段类型
	firstItem := items[0]
	val := reflect.ValueOf(firstItem)
	val = safeDeref(val)

	// 验证所有字段并初始化切片
	for _, selector := range selectors {
		field := val.FieldByName(selector.FieldName)
		if !field.IsValid() {
			return
		}

		// 检查 Result 是否是切片指针
		resultPtr := reflect.ValueOf(selector.Result)
		if resultPtr.Kind() != reflect.Ptr {
			return
		}

		sliceVal := resultPtr.Elem()
		if sliceVal.Kind() != reflect.Slice {
			return
		}

		// 确保切片元素类型匹配字段类型
		elemType := sliceVal.Type().Elem()
		if elemType != field.Type() {
			return
		}

		// 初始化或清空切片
		sliceVal.Set(reflect.MakeSlice(sliceVal.Type(), 0, len(items)))
	}

	// 填充数据
	for _, item := range items {
		itemVal := reflect.ValueOf(item)
		itemVal = safeDeref(itemVal)

		for _, selector := range selectors {
			field := itemVal.FieldByName(selector.FieldName)
			if !field.IsValid() {
				continue
			}

			// 获取对应的切片并追加
			slicePtr := reflect.ValueOf(selector.Result)
			sliceVal := slicePtr.Elem()
			sliceVal.Set(reflect.Append(sliceVal, field))
		}
	}

	return
}

// PluckFilter returns a slice of values from each item in the slice.
func PluckFilter[S any, T comparable](items []S, fieldName string, fn func(i S) bool) (result []T) {
	result = make([]T, 0)

	for _, item := range items {
		// 使用反射获取字段值
		val := reflect.ValueOf(item)
		if !val.IsValid() {
			continue
		}

		// 安全地解引用指针
		val = safeDeref(val)
		if !val.IsValid() {
			continue
		}

		if val.Kind() == reflect.Ptr {
			if val.IsNil() {
				continue
			}
			val = val.Elem()
		}

		field := val.FieldByName(fieldName)
		if !field.IsValid() {
			continue
		}

		// 类型断言
		fieldValue, ok := field.Interface().(T)
		if !ok {
			continue
		}

		if !fn(item) {
			continue
		}

		result = append(result, fieldValue)
	}

	return result
}

// PluckMap returns a slice of values from each item in the slice.
func PluckMap[S any, K, V comparable](items []S, keyField, valueField string) (result map[K]V) {
	result = make(map[K]V)
	if items == nil {
		return
	}

	for _, item := range items {
		val := reflect.ValueOf(item)

		// 处理指针
		if val.Kind() == reflect.Ptr {
			if val.IsNil() {
				continue // 跳过 nil 指针
			}
			val = val.Elem()
		}
		val = safeDeref(val)
		if !val.IsValid() {
			continue
		}

		// 确保是结构体
		if val.Kind() != reflect.Struct {
			return
		}

		// 获取 key 字段
		keyFieldVal := val.FieldByName(keyField)
		if !keyFieldVal.IsValid() {
			return
		}
		if !keyFieldVal.CanInterface() {
			return
		}

		// 获取 value 字段
		valueFieldVal := val.FieldByName(valueField)
		if !valueFieldVal.IsValid() {
			return
		}
		if !valueFieldVal.CanInterface() {
			return
		}

		// 类型断言
		key, keyOk := keyFieldVal.Interface().(K)
		value, valueOk := valueFieldVal.Interface().(V)

		if !keyOk || !valueOk {
			return
		}

		result[key] = value
	}

	return result
}

// PluckMapSlice returns a slice of values from each item in the slice.
func PluckMapSlice[S any, K, V comparable](items []S, keyField, valueField string) (result map[K][]V) {
	result = make(map[K][]V)
	if items == nil {
		return
	}

	for _, item := range items {
		val := reflect.ValueOf(item)

		// 处理指针
		if val.Kind() == reflect.Ptr {
			if val.IsNil() {
				continue // 跳过 nil 指针
			}
			val = val.Elem()
		}
		val = safeDeref(val)
		if !val.IsValid() {
			continue
		}

		// 确保是结构体
		if val.Kind() != reflect.Struct {
			return
		}

		// 获取 key 字段
		keyFieldVal := val.FieldByName(keyField)
		if !keyFieldVal.IsValid() {
			return
		}
		if !keyFieldVal.CanInterface() {
			return
		}

		// 获取 value 字段
		valueFieldVal := val.FieldByName(valueField)
		if !valueFieldVal.IsValid() {
			return
		}
		if !valueFieldVal.CanInterface() {
			return
		}

		// 类型断言
		key, keyOk := keyFieldVal.Interface().(K)
		value, valueOk := valueFieldVal.Interface().(V)

		if !keyOk || !valueOk {
			return
		}

		if _, ok := result[key]; !ok {
			result[key] = make([]V, 0)
		}
		result[key] = append(result[key], value)
	}

	return result
}

// PluckWithDefault returns a slice of values from each item in the slice, with a default value if the field is not found.
func PluckWithDefault[S any, T any](items []S, fieldName string, defaultValue T) []T {
	if items == nil {
		return []T{defaultValue}
	}

	result := make([]T, len(items))

	for i, item := range items {
		val := reflect.ValueOf(item)

		// 1. 检查值是否有效
		if !val.IsValid() {
			result[i] = defaultValue
			continue
		}

		// 2. 安全地解引用指针
		val = safeDeref(val)
		if !val.IsValid() {
			result[i] = defaultValue
			continue
		}

		// 3. 确保是结构体
		if val.Kind() != reflect.Struct {
			result[i] = defaultValue
			continue
		}

		// 4. 获取字段
		field := val.FieldByName(fieldName)
		if !field.IsValid() {
			result[i] = defaultValue
			continue
		}

		// 5. 关键：检查字段是否可导出
		if !field.CanInterface() {
			result[i] = defaultValue
			continue
		}

		// 6. 安全地获取接口值
		fieldInterface := field.Interface()

		// 7. 类型转换
		if fieldValue, ok := fieldInterface.(T); ok {
			result[i] = fieldValue
		} else {
			result[i] = defaultValue
		}
	}

	// 如果原始切片为空，返回包含默认值的切片
	if len(result) == 0 {
		return []T{defaultValue}
	}

	return result
}

// safeDeref 安全地解引用指针
func safeDeref(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return reflect.Value{}
		}
		v = v.Elem()
	}
	return v
}

// EqualIgnoreOrder 比较切片和字符串是否存在相等 忽略顺序
func EqualIgnoreOrder[T enum.SimpleType](str string, separator string, slice []T) bool {
	fmt.Println(str, separator, slice)

	// 分割字符串
	strParts := strings.Split(str, separator)
	if str == "" && len(slice) == 0 {
		return true
	}
	// 如果长度不同，直接返回false
	if len(strParts) != len(slice) {
		return false
	}

	// 创建映射统计元素频率
	freq := make(map[T]int, len(strParts))
	// 统计字符串中的元素
	for _, part := range strParts {
		// 去除空格
		trimmed := strings.TrimSpace(part)
		elem, err := parseToType[T](trimmed)
		if err != nil {
			return false // 转换失败，直接返回 false
		}
		freq[elem]++
	}

	// 统计切片中的元素并比较
	for _, item := range slice {
		if count, exists := freq[item]; exists && count > 0 {
			freq[item]--
		} else {
			return false
		}
	}

	// 检查所有元素都匹配
	for _, count := range freq {
		if count != 0 {
			return false
		}
	}

	return true
}

func parseToType[T enum.SimpleType](s string) (T, error) {
	var zero T
	var result interface{}
	var err error

	// 根据类型进行转换
	switch any(zero).(type) {
	case string:
		result = s
	case int:
		var val int
		val, err = strconv.Atoi(s)
		result = val
	case int8:
		var val int64
		val, err = strconv.ParseInt(s, 10, 8)
		result = int8(val)
	case int16:
		var val int64
		val, err = strconv.ParseInt(s, 10, 16)
		result = int16(val)
	case int32:
		var val int64
		val, err = strconv.ParseInt(s, 10, 32)
		result = int32(val)
	case int64:
		var val int64
		val, err = strconv.ParseInt(s, 10, 64)
		result = val
	case uint:
		var val uint64
		val, err = strconv.ParseUint(s, 10, 0)
		result = uint(val)
	case uint8:
		var val uint64
		val, err = strconv.ParseUint(s, 10, 8)
		result = uint8(val)
	case uint16:
		var val uint64
		val, err = strconv.ParseUint(s, 10, 16)
		result = uint16(val)
	case uint32:
		var val uint64
		val, err = strconv.ParseUint(s, 10, 32)
		result = uint32(val)
	case uint64:
		var val uint64
		val, err = strconv.ParseUint(s, 10, 64)
		result = val
	case float32:
		var val float64
		val, err = strconv.ParseFloat(s, 32)
		result = float32(val)
	case float64:
		var val float64
		val, err = strconv.ParseFloat(s, 64)
		result = val
	case bool:
		var val bool
		val, err = strconv.ParseBool(s)
		result = val
	default:
		return zero, nil
	}

	if err != nil {
		return zero, err
	}

	return result.(T), nil
}

// ParseIntegerSlice 解析字符串为数字类型切片
func ParseIntegerSlice[T enum.Number](str string, bitSize int) ([]T, error) {
	if str == "" {
		return []T{}, nil
	}

	parts := strings.Split(str, ",")
	result := make([]T, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}

		// 根据类型转换
		var num T
		switch any(num).(type) {
		case int, int8, int16, int32, int64:
			val, err := strconv.ParseInt(trimmed, 10, bitSize)
			if err != nil {
				return nil, err
			}
			result = append(result, T(val))
		case uint, uint8, uint16, uint32, uint64:
			val, err := strconv.ParseUint(trimmed, 10, bitSize)
			if err != nil {
				return nil, err
			}
			result = append(result, T(val))
		case float32, float64:
			val, err := strconv.ParseFloat(trimmed, bitSize)
			if err != nil {
				return nil, err
			}
			result = append(result, T(val))
			return result, nil
		default:
			return nil, strconv.ErrSyntax
		}
	}

	return result, nil
}

func PluckColumnMap[T any, S comparable](items []T, keyField string) (result map[S]T) {
	result = make(map[S]T)
	if len(items) == 0 {
		return
	}

	for _, item := range items {
		val := reflect.ValueOf(item)

		// 处理指针
		if val.Kind() == reflect.Ptr {
			if val.IsNil() {
				continue // 跳过 nil 指针
			}
			val = val.Elem()
		}
		val = safeDeref(val)
		if !val.IsValid() {
			continue
		}

		// 确保是结构体
		if val.Kind() != reflect.Struct {
			return
		}

		// 获取 key 字段
		keyFieldVal := val.FieldByName(keyField)
		if !keyFieldVal.IsValid() {
			return
		}
		if !keyFieldVal.CanInterface() {
			return
		}

		// 类型断言
		key, keyOk := keyFieldVal.Interface().(S)

		if !keyOk {
			return
		}

		result[key] = item
	}

	return result
}

func PluckColumnSliceMap[T any, S comparable](items []T, keyField string) (result map[S][]T) {
	result = make(map[S][]T)
	if len(items) == 0 {
		return
	}

	for _, item := range items {
		val := reflect.ValueOf(item)

		// 处理指针
		if val.Kind() == reflect.Ptr {
			if val.IsNil() {
				continue // 跳过 nil 指针
			}
			val = val.Elem()
		}
		val = safeDeref(val)
		if !val.IsValid() {
			continue
		}

		// 确保是结构体
		if val.Kind() != reflect.Struct {
			return
		}

		// 获取 key 字段
		keyFieldVal := val.FieldByName(keyField)
		if !keyFieldVal.IsValid() {
			return
		}
		if !keyFieldVal.CanInterface() {
			return
		}

		// 类型断言
		key, keyOk := keyFieldVal.Interface().(S)

		if !keyOk {
			return
		}

		result[key] = append(result[key], item)
	}

	return result
}
