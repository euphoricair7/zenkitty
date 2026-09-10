package tui

import (
	"strings"
	"testing"
	"time"

	"github.com/euphoricair7/zenkitty/internal/timer"
)

func TestViewShowsKitty(t *testing.T) {
	m := New()
	m.width = 80
	m.height = 24
	got := m.View()
	if !strings.Contains(got, `/\_/\`) {
		t.Fatalf("missing kitty:\n%s", got)
	}
	if !strings.Contains(got, "enter sit") {
		t.Fatalf("missing footer:\n%s", got)
	}
}

func TestDurationAndDoneViews(t *testing.T) {
	m := New()
	m.width = 80
	m.height = 24

	m.screen = screenDuration
	got := m.View()
	if !strings.Contains(got, "how long do you want to sit?") {
		t.Fatalf("duration:\n%s", got)
	}

	now := time.Unix(0, 0).UTC()
	m.screen = screenSit
	m.session.SetTotal(time.Second)
	m.session.Start(now)
	m.session.Tick(now.Add(time.Second))
	if m.session.Phase() != timer.Done {
		t.Fatalf("phase %v", m.session.Phase())
	}
	m.refreshCopy(true)
	got = m.View()
	if !strings.Contains(got, "stay a little longer?") {
		t.Fatalf("done:\n%s", got)
	}
}

func TestSittingViewShowsClock(t *testing.T) {
	m := New()
	m.width = 80
	m.height = 24
	now := time.Unix(0, 0).UTC()
	m.session.Start(now)
	m.session.Tick(now.Add(90 * time.Second))
	m.refreshCopy(true)
	got := m.View()
	if !strings.Contains(got, "1:30") {
		t.Fatalf("missing clock:\n%s", got)
	}
}
