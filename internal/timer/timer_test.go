package timer

import (
	"testing"
	"time"
)

func TestStartPauseResume(t *testing.T) {
	s := New()
	t0 := time.Unix(0, 0).UTC()

	s.Start(t0)
	s.Tick(t0.Add(10 * time.Second))
	if s.Phase() != Running {
		t.Fatalf("phase %v", s.Phase())
	}
	if s.Elapsed() != 10*time.Second {
		t.Fatalf("elapsed %v", s.Elapsed())
	}

	s.Pause(t0.Add(10 * time.Second))
	s.Tick(t0.Add(15 * time.Second))
	if s.Phase() != Paused {
		t.Fatalf("phase %v", s.Phase())
	}
	if s.Elapsed() != 10*time.Second {
		t.Fatalf("elapsed during pause %v", s.Elapsed())
	}
	if s.PauseElapsed() != 5*time.Second {
		t.Fatalf("pause elapsed %v", s.PauseElapsed())
	}

	s.Resume(t0.Add(15 * time.Second))
	s.Tick(t0.Add(18 * time.Second))
	if s.Phase() != Running {
		t.Fatalf("phase %v", s.Phase())
	}
	if s.Elapsed() != 13*time.Second {
		t.Fatalf("elapsed after resume %v", s.Elapsed())
	}
}

func TestSoftUntilAndStayLonger(t *testing.T) {
	s := New()
	t0 := time.Unix(0, 0).UTC()
	s.SetTotal(5 * time.Second)
	s.Start(t0)
	s.Tick(t0.Add(5 * time.Second))

	if s.Phase() != Done {
		t.Fatalf("phase %v", s.Phase())
	}
	if s.Elapsed() != 5*time.Second {
		t.Fatalf("elapsed %v", s.Elapsed())
	}

	s.StayLonger(t0.Add(8 * time.Second))
	if s.Phase() != Running {
		t.Fatalf("phase %v", s.Phase())
	}
	if s.HasDuration() {
		t.Fatal("expected open sit")
	}
	s.Tick(t0.Add(10 * time.Second))
	if s.Elapsed() != 7*time.Second {
		t.Fatalf("elapsed after staying %v", s.Elapsed())
	}
}

func TestNudgeOpenSit(t *testing.T) {
	s := New()
	t0 := time.Unix(0, 0).UTC()
	s.Start(t0)
	s.Tick(t0.Add(2 * time.Minute))
	s.Nudge(5 * time.Minute)
	rem, ok := s.Remaining()
	if !ok {
		t.Fatal("expected duration")
	}
	if rem != 5*time.Minute {
		t.Fatalf("remaining %v", rem)
	}
}

func TestNudgeClampsAboveElapsed(t *testing.T) {
	s := New()
	t0 := time.Unix(0, 0).UTC()
	s.SetTotal(10 * time.Minute)
	s.Start(t0)
	s.Tick(t0.Add(8 * time.Minute))
	s.Nudge(-10 * time.Minute)
	rem, ok := s.Remaining()
	if !ok {
		t.Fatal("expected duration")
	}
	if rem != time.Minute {
		t.Fatalf("remaining %v", rem)
	}
}

func TestStopKeepsDuration(t *testing.T) {
	s := New()
	t0 := time.Unix(0, 0).UTC()
	s.SetTotal(25 * time.Minute)
	s.Start(t0)
	s.Tick(t0.Add(time.Minute))
	s.Stop()
	if s.Phase() != Idle {
		t.Fatalf("phase %v", s.Phase())
	}
	if s.Elapsed() != 0 {
		t.Fatalf("elapsed %v", s.Elapsed())
	}
	if !s.HasDuration() || s.Duration() != 25*time.Minute {
		t.Fatalf("duration %v %v", s.HasDuration(), s.Duration())
	}
}

func TestSetRemainingFromMidSit(t *testing.T) {
	s := New()
	t0 := time.Unix(0, 0).UTC()
	s.Start(t0)
	s.Tick(t0.Add(10 * time.Minute))
	s.SetRemaining(15 * time.Minute)
	if s.Duration() != 25*time.Minute {
		t.Fatalf("duration %v", s.Duration())
	}
}

func TestFormat(t *testing.T) {
	if got := Format(4*time.Minute + 18*time.Second); got != "4:18" {
		t.Fatalf("got %q", got)
	}
	if got := Format(time.Hour + 2*time.Minute + 3*time.Second); got != "1:02:03" {
		t.Fatalf("got %q", got)
	}
}
