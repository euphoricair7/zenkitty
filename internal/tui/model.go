package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/euphoricair7/zenkitty/internal/copy"
	"github.com/euphoricair7/zenkitty/internal/kitty"
	"github.com/euphoricair7/zenkitty/internal/timer"
)

type screen int

const (
	screenSit screen = iota
	screenDuration
	screenConfirm
	screenHelp
)

type clockTickMsg time.Time
type animTickMsg time.Time

type durationChoice int

const (
	durLittle durationChoice = iota
	durStretch
	durDeep
	durOpen
)

var durationMins = map[durationChoice]time.Duration{
	durLittle:  15 * time.Minute,
	durStretch: 25 * time.Minute,
	durDeep:    50 * time.Minute,
	durOpen:    0,
}

type Model struct {
	width      int
	height     int
	session    *timer.Session
	screen     screen
	frame      int
	line       string
	mood       kitty.Mood
	picker     *copy.Picker
	farewell   string
	durCursor  durationChoice
	durInput   string
	listCursor int
	murmur     int
}

func New() *Model {
	s := timer.New()
	p := copy.New()
	return &Model{
		session:  s,
		picker:   p,
		line:     p.Line(kitty.Idle),
		mood:     kitty.Idle,
		farewell: copy.Leave(),
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(tickClock(), tickAnim())
}

func (m *Model) Goodbye() string {
	if m.farewell == "" {
		return copy.Leave()
	}
	return m.farewell
}

func tickClock() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return clockTickMsg(t)
	})
}

func tickAnim() tea.Cmd {
	return tea.Tick(520*time.Millisecond, func(t time.Time) tea.Msg {
		return animTickMsg(t)
	})
}

func (m *Model) refreshCopy(force bool) {
	mood := kitty.MoodOf(m.session)
	if force || mood != m.mood {
		if mood == kitty.Done && m.mood != kitty.Done {
			m.listCursor = 0
		}
		m.mood = mood
		m.line = m.picker.Line(mood)
		m.murmur = 0
		return
	}
	if mood == kitty.Present {
		m.murmur++
		if m.murmur >= 16 {
			m.murmur = 0
			m.line = m.picker.Line(mood)
		}
	}
}
