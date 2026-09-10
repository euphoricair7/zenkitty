package kitty

import (
	"time"

	"github.com/euphoricair7/zenkitty/internal/timer"
)

type Mood int

const (
	Idle Mood = iota
	Settling
	Present
	Paused
	Napping
	Winding
	Done
)

const (
	settlingFor   = 30 * time.Second
	nappingAfter  = 45 * time.Second
	windingWithin = 2 * time.Minute
)

func MoodOf(s *timer.Session) Mood {
	switch s.Phase() {
	case timer.Idle:
		return Idle
	case timer.Done:
		return Done
	case timer.Paused:
		if s.PauseElapsed() >= nappingAfter {
			return Napping
		}
		return Paused
	case timer.Running:
		if rem, ok := s.Remaining(); ok && rem > 0 && rem <= windingWithin {
			return Winding
		}
		if s.Elapsed() < settlingFor {
			return Settling
		}
		return Present
	default:
		return Idle
	}
}
