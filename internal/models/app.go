package models

import (
	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	mainView state = iota
	listView
	detailView
)

type App struct {
	state  state
	width  int
	height int
}

func NewApp() *App {
	return &App{
		state: mainView,
	}
}

func (a *App) Init() tea.Cmd {
	return nil
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return a, tea.Quit
		}
	}
	return a, nil
}

func (a *App) View() string {
	switch a.state {
	case mainView:
		return a.mainView()
	case listView:
		return a.listView()
	default:
		return "Unknown view"
	}
}

func (a *App) mainView() string {
	return "Welcome to your Bubble Tea app!"
}

func (a *App) listView() string {
	return "List view"
}
