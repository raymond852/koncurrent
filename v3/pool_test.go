package koncurrent

import (
	"context"
	"testing"
)

func TestPoolExecutor_Execute(t *testing.T) {
	underTest := NewPoolExecutor(10, 10)
	ret := underTest.Execute(context.Background(), func(ctx context.Context) error {
		return nil
	})
	resultErr := <- ret
	assertNil(t, resultErr)
}
