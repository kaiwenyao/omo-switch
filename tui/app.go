package tui

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/kaiwenyao/omo-switch/config"
)

// AppModel is the main Bubble Tea model for the TUI.
// It holds the config tree and navigation state.
type AppModel struct {
	Config   *config.Config
	Tree     []*TreeNode   // root provider nodes
	FlatList []TreeNode    // flattened tree for linear navigation
	Cursor   int           // current position in FlatList
	Viewport int           // visible rows
}

// NewAppModel creates a new AppModel from a config,
// building the tree structure and flattening it for navigation.
func NewAppModel(cfg *config.Config) AppModel {
	tree := BuildTree(cfg)
	flatList := Flatten(tree)

	return AppModel{
		Config:   cfg,
		Tree:     tree,
		FlatList: flatList,
		Cursor:   0,
		Viewport: 20, // default, will be updated on resize
	}
}

// Init initializes the Bubble Tea model.
// Returns no initial commands.
func (m AppModel) Init() tea.Cmd {
	return nil
}

// Update handles window resize events and updates the viewport.
// Returns the updated model and no command.
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Account for header and footer lines (6 total)
		m.Viewport = msg.Height - 6
		if m.Viewport < 1 {
			m.Viewport = 1
		}
	}
	return m, nil
}

// View renders the TUI by delegating to the package-level View function.
func (m AppModel) View() string {
	return View(m)
}
