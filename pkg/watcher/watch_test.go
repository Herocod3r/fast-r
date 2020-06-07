package watcher

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestListener_Listen_WithUniform(t *testing.T) {
	ctx, cac := context.WithCancel(context.Background())
	listn := NewListiner(cac)
	dataBytes := []int64{524288, 524288, 524288, 524288, 524288, 524288, 524288}

	for _, val := range dataBytes {
		listn.Listen(val)
		<-time.After(time.Millisecond * 5)
	}

	<-ctx.Done()
	fmt.Println("Completed")
}

func TestListener_Listen_WithNonUniform(t *testing.T) {
	ctx, cac := context.WithCancel(context.Background())
	listn := NewListiner(cac)
	dataBytes := []int64{524288, 524288, 524288, 524288, 524288, 524288, 524288}

	for _, val := range dataBytes {
		listn.Listen(val)
		<-time.After(time.Millisecond * 100)
	}

	select {
	case <-ctx.Done():
		t.Fail()
	case <-time.After(2 * time.Second):
	}

}
