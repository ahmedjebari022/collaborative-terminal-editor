package main

import (
	"bufio"
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



func readFile(path string)([]string, error){
	data, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err){
			return nil, fmt.Errorf("file doesn't exist: %s",err.Error())
		}
			return nil, fmt.Errorf("error while opening the file: %s",err.Error())
	}
	if data.IsDir(){
		return nil, fmt.Errorf("this is a directory please select a file")
	}
	file, err := os.Open(path)		
	if err != nil {
		return nil, fmt.Errorf("can't open the file : %s",err.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	fileContent := []string{}
	for scanner.Scan() {
		fileContent = append(fileContent, scanner.Text())
	}
	if err := scanner.Err() ; err != nil {
			return nil, fmt.Errorf("error while reading the line: %s", err.Error())
		}

	return fileContent, nil
}

func updateFile (path string, content []string)error{
		data, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err){
				return fmt.Errorf("file doesn't exist: %v",err)
			}
			return fmt.Errorf("error when looking for file: %v",err)
		}
		if data.IsDir(){
			return fmt.Errorf("this is a directory file needed for this operation")
		}
		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("error when opening the file :%v",err)
		}
		defer file.Close()
		scanner := bufio.NewWriter(file)
		for i, l := range content{
			_, err := scanner.WriteString(l + "\n")
			if err != nil {
				scanner.Flush()
				return fmt.Errorf("error while writing to file at line %d: %v",i, err)
			}
		}
		scanner.Flush()
		return nil
	}


