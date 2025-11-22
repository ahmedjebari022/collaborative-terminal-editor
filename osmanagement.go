package main

import (
	"fmt"
	"os"
	"path/filepath"
)

//creating a file in a specific path
func CreateFile(file string,dir string)error{
	file = filepath.Base(file)
	_, err := os.Stat(dir)
	if err != nil{
		if os.IsNotExist(err){
			return fmt.Errorf("directory doesn't exist")
		}
		return fmt.Errorf("error with directory: %s",err.Error())
	}

	fp := filepath.Join(dir,file)
	_, err = os.Stat(fp)
	if err == nil{
		return fmt.Errorf("file already exists")
	}
	if !os.IsNotExist(err){
		return fmt.Errorf("probleme when verifying the file existance: %s",err.Error())
	}
	nf,err := os.Create(fp)
	if err != nil {
		return fmt.Errorf("failed to create the file%s",err.Error())
	}
	defer nf.Close()
	return nil
}


//Get files/folder in a sepecific folder
func getFolderContent(dir string)([]string,error){
	
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err){
			return []string{},fmt.Errorf("folder doens't exist: %s",err.Error())
		}
		return  []string{},fmt.Errorf("error when looking for dir: %s",err.Error())
	}
	files, err := os.ReadDir(dir)
	if err != nil {
		return  []string{},fmt.Errorf("error displaying folder content: %s",err.Error())
	}
	folderNames := []string{"..."}
	for _, f := range files{
		folderNames = append(folderNames, filepath.Join(dir,f.Name()))
	}
	return folderNames,nil
}


func getParentFolder(dir string)(string, error){
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err){
			return "",fmt.Errorf("dir doesn't exist :%s",err.Error())
		}
		return "",fmt.Errorf("folder ")
	}
	parentDir := filepath.Dir(dir)
	return parentDir,nil
}