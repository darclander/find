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
    Name        string
    Path        string
    Debug       bool
    IncludeDirs bool
    ExactMatch  bool
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
    originalName := filename
    
    // Convert to lowercase for matching, unless exactmatch is set
    if (!flags.ExactMatch) {
        filename = strings.ToLower(filename)
    }
    
    // Handle 'find' rules
    if strings.Contains(filename, flags.Name) {
        b := originalName[:strings.Index(filename, flags.Name)]
        a := originalName[strings.Index(filename, flags.Name) + len(flags.Name):]
        r := b + colorTextRed(flags.Name) + a
        fmt.Printf("%s/%s\n", filepath.Dir(path), r)
    }
}

func main() {
    name := flag.String("name", "", "Name of file to search for.")
    rootPath := flag.String("path", "./", "Path in which to search for files. (Default current directory)")
    debug := flag.Bool("v", false, "Enable debug (verbose) output.")
    includeDirs := flag.Bool("d", false, "Include directories in the search.")
    exactMatch := flag.Bool("e", false, "Exact matching, e.g., case sensitive search.")

    flag.Parse()

    flags := Flags {
        Name: *name,
        Path: *rootPath,
        Debug: *debug,
        IncludeDirs: *includeDirs,
        ExactMatch: *exactMatch,
    }
    
    var wg sync.WaitGroup

    err := filepath.WalkDir(flags.Path, func(path string, d os.DirEntry, err error) error {
        if err != nil {
            return err
        }

        if d.IsDir() && !flags.IncludeDirs {
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

