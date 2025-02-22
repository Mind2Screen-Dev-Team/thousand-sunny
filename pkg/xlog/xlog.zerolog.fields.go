package xlog

import (
	"fmt"
	"net"
	"reflect"
	"time"

	"github.com/rs/zerolog"
)

type LogArrayMarshaler interface {
	MarshalZerologArray(*zerolog.Array)
}

type LogObjectMarshaler interface {
	MarshalZerologObject(*zerolog.Event)
}

// Generic function to handle pointer dereferencing
func deref(v any) any {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr && !rv.IsNil() {
		return rv.Elem().Interface()
	}
	return v
}

func AnyFieldsToContext(c zerolog.Context, m map[string]any) zerolog.Context {
	for key, value := range m {
		c = AnyFieldToZeroLogContext(c, key, value)
	}
	return c
}

func AnyFieldToZeroLogContext(c zerolog.Context, key string, value any) zerolog.Context {
	v := deref(value)

	switch v := v.(type) {
	case string:
		c = c.Str(key, v)
	case []string:
		c = c.Strs(key, v)

	// Numeric types
	case int, int8, int16, int32, int64:
		c = c.Int64(key, reflect.ValueOf(v).Int())
	case uint, uint8, uint16, uint32, uint64:
		c = c.Uint64(key, reflect.ValueOf(v).Uint())
	case float32:
		c = c.Float32(key, v)
	case float64:
		c = c.Float64(key, v)

	// Boolean types
	case bool:
		c = c.Bool(key, v)
	case []bool:
		c = c.Bools(key, v)

	// Time-related types
	case time.Time:
		c = c.Time(key, v)
	case []time.Time:
		c = c.Times(key, v)
	case time.Duration:
		c = c.Dur(key, v)
	case []time.Duration:
		c = c.Durs(key, v)

	// Byte-related types
	case []byte:
		c = c.Bytes(key, v)

	// Errors
	case error:
		c = c.AnErr(key, v)
	case []error:
		c = c.Errs(key, v)

	// Network types
	case net.IP:
		c = c.IPAddr(key, v)
	case net.IPNet:
		c = c.IPPrefix(key, v)
	case net.HardwareAddr:
		c = c.MACAddr(key, v)

	// Custom types
	case fmt.Stringer:
		c = c.Stringer(key, v)
	case LogArrayMarshaler:
		c = c.Array(key, v)
	case LogObjectMarshaler:
		c = c.Object(key, v)

	default:
		c = c.Interface(key, v)
	}

	return c
}

func AnyFieldsToEvent(c *zerolog.Event, m map[string]any) *zerolog.Event {
	for key, value := range m {
		c = AnyFieldToZeroLogEvent(c, key, value)
	}
	return c
}

func AnyFieldToZeroLogEvent(c *zerolog.Event, key string, value any) *zerolog.Event {
	v := deref(value)

	switch v := v.(type) {
	case string:
		c = c.Str(key, v)
	case []string:
		c = c.Strs(key, v)

	// Numeric types
	case int, int8, int16, int32, int64:
		c = c.Int64(key, reflect.ValueOf(v).Int())
	case uint, uint8, uint16, uint32, uint64:
		c = c.Uint64(key, reflect.ValueOf(v).Uint())
	case float32:
		c = c.Float32(key, v)
	case float64:
		c = c.Float64(key, v)

	// Boolean types
	case bool:
		c = c.Bool(key, v)
	case []bool:
		c = c.Bools(key, v)

	// Time-related types
	case time.Time:
		c = c.Time(key, v)
	case []time.Time:
		c = c.Times(key, v)
	case time.Duration:
		c = c.Dur(key, v)
	case []time.Duration:
		c = c.Durs(key, v)

	// Byte-related types
	case []byte:
		c = c.Bytes(key, v)

	// Errors
	case error:
		c = c.AnErr(key, v)
	case []error:
		c = c.Errs(key, v)

	// Network types
	case net.IP:
		c = c.IPAddr(key, v)
	case net.IPNet:
		c = c.IPPrefix(key, v)
	case net.HardwareAddr:
		c = c.MACAddr(key, v)

	// Custom types
	case fmt.Stringer:
		c = c.Stringer(key, v)
	case LogArrayMarshaler:
		c = c.Array(key, v)
	case LogObjectMarshaler:
		c = c.Object(key, v)

	default:
		c = c.Interface(key, v)
	}

	return c
}
