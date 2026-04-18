package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Theme defines the color scheme for the TUI.
type Theme struct {
	// Background colors
	Background lipgloss.Color
	Surface    lipgloss.Color

	// Text colors
	Text    lipgloss.Color
	Subtext lipgloss.Color

	// Node type colors
	Provider lipgloss.Color
	Agent    lipgloss.Color
	Category lipgloss.Color

	// UI element colors
	Cursor   lipgloss.Color
	Selected lipgloss.Color
	Border   lipgloss.Color
	Divider  lipgloss.Color
	Accent   lipgloss.Color
}

// DefaultTheme returns the default color theme.
func DefaultTheme() Theme {
	return Theme{
		Background: lipgloss.Color("#1A1A2E"),
		Surface:    lipgloss.Color("#2D2D44"),
		Text:       lipgloss.Color("#FAFAFA"),
		Subtext:    lipgloss.Color("#6B7280"),
		Provider:   lipgloss.Color("#7B68EE"),
		Agent:      lipgloss.Color("#98D8C8"),
		Category:   lipgloss.Color("#F7DC6F"),
		Cursor:     lipgloss.Color("#FF6B6B"),
		Selected:   lipgloss.Color("#98D8C8"),
		Border:     lipgloss.Color("#1A1A2E"),
		Divider:    lipgloss.Color("#6B7280"),
		Accent:     lipgloss.Color("#FF6B6B"),
	}
}

// Global theme instance
var theme = DefaultTheme()

// Style functions using theme colors
func headerStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Text).
		Background(theme.Background).
		Bold(true)
}

func providerStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Provider).
		Bold(true)
}

func agentStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Agent)
}

func categoryStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Category)
}

func modelStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Subtext).
		Italic(true)
}

func cursorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Cursor)
}

func footerStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Subtext)
}

func selectedStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Selected)
}

func editFieldStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Text).
		Background(theme.Surface).
		Bold(true)
}

func modelPickerStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Text).
		Background(theme.Background)
}

// Badge styles with background highlighting for type badges

func badgeProviderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Text).
		Background(theme.Provider).
		Bold(true)
}

func badgeAgentStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1A1A2E")).
		Background(theme.Agent).
		Bold(true)
}

func badgeCategoryStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1A1A2E")).
		Background(theme.Category).
		Bold(true)
}

func scrollIndicatorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Subtext)
}

func dividerStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Divider)
}

func modeBadgeStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)
}

func hoverRowStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(theme.Surface)
}

func errorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF6B6B"))
}

// buildFooter returns context-sensitive footer hints based on current mode and selection state.
func buildFooter(m AppModel) string {
	if m.ModelPickerMode {
		if m.ProviderPickerMode {
			return "j/k: navigate | Enter: select provider"
		}
		return "j/k: navigate | Enter: select model | Esc: back"
	}
	if m.EditMode {
		return "Enter: confirm | Esc: cancel"
	}
	// Normal mode
	selCount := m.SelectionCount
	if selCount == 0 {
		return "j/k: navigate | Enter: pick model | q: quit"
	}
	return fmt.Sprintf("j/k: navigate | Space: toggle | Enter: apply to %d node%s | x: clear | q: quit", selCount, plural(selCount))
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

// View renders the TUI with the config tree or model picker.
func View(m AppModel) string {
	var s strings.Builder

	// Model picker mode - two level: provider list + model preview
	if m.ModelPickerMode {
		s.WriteString(headerStyle().Render("  omo-switch "))
		s.WriteString(modeBadgeStyle().Render(" ◆ "))
		s.WriteString(headerStyle().Render("Select Model  "))
		s.WriteString("\n\n")

		if m.ProviderPickerMode {
			// Left side: provider list with viewport scrolling
			s.WriteString(modelPickerStyle().Render("  Providers"))
			s.WriteString("\n\n")

			maxVisible := m.Viewport - 6 // account for header + footer
			if maxVisible < 1 {
				maxVisible = 1
			}

			start := 0
			end := len(m.Providers)
			if len(m.Providers) > maxVisible {
				start = m.ProviderIdx - maxVisible/2
				if start < 0 {
					start = 0
				}
				end = start + maxVisible
				if end > len(m.Providers) {
					end = len(m.Providers)
					start = end - maxVisible
					if start < 0 {
						start = 0
					}
				}
			}

			// Scroll indicator: items above viewport
			if start > 0 {
				s.WriteString(scrollIndicatorStyle().Render(fmt.Sprintf("  ↑ %d more", start)))
				s.WriteString("\n")
			}

			for i := start; i < end; i++ {
				provider := m.Providers[i]
				cursor := "  "
				if i == m.ProviderIdx {
					cursor = cursorStyle().Render(" >")
				}
				s.WriteString(fmt.Sprintf("%s %s (%d models)\n", cursor, provider.Name, len(provider.Models)))
			}

			// Scroll indicator: items below viewport
			if end < len(m.Providers) {
				s.WriteString(scrollIndicatorStyle().Render(fmt.Sprintf("  ↓ %d more", len(m.Providers)-end)))
				s.WriteString("\n")
			}

			s.WriteString("\n")
			s.WriteString(footerStyle().Render("  j/k: navigate | Enter: select provider"))
			s.WriteString("\n")

			// Right side: preview of selected provider's models
			if m.ProviderIdx < len(m.Providers) {
				selectedProvider := m.Providers[m.ProviderIdx]
				s.WriteString(modelPickerStyle().Render(fmt.Sprintf("  Models: %s", selectedProvider.Name)))
				s.WriteString("\n\n")

				previewCount := 10
				if previewCount > len(selectedProvider.Models) {
					previewCount = len(selectedProvider.Models)
				}
				for i := 0; i < previewCount; i++ {
					model := selectedProvider.Models[i]
					s.WriteString(fmt.Sprintf("    %s\n", model))
				}
				if len(selectedProvider.Models) > previewCount {
					s.WriteString(fmt.Sprintf("    ... and %d more\n", len(selectedProvider.Models)-previewCount))
				}
			}
		} else {
			// Model selection within a provider
			if len(m.Providers) == 0 || m.ProviderIdx < 0 || m.ProviderIdx >= len(m.Providers) {
				s.WriteString(modelPickerStyle().Render("  No providers available"))
				s.WriteString("\n")
				return s.String()
			}
			provider := m.Providers[m.ProviderIdx]
			s.WriteString(modelPickerStyle().Render(fmt.Sprintf("  Provider: %s", provider.Name)))
			s.WriteString("\n\n")

			// Calculate visible range (viewport height)
			maxVisible := m.Viewport - 10 // account for header
			if maxVisible < 1 {
				maxVisible = 1
			}
			if maxVisible > len(provider.Models) {
				maxVisible = len(provider.Models)
			}

			start := 0
			end := maxVisible
			if len(provider.Models) > maxVisible {
				start = m.ModelIdx - maxVisible/2
				if start < 0 {
					start = 0
				}
				end = start + maxVisible
				if end > len(provider.Models) {
					end = len(provider.Models)
					start = end - maxVisible
					if start < 0 {
						start = 0
					}
				}
			}

			for i := start; i < end; i++ {
				model := provider.Models[i]
				cursor := "  "
				if i == m.ModelIdx {
					cursor = cursorStyle().Render(" >")
				}
				s.WriteString(fmt.Sprintf("%s %s\n", cursor, model))
			}

			s.WriteString("\n")
			s.WriteString(footerStyle().Render("  j/k: navigate | Enter: select model | Esc: back"))
			s.WriteString("\n")
		}

		return s.String()
	}

	// Normal tree view
	s.WriteString(headerStyle().Render("  omo-switch "))
	s.WriteString(modeBadgeStyle().Render(" ◆ "))
	s.WriteString(headerStyle().Render("Config Editor  "))
	s.WriteString("\n")
	s.WriteString(dividerStyle().Render("  " + strings.Repeat("─", 40) + "  "))
	s.WriteString("\n\n")

	// Edit mode indicator
	if m.EditMode {
		s.WriteString(headerStyle().Render("  omo-switch "))
		s.WriteString(modeBadgeStyle().Render(" ◆ "))
		s.WriteString(headerStyle().Render("Editing  "))
		s.WriteString("\n\n")
		s.WriteString(editFieldStyle().Render(fmt.Sprintf("  > %s", m.EditValue)))
		s.WriteString("\n")
		s.WriteString(footerStyle().Render("  Enter: confirm | Esc: cancel"))
		s.WriteString("\n\n")
	}

	// Calculate visible range
	start := 0
	end := len(m.FlatList)
	totalListLen := len(m.FlatList)
	isScrolling := totalListLen > m.Viewport
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

	// Scroll indicator: items above viewport
	if isScrolling && start > 0 {
		above := start
		s.WriteString(scrollIndicatorStyle().Render(fmt.Sprintf("  ↑ %d more above", above)))
		s.WriteString("\n")
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

			// Cursor indicator and row highlighting
			isCursor := actualIdx == m.Cursor
			if isCursor {
				cursor = cursorStyle().Render(" >")
			}

			// Selection marker
			selPrefix := "[ ]"
			if m.IsSelected(node.ID) {
				selPrefix = selectedStyle().Render("[x]")
			}

			// Indent based on depth
			if node.Depth > 0 {
				prefix = strings.Repeat("  ", node.Depth) + prefix
			}

			// Type badge and styling with model on next line (Phase 8)
			rowContent := prefix + cursor + selPrefix + " "
			if isCursor {
				rowContent = hoverRowStyle().Render(rowContent)
			}
			s.WriteString(rowContent)

			switch node.Type {
			case "provider":
				s.WriteString(badgeProviderStyle().Render("[P]") + " " + providerStyle().Render(node.Label))
			case "agent":
				s.WriteString(badgeAgentStyle().Render("[A]") + " " + agentStyle().Render(node.Label))
				if node.Model != "" {
					s.WriteString("\n" + strings.Repeat("  ", node.Depth+1) + cursorStyle().Render(" └─ ") + modelStyle().Render(node.Model))
				}
			case "category":
				s.WriteString(badgeCategoryStyle().Render("[C]") + " " + categoryStyle().Render(node.Label))
				if node.Model != "" {
					s.WriteString("\n" + strings.Repeat("  ", node.Depth+1) + cursorStyle().Render(" └─ ") + modelStyle().Render(node.Model))
				}
			}

			s.WriteString("\n")
		}
	}

	// Scroll indicator: items below viewport
	if isScrolling && end < totalListLen {
		below := totalListLen - end
		s.WriteString(scrollIndicatorStyle().Render(fmt.Sprintf("  ↓ %d more below", below)))
		s.WriteString("\n")
	}

	// Stats with position indicator
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
	if m.ErrorMsg != "" {
		s.WriteString(errorStyle().Render(fmt.Sprintf("  ⚠ %s", m.ErrorMsg)))
		s.WriteString("\n")
	}
	statsLine := fmt.Sprintf("  Agents: %d | Categories: %d | Selected: %d", agentCount, categoryCount, m.SelectionCount)
	if isScrolling {
		statsLine += fmt.Sprintf("  [%d/%d]", m.Cursor+1, totalListLen)
	}
	s.WriteString(footerStyle().Render(statsLine))
	s.WriteString("\n")
	s.WriteString(footerStyle().Render("  " + buildFooter(m)))
	s.WriteString("\n")

	return s.String()
}
