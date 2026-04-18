package tui

import (
	"github.com/charmbracelet/bubbletea"
)

// validEditChars contains characters allowed in model names.
// Includes alphanumeric, slash, hyphen, dot, underscore, colon.
var validEditChars = []rune{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
	'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
	'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'/', '-', '.', '_', ':',
}

// isValidEditChar returns true if the rune is valid for model input.
func isValidEditChar(r rune) bool {
	for _, c := range validEditChars {
		if c == r {
			return true
		}
	}
	return false
}

// HandleKey processes keyboard input for vim-style navigation and selection.
// It handles:
//   - j/k or Down/Up arrows: move cursor down/up
//   - h/l or Left/Right arrows: move to parent/child or collapse/expand
//   - Space: toggle selection on current node
//   - Enter: start edit mode (when not editing) or confirm edit (when editing)
//   - Escape: cancel edit mode
//   - Character keys: append to EditValue when in edit mode
func HandleKey(msg tea.KeyMsg, m *AppModel) (*AppModel, tea.Cmd) {
	// Handle character input in edit mode
	if m.EditMode && msg.Type == tea.KeyRunes {
		r := msg.Runes[0]
		if isValidEditChar(r) {
			m.EditValue += string(r)
		}
		return m, nil
	}

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
		// Move cursor down (only when not in edit mode)
		if !m.EditMode && m.Cursor < len(m.FlatList)-1 {
			m.Cursor++
		}

	case tea.KeyUp:
		// Move cursor up (only when not in edit mode)
		if !m.EditMode && m.Cursor > 0 {
			m.Cursor--
		}

	case tea.KeyRight:
		// Move to child node (first child of current provider)
		if !m.EditMode {
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
		}

	case tea.KeyLeft:
		// Move to parent node
		if !m.EditMode {
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
		}

	case tea.KeySpace:
		// Toggle selection on current node (only when not in edit mode)
		if !m.EditMode {
			current := m.FlatList[m.Cursor]
			m.ToggleSelection(current.ID)
		}

	case tea.KeyEnter:
		if m.EditMode {
			// Confirm edit and save
			m.ConfirmEdit()
		} else {
			// Start edit mode for agent/category
			current := m.FlatList[m.Cursor]
			if current.Type == "agent" || current.Type == "category" {
				m.StartEdit(current.ID, current.Model)
			}
		}

	case tea.KeyEscape:
		if m.EditMode {
			m.CancelEdit()
		}
	}

	return m, nil
}
