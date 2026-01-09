package slice

import (
	"reflect"
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
