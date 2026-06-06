package san_config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// LoadFromEnv populates a struct from environment variables using the `env` struct tag.
//
// Supported struct tag format:
//
//	`env:"ENV_VAR_NAME"`            — maps field to ENV_VAR_NAME
//	`env:"ENV_VAR_NAME,required"`   — returns error if ENV_VAR_NAME is not set
//
// Supported field types: string, bool, int, int8, int16, int32, int64,
// uint, uint8, uint16, uint32, uint64, float32, float64, and time.Duration.
//
// Fields without an `env` tag are skipped. Nested structs are traversed recursively.
// The data argument must be a non-nil pointer to a struct.
func LoadFromEnv(data any) error {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("san_config: LoadFromEnv requires a non-nil pointer to a struct")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("san_config: LoadFromEnv requires a pointer to a struct, got pointer to %s", v.Kind())
	}

	return loadStruct(v)
}

func loadStruct(v reflect.Value) error {
	t := v.Type()

	for i := range t.NumField() {
		field := t.Field(i)
		fieldVal := v.Field(i)

		// Recurse into embedded/nested structs.
		if field.Type.Kind() == reflect.Struct && field.Type != reflect.TypeOf(time.Duration(0)) {
			if err := loadStruct(fieldVal); err != nil {
				return err
			}
			continue
		}

		tag := field.Tag.Get("env")
		if tag == "" {
			continue
		}

		envKey, required := parseTag(tag)
		envVal, ok := os.LookupEnv(envKey)

		if !ok {
			if required {
				return fmt.Errorf("san_config: required environment variable %q is not set (field %s)", envKey, field.Name)
			}
			continue
		}

		if err := setField(fieldVal, envVal); err != nil {
			return fmt.Errorf("san_config: failed to set field %s from env %q: %w", field.Name, envKey, err)
		}
	}

	return nil
}

func parseTag(tag string) (key string, required bool) {
	parts := strings.Split(tag, ",")
	key = strings.TrimSpace(parts[0])
	for _, opt := range parts[1:] {
		if strings.TrimSpace(opt) == "required" {
			required = true
		}
	}
	return
}

func setField(fieldVal reflect.Value, envVal string) error {
	// Handle time.Duration specially.
	if fieldVal.Type() == reflect.TypeOf(time.Duration(0)) {
		d, err := time.ParseDuration(envVal)
		if err != nil {
			return fmt.Errorf("invalid duration %q: %w", envVal, err)
		}
		fieldVal.Set(reflect.ValueOf(d))
		return nil
	}

	switch fieldVal.Kind() {
	case reflect.String:
		fieldVal.SetString(envVal)

	case reflect.Bool:
		b, err := strconv.ParseBool(envVal)
		if err != nil {
			return fmt.Errorf("invalid bool %q: %w", envVal, err)
		}
		fieldVal.SetBool(b)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(envVal, 10, fieldVal.Type().Bits())
		if err != nil {
			return fmt.Errorf("invalid int %q: %w", envVal, err)
		}
		fieldVal.SetInt(n)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(envVal, 10, fieldVal.Type().Bits())
		if err != nil {
			return fmt.Errorf("invalid uint %q: %w", envVal, err)
		}
		fieldVal.SetUint(n)

	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(envVal, fieldVal.Type().Bits())
		if err != nil {
			return fmt.Errorf("invalid float %q: %w", envVal, err)
		}
		fieldVal.SetFloat(f)

	default:
		return fmt.Errorf("unsupported field type %s", fieldVal.Kind())
	}

	return nil
}
