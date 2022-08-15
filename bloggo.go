package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/domonda/go-types/date"
	fs "github.com/ungerik/go-fs"
)

var (
	author    = "Erik Unger"
	blogTitle = "Blog of Erik Unger"

	buildDirName = "docs"

	projectDir fs.File
	sourceDir  fs.File
	buildDir   fs.File
	routes     = make(map[string]interface{})

	pageTemplate *template.Template
	postTemplate *template.Template

	pages []*pageData
	posts []*postData
)

type pageData struct {
	RootPath string
	Path     string
	Title    string
	Index    int
	Author   string
	Body     template.HTML
}

type postData struct {
	RootPath string
	Path     string
	Title    string
	Date     date.Date
	Author   string
	Body     template.HTML
}

type rootData struct {
	Title  string
	Author string
	Pages  []*pageData
	Posts  []*postData
}

func handlePageDir(pageDir fs.File) error {
	return nil
}

func handlePostDir(postDir fs.File) error {
	dirName := postDir.Name()
	if len(dirName) < len("2001-12-31_x") {
		return fmt.Errorf("post directory name too short: '%s'", dirName)
	}
	postDate := date.Date(dirName[:len("2001-12-31")])
	postSlug := dirName[len("2001-12-31_"):]

	fmt.Println(postDate, postSlug)

	bodyFile := postDir.Join("body.html")
	body, err := bodyFile.ReadAllString()
	if err != nil {
		return err
	}

	buildPostDir := buildDir.Join(postSlug)
	err = buildPostDir.MakeDir()
	if err != nil {
		return err
	}

	indexFile := buildPostDir.Join("index.html")
	writer, err := indexFile.OpenWriter()
	if err != nil {
		return err
	}
	defer writer.Close()

	data := &postData{
		RootPath: "../",
		Path:     postSlug,
		Title:    postSlug,
		Date:     postDate,
		Author:   author,
		Body:     template.HTML(body),
	}

	err = postTemplate.Execute(writer, data)
	if err != nil {
		return err
	}

	posts = append(posts, data)

	return nil
}

func main() {
	projectDir = fs.File(".")
	if len(os.Args) > 1 {
		projectDir = fs.File(os.Args[1])
	}
	projectDir = projectDir.ToAbsPath()
	if !projectDir.Exists() {
		fmt.Println("project directory does not exist:", projectDir)
		os.Exit(1)
	}
	sourceDir = projectDir.Join("source")
	if !sourceDir.Exists() {
		fmt.Println("source directory does not exist:", sourceDir)
		os.Exit(1)
	}
	buildDir = projectDir.Join(buildDirName)
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

	cssFile := sourceDir.Join("style.css")
	if cssFile.Exists() {
		err = fs.CopyFile(cssFile, buildDir)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	pageTemplate = template.Must(template.ParseFiles(sourceDir.Join("template.html").Path()))
	postTemplate = pageTemplate

	pagesDir := sourceDir.Join("pages")
	postsDir := sourceDir.Join("posts")
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

	writer, err := buildDir.Join("index.html").OpenWriter()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rootTemplate := template.Must(template.ParseFiles(sourceDir.Join("root.html").Path()))
	data := &rootData{
		Title:  blogTitle,
		Author: author,
		Pages:  pages,
		Posts:  posts,
	}
	err = rootTemplate.Execute(writer, data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
