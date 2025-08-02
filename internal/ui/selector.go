package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"harvest-cli/internal/api"
)

type Selectable interface {
	GetID() string
	GetTitle() string
	GetDescription() string
}

type DataLoader[T Selectable] interface {
	Load() ([]T, error)
}

// Generic list item wrapper
type selectableItem[T Selectable] struct {
	item T
}

func (i selectableItem[T]) FilterValue() string { return i.item.GetTitle() }
func (i selectableItem[T]) Title() string       { return i.item.GetTitle() }
func (i selectableItem[T]) Description() string { return i.item.GetDescription() }

// Generic selector model
type selectorModel[T Selectable] struct {
	loading     bool
	spinner     spinner.Model
	list        list.Model
	items       []T
	selected    *T
	err         error
	quitting    bool
	loader      DataLoader[T]
	title       string
	emptyMsg    string
	loadingMsg  string
}

type itemsLoadedMsg[T Selectable] []T
type itemsErrorMsg error

// SelectorConfig holds configuration for the selector
type SelectorConfig struct {
	Title      string
	EmptyMsg   string
	LoadingMsg string
	Width      int
	Height     int
}

func NewSelector[T Selectable](loader DataLoader[T], config SelectorConfig) *selectorModel[T] {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	// Set defaults
	if config.Title == "" {
		config.Title = "Select an Item"
	}
	if config.EmptyMsg == "" {
		config.EmptyMsg = "No items found."
	}
	if config.LoadingMsg == "" {
		config.LoadingMsg = "Loading items..."
	}
	if config.Width == 0 {
		config.Width = 80
	}
	if config.Height == 0 {
		config.Height = 14
	}

	return &selectorModel[T]{
		loading:    true,
		spinner:    s,
		loader:     loader,
		title:      config.Title,
		emptyMsg:   config.EmptyMsg,
		loadingMsg: config.LoadingMsg,
	}
}

func (m *selectorModel[T]) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.loadItems)
}

func (m *selectorModel[T]) loadItems() tea.Msg {
	items, err := m.loader.Load()
	if err != nil {
		return itemsErrorMsg(err)
	}
	return itemsLoadedMsg[T](items)
}

func (m *selectorModel[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if !m.loading && len(m.items) > 0 {
				if selectedItem, ok := m.list.SelectedItem().(selectableItem[T]); ok {
					m.selected = &selectedItem.item
					m.quitting = true
					return m, tea.Quit
				}
			}
		}

	case itemsLoadedMsg[T]:
		m.loading = false
		m.items = []T(msg)
		
		// Create list items
		listItems := make([]list.Item, len(m.items))
		for i, item := range m.items {
			listItems[i] = selectableItem[T]{item: item}
		}
		
		// Setup list
		l := list.New(listItems, list.NewDefaultDelegate(), 80, 14)
		l.Title = m.title
		l.SetShowStatusBar(false)
		l.SetFilteringEnabled(true)
		l.Styles.Title = lipgloss.NewStyle().
			Foreground(lipgloss.Color("62")).
			Bold(true).
			Padding(0, 0, 1, 2)
		
		m.list = l
		return m, nil

	case itemsErrorMsg:
		m.loading = false
		m.err = error(msg)
		return m, nil

	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	if !m.loading {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *selectorModel[T]) View() string {
	if m.quitting {
		return ""
	}

	if m.err != nil {
		return fmt.Sprintf("Error loading items: %v\n", m.err)
	}

	if m.loading {
		return fmt.Sprintf("\n   %s %s\n\n", m.spinner.View(), m.loadingMsg)
	}

	if len(m.items) == 0 {
		return m.emptyMsg + "\n"
	}

	return "\n" + m.list.View() + "\n\nPress Enter to select, q/esc to quit\n"
}

// RunSelector runs the selector and returns the selected item
func RunSelector[T Selectable](loader DataLoader[T], config SelectorConfig) (*T, error) {
	model := NewSelector(loader, config)
	p := tea.NewProgram(model)
	
	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}
	
	if final, ok := finalModel.(*selectorModel[T]); ok {
		if final.selected != nil {
			return final.selected, nil
		}
		if final.err != nil {
			return nil, final.err
		}
	}
	
	return nil, fmt.Errorf("no item selected")
}

// Entry implementation of Selectable interface
type EntrySelectable struct {
	*api.Entry
}

func (e EntrySelectable) GetID() string {
	return e.Entry.ID
}

func (e EntrySelectable) GetTitle() string {
	return e.Entry.Title
}

func (e EntrySelectable) GetDescription() string {
	return fmt.Sprintf("ID: %s | Status: %s | Private: %t", e.Entry.ID, e.Entry.Status, e.Entry.Private)
}

// Task implementation of Selectable interface
type TaskSelectable struct {
	*api.Task
}

func (t TaskSelectable) GetTitle() string {
	return t.Task.Name
}

func (t TaskSelectable) GetDescription() string {
	return fmt.Sprintf("ID: %s | Project: %s", t.Task.ID)
}

type ProjectSelectable struct {
	*api.ProjectAssignment
}

func (p ProjectSelectable) GetID() string {
	return p.Project.Name
}

func (p ProjectSelectable) GetTitle() string {
	return p.Project.Name
}

func (p ProjectSelectable) GetDescription() string {
	return fmt.Sprintf("ID: %s | Client: %s | Status: %s", p.Project.ID)
}

// Entry loader implementation
type EntryLoader struct {
	client *api.Client
	params api.ListParams
}

func (el *EntryLoader) Load() ([]EntrySelectable, error) {
	entries, err := el.client.ListEntries(el.params)
	if err != nil {
		return nil, err
	}
	
	selectableEntries := make([]EntrySelectable, len(entries))
	for i, entry := range entries {
		selectableEntries[i] = EntrySelectable{Entry: entry}
	}
	
	return selectableEntries, nil
}

// Task loader implementation
type TaskLoader struct {
	client *api.Client
	params api.ListParams
}

// func (tl *TaskLoader) Load() ([]TaskSelectable, error) {
// 	tasks, err := tl.client.ListTasks()
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	selectableTasks := make([]TaskSelectable, len(tasks))
// 	for i, task := range tasks {
// 		selectableTasks[i] = TaskSelectable{Task: task}
// 	}
//
// 	return selectableTasks, nil
// }

// Project loader implementation
type ProjectLoader struct {
	client *api.Client
	params api.ListParams
}

func (pl *ProjectLoader) Load() ([]ProjectSelectable, error) {
	projects, err := pl.client.ListAssignedProjects(pl.params)
	if err != nil {
		return nil, err
	}
	
	selectableProjects := make([]ProjectSelectable, len(projects))
	for i, project := range projects {
		selectableProjects[i] = ProjectSelectable{ProjectAssignment: project}
	}
	
	return selectableProjects, nil
}

// Helper function to build list parameters
func buildListParams() api.ListParams {
	// Return the appropriate api.ListParams struct
	// You'll need to adjust this based on your actual api.ListParams structure
	return api.ListParams{
		// Add fields based on your api.ListParams struct
		// Examples (adjust field names to match your API):
		// Limit:  100,
		// Status: "",
		// Page:   1,
	}
}

// SelectEntryInteractively allows interactive selection of an entry
func SelectEntryInteractively(client *api.Client) (*api.Entry, error) {
	loader := &EntryLoader{
		client: client,
		params: buildListParams(),
	}
	config := SelectorConfig{
		Title:      "Select an Entry",
		EmptyMsg:   "No entries found.",
		LoadingMsg: "Loading entries...",
	}
	
	selected, err := RunSelector(loader, config)
	if err != nil {
		return nil, err
	}
	
	return selected.Entry, nil
}

// // SelectTaskInteractively allows interactive selection of a task
// func SelectTaskInteractively(client *api.Client) (*api.Task, error) {
// 	loader := &TaskLoader{client: client}
// 	config := SelectorConfig{
// 		Title:      "Select a Task",
// 		EmptyMsg:   "No tasks found.",
// 		LoadingMsg: "Loading tasks...",
// 	}
//
// 	selected, err := RunSelector(loader, config)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return selected.Task, nil
// }

func SelectProjectInteractively(client *api.Client) (*api.Project, error) {
	loader := &ProjectLoader{client: client}
	config := SelectorConfig{
		Title:      "Select a Project",
		EmptyMsg:   "No projects found.",
		LoadingMsg: "Loading projects...",
	}
	
	selected, err := RunSelector(loader, config)
	if err != nil {
		return nil, err
	}
	
	return &selected.Project, nil
}
