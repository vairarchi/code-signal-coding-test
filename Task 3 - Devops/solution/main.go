package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// FileInfo ...
type FileInfo struct {
	Name            string
	Size            int
	Command         string
	CompressedSize  int
	CompressionRate string
}

// FileArray ...
var FileArray []FileInfo

// ByCmd ...
type ByCmd []FileInfo

func (a ByCmd) Len() int           { return len(a) }
func (a ByCmd) Less(i, j int) bool { return a[i].Command < a[j].Command }
func (a ByCmd) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// ByName ...
type ByName []FileInfo

func (a ByName) Len() int           { return len(a) }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func sortFileArray(arr []FileInfo) {
	sort.Slice(arr, func(i, j int) bool {
		switch strings.Compare(arr[i].Command, arr[j].Command) {
		case -1:
			return true
		case 1:
			return false
		}
		return arr[i].Name < arr[j].Name
	})
	// sort.Sort(ByCmd(arr))
}
func main() {
	commandsf, err := os.Open("/root/devops/commands.txt")

	if err != nil {
		log.Fatalf("failed to open")

	}

	cmds := readCommands(commandsf)
	commandsf.Close()

	readFilesToZip("/root/devops/files/", cmds)

}
func readCommands(commandsf *os.File) []string {
	scanner := bufio.NewScanner(commandsf)

	scanner.Split(bufio.ScanLines)
	var commands []string

	for scanner.Scan() {
		commands = append(commands, scanner.Text())
	}
	return commands
}

func readFilesToZip(rootpath string, commands []string) {
	currentFile := &FileInfo{}
	sort.Strings(commands)
	for _, cmd := range commands {
		temp := []FileInfo{}
		err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				currentFile = &FileInfo{}
				currentFile.Name = info.Name()
				currentFile.Command = cmd
				currentFile.Size = int(info.Size())
				if err := checkAndExecuteCmd(currentFile, path); err == nil {
					temp = append(temp, *currentFile)

				}
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
		sort.Sort(ByName(temp))
		for _, file := range temp {
			// gzip -1,a.txt,253%
			fmt.Printf("%v,%v,%v\n", file.Command, file.Name, file.CompressionRate)
		}
	}

}
func checkAndExecuteCmd(currentFile *FileInfo, path string) error {
	cmds := strings.Split(currentFile.Command, " ")
	cmds = append(cmds, "-c")
	cmds = append(cmds, path)
	out, err := exec.Command(cmds[0], cmds[1:]...).Output()
	if err != nil {
		return err
	}

	currentFile.CompressedSize = len(out)
	percent := math.Round(float64(currentFile.CompressedSize*100) / float64(currentFile.Size))
	currentFile.CompressionRate = fmt.Sprintf("%v", percent) + "%"
	return nil
}

// func readCommands(commandsf *os.File) []string {
// 	scanner := bufio.NewScanner(commandsf)

// 	scanner.Split(bufio.ScanLines)
// 	var commands []string

// 	for scanner.Scan() {
// 		commands = append(commands, scanner.Text())
// 	}
// 	return commands
// }

// func readFilesToZip(rootpath string, commands []string) {
// 	currentFile := &FileInfo{}
//     // sort.Strings(commands)
// 	for _, cmd := range commands {
//         tempArray := []FileInfo{}
// 		err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
// 			if !info.IsDir() {
// 				currentFile = &FileInfo{}
// 				currentFile.Name = info.Name()
// 				currentFile.Command = cmd
// 				currentFile.Size = int(info.Size())
// 				if err := checkAndExecuteCmd(currentFile, path); err == nil {
// 					tempArray = append(tempArray, *currentFile)

// 				}
// 			}
// 			return nil
// 		})
//         sort.Sort(ByName(tempArray))
//         FileArray = append(FileArray, tempArray...)
// 		if err != nil {
// 			log.Fatalf(err.Error())
// 		}
// 	}

// }
// func checkAndExecuteCmd(currentFile *FileInfo, path string) error {
// 	cmds := strings.Split(currentFile.Command, " ")
// 	cmds = append(cmds, "-c")
// 	cmds = append(cmds, path)
// 	out, err := exec.Command(cmds[0], cmds[1:]...).Output()
// 	if err != nil {
// 		return err
// 	}

// 	currentFile.CompressedSize = len(out)
// 	percent := math.Round(float64(currentFile.CompressedSize*100) / float64(currentFile.Size))
// 	currentFile.CompressionRate = fmt.Sprintf("%v", percent) + "%"
// 	return nil
// }
