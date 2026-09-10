package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/euphoricair7/zenkitty/internal/tui"
)

func main() {
	p := tea.NewProgram(tui.New(), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if fm, ok := m.(*tui.Model); ok {
		bye := lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#8B5E5A", Dark: "#D4B4AE"}).
			Italic(true).
			Render(fm.Goodbye())
		fmt.Println()
		fmt.Println("  " + bye)
		fmt.Println()
	}
}
