package timer

import "time"

type Phase int

const (
	Idle Phase = iota
	Running
	Paused
	Done
)

// Session is a flexible sit: open-ended or softly timed, pause without penalty.
type Session struct {
	phase        Phase
	elapsed      time.Duration
	pauseElapsed time.Duration
	lastTick     time.Time
	duration     time.Duration
	hasDuration  bool
}

func New() *Session {
	return &Session{}
}

func (s *Session) Phase() Phase { return s.phase }

func (s *Session) Elapsed() time.Duration { return s.elapsed }

func (s *Session) PauseElapsed() time.Duration { return s.pauseElapsed }

func (s *Session) HasDuration() bool { return s.hasDuration }

func (s *Session) Duration() time.Duration { return s.duration }

func (s *Session) Remaining() (time.Duration, bool) {
	if !s.hasDuration {
		return 0, false
	}
	if s.elapsed >= s.duration {
		return 0, true
	}
	return s.duration - s.elapsed, true
}

func (s *Session) Sitting() bool {
	return s.phase == Running || s.phase == Paused
}

func (s *Session) Tick(now time.Time) {
	if s.lastTick.IsZero() {
		s.lastTick = now
		return
	}
	delta := now.Sub(s.lastTick)
	if delta < 0 {
		delta = 0
	}
	s.lastTick = now
	switch s.phase {
	case Running:
		s.elapsed += delta
		if s.hasDuration && s.elapsed >= s.duration {
			s.elapsed = s.duration
			s.phase = Done
		}
	case Paused:
		s.pauseElapsed += delta
	}
}

func (s *Session) Start(now time.Time) {
	if s.phase == Running {
		return
	}
	if s.phase == Paused {
		s.Resume(now)
		return
	}
	if s.phase == Done {
		return
	}
	s.phase = Running
	s.elapsed = 0
	s.pauseElapsed = 0
	s.lastTick = now
}

func (s *Session) Pause(now time.Time) {
	if s.phase != Running {
		return
	}
	s.Tick(now)
	if s.phase != Running {
		return
	}
	s.phase = Paused
	s.pauseElapsed = 0
}

func (s *Session) Resume(now time.Time) {
	if s.phase != Paused {
		return
	}
	s.Tick(now)
	s.phase = Running
	s.pauseElapsed = 0
	s.lastTick = now
}

func (s *Session) Stop() {
	*s = Session{
		duration:    s.duration,
		hasDuration: s.hasDuration,
	}
}

func (s *Session) StayLonger(now time.Time) {
	if s.phase != Done {
		return
	}
	s.phase = Running
	s.hasDuration = false
	s.duration = 0
	s.lastTick = now
}

func (s *Session) SetOpen() {
	s.hasDuration = false
	s.duration = 0
}

func (s *Session) SetTotal(d time.Duration) {
	if d <= 0 {
		s.SetOpen()
		return
	}
	s.hasDuration = true
	s.duration = d
	if s.phase == Running && s.elapsed >= s.duration {
		s.elapsed = s.duration
		s.phase = Done
	}
}

func (s *Session) SetRemaining(d time.Duration) {
	if d <= 0 {
		s.SetOpen()
		return
	}
	s.hasDuration = true
	s.duration = s.elapsed + d
}

func (s *Session) Nudge(delta time.Duration) {
	if s.phase == Done {
		return
	}
	if !s.hasDuration {
		if delta <= 0 {
			return
		}
		base := s.elapsed
		if s.phase == Idle {
			base = 0
		}
		s.hasDuration = true
		s.duration = base + delta
		return
	}
	s.duration += delta
	minDur := s.elapsed + time.Minute
	if s.phase == Idle {
		minDur = time.Minute
	}
	if s.duration < minDur {
		if s.phase == Idle && s.duration < time.Minute {
			s.SetOpen()
			return
		}
		s.duration = minDur
	}
}
