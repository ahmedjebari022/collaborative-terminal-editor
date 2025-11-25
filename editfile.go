package main

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/charmbracelet/bubbles/filepicker"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

var (
    editFileTitleStyle = lipgloss.NewStyle().
        Foreground(accentColor).
        Bold(true).
        Padding(1, 0)
    
    filePickerContainerStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(secondaryColor).
        Padding(2, 4)
)

type editFileModel struct{
    filePicker filepicker.Model
    selectedFile string
}

func initEditFileModel()tea.Model{
    home, _ := os.UserHomeDir()
    path := filepath.Join(home, "Workspace/github/ahmedjebari022/collaborative-terminal-editor")	
    fp := filepicker.New()
    fp.CurrentDirectory = path
    fp.AllowedTypes = nil
    fp.Height = 15  // ← ADD THIS LINE
    
    // Style the filepicker
    fp.Styles.Cursor = lipgloss.NewStyle().Foreground(primaryColor)
    fp.Styles.Symlink = lipgloss.NewStyle().Foreground(accentColor)
    fp.Styles.Directory = lipgloss.NewStyle().Foreground(accentColor).Bold(true)
    fp.Styles.File = lipgloss.NewStyle().Foreground(textColor)
    fp.Styles.Selected = lipgloss.NewStyle().Foreground(primaryColor).Bold(true)
    
    return editFileModel{
        filePicker: fp,
    }
}

func (efm editFileModel) Init() tea.Cmd{
    return efm.filePicker.Init()
}

func (efm editFileModel)View() string{
    title := editFileTitleStyle.Render("📂 SELECT FILE TO EDIT")
    
    currentDir := directoryStyle.Render(fmt.Sprintf("Current: %s", efm.filePicker.CurrentDirectory))
    
    pickerView := efm.filePicker.View()
    
    help := footerStyle.Render("enter=select | esc=back | ↑↓=navigate")
    
    content := lipgloss.JoinVertical(
        lipgloss.Left,
        title,
        "",
        currentDir,
        "",
        pickerView,
        "",
        help,
    )
    
    styledContent := filePickerContainerStyle.Render(content)
    
    return lipgloss.Place(
        100, 30,
        lipgloss.Center, lipgloss.Center,
        styledContent,
    )
}

func (efm editFileModel)Update(msg tea.Msg)(tea.Model,tea.Cmd){
    switch msg := msg.(type){
        case tea.KeyMsg :
            switch msg.String(){
                case "esc":
                    return efm, func() tea.Msg{
                        return NavigateBack{}
                    }
                case "ctrl+c", "q":
                    return efm, tea.Quit
            }
    }
    
    var cmd tea.Cmd
    efm.filePicker, cmd = efm.filePicker.Update(msg)
    
    if didSelect, path := efm.filePicker.DidSelectFile(msg); didSelect{
        return efm, func() tea.Msg{
            return NavigateToOpenFile{
                filePath: path,
            }
        }
    }

    return efm, cmd
}