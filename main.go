package main

import (
	"fmt"
	"os"
	

	tea "github.com/charmbracelet/bubbletea"
)


type model struct{
	currentScreen tea.Model
	previousScreens []tea.Model
}


func initModel()tea.Model{
	screen := initialMenu()
	return model{
		currentScreen:screen,
		previousScreens: []tea.Model{},
	}
}

func (m model)Init()tea.Cmd{
	return nil
}


func (m model)Update(msg tea.Msg)(tea.Model,tea.Cmd){
	
	switch msg.(type){
		case NavigateToHelpMsg:
			helpScreen := initHelpmodel()
			m.previousScreens = append(m.previousScreens, m.currentScreen)
			m.currentScreen = helpScreen
		case NavigateToCreateFileMsg:
			createFileScreen := initCreateFileModel()
			m.previousScreens = append(m.previousScreens, m.currentScreen)
			m.currentScreen = createFileScreen
		case NavigateToEditFileMsg:
			//navigate
		case NavigateBack:
			m.currentScreen = m.previousScreens[len(m.previousScreens)-1]
			m.previousScreens = m.previousScreens[:len(m.previousScreens)-1]
		default :
			um, cmd := m.currentScreen.Update(msg)
			m.currentScreen = um
			return m, cmd
	}
		return m, nil
}


func (m model)View()string{
	return m.currentScreen.View()
}


func main() {
    p := tea.NewProgram(initModel())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}