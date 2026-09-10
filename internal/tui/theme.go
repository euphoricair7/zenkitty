package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	bg          = lipgloss.AdaptiveColor{Light: "#F4E8D0", Dark: "#1A1410"}
	kittyColor  = lipgloss.AdaptiveColor{Light: "#5C3D32", Dark: "#F3D3B8"}
	timeColor   = lipgloss.AdaptiveColor{Light: "#4F6F52", Dark: "#A8C09A"}
	copyColor   = lipgloss.AdaptiveColor{Light: "#8B5E5A", Dark: "#D4B4AE"}
	mutedColor  = lipgloss.AdaptiveColor{Light: "#A89888", Dark: "#6B5E55"}
	dimColor    = lipgloss.AdaptiveColor{Light: "#C4B8A8", Dark: "#4A403A"}
	selectColor = lipgloss.AdaptiveColor{Light: "#6B3F36", Dark: "#F0D0C4"}
	sparkColor  = lipgloss.AdaptiveColor{Light: "#8A7340", Dark: "#E6D39A"}
	grassColor  = lipgloss.AdaptiveColor{Light: "#4F6F52", Dark: "#7E9A78"}
)

func screenStyle() lipgloss.Style {
	return lipgloss.NewStyle().Background(bg)
}

func kittyStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(kittyColor).Background(bg)
}

func timeStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(timeColor).Background(bg).Faint(true)
}

func copyStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(copyColor).Background(bg).Italic(true)
}

func mutedStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(mutedColor).Background(bg)
}

func dimStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(dimColor).Background(bg)
}

func selectStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(selectColor).Background(bg)
}

func sparkStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(sparkColor).Background(bg)
}

func grassStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(grassColor).Background(bg).Faint(true)
}

func footerStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(dimColor).Background(bg)
}
