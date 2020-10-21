package koncurrent

import (
	"context"
	"testing"
)

func TestAsyncExecutor_Execute(t *testing.T) {
	underTest := NewAsyncExecutor()
	ret := underTest.Execute(context.Background(), func(ctx context.Context) error {
		return nil
	})
	resultErr := ret.Get()
	assertNil(t, resultErr)
}
