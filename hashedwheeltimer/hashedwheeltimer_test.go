package hwtimer

import (
	"fmt"
	"testing"
	"time"
)

func TestNewTimer(t *testing.T) {
	resultC := make(chan time.Time)
	tickDuration := 300 * time.Millisecond

	hw := NewTimer(100, 4)

	hw.AfterFunc(tickDuration, func() {
		resultC <- time.Now().UTC()
	})

	start := time.Now().UTC()

	got := (<-resultC).Truncate(time.Millisecond)
	want := start.Add(tickDuration).Truncate(time.Millisecond)

	//fmt.Println(got.Sub(want))
	fmt.Println(want)
	fmt.Println(got)

	delta := 5 * time.Millisecond
	if got.Before(want) || got.After(want.Add(delta)) {
		t.Fatalf("expected [%s, %s], but got %s\n", want, want.Add(delta), got)
	}

}

func TestCancelTask(t *testing.T) {
	//hw := NewTimer(100, 10)

}
