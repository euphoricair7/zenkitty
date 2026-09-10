package tui

import (
	"strconv"
	"time"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/euphoricair7/zenkitty/internal/copy"
	"github.com/euphoricair7/zenkitty/internal/timer"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case clockTickMsg:
		m.session.Tick(time.Time(msg))
		if m.session.Phase() == timer.Done {
			m.screen = screenSit
		}
		m.refreshCopy(false)
		return m, tickClock()

	case animTickMsg:
		m.frame++
		return m, tickAnim()

	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	now := time.Now()
	m.session.Tick(now)

	if msg.String() == "ctrl+c" {
		m.farewell = copy.Leave()
		return m, tea.Quit
	}

	switch m.screen {
	case screenDuration:
		return m.handleDuration(msg, now)
	case screenConfirm:
		return m.handleConfirm(msg)
	case screenHelp:
		return m.handleHelp(msg)
	default:
		return m.handleSit(msg, now)
	}
}

func (m *Model) handleSit(msg tea.KeyMsg, now time.Time) (tea.Model, tea.Cmd) {
	if m.session.Phase() == timer.Done {
		return m.handleDone(msg, now)
	}

	switch msg.String() {
	case "q":
		m.farewell = copy.Leave()
		return m, tea.Quit
	case "?":
		m.screen = screenHelp
		return m, nil
	case "d":
		m.screen = screenDuration
		m.durInput = ""
		m.durCursor = durStretch
		return m, nil
	case "n":
		if m.session.Sitting() {
			m.screen = screenConfirm
			m.listCursor = 0
			return m, nil
		}
		return m, nil
	case "enter", "s":
		if m.session.Phase() == timer.Idle {
			m.session.Start(now)
			m.refreshCopy(true)
		}
		return m, nil
	case " ":
		switch m.session.Phase() {
		case timer.Idle:
			m.session.Start(now)
			m.refreshCopy(true)
		case timer.Running:
			m.session.Pause(now)
			m.refreshCopy(true)
		case timer.Paused:
			m.session.Resume(now)
			m.refreshCopy(true)
		}
		return m, nil
	case "+", "=":
		m.session.Nudge(5 * time.Minute)
		m.refreshCopy(true)
		return m, nil
	case "-", "_":
		m.session.Nudge(-5 * time.Minute)
		m.refreshCopy(true)
		return m, nil
	}
	return m, nil
}

func (m *Model) handleDone(msg tea.KeyMsg, now time.Time) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		m.farewell = copy.Leave()
		return m, tea.Quit
	case "left", "h", "right", "l", "tab":
		m.listCursor = 1 - m.listCursor
		return m, nil
	case "enter", "s":
		if m.listCursor == 0 {
			m.session.StayLonger(now)
			m.refreshCopy(true)
		} else {
			m.session.Stop()
			m.refreshCopy(true)
		}
		return m, nil
	case "n":
		m.session.Stop()
		m.refreshCopy(true)
		return m, nil
	case " ":
		m.session.StayLonger(now)
		m.refreshCopy(true)
		return m, nil
	}
	return m, nil
}

func (m *Model) handleDuration(msg tea.KeyMsg, now time.Time) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.screen = screenSit
		m.durInput = ""
		return m, nil
	case "q":
		m.farewell = copy.Leave()
		return m, tea.Quit
	case "up", "k":
		if m.durCursor > durLittle {
			m.durCursor--
		}
		return m, nil
	case "down", "j":
		if m.durCursor < durOpen {
			m.durCursor++
		}
		return m, nil
	case "backspace", "ctrl+h":
		if m.durInput != "" {
			m.durInput = m.durInput[:len(m.durInput)-1]
		}
		return m, nil
	case "enter":
		m.applyDuration(now)
		m.screen = screenSit
		m.durInput = ""
		m.refreshCopy(true)
		return m, nil
	}

	if len(msg.Runes) == 1 && unicode.IsDigit(msg.Runes[0]) {
		if len(m.durInput) < 4 {
			m.durInput += string(msg.Runes)
		}
		return m, nil
	}
	return m, nil
}

func (m *Model) applyDuration(now time.Time) {
	var d time.Duration
	if m.durInput != "" {
		mins, err := strconv.Atoi(m.durInput)
		if err != nil || mins < 0 {
			return
		}
		if mins == 0 {
			d = 0
		} else {
			d = time.Duration(mins) * time.Minute
		}
	} else {
		d = durationMins[m.durCursor]
	}

	sitting := m.session.Sitting()
	if sitting {
		if d == 0 {
			m.session.SetOpen()
		} else {
			m.session.SetRemaining(d)
		}
		return
	}
	if d == 0 {
		m.session.SetOpen()
	} else {
		m.session.SetTotal(d)
	}
	if m.session.Phase() == timer.Idle {
		m.session.Start(now)
	}
}

func (m *Model) handleConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.screen = screenSit
		return m, nil
	case "q":
		m.farewell = copy.Leave()
		return m, tea.Quit
	case "up", "k", "down", "j", "tab":
		m.listCursor = 1 - m.listCursor
		return m, nil
	case "enter":
		if m.listCursor == 1 {
			m.session.Stop()
			m.refreshCopy(true)
		}
		m.screen = screenSit
		return m, nil
	}
	return m, nil
}

func (m *Model) handleHelp(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "?", "q", "enter":
		if msg.String() == "q" {
			m.farewell = copy.Leave()
			return m, tea.Quit
		}
		m.screen = screenSit
	}
	return m, nil
}
