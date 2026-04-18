package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Style definitions using lipgloss
var (
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#1A1A2E")).
			Bold(true).
			Width(60)

	providerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7B68EE")). // Medium slate blue
			Bold(true)

	agentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#98D8C8")) // Seafoam green

	categoryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F7DC6F")) // Yellow

	modelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")). // Gray
			Italic(true)

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")) // Coral red

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")) // Gray
)

// View renders the TUI with the config tree.
func View(m AppModel) string {
	var s strings.Builder

	// Header
	s.WriteString(headerStyle.Render("  omo-switch - Config Editor  "))
	s.WriteString("\n\n")

	// Calculate visible range
	start := 0
	end := len(m.FlatList)
	if end > m.Viewport {
		// Center cursor in viewport if possible
		start = m.Cursor - m.Viewport/2
		if start < 0 {
			start = 0
		}
		end = start + m.Viewport
		if end > len(m.FlatList) {
			end = len(m.FlatList)
			start = end - m.Viewport
			if start < 0 {
				start = 0
			}
		}
	}

	// Render visible nodes
	visibleNodes := m.FlatList
	if len(visibleNodes) == 0 {
		s.WriteString("  No agents or categories configured.\n")
		s.WriteString("  Edit ~/.config/opencode/oh-my-openagent.json to add providers.\n")
	} else {
		for i, node := range visibleNodes[start:end] {
			actualIdx := start + i
			prefix := "  "
			cursor := "  "

			// Cursor indicator
			if actualIdx == m.Cursor {
				cursor = cursorStyle.Render(" >")
			}

			// Indent based on depth
			if node.Depth > 0 {
				prefix = strings.Repeat("  ", node.Depth) + prefix
			}

			// Type icon and styling
			switch node.Type {
			case "provider":
				s.WriteString(prefix + cursor + providerStyle.Render("[P]") + " " + providerStyle.Render(node.Label))
			case "agent":
				s.WriteString(prefix + cursor + agentStyle.Render("[A]") + " " + agentStyle.Render(node.Label) + " " + modelStyle.Render(node.Model))
			case "category":
				s.WriteString(prefix + cursor + categoryStyle.Render("[C]") + " " + categoryStyle.Render(node.Label) + " " + modelStyle.Render(node.Model))
			}

			s.WriteString("\n")
		}
	}

	// Stats
	agentCount := 0
	categoryCount := 0
	for _, node := range m.FlatList {
		if node.Type == "agent" {
			agentCount++
		} else if node.Type == "category" {
			categoryCount++
		}
	}

	s.WriteString("\n")
	s.WriteString(footerStyle.Render(fmt.Sprintf("  Agents: %d | Categories: %d", agentCount, categoryCount)))
	s.WriteString("\n")
	s.WriteString(footerStyle.Render("  h/j/k/l: navigate | Enter: edit | Space: select"))
	s.WriteString("\n")

	return s.String()
}
