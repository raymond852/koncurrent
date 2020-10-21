package koncurrent

import (
	"reflect"
	"testing"
)

func assertEqual(t *testing.T, a, b interface{}) {
	if a != b {
		t.Errorf("unexpected not equal, %+v != %+v", a, b)
	}
}

func assertNil(t *testing.T, v interface{}) {
	if v != nil && !reflect.ValueOf(v).IsNil() {
		t.Errorf("unexpected not nil value %+v", v)
	}
}

func assertNotNil(t *testing.T, v interface{}) {
	if v == nil || reflect.ValueOf(v).IsNil() {
		t.Error("unexpected nil value")
	}
}

func assertTrue(t *testing.T, v bool) {
	if !v {
		t.Error("unexpected false value")
	}
}
