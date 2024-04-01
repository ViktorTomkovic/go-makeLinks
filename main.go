// Package main contains the main function which
// provides a script making symlinks from a definition
// in a file in JSON format { "target" : [ "link1", "link2"] }
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// thank you https://stackoverflow.com/a/61837617

// StringSlice is a slice of strings used for more compact form
type StringSlice []string

// UnmarshalJSON is a custom unmarshalling implementation for string slices
func (ss *StringSlice) UnmarshalJSON(data []byte) error {
	fmt.Println(string(data))
	if data[0] == '[' {
		return json.Unmarshal(data, (*[]string)(ss))
	} else if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		*ss = append(*ss, s)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s linkFile\n", os.Args[0])
		return
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Working directory: %s\n", wd)
	sourceFile := os.Args[1]
	sourceBytes, err := os.ReadFile(sourceFile)
	if err != nil {
		panic(err)
	}
	var linkMap map[string]StringSlice
	err = json.Unmarshal([]byte(sourceBytes), &linkMap)
	if err != nil {
		panic(err)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	for target, links := range linkMap {
		fmt.Printf("%s %s\n", target, links)
		target = filepath.Join(wd, target)
		for _, link := range links {
			if strings.HasPrefix(link, "~/") {
				link = filepath.Join(home, link[2:])
			}
			lstat, err := os.Lstat(link)
			if errors.Is(err, os.ErrNotExist) {
				fmt.Println("Link location doesn't exist.")
				dir := filepath.Dir(link)
				err := os.MkdirAll(dir, os.FileMode(0755))
				if err != nil {
					panic(err)
				}
			} else if err != nil {
				panic(err)
			} else {
				// Assuming err is nil and lstat is not nil
				//fmt.Printf("%+v\n", lstat)
				if lstat.Mode().Type() == fs.ModeSymlink {
					fmt.Printf("Removing symlink '%s'\n", link)
					err := os.Remove(link)
					if err != nil {
						panic(err)
					}
				}
				if lstat.Mode().IsDir() || lstat.Mode().IsRegular() {
					fmt.Printf("Moving '%s' -> '%s'", link, link+".bak\n")
					err = os.Rename(link, link+".bak")
					if err != nil {
						panic(err)
					}
				}
			}
			fmt.Printf("Creating link '%s' -> '%s'\n", link, target)
			err = os.Symlink(target, link)
			if err != nil {
				fmt.Printf("%+v\n", err)
			}
		}
	}
}

