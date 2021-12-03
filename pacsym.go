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

// Arguments
var pacsymArguments = []string{"-s", "--separate-builddir", "--local", "-l"}

// Common variabes
var (
	makearguments string
)

// Variables for sync
var (
	packagePath  = []string{}
	filePath     = []string{}
	trueFilePath = []string{}
)

// Variables for build
var (
	tarName        string
	tarPath        string
	url            string
	buildarguments = []string{}
)

// Variables for install

var (
	packageBuiltDir  string
	installarguments = []string{}
)

func help() {
	fmt.Printf("sync - Symlinks all installed packages in /usr/pkg\nbuild <URL> [OPTIONS] [MAKEFLAGS] - Compiles a package from url or local file\n--local -l - Compiles from local file, takes a filepath instead of url\n--separate-builddir -s - Compiles in a separate build directory\ninstall <PACKAGE NAME> <PACKAGE VERSION> [MAKEFLAGS] - Installs previously built package and gives it a name and a version.\nclean - Removes leftover sourcecode.")
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

func inList(x []string, y string) bool {
	for i := 0; i < len(x); i++ {
		if strings.Compare(y, x[i]) == 0 {
			return true
		}
	}
	return false
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
				packagePath = append(packagePath, "/usr/pkg/"+f.Name())
			}

			/// Reads package version
			for i := 0; i < len(packagePath); i++ {
				files, err := ioutil.ReadDir(packagePath[i])
				if err != nil {
					log.Fatal(err)
				}

				for _, f := range files {
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
				cmd := exec.Command("mkdir", "--verbose", "--parents", strings.Join(strings.Split(trueFilePath[i], "/")[:len(strings.Split(trueFilePath[i], "/"))-1], "/"))
				stdout, err := cmd.Output()
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println(string(stdout))
				cmd = exec.Command("ln", "-svf", filePath[i], trueFilePath[i])
				stdout, err = cmd.Output()
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println(string(stdout))
			}

			fmt.Println("Done.")

		// Downloads packages and compiles it with set flags
		case "build":
			buildarguments = os.Args[3:]

			// Checks if its a local package
			if inList(buildarguments, "--local") || inList(buildarguments, "-l") {

				// Finds and indexes package
				tarPath = os.Args[2]
				tarName = (strings.Split(tarPath, "/")[len(strings.Split(tarPath, "/"))-1])
			} else {

				// Downloads, finds and indexes package
				url = os.Args[2]
				cmd := exec.Command("wget", url)
				stdout, err := cmd.Output()
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println(string(stdout))
				tarName = (strings.Split(url, "/")[len(strings.Split(url, "/"))-1])
				tarPath = "./" + tarName
			}
			fmt.Println(tarName)

			// Moves package to /usr/pkgsrc/, aka build location.
			cmd := exec.Command("cp", "-v", tarPath, "/usr/pkgsrc/")
			stdout, err := cmd.Output()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println(string(stdout))

			// Extract tarball
			cmd = exec.Command("tar", "xvf", tarName, "-C", "/usr/pkgsrc/")
			stdout, err = cmd.Output()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			builddir := "/usr/pkgsrc/" + tarName[:len(tarName)-7]

			if inList(buildarguments, "--separate-builddir") || inList(buildarguments, "-s") {
				// Create separate builddir
				builddir = "/usr/pkgsrc/" + tarName[:len(tarName)-7] + "/build"
				cmd = exec.Command("mkdir", "-v", builddir)
				stdout, err = cmd.Output()
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println(string(stdout))
			}

			// Changes dir into builddir
			os.Chdir(builddir)

			// Removes pacsym specific arguments
			for i := 0; i < len(buildarguments); i++ {
				if inList(pacsymArguments, buildarguments[i]) == false {
					makearguments = makearguments + " " + buildarguments[i]
				}
			}

			fmt.Println("Configuring...")
			if inList(buildarguments, "--separate-builddir") || inList(buildarguments, "-s") {
				// Configures package if in builddir
				cmd = exec.Command("../configure", makearguments)
				stdout, err = cmd.Output()
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println(string(stdout))
			} else {
				// Configures package
				cmd = exec.Command("./configure", makearguments)
				stdout, err = cmd.Output()
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println(string(stdout))
			}

			// Builds package
			fmt.Println("Compiling...")
			cmd = exec.Command("make")
			stdout, err = cmd.Output()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println(string(stdout))

		case "install":
			packageName := os.Args[1]
			packageVer := os.Args[2]
			installarguments = os.Args[3:]

			// Discoveres built packages
			files, err := ioutil.ReadDir("/usr/pkgsrc")
			if err != nil {
				log.Fatal(err)
			}
			for _, f := range files {
				packageBuiltDir = "/usr/pkgsrc/" + f.Name()
			}

			os.Chdir(packageBuiltDir)

			// Attempts to chdir into built
			os.Chdir("build")

			for i := 0; i < len(installarguments); i++ {
				makearguments = makearguments + " " + installarguments[i]
			}

			fmt.Println("Installing...")
			destdir := "DESTDIR=/usr/pkg/" + packageName + "/" + packageVer + "/"
			cmd := exec.Command("make", "install", destdir, makearguments)
			stdout, err := cmd.Output()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println(string(stdout))
			fmt.Println("Done.")

		case "clean":
			cmd := exec.Command("rm -r * /usr/pkgsrc/*")
			stdout, err := cmd.Output()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println(string(stdout))
			fmt.Println("Done.")

		default:
			help()
		}
	} else {
		help()
	}
}
