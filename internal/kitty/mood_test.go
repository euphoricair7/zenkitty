package kitty

import (
	"testing"
	"time"

	"github.com/euphoricair7/zenkitty/internal/timer"
)

func TestMoods(t *testing.T) {
	t0 := time.Unix(0, 0).UTC()

	idle := timer.New()
	if MoodOf(idle) != Idle {
		t.Fatal("idle")
	}

	s := timer.New()
	s.Start(t0)
	s.Tick(t0.Add(10 * time.Second))
	if MoodOf(s) != Settling {
		t.Fatalf("got %v", MoodOf(s))
	}

	s.Tick(t0.Add(60 * time.Second))
	if MoodOf(s) != Present {
		t.Fatalf("got %v", MoodOf(s))
	}

	s.Pause(t0.Add(60 * time.Second))
	s.Tick(t0.Add(70 * time.Second))
	if MoodOf(s) != Paused {
		t.Fatalf("got %v", MoodOf(s))
	}

	s.Tick(t0.Add(120 * time.Second))
	if MoodOf(s) != Napping {
		t.Fatalf("got %v", MoodOf(s))
	}

	w := timer.New()
	w.SetTotal(10 * time.Minute)
	w.Start(t0)
	w.Tick(t0.Add(9*time.Minute + 30*time.Second))
	if MoodOf(w) != Winding {
		t.Fatalf("got %v", MoodOf(w))
	}

	w.Tick(t0.Add(10 * time.Minute))
	if MoodOf(w) != Done {
		t.Fatalf("got %v", MoodOf(w))
	}
}

func TestFrameStableHeight(t *testing.T) {
	h := stringsCount(Frame(Idle, 0))
	for mood := Idle; mood <= Done; mood++ {
		for i := 0; i < 8; i++ {
			if stringsCount(Frame(mood, i)) != h {
				t.Fatalf("mood %d frame %d height changed", mood, i)
			}
		}
	}
}

func stringsCount(s string) int {
	n := 1
	for _, r := range s {
		if r == '\n' {
			n++
		}
	}
	return n
}
