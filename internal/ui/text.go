package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func WriteTextStep(key string, value string) {
	var b strings.Builder
	textKeyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("51"))
	b.WriteString(textKeyStyle.Render(key+": ") + value)
	fmt.Println(b.String())
}

func WriteTextSuccess(text string) {
	var b strings.Builder
	textSuccessStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
	b.WriteString(textSuccessStyle.Render(text))
	fmt.Println(b.String())
}

func WriteTextError(text string) {
	var b strings.Builder
	textErrorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("124"))
	b.WriteString(textErrorStyle.Render(text))
	fmt.Println(b.String())
}
