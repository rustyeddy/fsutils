package fsutils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Scan this path and write an error if something goes wrong
func Scan(path string) {

	//chfi := make(chan os.FileInfo)
	//chdir := make(chan os.FileInfo)
	//go func() { chdir <- os.Args[1:] }()
	//doDirectories(path, chfi, chdir)

}

func doDirectory(d string, chfi, chdir chan os.FileInfo) *Error {
	files, err := ioutil.ReadDir(d)
	if err != nil {
		return NewError(d, err, "iotuil.ReadDir()")
	}

	for _, fi := range files {
		mode := fi.Mode()
		switch {
		case mode.IsDir():
			fmt.Printf("  dir: %s\n", fi.Name())
			chfi <- fi
		case mode.IsRegular():
			fmt.Printf("  reg: %s\n", fi.Name())
			chdir <- fi
		default:
			fmt.Printf(" Hmmm: %s probably perms .?.\n", fi.Name())
		}
	}
	return nil
}

// fatalError fail if there is an Error
func FatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func FatalNil(obj interface{}) {
	if obj == nil {
		log.Fatalf("obj %T was nil", obj)
	}
}

func doFile(fi os.FileInfo) {
	fmt.Printf("We have a file: \n\t%+v\n", fi)
}

func doDir(fi os.FileInfo) {
	fmt.Printf("We have a directory: \n")
}
