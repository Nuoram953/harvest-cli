package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

type DateInputModel struct {
	textInput   textinput.Model
	err         error
	validDate   *time.Time
	submitted   bool
	title       string
	placeholder string
}

func NewDateInput(title string) DateInputModel {
	ti := textinput.New()
	ti.Placeholder = "YYYY-MM-DD"
	ti.Focus()
	ti.CharLimit = 10
	ti.Width = 20

	return DateInputModel{
		textInput:   ti,
		title:       title,
		placeholder: "YYYY-MM-DD",
	}
}

func (m *DateInputModel) SetPlaceholder(placeholder string) {
	m.placeholder = placeholder
	m.textInput.Placeholder = placeholder
}

func (m DateInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m DateInputModel) Update(msg tea.Msg) (DateInputModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			dateStr := strings.TrimSpace(m.textInput.Value())
			
			var parsedDate time.Time
			var err error
			
			if dateStr == "" {
				parsedDate = time.Now()
			} else {
				parsedDate, err = time.Parse("2006-01-02", dateStr)
				if err != nil {
					m.err = fmt.Errorf("invalid date format. Use YYYY-MM-DD")
					return m, nil
				}
			}

			m.validDate = &parsedDate
			m.err = nil
			m.submitted = true
			return m, nil

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, nil
		}

	case error:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m DateInputModel) View() string {
	var b strings.Builder

	if m.submitted && m.validDate != nil {
		return b.String()
	}

	b.WriteString("Enter a date:")
	b.WriteString("\n\n")

	b.WriteString(inputStyle.Render(m.textInput.View()))
	b.WriteString("\n\n")

	if m.err != nil {
		b.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		b.WriteString("\n\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Press Enter to submit (empty = today's date)"))

	return b.String()
}

func (m DateInputModel) IsSubmitted() bool {
	return m.submitted && m.validDate != nil
}

func (m DateInputModel) GetDate() *time.Time {
	if m.submitted && m.validDate != nil {
		return m.validDate
	}
	return nil
}

func (m DateInputModel) GetDateISO8601() string {
	if m.submitted && m.validDate != nil {
		return m.validDate.Format("2006-01-02")
	}
	return ""
}

func (m *DateInputModel) Reset() {
	m.textInput.SetValue("")
	m.err = nil
	m.validDate = nil
	m.submitted = false
	m.textInput.Focus()
}

func TextInputDate(title string) (string, error) {
	model := NewDateInput(title)
	
	program := tea.NewProgram(dateInputWrapper{model})
	finalModel, err := program.Run()
	if err != nil {
		return "", err
	}

	if wrapper, ok := finalModel.(dateInputWrapper); ok {
		return wrapper.model.GetDateISO8601(), nil
	}
	
	return "", fmt.Errorf("failed to get date input")
}

func TextInputDateAsTime(title string) (*time.Time, error) {
	model := NewDateInput(title)
	
	program := tea.NewProgram(dateInputWrapper{model})
	finalModel, err := program.Run()
	if err != nil {
		return nil, err
	}

	if wrapper, ok := finalModel.(dateInputWrapper); ok {
		return wrapper.model.GetDate(), nil
	}
	
	return nil, fmt.Errorf("failed to get date input")
}

type dateInputWrapper struct {
	model DateInputModel
}

func (w dateInputWrapper) Init() tea.Cmd {
	return w.model.Init()
}

func (w dateInputWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return w, tea.Quit
		case tea.KeyEnter:
			newModel, cmd := w.model.Update(msg)
			w.model = newModel
			if w.model.IsSubmitted() {
				return w, tea.Quit
			}
			return w, cmd
		}
	}

	newModel, cmd := w.model.Update(msg)
	w.model = newModel
	return w, cmd
}

func (w dateInputWrapper) View() string {
	view := w.model.View()
	if !w.model.IsSubmitted() {
		view += "\n"
	}
	return view
}
