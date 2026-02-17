package clock

import "time"

type Clock interface {
	Now() time.Time
}

type realClock struct{}

func New() Clock {
	return &realClock{}
}

func (c *realClock) Now() time.Time {
	return time.Now().UTC()
}

type MockClock struct {
	FixedTime time.Time
}

func NewMock(t time.Time) *MockClock {
	return &MockClock{FixedTime: t}
}

func (c *MockClock) Now() time.Time {
	return c.FixedTime
}

func (c *MockClock) Set(t time.Time) {
	c.FixedTime = t
}
