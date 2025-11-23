package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)



type editFileModel struct{
	filePicker filepicker.Model
	selectedFile string
	
}

func initEditFileModel()tea.Model{
	// home, _ := os.UserHomeDir()
	// path := filepath.Join(home, "Workspace/github/ahmedjebari022/collaborative-terminal-editor")	
	fp :=  filepicker.New()
	fp.CurrentDirectory, _ = os.UserHomeDir()

	
	fp.AllowedTypes = nil
	return editFileModel{
		filePicker: fp,
		
	}
}

func (efm editFileModel) Init() tea.Cmd{
	return efm.filePicker.Init()
}
func (efm editFileModel)View() string{

	s := "Select A file:\n"
    s += fmt.Sprintf("Current Dir: %s\n", efm.filePicker.CurrentDirectory)
	s += fmt.Sprintf("%s",efm.filePicker.View())
	return s


}


func (efm editFileModel)Update(msg tea.Msg)(tea.Model,tea.Cmd){
		
			switch  msg := msg.(type){
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