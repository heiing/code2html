// code2html project main.go
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

var script, input, output, filter, vimfile string
var input_length int

func usage() {
	fmt.Println("Usage: ", os.Args[0], " ScriptFile FileFilter InputPath OutputPath")
	fmt.Println("ScriptFile   for example: D:\\tohtml.vim")
	fmt.Println("FileFilter   for example: *.c or *.c,*.h or *")
	fmt.Println("InputPath    for example: D:\\test")
	fmt.Println("OutputPath   for example: D:\\testout")
	os.Exit(0)
}

func err(msg string) {
	fmt.Println("Error: ", msg, "\n")
}

func fileexist(f string) (bool, bool) {
	fi, e := os.Stat(f)
	if (e == nil) && (fi != nil) {
		return true, fi.IsDir()
	}
	return false, false
}

func filemove(from string, to string) error {
	content, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}
	werr := ioutil.WriteFile(to, content, os.ModeAppend)
	if werr != nil {
		return werr
	}
	return nil
}

func check_args() {
	if len(os.Args) < 5 {
		err("missing arguments")
		usage()
	}

	script_exists, script_isdir := fileexist(os.Args[1])
	if !script_exists || script_isdir {
		err(os.Args[1] + " not exists or not file")
		usage()
	}

	inputdir_exists, inputdir_isdir := fileexist(os.Args[3])
	if !inputdir_exists || !inputdir_isdir {
		err(os.Args[3] + " not exists or not dir")
		usage()
	}

	if (len(os.Args[4]) > len(os.Args[3])) &&
		(os.Args[4][0:len(os.Args[3])] == os.Args[3]) {
		err(os.Args[4] + " can not in " + os.Args[3])
		os.Exit(1)
	}

	outputdir_exists, outputdir_isdir := fileexist(os.Args[4])
	if outputdir_exists && !outputdir_isdir {
		err(os.Args[4] + " not dir")
		usage()
	} else if !outputdir_exists {
		mkdirerror := os.MkdirAll(os.Args[4], os.ModeDir)
		if mkdirerror != nil {
			err("create dir " + os.Args[4] + " faild")
			os.Exit(1)
		}
	}
}

func dirwalker(path string, info os.FileInfo, wk_err error) error {
	if wk_err != nil {
		return wk_err
	}
	if info.IsDir() {
		if len(path) > input_length {
			target_dir := output + path[input_length:]
			exists, _ := fileexist(target_dir)
			if !exists {
				mkdir_err := os.Mkdir(target_dir, os.ModeDir)
				if mkdir_err != nil {
					err("create dir " + target_dir + " faild")
					os.Exit(1)
				}
			}
		}
	} else if len(path) > input_length {
		cmd := exec.Command(vimfile, "-S", script, path)
		run_err := cmd.Run()
		if run_err == nil {
			htmlfile := path + ".html"
			if htmlexists, _ := fileexist(htmlfile); htmlexists {
				moveerr := os.Rename(htmlfile, output+path[input_length:]+".html")
				if moveerr == nil {
					fmt.Println("convert " + path + " OK.")
				}
			}
		} else {
			fmt.Println(run_err)
		}
	}

	return nil
}

func deal() {
	vimfile = "vim"
	if len(os.Args) >= 6 {
		if vimfile_exists, vimfile_isdir := fileexist(os.Args[5]); vimfile_exists && !vimfile_isdir {
			vimfile = os.Args[5]
		}
	}
	script = os.Args[1]
	filter = os.Args[2]
	input = os.Args[3]
	output = os.Args[4]

	input_length = len(input)

	filepath.Walk(input, dirwalker)
}

func main() {
	check_args()
	deal()
}
