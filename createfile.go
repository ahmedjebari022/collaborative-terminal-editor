package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	createFileTitleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(1, 0)

	directoryStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)

	fileItemStyle = lipgloss.NewStyle().
			Foreground(textColor).
			PaddingLeft(2)

	selectedFileItemStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Background(secondaryColor).
			PaddingLeft(1)

	inputContainerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accentColor).
			Padding(0, 1).
			MarginBottom(1)

	createFileContainerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(2, 4)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)
)

type createFileModel struct {
	currentDir string
	dirEntries []string
	cursor     int
	textInput  textinput.Model
	isTyping   bool
	err        error
}

func initCreateFileModel() tea.Model {
	ti := textinput.New()
	ti.Placeholder = "file name"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 40
	ti.PromptStyle = lipgloss.NewStyle().Foreground(primaryColor)
	ti.TextStyle = lipgloss.NewStyle().Foreground(textColor)

	homeDir, _ := os.UserHomeDir()
	startingDir := filepath.Join(homeDir, "Workspace")
	homeDirEntries, _ := getFolderContent(startingDir)
	return createFileModel{
		currentDir: startingDir,
		dirEntries: homeDirEntries,
		textInput:  ti,
		isTyping:   false,
	}
}

func (cfm createFileModel) Init() tea.Cmd {
	return textinput.Blink
}

func (cfm createFileModel) View() string {
	title := createFileTitleStyle.Render("📝 CREATE NEW FILE")

	currentDir := directoryStyle.Render(fmt.Sprintf("📁 %s", cfm.currentDir))

	// Input box
	inputLabel := lipgloss.NewStyle().Foreground(accentColor).Render("Enter filename:")
	var inputBox string
	if cfm.isTyping {
		inputBox = inputContainerStyle.Render(cfm.textInput.View())
	} else {
		inputBox = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#666666")).
				Render("Press 'i' to enter filename")
	}

	// Error message
	errorMsg := ""
	if cfm.err != nil {
		errorMsg = errorStyle.Render(fmt.Sprintf("❌ %s", cfm.err.Error()))
	}

	// File list
	var fileItems []string
	for i, file := range cfm.dirEntries {
		if cfm.cursor == i {
			fileItems = append(fileItems, selectedFileItemStyle.Render("▶ "+file))
		} else {
			fileItems = append(fileItems, fileItemStyle.Render("  "+file))
		}
	}
	fileList := lipgloss.JoinVertical(lipgloss.Left, fileItems...)

	// Help
	help := footerStyle.Render("i=input | enter=navigate | esc=back | ↑↓/jk=move")

	// Combine all
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		currentDir,
		"",
		inputLabel,
		inputBox,
		errorMsg,
		"",
		fileList,
		"",
		help,
	)

	styledContent := createFileContainerStyle.Render(content)

	return lipgloss.Place(
		100, 30,
		lipgloss.Center, lipgloss.Center,
		styledContent,
	)
}

func (cfm createFileModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !cfm.isTyping {
			switch msg.String() {
			case "q", "ctrl+c":
				return cfm, tea.Quit
			case "up", "k":
				if cfm.cursor > 0 {
					cfm.cursor--
				}
			case "down", "j":
				if cfm.cursor < len(cfm.dirEntries)-1 {
					cfm.cursor++
				}
			case "enter", " ":
				if cfm.cursor == 0 {
					cfm.currentDir, _ = getParentFolder(cfm.currentDir)
					cfm.dirEntries, _ = getFolderContent(cfm.currentDir)
				} else {
					cfm.currentDir = cfm.dirEntries[cfm.cursor]
					cfm.dirEntries, _ = getFolderContent(cfm.currentDir)
				}
			case "i":
				cfm.isTyping = true
			case "esc", "b":
				return cfm, func() tea.Msg {
					return NavigateBack{}
				}
			}
		} else {
			switch msg.String() {
			case "enter":
				err := CreateFile(cfm.textInput.Value(), cfm.currentDir)
				cfm.err = err
				if err == nil {
					cfm.textInput.SetValue("")
					cfm.isTyping = false
				}
				cfm.dirEntries, _ = getFolderContent(cfm.currentDir)
			case "esc":
				cfm.isTyping = false
				cfm.textInput.SetValue("")
			}
		}
	}

	cfm.textInput, cmd = cfm.textInput.Update(msg)
	return cfm, cmd
}


