package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	textInputTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#7D56F4")).
				Padding(0, 1)

	textInputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205"))

	textInputErrorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("196"))

	textKeyStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("46"))

	textInputHelpStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241"))
)

type TextInputOptions struct {
	Title        string
	Prompt       string
	Placeholder  string
	Required     bool
	Password     bool
	CharLimit    int
	Width        int
	ValidateFunc func(string) error
	DefaultValue string
}

type TextInputModel struct {
	textInput textinput.Model
	options   TextInputOptions
	err       error
	submitted bool
	cancelled bool
}

func NewTextInput(options TextInputOptions) TextInputModel {
	ti := textinput.New()

	if options.Prompt == "" {
		options.Prompt = "Enter text:"
	}
	if options.Width == 0 {
		options.Width = 50
	}
	if options.CharLimit == 0 {
		options.CharLimit = 200
	}

	ti.Placeholder = options.Placeholder
	ti.Width = options.Width
	ti.CharLimit = options.CharLimit
	ti.Focus()

	if options.Password {
		ti.EchoMode = textinput.EchoPassword
		ti.EchoCharacter = '•'
	}

	if options.DefaultValue != "" {
		ti.SetValue(options.DefaultValue)
	}

	return TextInputModel{
		textInput: ti,
		options:   options,
	}
}

func (m TextInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m TextInputModel) Update(msg tea.Msg) (TextInputModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			value := strings.TrimSpace(m.textInput.Value())

			if m.options.Required && value == "" {
				m.err = fmt.Errorf("this field is required")
				return m, nil
			}

			if m.options.ValidateFunc != nil {
				if err := m.options.ValidateFunc(value); err != nil {
					m.err = err
					return m, nil
				}
			}

			m.err = nil
			m.submitted = true
			return m, nil

		case tea.KeyCtrlC, tea.KeyEsc:
			m.cancelled = true
			return m, nil
		}

	case error:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)

	if m.err != nil && msg != nil {
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyRunes {
			m.err = nil
		}
	}

	return m, cmd
}

func (m TextInputModel) View() string {
	var b strings.Builder

	if m.submitted {
		return b.String()
	}

	if m.cancelled {
		b.WriteString(textInputErrorStyle.Render("✗ Input cancelled"))
		return b.String()
	}

	prompt := m.options.Prompt
	if m.options.Required {
		prompt += " *"
	}
	b.WriteString(prompt)
	b.WriteString("\n\n")

	b.WriteString(textInputStyle.Render(m.textInput.View()))
	b.WriteString("\n\n")

	if m.err != nil {
		b.WriteString(textInputErrorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		b.WriteString("\n\n")
	}


	return b.String()
}

func (m TextInputModel) IsSubmitted() bool {
	return m.submitted
}

func (m TextInputModel) IsCancelled() bool {
	return m.cancelled
}

func (m TextInputModel) GetValue() string {
	if m.submitted {
		return strings.TrimSpace(m.textInput.Value())
	}
	return ""
}

func (m *TextInputModel) Reset() {
	m.textInput.SetValue(m.options.DefaultValue)
	m.err = nil
	m.submitted = false
	m.cancelled = false
	m.textInput.Focus()
}

func (m *TextInputModel) SetValue(value string) {
	m.textInput.SetValue(value)
}

func (m *TextInputModel) Focus() {
	m.textInput.Focus()
}

func (m *TextInputModel) Blur() {
	m.textInput.Blur()
}

func TextInput(options TextInputOptions) (string, error) {
	model := NewTextInput(options)

	program := tea.NewProgram(textInputWrapper{model})
	finalModel, err := program.Run()
	if err != nil {
		return "", err
	}

	if wrapper, ok := finalModel.(textInputWrapper); ok {
		if wrapper.model.IsCancelled() {
			return "", fmt.Errorf("input cancelled")
		}
		return wrapper.model.GetValue(), nil
	}

	return "", fmt.Errorf("failed to get text input")
}

func SimpleTextInput(title, prompt string) (string, error) {
	return TextInput(TextInputOptions{
		Title:  title,
		Prompt: prompt,
	})
}

func RequiredTextInput(title, prompt string) (string, error) {
	return TextInput(TextInputOptions{
		Title:    title,
		Prompt:   prompt,
		Required: true,
	})
}

func PasswordInput(title, prompt string) (string, error) {
	return TextInput(TextInputOptions{
		Title:    title,
		Prompt:   prompt,
		Password: true,
		Required: true,
	})
}

func ValidatedTextInput(title, prompt string, validateFunc func(string) error) (string, error) {
	return TextInput(TextInputOptions{
		Title:        title,
		Prompt:       prompt,
		Required:     true,
		ValidateFunc: validateFunc,
	})
}

type textInputWrapper struct {
	model TextInputModel
}

func (w textInputWrapper) Init() tea.Cmd {
	return w.model.Init()
}

func (w textInputWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			newModel, _ := w.model.Update(msg)
			w.model = newModel
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

func (w textInputWrapper) View() string {
	return w.model.View()
}
