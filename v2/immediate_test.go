package koncurrent

import (
	"context"
	"testing"
)

func TestImmediateExecutor_Execute(t *testing.T) {
	underTest := NewImmediateExecutor()
	ret := underTest.Execute(context.Background(), func(ctx context.Context) error {
		return nil
	})
	err := ret.Get()
	assertNil(t, err)
}
