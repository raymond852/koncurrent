package koncurrent

import (
	"context"
	"testing"
)

func TestPoolExecutor_Execute(t *testing.T) {
	underTest := NewPoolExecutor(10, 10)
	ret := make(chan TaskResult)
	underTest.Execute(context.Background(), func(ctx context.Context) error {
		return nil
	}, 0, ret)
	result := <- ret
	assertNil(t, result.err)
}
