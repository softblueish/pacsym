package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	packagePath  = []string{}
	filePath     = []string{}
	trueFilePath = []string{}
)

func help() {
	fmt.Printf("sync - Symlinks all installed packages in /usr/pkg\n")
}

func walk(s string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if !d.IsDir() {
		filePath = append(filePath, s)
	}
	return nil
}

func main() {
	if cap(os.Args) != 1 {
		switch os.Args[1] {

		// Creates symbolic links for every package
		case "sync":

			// Reads packages available
			files, err := ioutil.ReadDir("/usr/pkg")
			if err != nil {
				log.Fatal(err)
			}

			for _, f := range files {
				fmt.Println(f.Name())
				packagePath = append(packagePath, "/usr/pkg/"+f.Name())
			}

			// Reads package version
			for i := 0; i < len(packagePath); i++ {
				files, err := ioutil.ReadDir(packagePath[i])
				if err != nil {
					log.Fatal(err)
				}

				for _, f := range files {
					fmt.Println(f.Name())
					packagePath[i] = packagePath[i] + "/" + f.Name()

					// Reads binary file
					filepath.WalkDir(packagePath[i], walk)
				}
			}

			// Turn into true paths
			for i := 0; i < len(filePath); i++ {
				trueFilePath = append(trueFilePath, "/"+strings.Join(strings.Split(filePath[i], "/")[5:], "/"))
			}

			// Create symlinks
			for i := 0; i < len(filePath); i++ {
				exec.Command("ln -svf " + filePath[i] + " " + trueFilePath[i])
				fmt.Println("Executed: " + "ln -svf " + filePath[i] + " " + trueFilePath[i])
			}

			fmt.Println("Done.")

		default:
			help()
		}
	} else {
		help()
	}
}
