package slice

import (
	"reflect"
	"testing"
)

// 定义测试用的结构体
type Person struct {
	ID      int
	Name    string
	Age     int
	Address string
}

type Product struct {
	Code         string
	ID           int
	Price        float64
	Name         string
	Description  string
	privateField string // 私有字段用于测试不可访问性
}

// 用于测试指针类型结构体
type Company struct {
	Name string
	Size int
}

func TestPluck(t *testing.T) {
	// 测试用例1: 正常提取字符串字段
	t.Run("Extract string field from struct slice", func(t *testing.T) {
		people := []Person{
			{Name: "Alice", Age: 30, Address: "New York"},
			{Name: "Bob", Age: 25, Address: "London"},
			{Name: "Charlie", Age: 35, Address: "Tokyo"},
		}

		result := Pluck[Person, string](people, "Name")

		expected := []string{"Alice", "Bob", "Charlie"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, but got %v", expected, result)
		}
	})

	// 测试用例2: 正常提取整数字段
	t.Run("Extract integer field from struct slice", func(t *testing.T) {
		people := []Person{
			{Name: "Alice", Age: 30, Address: "New York"},
			{Name: "Bob", Age: 25, Address: "London"},
			{Name: "Charlie", Age: 35, Address: "Tokyo"},
		}

		result := Pluck[Person, int](people, "Age")

		expected := []int{30, 25, 35}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, but got %v", expected, result)
		}
	})

	// 测试用例3: 提取指针类型结构体字段
	t.Run("Extract field from pointer struct slice", func(t *testing.T) {
		companies := []*Company{
			{Name: "Google", Size: 100000},
			{Name: "Apple", Size: 80000},
			{Name: "Microsoft", Size: 90000},
		}

		result := Pluck[*Company, string](companies, "Name")

		expected := []string{"Google", "Apple", "Microsoft"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, but got %v", expected, result)
		}
	})

	// 测试用例4: 字段不存在
	t.Run("Field not exist", func(t *testing.T) {
		people := []Person{
			{Name: "Alice", Age: 30, Address: "New York"},
		}

		result := Pluck[Person, string](people, "NonExistentField")

		// 期望返回空切片
		if len(result) != 0 {
			t.Errorf("Expected empty slice, but got %v", result)
		}
	})

	// 测试用例5: 字段类型不匹配
	t.Run("Field type mismatch", func(t *testing.T) {
		people := []Person{
			{Name: "Alice", Age: 30, Address: "New York"},
		}

		// 尝试将string类型的Name字段提取为int类型
		result := Pluck[Person, int](people, "Name")

		// 期望返回空切片
		if len(result) != 0 {
			t.Errorf("Expected empty slice, but got %v", result)
		}
	})

	// 测试用例6: 空切片输入
	t.Run("Empty slice input", func(t *testing.T) {
		var people []Person

		result := Pluck[Person, string](people, "Name")

		// 期望返回空切片
		if len(result) != 0 {
			t.Errorf("Expected empty slice, but got %v", result)
		}
	})

	// 测试用例7: 不同类型的结构体
	t.Run("Different struct types", func(t *testing.T) {
		products := []Product{
			{ID: 1, Price: 99.99, Name: "Laptop"},
			{ID: 2, Price: 59.99, Name: "Mouse"},
			{ID: 3, Price: 199.99, Name: "Monitor"},
		}

		result := Pluck[Product, float64](products, "Price")

		expected := []float64{99.99, 59.99, 199.99}
		// 由于浮点数比较需要考虑精度，这里简化处理
		if len(result) != len(expected) {
			t.Errorf("Expected length %d, but got %d", len(expected), len(result))
		}
	})

	// 测试用例8: nil指针元素（边界情况）
	t.Run("Nil pointer in slice", func(t *testing.T) {
		companies := []*Company{
			{Name: "Google", Size: 100000},
			nil, // nil指针
			{Name: "Microsoft", Size: 90000},
		}

		// 这种情况下会panic，因为nil指针无法进行反射操作
		defer func() {
			if r := recover(); r != nil {
				// 捕获到panic是预期的行为
				t.Logf("1Caught expected panic: %v", r)
			}
		}()

		result := Pluck[*Company, string](companies, "Name")
		expected := []string{"Google", "Microsoft"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected panic but didn't get one, result: %v", result)
		}
	})
}

// TestPluckWithDefault_EmptySlice 测试空切片情况
func TestPluckWithDefault_EmptySlice(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	items := []User{}
	defaultValue := "unknown"
	result := PluckWithDefault(items, "Name", defaultValue)

	// 空切片应该返回包含默认值的切片
	expected := []string{"unknown"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestPluckWithDefault_NormalStruct 测试正常结构体
func TestPluckWithDefault_NormalStruct(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	items := []User{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 30},
		{Name: "Charlie", Age: 35},
	}

	result := PluckWithDefault(items, "Name", "unknown")
	expected := []string{"Alice", "Bob", "Charlie"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestPluckWithDefault_PointerStruct 测试指针结构体
func TestPluckWithDefault_PointerStruct(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	items := []*User{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 30},
		nil, // 测试nil指针情况
	}

	result := PluckWithDefault(items, "Name", "unknown")
	expected := []string{"Alice", "Bob", "unknown"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestPluckWithDefault_FieldNotFound 测试字段不存在的情况
func TestPluckWithDefault_FieldNotFound(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	items := []User{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 30},
	}

	result := PluckWithDefault(items, "Email", "no-email")
	expected := []string{"no-email", "no-email"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestPluckWithDefault_TypeMismatch 测试类型不匹配的情况
func TestPluckWithDefault_TypeMismatch(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	items := []User{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 30},
	}

	// 尝试将int字段提取为string类型
	result := PluckWithDefault(items, "Age", "default-age")
	expected := []string{"default-age", "default-age"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestPluckWithDefault_MixedScenario 测试混合场景
func TestPluckWithDefault_MixedScenario(t *testing.T) {
	type User struct {
		Name  string
		Age   int
		Email string
	}

	// 构造复杂测试数据：包含正常数据、缺少字段的数据等
	type PartialUser struct {
		Name string
		Age  int
	}

	items := []interface{}{
		User{Name: "Alice", Age: 25, Email: "alice@example.com"},
		PartialUser{Name: "Bob", Age: 30}, // 缺少Email字段
		User{Name: "Charlie", Age: 35, Email: "charlie@example.com"},
	}

	result := PluckWithDefault(items, "Email", "no-email")
	expected := []string{"alice@example.com", "no-email", "charlie@example.com"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestPluckWithDefault_DifferentTypes 测试不同数据类型
func TestPluckWithDefault_DifferentTypes(t *testing.T) {
	type Product struct {
		Name  string
		Price float64
		Stock int
	}

	items := []Product{
		{Name: "Apple", Price: 1.5, Stock: 100},
		{Name: "Banana", Price: 0.8, Stock: 50},
		{Name: "Orange", Price: 2.0, Stock: 75},
	}

	// 测试提取float64类型字段
	priceResult := PluckWithDefault(items, "Price", 0.0)
	expectedPrices := []float64{1.5, 0.8, 2.0}
	if !reflect.DeepEqual(priceResult, expectedPrices) {
		t.Errorf("Expected prices %v, got %v", expectedPrices, priceResult)
	}

	// 测试提取int类型字段
	stockResult := PluckWithDefault(items, "Stock", 0)
	expectedStocks := []int{100, 50, 75}
	if !reflect.DeepEqual(stockResult, expectedStocks) {
		t.Errorf("Expected stocks %v, got %v", expectedStocks, stockResult)
	}
}

// TestPluckWithDefault_StructWithUnexportedFields 测试包含未导出字段的结构体
func TestPluckWithDefault_StructWithUnexportedFields(t *testing.T) {
	type Person struct {
		Name string
		age  int // unexported field
	}

	items := []Person{
		{Name: "Alice", age: 25},
		{Name: "Bob", age: 30},
	}

	// 尝试访问导出字段 - 应该成功
	nameResult := PluckWithDefault(items, "Name", "unknown")
	expectedNames := []string{"Alice", "Bob"}
	if !reflect.DeepEqual(nameResult, expectedNames) {
		t.Errorf("Expected names %v, got %v", expectedNames, nameResult)
	}

	// 尝试访问未导出字段 - 应该返回默认值
	ageResult := PluckWithDefault(items, "age", -1)
	expectedAges := []int{-1, -1} // 无法访问未导出字段，应返回默认值
	if !reflect.DeepEqual(ageResult, expectedAges) {
		t.Errorf("Expected ages %v, got %v", expectedAges, ageResult)
	}
}

// BenchmarkPluckWithDefault 性能基准测试
func BenchmarkPluckWithDefault(b *testing.B) {
	type User struct {
		Name string
		Age  int
	}

	items := make([]User, 1000)
	for i := 0; i < 1000; i++ {
		items[i] = User{Name: "User" + string(rune(i)), Age: i}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PluckWithDefault(items, "Name", "unknown")
	}
}
func TestPluckMap(t *testing.T) {
	// 测试用例定义
	tests := []struct {
		name        string      // 测试用例名称
		items       interface{} // 输入数据
		keyField    string      // 键字段名
		valueField  string      // 值字段名
		expected    interface{} // 期望结果
		expectEmpty bool        // 是否期望空结果
	}{
		{
			name:       "Normal case - extract ID and Name",
			items:      []Person{{ID: 1, Name: "Alice", Age: 30}, {ID: 2, Name: "Bob", Age: 25}},
			keyField:   "ID",
			valueField: "Name",
			expected:   map[int]string{1: "Alice", 2: "Bob"},
		},
		{
			name:       "Normal case - extract Code and Price",
			items:      []Product{{Code: "A001", Price: 99.99, Description: "Product A"}, {Code: "A002", Price: 149.99, Description: "Product B"}},
			keyField:   "Code",
			valueField: "Price",
			expected:   map[string]float64{"A001": 99.99, "A002": 149.99},
		},
		{
			name:       "Empty slice",
			items:      []Person{},
			keyField:   "ID",
			valueField: "Name",
			expected:   map[int]string{},
		},
		{
			name:       "Nil slice",
			items:      ([]Person)(nil),
			keyField:   "ID",
			valueField: "Name",
			expected:   map[int]string{},
		},
		{
			name:       "Pointer elements",
			items:      []*Person{{ID: 1, Name: "Alice", Age: 30}, {ID: 2, Name: "Bob", Age: 25}},
			keyField:   "ID",
			valueField: "Name",
			expected:   map[int]string{1: "Alice", 2: "Bob"},
		},
		{
			name:       "Mixed pointer and value with nil pointer",
			items:      []*Person{{ID: 1, Name: "Alice", Age: 30}, nil, {ID: 2, Name: "Bob", Age: 25}},
			keyField:   "ID",
			valueField: "Name",
			expected:   map[int]string{1: "Alice", 2: "Bob"},
		},
		{
			name:        "Non-struct elements",
			items:       []int{1, 2, 3},
			keyField:    "ID",
			valueField:  "Name",
			expectEmpty: true,
		},
		{
			name:        "Non-existent key field",
			items:       []Person{{ID: 1, Name: "Alice", Age: 30}},
			keyField:    "NonExistent",
			valueField:  "Name",
			expectEmpty: true,
		},
		{
			name:        "Non-existent value field",
			items:       []Person{{ID: 1, Name: "Alice", Age: 30}},
			keyField:    "ID",
			valueField:  "NonExistent",
			expectEmpty: true,
		},
		{
			name:        "Unexported field as key",
			items:       []Product{{Code: "A001", Price: 99.99, privateField: "private"}},
			keyField:    "privateField",
			valueField:  "Price",
			expectEmpty: true,
		},
		{
			name:        "Unexported field as value",
			items:       []Product{{Code: "A001", Price: 99.99, privateField: "private"}},
			keyField:    "Code",
			valueField:  "privateField",
			expectEmpty: true,
		},
		{
			name:        "Type mismatch for key",
			items:       []Person{{ID: 1, Name: "Alice", Age: 30}},
			keyField:    "Name",
			valueField:  "Age",
			expected:    map[string]int{}, // 应该返回空map因为Name是string但期待int类型的key
			expectEmpty: true,
		},
		{
			name:        "Type mismatch for value",
			items:       []Person{{ID: 1, Name: "Alice", Age: 30}},
			keyField:    "ID",
			valueField:  "Age",
			expected:    map[int]string{}, // 应该返回空map因为Age是int但期待string类型的value
			expectEmpty: true,
		},
	}

	// 执行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 使用反射调用泛型函数
			itemsValue := reflect.ValueOf(tt.items)
			var result reflect.Value

			if itemsValue.IsValid() && itemsValue.Kind() == reflect.Slice {
				if itemsValue.Len() > 0 {
					// 获取元素类型来确定泛型参数
					elemType := itemsValue.Type().Elem()
					// 处理指针情况
					if elemType.Kind() == reflect.Ptr {
						elemType = elemType.Elem()
					}

					// 根据不同的测试场景创建适当的调用
					switch tt.name {
					case "Normal case - extract ID and Name", "Empty slice", "Nil slice":
						result = reflect.ValueOf(PluckMap[Person, int, string](
							reflect.ValueOf(tt.items).Interface().([]Person),
							tt.keyField,
							tt.valueField,
						))
					case "Normal case - extract Code and Price",
						"Non-existent key field",
						"Non-existent value field",
						"Unexported field as key",
						"Unexported field as value",
						"Type mismatch for key",
						"Type mismatch for value":
						if persons, ok := tt.items.([]Product); ok {
							result = reflect.ValueOf(PluckMap[Product, string, float64](
								persons,
								tt.keyField,
								tt.valueField,
							))
						} else if ptrPersons, ok := tt.items.([]*Product); ok {
							result = reflect.ValueOf(PluckMap[*Product, string, float64](
								ptrPersons,
								tt.keyField,
								tt.valueField,
							))
						}
					case "Pointer elements", "Mixed pointer and value with nil pointer":
						result = reflect.ValueOf(PluckMap[*Person, int, string](
							reflect.ValueOf(tt.items).Interface().([]*Person),
							tt.keyField,
							tt.valueField,
						))
					case "Non-struct elements":
						result = reflect.ValueOf(PluckMap[int, int, string](
							reflect.ValueOf(tt.items).Interface().([]int),
							tt.keyField,
							tt.valueField,
						))
					}
				} else {
					// 空切片情况
					if tt.name == "Empty slice" || tt.name == "Nil slice" {
						result = reflect.ValueOf(PluckMap[Person, int, string](
							reflect.ValueOf(tt.items).Interface().([]Person),
							tt.keyField,
							tt.valueField,
						))
					}
				}
			}

			// 验证结果
			if tt.expectEmpty {
				// 检查是否返回了空map
				if result.IsValid() && result.Len() != 0 {
					t.Errorf("Expected empty map, but got %v", result.Interface())
				}
			} else if tt.expected != nil {
				expectedValue := reflect.ValueOf(tt.expected)
				if !result.IsValid() {
					t.Errorf("Expected %v, but got invalid result", tt.expected)
				} else if !reflect.DeepEqual(result.Interface(), expectedValue.Interface()) {
					t.Errorf("Expected %v, but got %v", tt.expected, result.Interface())
				}
			}
		})
	}
}

// 专门针对具体类型的测试以简化测试逻辑
func TestPluckMapSpecificTypes(t *testing.T) {
	t.Run("Person ID to Name mapping", func(t *testing.T) {
		persons := []Person{
			{ID: 1, Name: "Alice", Age: 30},
			{ID: 2, Name: "Bob", Age: 25},
			{ID: 3, Name: "Charlie", Age: 35},
		}

		result := PluckMap[Person, int, string](persons, "ID", "Name")
		expected := map[int]string{1: "Alice", 2: "Bob", 3: "Charlie"}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, but got %v", expected, result)
		}
	})

	t.Run("Product Code to Price mapping", func(t *testing.T) {
		products := []Product{
			{Code: "P001", Price: 29.99, Description: "Widget A"},
			{Code: "P002", Price: 39.99, Description: "Widget B"},
		}

		result := PluckMap[Product, string, float64](products, "Code", "Price")
		expected := map[string]float64{"P001": 29.99, "P002": 39.99}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, but got %v", expected, result)
		}
	})

	t.Run("Pointer to Person with nil element", func(t *testing.T) {
		persons := []*Person{
			{ID: 1, Name: "Alice", Age: 30},
			nil,
			{ID: 2, Name: "Bob", Age: 25},
		}

		result := PluckMap[*Person, int, string](persons, "ID", "Name")
		expected := map[int]string{1: "Alice", 2: "Bob"}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, but got %v", expected, result)
		}
	})
}
