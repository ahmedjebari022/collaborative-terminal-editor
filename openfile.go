package main

import (
	"fmt"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)
type Mode string
const (
	NormalMode = "Normal"
	InsertMode = "Insert" 
)

var cursorStyle = lipgloss.NewStyle().
	Background(secondaryColor)

type fileModel struct{
	path string
	fileContent []string
	cp cursorPos
	vp viewPort
	m Mode
}
type cursorPos struct{
	x int
	y int 
}

type viewPort struct{
	ch int
	h int
}

func initFileModel(path string)tea.Model{
	fc, _ := readFile(path)
	for i,l := range fc{
		fc[i] = replaceTabulation(l)
	}
	vp := viewPort{
		ch: 0,
		h: 10,
	}
	if len(fc)< 10{
		vp.h = len(fc)
	}
	return fileModel{
		fileContent: fc,
		cp : cursorPos{
			x: 0,
			y: 0,
		},
		vp : vp,
		m : NormalMode,
	}
}

func (fm fileModel) Init() tea.Cmd{
	return nil
}

func Min(a, b int)int{
	if a < b {
		return a
	}
	return b	
}


func (fm fileModel) View() string{
    rl := len(fm.fileContent)-1-fm.vp.ch
    h := fm.vp.ch + fm.vp.h
    s := fmt.Sprintf("height :%d total lines :%d \n",h,len(fm.fileContent))
    min := Min(rl, fm.vp.h)
    for i := fm.vp.ch ; i < fm.vp.ch+min ; i++{
        if fm.cp.y == i{
            // Always get FRESH line data for cursor line
            line := fm.fileContent[i]
            
            if len(line) == 0 {
                // Empty line - render styled space
                c := cursorStyle.Render(" ")
                s += fmt.Sprintf("%s\n",c)
            } else if fm.cp.x >= len(line) {
                // Cursor past end of line - render line + styled space
                s += fmt.Sprintf("%s%s\n", line, cursorStyle.Render(" "))
            } else {
                // Normal case - cursor within line
                l := []rune(line)
                bc := string(l[:fm.cp.x])
                c := cursorStyle.Render(string(l[fm.cp.x]))
                ac := ""
                if fm.cp.x+1 < len(l) {
                    ac = string(l[fm.cp.x+1:])
                }
                s += fmt.Sprintf("%s%s%s\n", bc, c, ac)
            }
        } else {
            // Non-cursor line
            s += fmt.Sprintf("%s\n", fm.fileContent[i])
        }
    }
    return s
}

func (fm fileModel) Update(msg tea.Msg)(tea.Model, tea.Cmd){
			
	switch  msg := msg.(type){
		case tea.KeyMsg:
			if fm.m == NormalMode{
				switch msg.String(){
					case "ctrl+c":
						return fm,tea.Quit
					case "h":
						fm.MoveLeft()
					case "j":
						fm.MoveDown()
					case "k":
						fm.MoveUp()
					case "l":
						fm.MoveRight()
					case "i" : 
						fm.m = InsertMode
				}
			}else{
				switch msg.String(){
					case "esc":
						fm.m = NormalMode
					case "ctrl+c":
						return fm,tea.Quit
					case "backspace":
						fm.Delete()
					case "enter", " ":
						fm.AddLine()
					default :
						fm.Write(msg.String())
					}

			}

			}
			return fm,nil
}


func (fm *fileModel) MoveDown() {
	if fm.cp.y < len(fm.fileContent) - 1 {
		if fm.cp.y == fm.vp.ch + fm.vp.h - 1{
			fm.vp.ch += 1
		}
		if fm.cp.x >= len(fm.fileContent[fm.cp.y + 1]){
			if len(fm.fileContent[fm.cp.y + 1]) == 0{
				fm.cp.x = 0
			}else{
				fm.cp.x = len(fm.fileContent[fm.cp.y + 1]) - 1
			}
		}
		fm.cp.y += 1
		
	}
}

func (fm *fileModel) MoveRight(){
	if fm.cp.x < len(fm.fileContent[fm.cp.y]) - 1{
		fm.cp.x += 1
	}else {
		fm.MoveDown()
		fm.cp.x = 0
	}
}


func (fm *fileModel) MoveLeft(){
		if fm.cp.x > 0 {
			fm.cp.x -= 1
		}else{
			fm.MoveUp()
			fm.cp.x = len(fm.fileContent[fm.cp.y]) - 1
		}
}

func (fm *fileModel) MoveUp(){
	if fm.cp.y > 0 {
		
		if fm.cp.y == fm.vp.ch {
			fm.vp.ch -= 1
		}
		if fm.cp.x >= len(fm.fileContent[fm.cp.y - 1]){
			if len(fm.fileContent[fm.cp.y - 1]) == 0{
				fm.cp.x = 0
			}else{
				fm.cp.x = len(fm.fileContent[fm.cp.y - 1]) - 1
			}
		}
		fm.cp.y -= 1
	}

}

func (fm *fileModel)Write(input string) {
	if len(fm.fileContent) == 0{
		fm.fileContent = append(fm.fileContent, input)  
        fm.cp.x = 1  
        return
	}

	l := []rune(fm.fileContent[fm.cp.y])
	bi := l[:fm.cp.x]
	ai := l[fm.cp.x:]
	NewLine := string(bi) + input + string(ai)
	fm.fileContent[fm.cp.y] = NewLine
	fm.cp.x += 1	
}

func (fm *fileModel)Delete(){
	l := []rune(fm.fileContent[fm.cp.y])
	
	if len(l) == 0{
		fm.fileContent = slices.Delete(fm.fileContent,fm.cp.y,fm.cp.y+1)
		fm.MoveLeft()
		return
	}
	if fm.cp.x == 0{
		//delete the line
		//add the last line to the last line 
		ol := fm.fileContent[fm.cp.y]
		fm.MoveLeft()
		fm.fileContent = slices.Delete(fm.fileContent,fm.cp.y+1,fm.cp.y+2)
		fm.fileContent[fm.cp.y] += ol
		return
	}
	if fm.cp.x == len(fm.fileContent[fm.cp.y]){
		
	}
	bc := l[:fm.cp.x]
	ac := l[fm.cp.x + 1 :]
	NewLine := string(bc) + string(ac)
	fm.fileContent[fm.cp.y] = NewLine
	fm.MoveLeft()
}

func (fm *fileModel)AddLine(){
	newl := []string{""}
	fm.fileContent = append(fm.fileContent[:fm.cp.y+1],append(newl,fm.fileContent[fm.cp.y+1:]...)...)
	fm.MoveDown()
}


func replaceTabulation (line string)string{

	spacePerTab := 8
	replacement := strings.Repeat(" ", spacePerTab)

	return strings.ReplaceAll(line,"\t",replacement)

}