package tui

import (
	"github.com/charmbracelet/bubbletea"
)

// HandleKey processes keyboard input for vim-style navigation and selection.
// It handles:
//   - j/k or Down/Up arrows: move cursor down/up
//   - h/l or Left/Right arrows: move to parent/child or collapse/expand
//   - Space: toggle selection on current node
//   - Enter: start model picker (when not editing) or confirm selection (when in picker)
//   - Escape: cancel model picker or quit app
func HandleKey(msg tea.KeyMsg, m *AppModel) (*AppModel, tea.Cmd) {
	// Handle rune input (j/k for navigation)
	if msg.Type == tea.KeyRunes {
		r := msg.Runes[0]

		// j/k navigation in provider picker mode
		if m.ModelPickerMode && m.ProviderPickerMode {
			switch r {
			case 'j':
				if m.ProviderIdx < len(m.Providers)-1 {
					m.ProviderIdx++
				}
			case 'k':
				if m.ProviderIdx > 0 {
					m.ProviderIdx--
				}
			case 'q':
				m.CancelModelPicker()
				return m, nil
			}
			return m, nil
		}

		// j/k navigation in model selection mode
		if m.ModelPickerMode && !m.ProviderPickerMode {
			switch r {
			case 'j':
				if len(m.Providers) > 0 && m.ProviderIdx < len(m.Providers) && m.ModelIdx < len(m.Providers[m.ProviderIdx].Models)-1 {
					m.ModelIdx++
				}
			case 'k':
				if m.ModelIdx > 0 {
					m.ModelIdx--
				}
			case 'q':
				m.ProviderPickerMode = true
				m.ModelIdx = 0
				return m, nil
			}
			return m, nil
		}

		// Edit mode - type model name
		if m.EditMode {
			if isValidEditChar(r) {
				m.EditValue += string(r)
			}
			return m, nil
		}

		// Normal tree view - vim navigation
		switch r {
		case 'j':
			if len(m.FlatList) > 0 && m.Cursor < len(m.FlatList)-1 {
				m.Cursor++
			}
		case 'k':
			if m.Cursor > 0 {
				m.Cursor--
			}
		case 'l':
			if len(m.FlatList) == 0 {
				return m, nil
			}
			current := m.FlatList[m.Cursor]
			if current.Type == "provider" && len(current.Children) > 0 {
				for i, node := range m.FlatList {
					if node.Parent != nil && node.Parent.ID == current.ID {
						m.Cursor = i
						break
					}
				}
			}
		case 'h':
			if len(m.FlatList) == 0 {
				return m, nil
			}
			current := m.FlatList[m.Cursor]
			if current.Parent != nil {
				for i, node := range m.FlatList {
					if node.ID == current.Parent.ID {
						m.Cursor = i
						break
					}
				}
			}
		case 'x':
			if m.SelectionCount > 0 {
				m.ClearSelection()
			}
		case 'q':
			return m, tea.Quit
		}
		return m, nil
	}

	// Non-rune key handling
	switch msg.Type {
	case tea.KeyBackspace, tea.KeyDelete:
		if m.EditMode && len(m.EditValue) > 0 {
			m.EditValue = m.EditValue[:len(m.EditValue)-1]
		}

	case tea.KeyDown:
		if m.ModelPickerMode && m.ProviderPickerMode {
			if m.ProviderIdx < len(m.Providers)-1 {
				m.ProviderIdx++
			}
		} else if m.ModelPickerMode && !m.ProviderPickerMode {
			if len(m.Providers) > 0 && m.ProviderIdx < len(m.Providers) && m.ModelIdx < len(m.Providers[m.ProviderIdx].Models)-1 {
				m.ModelIdx++
			}
		} else if !m.EditMode && !m.ModelPickerMode && len(m.FlatList) > 0 && m.Cursor < len(m.FlatList)-1 {
			m.Cursor++
		}

	case tea.KeyUp:
		if m.ModelPickerMode && m.ProviderPickerMode {
			if m.ProviderIdx > 0 {
				m.ProviderIdx--
			}
		} else if m.ModelPickerMode && !m.ProviderPickerMode {
			if m.ModelIdx > 0 {
				m.ModelIdx--
			}
		} else if !m.EditMode && !m.ModelPickerMode && len(m.FlatList) > 0 && m.Cursor > 0 {
			m.Cursor--
		}

	case tea.KeyRight:
		if !m.EditMode && !m.ModelPickerMode && len(m.FlatList) > 0 && m.Cursor >= 0 && m.Cursor < len(m.FlatList) {
			current := m.FlatList[m.Cursor]
			if current.Type == "provider" && len(current.Children) > 0 {
				for i, node := range m.FlatList {
					if node.Parent != nil && node.Parent.ID == current.ID {
						m.Cursor = i
						break
					}
				}
			}
		}

	case tea.KeyLeft:
		if !m.EditMode && !m.ModelPickerMode && len(m.FlatList) > 0 && m.Cursor >= 0 && m.Cursor < len(m.FlatList) {
			current := m.FlatList[m.Cursor]
			if current.Parent != nil {
				for i, node := range m.FlatList {
					if node.ID == current.Parent.ID {
						m.Cursor = i
						break
					}
				}
			}
		}

	case tea.KeySpace:
		if !m.EditMode && !m.ModelPickerMode && len(m.FlatList) > 0 && m.Cursor >= 0 && m.Cursor < len(m.FlatList) {
			current := m.FlatList[m.Cursor]
			m.ToggleSelection(current.ID)
		}

	case tea.KeyEnter:
		if m.ModelPickerMode {
			m.ConfirmModelPicker()
		} else if m.EditMode {
			m.ConfirmEdit()
		} else if len(m.FlatList) > 0 && m.Cursor >= 0 && m.Cursor < len(m.FlatList) {
			current := m.FlatList[m.Cursor]
			if current.Type == "agent" || current.Type == "category" {
				m.StartModelPicker(current.ID, current.Model)
			}
		}

	case tea.KeyEscape:
		if m.ModelPickerMode {
			if !m.ProviderPickerMode {
				// Go back to provider selection
				m.ProviderPickerMode = true
				m.ModelIdx = 0
			} else {
				m.CancelModelPicker()
			}
		} else if m.EditMode {
			m.CancelEdit()
		} else {
			return m, tea.Quit
		}
	}

	return m, nil
}

// validEditChars contains characters allowed in model names.
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
