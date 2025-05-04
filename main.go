package main

import (
    "fmt"
    "flag"
    "os"
    "sync"
    "path/filepath"
    "strings"
)

type Flags struct {
    Name    string
    Path    string
    Debug   bool
}

func colorTextRed(text string) string {
    const red = "\033[31m"
    const reset = "\033[0m"

    return red + text + reset
}

func process(path string, wg *sync.WaitGroup, flags Flags) {
    defer wg.Done()

    // Debug
    if flags.Debug {
        fmt.Println("Processing path: ", path)
    }

    // Retrieve the name of the file
    filename := filepath.Base(path)
    
    // Convert to lowercase for matching, unless TODO: flags are set
    filename = strings.ToLower(filename) 
    
    // Handle 'find' rules
    if strings.Contains(filename, flags.Name) {
        b := filename[:strings.Index(filename, flags.Name)]
        a := filename[strings.Index(filename, flags.Name) + len(flags.Name):]
        r := b + colorTextRed(flags.Name) + a
        fmt.Printf("%s/%s\n", filepath.Dir(path), r)
    }
}

func main() {
    name := flag.String("name", "", "Name of file to search for")
    rootPath := flag.String("path", "./", "Path in which to search for files")
    debug := flag.Bool("debug", false, "Enable debug output")
    flag.Parse()

    flags := Flags {
        Name: *name,
        Path: *rootPath,
        Debug: *debug,
    }
    
    var wg sync.WaitGroup

    err := filepath.WalkDir(flags.Path, func(path string, d os.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if d.IsDir() {
            return nil
        }

        wg.Add(1)
        go process(path, &wg, flags)

        return nil
    })

    if err != nil {
        fmt.Println("Error walking directory: ", err)
    }

    wg.Wait()
}

