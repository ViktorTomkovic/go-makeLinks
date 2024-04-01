package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

// thank you https://stackoverflow.com/a/61837617
type StringSlice []string

func (ss *StringSlice) UnmarshalJSON(data []byte) error {
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
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fmt.Println(exPath)
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(wd)
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s linkFile", os.Args[0])
		return
	}
	sourceFile := os.Args[1]
	fmt.Println(sourceFile)
	sourcePath := path.Join(wd, sourceFile)
	fmt.Println(sourcePath)
	sourceBytes, err := os.ReadFile(sourcePath)
	fmt.Println(string(sourceBytes))
	var linkMap map[string]StringSlice
	err = json.Unmarshal([]byte(sourceBytes), &linkMap)
	if err != nil {
		panic(err)
	}
	for key, value := range linkMap {
		fmt.Printf("%s %s\n", key, value)
	}
}
