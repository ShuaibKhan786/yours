package workerpool

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
)

type Mock struct {
	Buff    *bytes.Buffer
	ID      string
	HaveErr bool
}

func (m Mock) Execute() error {
	if m.HaveErr {
		return errors.New("mock error")
	}
	fmt.Fprintf(m.Buff, m.ID)
	return nil
}

func TestQueue(t *testing.T) {
	buffer := make([]byte, 0)
	buff := bytes.NewBuffer(buffer)

	mockIDs := []string{"mock1", "mock2", "mock1", "mock1", "mock2", "mock2"}

	q := newQueue()
	for _, id := range mockIDs {
		m := &Mock{
			Buff: buff,
			ID:   id,
		}
		q.enqueue(m)
	}

	expectedLength := len(mockIDs)
	if expectedLength != q.len() {
		t.Errorf("Queue Length, Expected: %d but Got:%d", expectedLength, q.length)
	}

	for range mockIDs {
		m := q.dequeue()
		m.Execute()
	}

	if q.length != 0 {
		t.Errorf("Queue Length, Expected: 0 but Got:%d", q.length)
	}

	expected := strings.Join(mockIDs, "")
	got := buff.String()
	if expected != got {
		t.Errorf("Queue Result, Expected: %s but Got:%s", expected, got)
	}

	if q.dequeue() != nil {
		t.Errorf("Queue, Expected a nil result when dequeue")
	}

	q.enqueue(
		&Mock{
			Buff:    buff,
			ID:      "mock error",
			HaveErr: true,
		},
	)
	task := q.dequeue()
	err := task.Execute()

	if err == nil {
		t.Errorf("Expected an error but not none")
	}

	if err.Error() != "mock error" {
		t.Errorf("Expected an error: mock error but got: %s", err.Error())
	}
}
