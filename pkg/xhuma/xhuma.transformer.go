package xhuma

import (
	"reflect"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
	"github.com/danielgtaylor/huma/v2"
)

func TransformerTraceIdSetter(ctx huma.Context, status string, v any) (any, error) {
	tid, ok := ctx.Context().Value(xlog.XLOG_REQ_TRACE_ID_CTX_KEY).(string)
	if !ok {
		return v, nil
	}

	var (
		rv = reflect.ValueOf(v)
		rt = reflect.TypeOf(v)
	)

	// Case 1: Struct by value (not pointer)
	if rv.Kind() == reflect.Struct {
		ptr := reflect.New(rt) // *T
		ptr.Elem().Set(rv)     // copy original into pointer
		field := ptr.Elem().FieldByName("TraceID")
		if field.IsValid() && field.CanSet() && field.Kind() == reflect.String {
			field.SetString(tid)
		}
		return ptr.Elem().Interface(), nil // return updated value
	}

	// Case 2: Pointer to struct
	if rv.Kind() == reflect.Ptr && rv.Elem().Kind() == reflect.Struct {
		elem := rv.Elem()
		field := elem.FieldByName("TraceID")
		if field.IsValid() && field.CanSet() && field.Kind() == reflect.String {
			field.SetString(tid)
		}
		return v, nil // mutation is in-place
	}

	// Other types — do nothing
	return v, nil
}
