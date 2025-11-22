package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type createFileModel struct{
	currentDir string
	dirEntries []string
	cursor int
	textInput textinput.Model
	isTyping bool
	err error
}
func initCreateFileModel()tea.Model{
	ti := textinput.New()
	ti.Placeholder = "file name"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 20
	homeDir, _ := os.UserHomeDir()
	startingDir := filepath.Join(homeDir,"Workspace")
	homeDirEntries, _ := getFolderContent(startingDir)
	return createFileModel{
		currentDir: startingDir,
		dirEntries: homeDirEntries,
		textInput: ti,
		isTyping: false,
	}
}

func (cfm createFileModel) Init() tea.Cmd{
	return textinput.Blink
}


func (cfm createFileModel) View() string{
	s := "Files:\n"
	s += fmt.Sprintf("\n%s\n",cfm.textInput.View())
	for i, file := range cfm.dirEntries{
		cursor := ""
		if cfm.cursor == i{
			cursor = "<"
		}
		s += fmt.Sprintf("%s %s \n",file,cursor)
	}
	
	s += "\nescape to go back to menu\n" 
	return s
}


func (cfm createFileModel) Update(msg tea.Msg)(tea.Model,tea.Cmd){
	 	var cmd tea.Cmd
		switch msg := msg.(type){
			case tea.KeyMsg :
				if !cfm.isTyping{
					switch msg.String(){
						case "q", "ctrl+c":
							return cfm, tea.Quit
						case "up", "k":
							if cfm.cursor > 0{
								cfm.cursor --
							}
						case "down", "j":
							if cfm.cursor < len(cfm.dirEntries)-1{
								cfm.cursor ++
							}
						case "enter", " ":
							if cfm.cursor == 0 {
								cfm.currentDir, _ = getParentFolder(cfm.currentDir) 
								cfm.dirEntries, _= getFolderContent(cfm.currentDir)
							}else{
								cfm.currentDir = cfm.dirEntries[cfm.cursor]
								cfm.dirEntries, _ = getFolderContent(cfm.currentDir)
							}
						case "i" : 
							cfm.isTyping = true
						case "esc", "b":
							return cfm, func () tea.Msg{
								return NavigateBack{}
							}
					}
				}else{
					switch msg.String() {
						case "enter", " ":
							err := CreateFile(cfm.textInput.Value(),cfm.currentDir)
							cfm.err = err
							cfm.dirEntries, _ = getFolderContent(cfm.currentDir)

						case "esc":
							cfm.isTyping = false
				}
			}
		}
		
		cfm.textInput, cmd = cfm.textInput.Update(msg)
		return cfm, cmd
}


