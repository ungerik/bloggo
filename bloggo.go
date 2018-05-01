package main

import (
	"fmt"
	"os"

	"github.com/domonda/Domonda/go/types/date"
	"github.com/pkg/errors"
	fs "github.com/ungerik/go-fs"
)

var (
	projectDir fs.File
	sourceDir  fs.File
	buildDir   fs.File
	routes     = make(map[string]interface{})
)

func handlePageDir(pageDir fs.File) error {
	return nil
}

func handlePostDir(postDir fs.File) error {
	dirName := postDir.Name()
	if len(dirName) < len("2001-12-31_x") {
		return errors.Errorf("post directory name too short: '%s'", dirName)
	}
	postDate := date.Date(dirName[:len("2001-12-31")])
	postSlug := dirName[len("2001-12-31_"):]

	fmt.Println(postDate, postSlug)

	buildPostDir := buildDir.Relative(postSlug)
	err := buildPostDir.MakeDir()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	projectDir = fs.File(".")
	if len(os.Args) > 1 {
		projectDir = fs.File(os.Args[1])
	}
	projectDir = projectDir.MakeAbsolute()
	if !projectDir.Exists() {
		fmt.Println("project directory does not exist:", projectDir)
		os.Exit(1)
	}
	sourceDir = projectDir.Relative("source")
	if !sourceDir.Exists() {
		fmt.Println("source directory does not exist:", sourceDir)
		os.Exit(1)
	}
	buildDir = projectDir.Relative("build")
	err := buildDir.MakeAllDirs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = buildDir.RemoveDirContentsRecursive()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("bloggo project dir:", projectDir.Path())

	pagesDir := sourceDir.Relative("pages")
	postsDir := sourceDir.Relative("posts")
	if !pagesDir.Exists() && !postsDir.Exists() {
		fmt.Println("source/pages/ or source/posts/ must exist")
		os.Exit(1)
	}

	err = pagesDir.ListDir(handlePageDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = postsDir.ListDir(handlePostDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
