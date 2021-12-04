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

// pacsym-specific arguments
var pacsymArguments = []string{"-s", "--separate-builddir", "--local", "-l"}

// Common variabes among all options
var (
	makearguments = []string{}
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

// When this function is called, pacsym will display the help message.
func help() {
	fmt.Printf("\n" +
		"Commands\n" +
		"======================\n" +
		"sync - Symlinks all installed packages in /usr/pkg\n" +
		"build <URL> [OPTIONS] [MAKEFLAGS] - Compiles a package from url or local file\n" +
		"install <PACKAGE NAME> <PACKAGE VERSION> [MAKEFLAGS] - Installs a built package and gives it a name and a version.\n" +
		"clean - Removes leftover sourcecode.\n" +
		"\nFlags\n" +
		"======================\n" +
		"--local -l - Compiles from local file, takes a filepath instead of url\n" +
		"--separate-builddir -s - Compiles in a separate build directory\n" +
		"\n")
}

// This function is made to find all directories in a directory
func walk(s string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if !d.IsDir() {
		filePath = append(filePath, s)
	}
	return nil
}

// This will check if an element is in an array/slice and return a boolean output
func inList(x []string, y string) bool {
	for i := 0; i < len(x); i++ {
		if strings.Compare(y, x[i]) == 0 {
			return true
		}
	}
	return false
}

func main() {

	// This will make sure that the person executing pacsym put in options.
	if cap(os.Args) != 1 {

		// This switch statement activates depending on what they type next
		switch os.Args[1] {

		// When the user types sync, it is to find all packages in the /usr/pkg/ directory and symlink them according to its hiearchy
		case "sync":

			// This will read all directories in /usr/pkg, which will be used later
			files, err := ioutil.ReadDir("/usr/pkg")
			if err != nil {
				log.Fatal(err)
			}

			// This for loop will make sure the directory of every package in /usr/pkg is catalouged in the packagePath array
			for _, f := range files {
				packagePath = append(packagePath, "/usr/pkg/"+f.Name())
			}

			// This nested for loop will look in every package with the help of the packagePath array and add the version of the package to the packagePath array
			// An example would be /usr/pkg/foobar/1.1, which would be an element of /usr/pkg
			for i := 0; i < len(packagePath); i++ {
				files, err := ioutil.ReadDir(packagePath[i])
				if err != nil {
					log.Fatal(err)
				}

				for _, f := range files {
					packagePath[i] = packagePath[i] + "/" + f.Name()

					// When we reach this point in the for loop, we will index every file in the directory who are elements in the packagePath array and place them in the filepath array
					filepath.WalkDir(packagePath[i], walk)
				}
			}

			// The for loop described under will tear appart every element in filePath, remove /usr/pkg/<PACKAGENAME>/<PACKAGEVERSION>/ and make it an element of trueFilePath
			// An example would be the element /usr/pkg/foobar/1.1/usr/bin/dependency.sh would turn into just /usr/bin/dependency.sh
			for i := 0; i < len(filePath); i++ {
				trueFilePath = append(trueFilePath, "/"+strings.Join(strings.Split(filePath[i], "/")[5:], "/"))
			}

			// Now it's time for the main attraction, linking every element in filePath to trueFilePath, this will recursively repeat for every element present in the array filePath
			// An example of a symlink would be /usr/bin/dependency.sh would be symlinked to /usr/pkg/foobar/1.1/usr/bin/dependency.sh
			for i := 0; i < len(filePath); i++ {
				cmd := exec.Command("mkdir", "--verbose", "--parents", strings.Join(strings.Split(trueFilePath[i], "/")[:len(strings.Split(trueFilePath[i], "/"))-1], "/"))
				stdout, err := cmd.CombinedOutput()
				if err != nil {
					fmt.Println(fmt.Sprint(err) + ": " + string(stdout))
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

			// At this point the pacsym program will end.
			fmt.Println("Done.")

		// If the "build" option was chosen, the package manager will attempt to place a tarball in /usr/pkgsrc, extract it, configure it and finally, build it.
		case "build":

			// When building a package, it might require certian arguments to work correctly, that means everything written after pacsym build <tarball> will be counted as a build argument, unless it's a pacsym argument
			buildarguments = os.Args[3:]

			// Here's the first usecase of the pacsym arguments, "--local" or "-l", this will check if it exists in the buildarguments array. If it is present, it will just index the tarball, if its not, it will download the url present with wget, then index it.
			if inList(buildarguments, "--local") || inList(buildarguments, "-l") {

				// This will check os.Args[2], which the program assumes to be the location of the tarball, and index it to tarPath.
				tarPath = os.Args[2]

				// This will take tarPath and take the last word before it finds a /, which will leave just the name of the tarball and index it to tarName.
				tarName = (strings.Split(tarPath, "/")[len(strings.Split(tarPath, "/"))-1])
			} else {

				// If the --local option wasnt specified, os.Args will assume that the thing after "build" written is the url to the tarball, then index that url to the url variable.
				url = os.Args[2]

				// After which, a wget command will be executed and download the tarball.
				cmd := exec.Command("wget", url)
				stdout, err := cmd.CombinedOutput()
				if err != nil {
					fmt.Println(fmt.Sprint(err) + ": " + string(stdout))
					return
				}
				fmt.Println(string(stdout))

				// This will take the url, find the last segment of it which would be the name of the tarball and index it.
				tarName = (strings.Split(url, "/")[len(strings.Split(url, "/"))-1])

				// This is self explanitory, it takes the name of the tarball and adds ./, which means current directory.
				tarPath = "./" + tarName
			}

			// This executes a cp command, which will copy the tarball to the /usr/pkgsrc directory.
			cmd := exec.Command("cp", "-v", tarPath, "/usr/pkgsrc/")
			stdout, err := cmd.Output()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println(string(stdout))

			// After which, the tarball will be extracted into /usr/pkgsrc.
			cmd = exec.Command("tar", "xvf", tarName, "-C", "/usr/pkgsrc/")
			stdout, err = cmd.Output()
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			// The tarball will be deleted to not waste space on the drive and keep pacsym functioning.
			pkgSrcTarDir := "/usr/pkgsrc/" + tarName
			cmd = exec.Command("rm", pkgSrcTarDir)
			stdout, err = cmd.Output()
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			// The directory it was extracted to will be indexed into buildddir, pacsym assumes the length of the tarball extensions is 7 characters (example: .tar.xz and .tar.gz)
			extensionLength := 7
			builddir := "/usr/pkgsrc/" + tarName[:len(tarName)-extensionLength]

			// This part will find the directory of the configure file, which is used for creating the makefile.
			conflocation := builddir + "/configure"

			// Now pacsym will check if the package requires for it to be built in a separate build directory, which would require --separate-builddir to be passed, or -s
			if inList(buildarguments, "--separate-builddir") || inList(buildarguments, "-s") {

				//This will create a folder called "build" in the directory the tarball created. This folder will now be the path for builddir.
				builddir = "/usr/pkgsrc/" + tarName[:len(tarName)-extensionLength] + "/build"
				cmd = exec.Command("mkdir", "-v", builddir)
				stdout, err = cmd.CombinedOutput()
				if err != nil {
					fmt.Println(fmt.Sprint(err) + ": " + string(stdout))
					return
				}
				fmt.Println(string(stdout))
			}

			// At this point pacsym will change the current directory to the build directory, be it the tarball folder or the "build" folder
			os.Chdir(builddir)

			// Now pacsym specific arguments will be removed from buildarguments as they are no longer used.
			fmt.Println("Configuring...")
			for i := 0; i < len(buildarguments); i++ {
				if inList(pacsymArguments, buildarguments[i]) == false {
					makearguments = append(makearguments, buildarguments[i])
				}
			}

			// This will execute the configure script and pass all the elements of makearguments into it.
			cmd = exec.Command(conflocation, makearguments...)
			stdout, err = cmd.CombinedOutput()
			if err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + string(stdout))
				return
			}
			fmt.Println(string(stdout))

			// Now the fun part begins, a make command will be passed, which will start compiling the software.
			fmt.Println("Compiling...")
			cmd = exec.Command("make")
			stdout, err = cmd.Output()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println(string(stdout))

			// At this point the package is configured but not installed, this is important if the user needs to make changes before it's compiled. Pacsym ends here.
			fmt.Println("Done.")

		// When the install option is passed, it's going to discover the previously built package and put it into place, give it a name and a version.
		case "install":

			// When executing install, it requires a <PACKAGENAME> and <PACKAGEVERSION>. It will take these and index them. It will also assume everything after written are install arguments.
			packageName := os.Args[2]
			packageVer := os.Args[3]
			installarguments = os.Args[4:]

			// This if and for loop will look in /usr/pkgsrc for the extracted tarball and index it in packageBuiltDir
			files, err := ioutil.ReadDir("/usr/pkgsrc")
			if err != nil {
				log.Fatal(err)
			}

			for _, f := range files {
				packageBuiltDir = "/usr/pkgsrc/" + f.Name()
			}

			// Now pacsym will enter the package directory.
			os.Chdir(packageBuiltDir)

			// This will read every file in the package directory to look for the build folder, which is important if you made a separate build directory.
			files, err = ioutil.ReadDir(packageBuiltDir)
			if err != nil {
				log.Fatal(err)
			}

			for _, f := range files {
				// If "build" is found, it will index it into packageBuiltBuildDir and enter that directory.
				if f.Name() == "build" {
					packageBuiltBuildDir := packageBuiltDir + "/build"
					os.Chdir(packageBuiltBuildDir)
				}
			}

			// Now any installarguments passed that are pacsym specific will be removed. This code is currently useless but good for futureproofing pacsym.
			for i := 0; i < len(installarguments); i++ {
				if inList(pacsymArguments, installarguments[i]) == false {
					makearguments = append(makearguments, installarguments[i])
				}
			}

			// Pacsym will now create the directory where package will be installed into.
			// Example, if the package name was set to "Foobar" and the version was set to 1.1, it will create the directory /usr/pkg/Foobar/1.1/
			packageInstallDir := "/usr/pkg/" + packageName + "/" + packageVer + "/"
			cmd := exec.Command("mkdir", "--parents", "--verbose", packageInstallDir)
			stdout, err := cmd.Output()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println(string(stdout))

			// Now the packages will be installed, Where one of the default paramaters added by pacsym is DESTDIR=/usr/pkg/packageName/packageVersion, which is the core of pacsyms functions.
			// It will execute make install DESTDIR=/usr/pkg/packageName/packageVersion/ and then with the parameters passed with makearguments.
			fmt.Println("Installing...")
			param := "DESTDIR=/usr/pkg/" + packageName + "/" + packageVer + "/"
			makeinstall := append([]string{"install", param}, makearguments...)
			cmd = exec.Command("make", makeinstall...)
			stdout, err = cmd.CombinedOutput()
			if err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + string(stdout))
				return
			}
			fmt.Println(string(stdout))

			// Here pacsym will end, it's installed but it cant be used yet, it needs to be symlinked into place with pacsym sync
			fmt.Println("Done.")

		// If the clean option is passed, it will remove everything in /usr/pkgsrc, which is important as pacsym can't build a new package if the /usr/pkgsrc directory inst empty.
		case "clean":

			// Here pacsym will discover every file in /usr/pkgsrc and execute rm -r /usr/pkgsrc/sourcecode
			files, err := ioutil.ReadDir("/usr/pkgsrc")
			if err != nil {
				log.Fatal(err)
			}
			for _, f := range files {
				toRemove := "/usr/pkgsrc/" + f.Name()
				cmd := exec.Command("rm", "-r", toRemove)
				stdout, err := cmd.CombinedOutput()
				if err != nil {
					fmt.Println(fmt.Sprint(err) + ": " + string(stdout))
					return
				}
				fmt.Println(string(stdout))
			}

			// Here pacsym will end.
			fmt.Println("Done.")

		default:
			// If no option was passed, it will print the help menu
			help()
		}
	} else {
		// If no arguments were passed at all, it will print the help menu
		help()
	}
}
