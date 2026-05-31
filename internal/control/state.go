package control

import (
	"sync"
	"time"
)

// axis identifies which stick component a key controls.
type axis int

const (
	AxisLR axis = iota
	AxisFB
	AxisUD
	AxisYaw
)

// keyState keeps a per-axis timer so we can zero the axis when the key is released.
type keyState struct {
	mu             sync.Mutex
	timers         map[axis]*time.Timer
	pending        map[axis]int // the active non-zero value for each axis
	keyHoldTimeout time.Duration
}

func NewKeyState(keyHoldTimeout time.Duration) *keyState {
	return &keyState{
		timers:         make(map[axis]*time.Timer),
		pending:        make(map[axis]int),
		keyHoldTimeout: keyHoldTimeout,
	}
}

// Press is called on every key-Press (and auto-repeat) event.
// It sets value on the given axis and arms a release timer.
func (ks *keyState) Press(st *Sticks, ax axis, value int) {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	ks.pending[ax] = value

	// Reset existing timer if any.
	if t, ok := ks.timers[ax]; ok {
		t.Stop()
	}

	// When the timer fires the key has been released: zero out that axis.
	ks.timers[ax] = time.AfterFunc(ks.keyHoldTimeout, func() {
		ks.mu.Lock()
		delete(ks.pending, ax)
		delete(ks.timers, ax)
		ks.mu.Unlock()

		switch ax {
		case AxisLR:
			v := 0
			st.Set(&v, nil, nil, nil)
		case AxisFB:
			v := 0
			st.Set(nil, &v, nil, nil)
		case AxisUD:
			v := 0
			st.Set(nil, nil, &v, nil)
		case AxisYaw:
			v := 0
			st.Set(nil, nil, nil, &v)
		}
	})

	// Apply immediately.
	switch ax {
	case AxisLR:
		st.Set(&value, nil, nil, nil)
	case AxisFB:
		st.Set(nil, &value, nil, nil)
	case AxisUD:
		st.Set(nil, nil, &value, nil)
	case AxisYaw:
		st.Set(nil, nil, nil, &value)
	}
}
