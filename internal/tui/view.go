package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/euphoricair7/zenkitty/internal/kitty"
	"github.com/euphoricair7/zenkitty/internal/timer"
)

func (m *Model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	scene := kitty.Compose(m.mood, m.frame)
	cat := join(
		sparkStyle().Render(scene.Sky),
		kittyStyle().Render(scene.Cat),
		grassStyle().Render(scene.Ground),
	)
	body := m.body(cat)
	footer := footerStyle().Render(m.footer())

	contentH := m.height - 1
	if contentH < 1 {
		contentH = m.height
	}

	placed := lipgloss.Place(
		m.width,
		contentH,
		lipgloss.Center,
		lipgloss.Center,
		body,
		lipgloss.WithWhitespaceBackground(bg),
	)
	foot := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, footer, lipgloss.WithWhitespaceBackground(bg))

	out := lipgloss.JoinVertical(lipgloss.Left, placed, foot)
	return screenStyle().Width(m.width).Height(m.height).Render(out)
}

func (m *Model) body(cat string) string {
	switch m.screen {
	case screenDuration:
		return join(cat, m.durationView())
	case screenConfirm:
		return join(cat, m.confirmView())
	case screenHelp:
		return join(cat, m.helpView())
	default:
		return join(cat, m.sitView())
	}
}

func join(parts ...string) string {
	return lipgloss.JoinVertical(lipgloss.Center, parts...)
}

func (m *Model) sitView() string {
	clock := timeStyle().Render(m.clock())
	line := m.line
	if strings.TrimSpace(line) == "" {
		line = " "
	}
	words := copyStyle().Render(line)

	bits := []string{""}
	if m.session.Phase() != timer.Idle || m.session.HasDuration() {
		bits = append(bits, clock, "")
	}
	bits = append(bits, words)

	if m.session.Phase() == timer.Done {
		bits = append(bits, "", m.doneChoices())
	} else if trail := m.trail(); trail != "" {
		bits = append(bits, "", dimStyle().Render(trail))
	}

	return lipgloss.JoinVertical(lipgloss.Center, bits...)
}

func (m *Model) clock() string {
	s := m.session
	if s.Phase() == timer.Idle {
		if s.HasDuration() {
			return timer.Format(s.Duration())
		}
		return " "
	}
	if rem, ok := s.Remaining(); ok && s.Phase() != timer.Done {
		return timer.Format(rem)
	}
	return timer.Format(s.Elapsed())
}

func (m *Model) trail() string {
	rem, ok := m.session.Remaining()
	if !ok || m.session.Phase() == timer.Idle {
		return ""
	}
	total := m.session.Duration()
	if total <= 0 {
		return ""
	}
	const n = 7
	filled := int(float64(n) * (1 - float64(rem)/float64(total)))
	if filled < 0 {
		filled = 0
	}
	if filled > n-1 {
		filled = n - 1
	}
	dots := make([]string, n)
	for i := range dots {
		if i == filled {
			dots[i] = "•"
		} else {
			dots[i] = "·"
		}
	}
	return strings.Join(dots, " ")
}

func (m *Model) doneChoices() string {
	stay := "stay a little longer?"
	done := "all done"
	if m.listCursor == 0 {
		stay = selectStyle().Render("› " + stay)
		done = mutedStyle().Render(done)
	} else {
		stay = mutedStyle().Render(stay)
		done = selectStyle().Render("› " + done)
	}
	return stay + mutedStyle().Render("   ·   ") + done
}

func (m *Model) durationView() string {
	title := copyStyle().Render("how long do you want to sit?")
	if m.session.Sitting() {
		title = copyStyle().Render("change the until?")
	}

	opts := []string{
		"a little while     (~15)",
		"a good stretch     (~25)",
		"a deep dive        (~50)",
		"until we're done   (open)",
	}
	var rows []string
	for i, opt := range opts {
		if durationChoice(i) == m.durCursor && m.durInput == "" {
			rows = append(rows, selectStyle().Render("› "+opt))
		} else {
			rows = append(rows, mutedStyle().Render("  "+opt))
		}
	}

	typed := "or type minutes"
	if m.durInput != "" {
		typed = fmt.Sprintf("%s min", m.durInput)
		typed = selectStyle().Render(typed)
	} else {
		typed = dimStyle().Render(typed)
	}

	return lipgloss.JoinVertical(lipgloss.Center, append([]string{
		"",
		title,
		"",
	}, append(rows, "", typed)...)...)
}

func (m *Model) confirmView() string {
	title := copyStyle().Render("leave this one?")
	stay := "stay a little more"
	over := "start over"
	if m.listCursor == 0 {
		stay = selectStyle().Render("› " + stay)
		over = mutedStyle().Render("  " + over)
	} else {
		stay = mutedStyle().Render("  " + stay)
		over = selectStyle().Render("› " + over)
	}
	return lipgloss.JoinVertical(lipgloss.Center, "", title, "", stay, over)
}

func (m *Model) helpView() string {
	rows := []string{
		"space   pause, or come back",
		"enter   start sitting",
		"n       new sit",
		"d       how long",
		"+ -     a little more, a little less",
		"?       this note",
		"q       bye",
	}
	out := make([]string, 0, len(rows)+1)
	out = append(out, "", copyStyle().Render("a little note"))
	for _, r := range rows {
		out = append(out, mutedStyle().Render(r))
	}
	return lipgloss.JoinVertical(lipgloss.Left, out...)
}

func (m *Model) footer() string {
	switch m.screen {
	case screenDuration:
		return "enter sit  ·  esc back"
	case screenConfirm:
		return "enter  ·  esc stay"
	case screenHelp:
		return "esc back"
	}
	if m.session.Phase() == timer.Done {
		return "enter stay  ·  n all done  ·  q bye"
	}
	if m.session.Sitting() {
		return "space pause  ·  n new  ·  d duration  ·  q home"
	}
	return "enter sit  ·  d duration  ·  ? note  ·  q bye"
}
