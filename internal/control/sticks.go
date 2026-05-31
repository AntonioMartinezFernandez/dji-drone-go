package control

import "sync"

type Sticks struct {
	mu sync.Mutex
	lr int // left/right   (-100 left  … +100 right)
	fb int // fwd/backward (-100 back  … +100 forward)
	ud int // up/down      (-100 down  … +100 up)
	y  int // yaw          (-100 ccw   … +100 cw)
}

func (s *Sticks) Get() (lr, fb, ud, y int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lr, s.fb, s.ud, s.y
}

func (s *Sticks) Set(lr, fb, ud, y *int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if lr != nil {
		s.lr = *lr
	}
	if fb != nil {
		s.fb = *fb
	}
	if ud != nil {
		s.ud = *ud
	}
	if y != nil {
		s.y = *y
	}
}
