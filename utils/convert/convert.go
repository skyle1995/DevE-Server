package convert

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// ToString 将任意类型转换为字符串
func ToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case time.Time:
		return v.Format(time.RFC3339)
	case []byte:
		return string(v)
	default:
		if v == nil {
			return ""
		}
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Sprintf("%v", value)
		}
		return string(data)
	}
}

// ToInt 将任意类型转换为int
func ToInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint64:
		return int(v), nil
	case float32:
		return int(v), nil
	case float64:
		return int(v), nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, fmt.Errorf("cannot convert %T to int", value)
	}
}

// ToInt64 将任意类型转换为int64
func ToInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to int64", value)
	}
}

// ToFloat64 将任意类型转换为float64
func ToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", value)
	}
}

// ToBool 将任意类型转换为bool
func ToBool(value interface{}) (bool, error) {
	switch v := value.(type) {
	case bool:
		return v, nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(value).Int() != 0, nil
	case float32, float64:
		return reflect.ValueOf(value).Float() != 0, nil
	case string:
		v = strings.ToLower(v)
		return v == "true" || v == "yes" || v == "1" || v == "on", nil
	default:
		return false, fmt.Errorf("cannot convert %T to bool", value)
	}
}

// ToStringSlice 将任意类型转换为字符串切片
func ToStringSlice(value interface{}) ([]string, error) {
	switch v := value.(type) {
	case []string:
		return v, nil
	case string:
		return strings.Split(v, ","), nil
	case []interface{}:
		result := make([]string, len(v))
		for i, item := range v {
			result[i] = ToString(item)
		}
		return result, nil
	default:
		return nil, fmt.Errorf("cannot convert %T to []string", value)
	}
}

// ToIntSlice 将任意类型转换为整数切片
func ToIntSlice(value interface{}) ([]int, error) {
	switch v := value.(type) {
	case []int:
		return v, nil
	case []interface{}:
		result := make([]int, len(v))
		for i, item := range v {
			intVal, err := ToInt(item)
			if err != nil {
				return nil, err
			}
			result[i] = intVal
		}
		return result, nil
	case string:
		parts := strings.Split(v, ",")
		result := make([]int, len(parts))
		for i, part := range parts {
			intVal, err := strconv.Atoi(strings.TrimSpace(part))
			if err != nil {
				return nil, err
			}
			result[i] = intVal
		}
		return result, nil
	default:
		return nil, fmt.Errorf("cannot convert %T to []int", value)
	}
}

// ToMap 将任意类型转换为map[string]interface{}
func ToMap(value interface{}) (map[string]interface{}, error) {
	switch v := value.(type) {
	case map[string]interface{}:
		return v, nil
	case string:
		result := make(map[string]interface{})
		err := json.Unmarshal([]byte(v), &result)
		return result, err
	default:
		data, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}
		result := make(map[string]interface{})
		err = json.Unmarshal(data, &result)
		return result, err
	}
}

// ToStruct 将map或JSON字符串转换为结构体
func ToStruct(value interface{}, result interface{}) error {
	switch v := value.(type) {
	case string:
		return json.Unmarshal([]byte(v), result)
	default:
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		return json.Unmarshal(data, result)
	}
}

// ToTime 将任意类型转换为time.Time
func ToTime(value interface{}, layouts ...string) (time.Time, error) {
	switch v := value.(type) {
	case time.Time:
		return v, nil
	case int, int64, uint, uint64:
		return time.Unix(reflect.ValueOf(value).Int(), 0), nil
	case string:
		if len(layouts) > 0 {
			for _, layout := range layouts {
				t, err := time.Parse(layout, v)
				if err == nil {
					return t, nil
				}
			}
		}
		// 尝试常见的时间格式
		formats := []string{
			time.RFC3339,
			"2006-01-02 15:04:05",
			"2006-01-02",
			"2006/01/02 15:04:05",
			"2006/01/02",
			"01/02/2006",
			"01/02/2006 15:04:05",
		}
		for _, format := range formats {
			t, err := time.Parse(format, v)
			if err == nil {
				return t, nil
			}
		}
		return time.Time{}, fmt.Errorf("cannot parse time: %s", v)
	default:
		return time.Time{}, fmt.Errorf("cannot convert %T to time.Time", value)
	}
}

// StructToMap 将结构体转换为map
func StructToMap(obj interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	result := make(map[string]interface{})
	err = json.Unmarshal(data, &result)
	return result, err
}

// MapToStruct 将map转换为结构体
func MapToStruct(m map[string]interface{}, obj interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, obj)
}
