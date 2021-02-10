package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	// pflag because --flag is more "normal" than -flag

	flag "github.com/spf13/pflag"

	"gopkg.in/yaml.v3"
)

func main() {

	// arguments/flags
	path := flag.String("path", "", "relative or absolute path to read files from")
	indent := flag.Int("indent", 2, "output yaml indentation level")
	hidden := flag.Bool("hidden", false, "including looking for hidden files")
	excluded := flag.String("exclude", "", "files to exclude")
	flag.Parse()

	// show a human friendly message if --path isn't specified

	if !isFlagPassed("path") {
		fmt.Println("# error, --path must be provided")
		fmt.Printf("\nUsage of yaml-concate:\n")
		flag.PrintDefaults()
		os.Exit(102)
	}

	matches, err := findYaml(*path, *excluded, *hidden)

	// sort them
	sort.Strings(matches)

	if err != nil {
		fmt.Printf("# error?\n")
		os.Exit(101)
	}
	for _, i := range matches {
		// TODO: check formatFile works BEFORE inserting "---" into the stream, avoiding potential for empty docs
		fmt.Println("---")
		fmt.Fprintf(os.Stderr, "# reading : %s\n", i)
		formatFile(i, *indent, false)
	}
}

func formatFile(f string, indent int, overwrite bool) {
	r, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	var out bytes.Buffer
	if e := formatStream(r, &out, indent); e != nil {
		log.Fatalf("Failed formatting YAML stream: %v", e)
	}
	r.Close()
	if e := dumpStream(&out, f, overwrite); e != nil {
		log.Fatalf("Cannot overwrite: %v", e)
	}
}

func formatStream(r io.Reader, out io.Writer, indent int) error {
	d := yaml.NewDecoder(r)
	in := yaml.Node{}
	err := d.Decode(&in)
	for err == nil {
		e := yaml.NewEncoder(out)
		e.SetIndent(indent)
		if err := e.Encode(&in); err != nil {
			log.Fatal(err)
		}
		e.Close()

		if err = d.Decode(&in); err == nil {
			fmt.Fprintln(out, "---")
		}
	}

	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

func dumpStream(out *bytes.Buffer, f string, overwrite bool) error {
	if overwrite {
		return ioutil.WriteFile(f, out.Bytes(), 0744)
	}
	_, err := io.Copy(os.Stdout, out)
	return err
}

func isHidden(path string) bool {
	// get just the filename without base path
	file := filepath.Base(path)
	// return true if its first character is a fullstop
	return strings.HasPrefix(file, ".")
}

func findYaml(root, excluded string, includeHidden bool) ([]string, error) {

	fmt.Fprintf(os.Stderr, "# searching path: %s\n", root)

	// acceptable yaml extension patterns
	yamlPatterns := []string{
		"*.yaml",
		"*.yml",
	}

	// we'll store a slice of files matching the pattern here
	var matches []string

	// attempt to walk a folder looking for files that match the extensions
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// TODO: feels bad :/ please improve me
		for _, i := range yamlPatterns {
			matched, err := filepath.Match(i, filepath.Base(path))
			if err != nil {
				// TODO: unsure how to replicate this situation, test me and improve the reporting..
				fmt.Printf("\n # bad error :/\n")
				os.Exit(100)
			}

			// the level of if's here is terribad.. please david fix me and write tests :/
			if matched {
				if includeHidden {
					matches = append(matches, path)
				} else {
					if !isHidden(path) {
						matches = append(matches, path)
					}
				}
			}
		}
		// finished..
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
