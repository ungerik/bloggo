package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/domonda/Domonda/go/types/date"
	"github.com/pkg/errors"
	fs "github.com/ungerik/go-fs"
)

var (
	author = "Erik Unger"

	projectDir fs.File
	sourceDir  fs.File
	buildDir   fs.File
	routes     = make(map[string]interface{})

	indexTemplate *template.Template
)

type pageData struct {
	RootPath string
	Title    string
	Author   string
	Body     template.HTML
}

type postData struct {
	RootPath string
	Title    string
	Date     date.Date
	Author   string
	Body     template.HTML
}

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

	bodyFile := postDir.Relative("body.html")
	body, err := bodyFile.ReadAllString()
	if err != nil {
		return err
	}

	buildPostDir := buildDir.Relative(postSlug)
	err = buildPostDir.MakeDir()
	if err != nil {
		return err
	}

	indexFile := buildPostDir.Relative("index.html")
	writer, err := indexFile.OpenWriter()
	if err != nil {
		return err
	}
	defer writer.Close()

	data := &postData{
		RootPath: "../",
		Title:    postSlug,
		Date:     postDate,
		Author:   author,
		Body:     template.HTML(body),
	}

	err = indexTemplate.Execute(writer, data)
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

	cssFile := sourceDir.Relative("style.css")
	if cssFile.Exists() {
		err = fs.CopyFile(cssFile, buildDir)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	indexTemplate = template.Must(template.ParseFiles(sourceDir.Relative("template.html").Path()))

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
