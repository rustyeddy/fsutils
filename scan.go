package fsutils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func Scan(dirs []string) {
	for _, d := range os.Args[1:] {
		fmt.Printf("Walking Directory: %s\n", d)
		files, err := ioutil.ReadDir(d)
		fatalError(err)

		fmt.Printf("files %v", files)
	}
}

// fatalError fail if there is an Error
func fatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
