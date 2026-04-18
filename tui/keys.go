package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
)

// HandleKey processes keyboard input for vim-style navigation and selection.
// It handles:
//   - j/k or Down/Up arrows: move cursor down/up
//   - h/l or Left/Right arrows: move to parent/child or collapse/expand
//   - Space: toggle selection on current node
//   - Enter: print selected model (edit mode comes later)
func HandleKey(msg tea.KeyMsg, m *AppModel) (*AppModel, tea.Cmd) {
	// Handle letter keys (j/k/h/l) via runes
	if msg.Type == tea.KeyRunes {
		switch msg.Runes[0] {
		case 'j':
			// Move cursor down
			if m.Cursor < len(m.FlatList)-1 {
				m.Cursor++
			}
			return m, nil
		case 'k':
			// Move cursor up
			if m.Cursor > 0 {
				m.Cursor--
			}
			return m, nil
		case 'l':
			// Move to child node (first child of current provider)
			current := m.FlatList[m.Cursor]
			if current.Type == "provider" && len(current.Children) > 0 {
				// Find the first child in FlatList
				for i, node := range m.FlatList {
					if node.Parent != nil && node.Parent.ID == current.ID {
						m.Cursor = i
						break
					}
				}
			}
			return m, nil
		case 'h':
			// Move to parent node
			current := m.FlatList[m.Cursor]
			if current.Parent != nil {
				// Find parent in FlatList
				for i, node := range m.FlatList {
					if node.ID == current.Parent.ID {
						m.Cursor = i
						break
					}
				}
			}
			return m, nil
		}
	}

	switch msg.Type {
	case tea.KeyDown:
		// Move cursor down
		if m.Cursor < len(m.FlatList)-1 {
			m.Cursor++
		}

	case tea.KeyUp:
		// Move cursor up
		if m.Cursor > 0 {
			m.Cursor--
		}

	case tea.KeyRight:
		// Move to child node (first child of current provider)
		current := m.FlatList[m.Cursor]
		if current.Type == "provider" && len(current.Children) > 0 {
			// Find the first child in FlatList
			for i, node := range m.FlatList {
				if node.Parent != nil && node.Parent.ID == current.ID {
					m.Cursor = i
					break
				}
			}
		}

	case tea.KeyLeft:
		// Move to parent node
		current := m.FlatList[m.Cursor]
		if current.Parent != nil {
			// Find parent in FlatList
			for i, node := range m.FlatList {
				if node.ID == current.Parent.ID {
					m.Cursor = i
					break
				}
			}
		}

	case tea.KeySpace:
		// Toggle selection on current node
		current := m.FlatList[m.Cursor]
		m.ToggleSelection(current.ID)

	case tea.KeyEnter:
		// Print selected model for now (edit mode comes in Plan 04)
		current := m.FlatList[m.Cursor]
		if current.Type == "agent" || current.Type == "category" {
			fmt.Printf("Selected model: %s\n", current.Model)
		}
	}

	return m, nil
}
