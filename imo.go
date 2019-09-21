/*
 * Image Organizer
 * @author Benjamin Lee
 * @description Extract images from sub-folders into a single directory
 */

package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// version
const VER_MAJ int = 1 // major
const VER_MIN int = 0 // minor
const VER_REV int = 0 // revision

// options
var optIn string       // input directory
var optOut string      // output directory
var optExt string      // file extensions
var optDepth int       // search depth
var optVerboseErr bool // show error messages
var optVerboseAll bool // show all messages
var optScanOnly bool   // scan without copy

// runtime variables
var id int = 0      // image ID
var found int = 0   // qualified files
var copied int = 0  // files copied
var extArr []string // split optExt into string array

// error counters
var failed int = 0            // failed operations
var dirError int = 0          // failed to read from directory
var copyError int = 0         // failed to copy
var depthLimitReached int = 0 // stopped by maximum depth, you may want to raise the value of -d to do a deeper search

/*
 * Initialize options
 * set default values and help messages for options using package flag
 * @see https://golang.org/pkg/flag/
 */
func initOpts() {
	flag.StringVar(&optIn, "i", ".", "input directory")
	flag.StringVar(&optOut, "o", "image-organizer", "output directory")
	flag.StringVar(&optExt, "e", "jpg|jpeg|png|bmp", "file extensions")
	flag.IntVar(&optDepth, "d", 10, "search depth")
	flag.BoolVar(&optVerboseErr, "v", false, "show error log")
	flag.BoolVar(&optVerboseAll, "vv", false, "show error and message logs")
	flag.BoolVar(&optScanOnly, "s", false, "search without copy")
}

/*
 * Process a given directory
 * @param from	search this directory for images
 * @param to    once found, copy image to this directory
 * @param depth stop when exceeding optDepth
 */
func processDir(from string, to string, depth int) {
	// stop if we've reached maximum depth
	if depth > optDepth {
		depthLimitReached++ // record this incident
		return
	}
	// don't copy to itself
	if from == to {
		return
	}
	// scan directory specified by from
	// @see https://golang.org/pkg/io/ioutil/#ReadDir
	files, err := ioutil.ReadDir(from)
	// if we encounter an directory error, this would likely to be
	// 1. directory not exist
	// 2. directory permissions
	// TODO: show suggestions depending on different errors
	if err != nil {
		dirError++ // record this incident
		failed++
		if optVerboseErr || optVerboseAll { // TODO: replace by log level in integer
			fmt.Fprintln(os.Stderr, err.Error())
		}
		return
	}
	// if we successfully read the directory,
	// parse its files/sub-directories
	for _, file := range files {
		if file.IsDir() { // if we find a directory, search it
			processDir(filepath.Join(from, file.Name()), to, depth+1)
		} else { // if we find a file, get its properties
			var filename string = file.Name()                        // get filename
			var ext string = strings.ToLower(filepath.Ext(filename)) // convert extension to lowercase for easier filtering
			// exclude system files
			if filename == ".DS_STORE" || filename == "thumb.db" || filename == "Thumb.db" {
				continue
			}
			// filter extension
			var validExt bool = false // valid extension flag
			for i := 0; i < len(extArr); i++ {
				if "."+extArr[i] == ext {
					validExt = true
					break // don't need to check the rest if we've got a correct one
				}
			}
			if validExt { // if extension is valid
				found++ // record this incident
			} else {
				continue
			}
			if optScanOnly { // skip copy if -s is enabled
				if optVerboseAll { // TODO: replace by log level
					fmt.Println(filepath.Join(from, filename))
				}
				continue
			}
			// copy file
			var cpFrom string = filepath.Join(from, filename) // copy from
			id++
			var cpTo = filepath.Join(to, strconv.Itoa(id)+ext) // copy to
			if optVerboseAll {                                 // TODO: replace by log level
				fmt.Println("\"" + cpFrom + "\",\"" + cpTo + "\"")
			}
			var err = copy(cpFrom, cpTo) // copy
			if err != nil {              // if we encounter an error in copy process
				failed++ // record this incident
				copyError++
				if optVerboseErr || optVerboseAll { // TODO: replace by log level
					fmt.Fprintln(os.Stderr, err.Error())
				}
			} else {
				copied++ // record how many files were copied
			}
		}
	}
}

/*
 * Copy a single file from one place to another
 */
func copy(from string, to string) error {
	in, err := os.Open(from)

	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(to)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func main() {
	// initialize options
	initOpts()
	// parse options
	flag.Parse()
	if !flag.Parsed() { // if flag failed to parse options
		fmt.Fprintln(os.Stderr, "failed to parse options")
		os.Exit(1)
	}
	// parse extension string specified in -e
	extArr = strings.Split(optExt, "|")
	if len(extArr) == 0 { // if we've got an empty string
		fmt.Fprintln(os.Stderr, "failed to prase extension string")
		os.Exit(2)
	}
	// convert pathes given by -i and -o to absolute pathes
	absIn, errIn := filepath.Abs(optIn)
	if errIn != nil {
		fmt.Fprintln(os.Stderr, errIn.Error())
		os.Exit(3)
	}
	absOut, errOut := filepath.Abs(optOut)
	if errOut != nil {
		fmt.Fprintln(os.Stderr, errOut.Error())
		os.Exit(4)
	}
	// create output directory if not exists
	os.Mkdir(absOut, os.ModePerm)
	// process directory
	processDir(absIn, absOut, 0)
	// show result
	fmt.Println("")
	fmt.Printf("Image Organizer v%d.%d.%d    ", VER_MAJ, VER_MIN, VER_REV)
	fmt.Println("")
	fmt.Println("")
	fmt.Println("Found", found, "files with extension", optExt, "under directory")
	fmt.Println(absIn)
	if copied != 0 {
		fmt.Println("Copied", copied, "files to directory")
		fmt.Println(absOut)
	}
	if failed != 0 {
		fmt.Println("Encountered", failed, "failures, including", copyError, "copy failures and", dirError, "directory failures")
	}
	if depthLimitReached != 0 {
		fmt.Println("Stopped at maximum depth", optDepth, "for", depthLimitReached, "times ")
	}
	fmt.Println("")
	fmt.Println("\"imo -h\" for help")
	fmt.Println("")
	os.Exit(0)
}
