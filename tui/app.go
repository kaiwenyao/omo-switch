package tui

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/kaiwenyao/omo-switch/config"
)

// AppModel is the main Bubble Tea model for the TUI.
// It holds the config tree and navigation state.
type AppModel struct {
	Config         *config.Config
	Tree           []*TreeNode     // root provider nodes
	FlatList       []TreeNode      // flattened tree for linear navigation
	Cursor         int             // current position in FlatList
	Viewport       int             // visible rows
	Selected       map[string]bool // node ID -> selected
	SelectionCount int             // cached count of selected nodes
	// Edit mode state
	EditMode      bool   // true when editing
	EditNodeID    string // "agent/name" or "category/name"
	EditValue     string // current input value
	OriginalModel string // for Escape to restore
	ConfigPath    string // path to config file for auto-save
	ErrorMsg      string // error message to display in TUI (empty = no error)
	// Model picker state - two level: provider list then model list
	ModelPickerMode    bool       // true when selecting from available models
	ProviderPickerMode bool       // true when in provider selection (first level)
	ProviderIdx        int        // cursor position in providers list
	ModelIdx           int        // cursor position in models list (within selected provider)
	Providers          []Provider // all providers with their models
	AllModels          []string   // flat list of all models for compatibility
}

// NewAppModel creates a new AppModel from a config,
// building the tree structure and flattening it for navigation.
func NewAppModel(cfg *config.Config) AppModel {
	tree := BuildTree(cfg)
	flatList := Flatten(tree)
	allModels := GetAvailableModels(cfg)
	providers := GetProvidersWithModels(cfg)

	return AppModel{
		Config:             cfg,
		Tree:               tree,
		FlatList:           flatList,
		Cursor:             0,
		Viewport:           20, // default, will be updated on resize
		Selected:           make(map[string]bool),
		SelectionCount:     0,
		EditMode:           false,
		EditNodeID:         "",
		EditValue:          "",
		OriginalModel:      "",
		ConfigPath:         "", // set by main.go before creating model
		ModelPickerMode:    false,
		ProviderPickerMode: true, // start in provider selection mode
		ProviderIdx:        0,
		ModelIdx:           0,
		Providers:          providers,
		AllModels:          allModels,
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

// ClearSelection clears all selections.
func (m *AppModel) ClearSelection() {
	m.Selected = make(map[string]bool)
	m.updateSelectionCount()
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
	m.ModelPickerMode = false
}

// StartModelPicker initiates the model picker for the current node.
// Starts in provider selection mode.
// If the cursor node is not in the current selection, the selection is cleared
// so the edit targets the cursor node alone rather than the stale selection.
func (m *AppModel) StartModelPicker(nodeID string, currentModel string) {
	if m.SelectionCount > 0 && !m.IsSelected(nodeID) {
		m.ClearSelection()
	}
	m.ModelPickerMode = true
	m.ProviderPickerMode = true
	m.EditNodeID = nodeID
	m.EditValue = currentModel
	m.OriginalModel = currentModel
	m.ProviderIdx = 0
	m.ModelIdx = 0
	m.Providers = GetAllProviders(m.Config, nodeID)
}

// ConfirmModelPicker confirms the selected model from the picker.
func (m *AppModel) ConfirmModelPicker() error {
	if !m.ModelPickerMode {
		return nil
	}

	if len(m.Providers) == 0 {
		m.CancelModelPicker()
		return nil
	}

	var selectedModel string
	if m.ProviderPickerMode {
		if m.ProviderIdx < 0 || m.ProviderIdx >= len(m.Providers) {
			return nil
		}
		m.ProviderPickerMode = false
		m.ModelIdx = 0
		return nil
	}

	if m.ProviderIdx < 0 || m.ProviderIdx >= len(m.Providers) {
		return nil
	}
	provider := m.Providers[m.ProviderIdx]
	if m.ModelIdx < 0 || m.ModelIdx >= len(provider.Models) {
		return nil
	}

	if provider.Name == "recently" {
		if m.ModelIdx < len(provider.FullPaths) {
			selectedModel = provider.FullPaths[m.ModelIdx]
		} else {
			m.CancelModelPicker()
			return nil
		}
	} else {
		selectedModel = provider.Name + "/" + provider.Models[m.ModelIdx]
	}

	m.EditValue = selectedModel

	m.Config.AddRecentModel(m.EditNodeID, selectedModel)

	return m.ConfirmEdit()
}

// CancelModelPicker cancels the model picker and restores the original value.
func (m *AppModel) CancelModelPicker() {
	m.ModelPickerMode = false
	m.ProviderPickerMode = true
	m.EditNodeID = ""
	m.EditValue = ""
	m.OriginalModel = ""
	m.ProviderIdx = 0
	m.ModelIdx = 0
}

// ConfirmEdit applies the edited value and saves the config.
//
// If auto-save fails (disk full, permission denied, read-only filesystem, …),
// the in-memory mutation is rolled back from a pre-mutation JSON snapshot,
// the picker/edit mode and selection are preserved so the user can retry or
// cancel, and the error surfaces on screen. This prevents the silent
// memory-vs-disk divergence where quitting after a failed save would discard
// the change without warning.
func (m *AppModel) ConfirmEdit() error {
	if !m.EditMode && !m.ModelPickerMode {
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

	// Snapshot for rollback on save failure. If snapshotting itself fails,
	// we abort before mutating so memory stays consistent with disk.
	var snapshot []byte
	if m.ConfigPath != "" {
		var snapErr error
		snapshot, snapErr = json.Marshal(m.Config)
		if snapErr != nil {
			m.ErrorMsg = fmt.Sprintf("snapshot failed: %v", snapErr)
			return snapErr
		}
	}

	// Apply change to each node (batch mode updates fallback_models)
	batch := len(m.Selected) > 0
	for _, nodeID := range nodeIDs {
		ApplyChange(m.Config, nodeID, m.EditValue, batch)
	}

	// Auto-save to JSON. On failure, roll back and keep the picker/selection
	// open so the user can retry the save or Esc to cancel cleanly.
	if m.ConfigPath != "" {
		if err := SaveConfig(m.Config, m.ConfigPath); err != nil {
			var restored config.Config
			if rbErr := json.Unmarshal(snapshot, &restored); rbErr == nil {
				*m.Config = restored
			}
			m.Tree = BuildTree(m.Config)
			m.FlatList = Flatten(m.Tree)
			if len(m.FlatList) == 0 {
				m.Cursor = 0
			} else if m.Cursor >= len(m.FlatList) {
				m.Cursor = len(m.FlatList) - 1
			}
			m.ErrorMsg = fmt.Sprintf("save failed (changes reverted): %v", err)
			return err
		}
		m.ErrorMsg = ""
	}

	// Rebuild tree to reflect changes
	m.Tree = BuildTree(m.Config)
	m.FlatList = Flatten(m.Tree)
	// Clamp cursor to new list bounds
	if len(m.FlatList) == 0 {
		m.Cursor = 0
	} else if m.Cursor >= len(m.FlatList) {
		m.Cursor = len(m.FlatList) - 1
	} else if m.Cursor < 0 {
		m.Cursor = 0
	}
	// Refresh available models
	m.AllModels = GetAvailableModels(m.Config)
	// Refresh providers list
	if m.EditNodeID != "" {
		m.Providers = GetAllProviders(m.Config, m.EditNodeID)
	}

	// Exit edit/picker mode
	m.EditMode = false
	m.ModelPickerMode = false
	m.EditNodeID = ""
	m.EditValue = ""
	m.OriginalModel = ""
	m.ClearSelection()

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
		_, cmd := HandleKey(msg, &m)
		return m, cmd
	}
	return m, nil
}

// View renders the TUI by delegating to the package-level View function.
func (m AppModel) View() string {
	return View(m)
}
