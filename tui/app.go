package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/kaiwenyao/omo-switch/config"
)

// AppModel is the main Bubble Tea model for the TUI.
// It holds the config tree and navigation state.
type AppModel struct {
	Config         *config.Config
	Tree           []*TreeNode   // root provider nodes
	FlatList       []TreeNode    // flattened tree for linear navigation
	Cursor         int           // current position in FlatList
	Viewport       int           // visible rows
	Selected       map[string]bool // node ID -> selected
	SelectionCount int           // cached count of selected nodes
	// Edit mode state
	EditMode       bool   // true when editing
	EditNodeID     string // "agent/name" or "category/name"
	EditValue      string // current input value
	OriginalModel  string // for Escape to restore
	ConfigPath     string // path to config file for auto-save
}

// NewAppModel creates a new AppModel from a config,
// building the tree structure and flattening it for navigation.
func NewAppModel(cfg *config.Config) AppModel {
	tree := BuildTree(cfg)
	flatList := Flatten(tree)

	return AppModel{
		Config:         cfg,
		Tree:           tree,
		FlatList:       flatList,
		Cursor:         0,
		Viewport:       20, // default, will be updated on resize
		Selected:       make(map[string]bool),
		SelectionCount: 0,
		EditMode:       false,
		EditNodeID:     "",
		EditValue:      "",
		OriginalModel:  "",
		ConfigPath:     "", // set by main.go before creating model
	}
}

// ToggleSelection toggles the selection state of a node by ID.
func (m *AppModel) ToggleSelection(nodeID string) {
	if m.Selected[nodeID] {
		delete(m.Selected, nodeID)
	} else {
		m.Selected[nodeID] = true
	}
	m.updateSelectionCount()
}

// IsSelected returns true if the node is selected.
func (m *AppModel) IsSelected(nodeID string) bool {
	return m.Selected[nodeID]
}

// updateSelectionCount recalculates the selection count.
func (m *AppModel) updateSelectionCount() {
	m.SelectionCount = len(m.Selected)
}

// StartEdit initiates edit mode for the specified node.
func (m *AppModel) StartEdit(nodeID string, currentModel string) {
	m.EditMode = true
	m.EditNodeID = nodeID
	m.EditValue = currentModel
	m.OriginalModel = currentModel
}

// ConfirmEdit applies the edited value and saves the config.
func (m *AppModel) ConfirmEdit() error {
	if !m.EditMode {
		return nil
	}

	// Apply the change to all selected nodes (or the single edited node)
	nodeIDs := make([]string, 0, len(m.Selected)+1)
	if len(m.Selected) > 0 {
		for id := range m.Selected {
			nodeIDs = append(nodeIDs, id)
		}
	} else if m.EditNodeID != "" {
		nodeIDs = append(nodeIDs, m.EditNodeID)
	}

	// Apply change to each node
	for _, nodeID := range nodeIDs {
		ApplyChange(m.Config, nodeID, m.EditValue)
	}

	// Auto-save to JSON
	if m.ConfigPath != "" {
		if err := SaveConfig(m.Config, m.ConfigPath); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
		}
	}

	// Rebuild tree to reflect changes
	m.Tree = BuildTree(m.Config)
	m.FlatList = Flatten(m.Tree)

	// Exit edit mode
	m.EditMode = false
	m.EditNodeID = ""
	m.EditValue = ""
	m.OriginalModel = ""

	return nil
}

// CancelEdit restores the original value and exits edit mode.
func (m *AppModel) CancelEdit() {
	m.EditMode = false
	m.EditNodeID = ""
	m.EditValue = ""
	m.OriginalModel = ""
}

// Init initializes the Bubble Tea model.
// Returns no initial commands.
func (m AppModel) Init() tea.Cmd {
	return nil
}

// Update handles window resize events and keyboard input.
// Returns the updated model and no command.
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Account for header and footer lines (6 total)
		m.Viewport = msg.Height - 6
		if m.Viewport < 1 {
			m.Viewport = 1
		}
	case tea.KeyMsg:
		// Handle vim-style navigation and selection keys
		HandleKey(msg, &m)
	}
	return m, nil
}

// View renders the TUI by delegating to the package-level View function.
func (m AppModel) View() string {
	return View(m)
}
