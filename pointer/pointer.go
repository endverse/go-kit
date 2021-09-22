package pointer

import "time"

// Int32Ptr returns a pointer to an int32
func Int32Ptr(i int32) *int32 {
	return &i
}

// Int64Ptr returns a pointer to an int64
func Int64Ptr(i int64) *int64 {
	return &i
}

// IntPtr returns a pointer to an int
func IntPtr(i int) *int {
	return &i
}

// Int32PtrDerefOr dereference the int32 ptr and returns it if not nil,
// else returns def.
func Int32PtrDerefOr(ptr *int32, def int32) int32 {
	if ptr != nil {
		return *ptr
	}
	return def
}

// BoolPtr returns a pointer to a bool
func BoolPtr(b bool) *bool {
	return &b
}

// TimePtr returns a pointer to a time.Time
func TimePtr(t time.Time) *time.Time {
	return &t
}

// StringPtr returns a pointer to the passed string.
func StringPtr(s string) *string {
	return &s
}

// Float32Ptr returns a pointer to the passed float32.
func Float32Ptr(i float32) *float32 {
	return &i
}

// Float64Ptr returns a pointer to the passed float64.
func Float64Ptr(i float64) *float64 {
	return &i
}
