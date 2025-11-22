package main

import tea "github.com/charmbracelet/bubbletea"



type helpmodel struct{
	helpText string
}


func initHelpmodel()helpmodel{
	return helpmodel{
		helpText: "Help",
	}
}

func (h helpmodel)Init()tea.Cmd{
	return nil
}

func (h helpmodel)Update(msg tea.Msg )(tea.Model, tea.Cmd){
	switch msg := msg.(type){
		case tea.KeyMsg:
			switch msg.String(){
				case "esc", "b":
					return h, func() tea.Msg{
						return NavigateBack{}
					}
			}		
	}
	return h,nil
}

func(h helpmodel)View()string{
	return h.helpText
}