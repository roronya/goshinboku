package main

import (
	"archive/zip"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/roronya/goshinboku/jira"
	"github.com/roronya/goshinboku/xmind"
	"io"
	"log"
	"os"
)

var dryrun bool

/**
This program is rewrite a content.json in Xmind file.
The process is below.
1. Unzip a xmind file to get files. Xmind file is actually zip file.
2. Find a content.json in the files, create a ticket and rewrite the content.json.
3. Create a new zip file using the content.json and rest of the files.
4. Remove an old xmind file and rename the new zip file.
*/
func main() {
	flag.BoolVar(&dryrun, "dryrun", false, "skip creating jira tickets if dryrun option is true")
	flag.Parse()
	in := flag.Arg(0)
	if in == "" {
		flag.Usage()
		os.Exit(1)
	}

	// 1. Unzip a xmind file to get files. Xmind file is actually zip file.
	zr, err := zip.OpenReader(in)
	if err != nil {
		log.Fatal(err)
	}
	defer zr.Close()

	// 2. Find a content.json in the files, create a ticket and rewrite the content.json.
	f, err := findContentJsonFile(zr.File)
	if err != nil {
		log.Fatal(err)
	}
	fr, err := f.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer fr.Close()

	dec := json.NewDecoder(fr)
	var c xmind.Contents
	for {
		if err := dec.Decode(&c); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
	}

	// We are coding a magic number 0.
	// Because the content.json's root object is array that just has one element.
	r := c[0].RootTopic
	r.ParseTitle()
	// It is a specification that a project and an epic must be set RootObject.
	if r.Project == "" || r.Epic == "" {
		log.Fatal("RootTopic must be set project and epic")
	}

	leaves := r.FindLeaves()

	if dryrun == true {
		fmt.Println("skipped creating jira tickets because dryrun option is true")
		fmt.Println("below tickets will be create")
		fmt.Printf("project:%s\ncomponent:%s\nepic:%s\n", r.Project, r.Component, r.Epic)
		for i, leaf := range leaves {
			fmt.Printf("ticket %d: %s\n", i, leaf.Title)
		}
	} else {
		user := os.Getenv("JIRA_USER")
		pass := os.Getenv("JIRA_PASSWORD")
		server := os.Getenv("JIRA_SERVER")
		client, err := jira.NewClient(user, pass, server)
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < len(leaves); i++ {
			leaves[i].ParseTitle()
			if leaves[i].IssueId == "" {
				issue, url, err := jira.IssueCreate(client, r.Project, r.Component, r.Epic, "Task", leaves[i].Title, "", "", "")
				if err != nil {
					log.Printf("got an error creating a ticket titled \"%s\". error is below:\n%s", leaves[i].Title, err)
					continue
				}
				leaves[i].SetIssueId(issue.Key)
				fmt.Printf("created a ticket\n\"%s\"\n %s\n", leaves[i].Title, url)
			} else {
				issue, err := jira.GetIssue(client, leaves[i].IssueId)
				if err != nil {
					fmt.Printf("issue %s is not found. Please remove IssueId property and rerun to remake a ticket.", leaves[i].IssueId)
				}
				// FIXME: Each jira projects have different status and workflow. We should use config to handle it.
				switch issue.Fields.Status.Name {
				case "オープン", "再オープン":
					leaves[i].SetMakerAsTodo()
				case "進行中", "In Review":
					leaves[i].SetMarkerAsProgress()
				case "解決済み", "クローズ":
					leaves[i].SetMakerAsDone()
				default:
					log.Fatalf("invalid ticket status %s", issue.Fields.Status.Name)
				}
				fmt.Printf("issue %s is updated the status.\n", issue.Key)
			}
		}
	}

	// convert struct to json
	j, err := json.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}

	// 3. Create a new zip file using the content.json and rest of the files.
	z, err := os.CreateTemp("", "new.xmind")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(z.Name())

	if err := save(zr.File, j, z); err != nil {
		log.Fatal(err)
	}

	// 4. Remove an old xmind file and rename the new zip file.
	zr.Close() // close before removing
	if err := os.Remove(in); err != nil {
		log.Fatal(err)
	}
	if err := os.Rename(z.Name(), in); err != nil {
		log.Fatal(err)
	}
}

func findContentJsonFile(files []*zip.File) (*zip.File, error) {
	for _, f := range files {
		if f.Name == "content.json" {
			return f, nil
		}
	}
	return nil, fmt.Errorf("cannot find content.json")
}

func save(files []*zip.File, c []byte, z *os.File) error {
	zw := zip.NewWriter(z)
	defer zw.Close()

	fw, err := zw.Create("content.json")
	if err != nil {
		return err
	}
	if _, err := fw.Write(c); err != nil {
		return err
	}
	for _, file := range files {
		if file.Name == "content.json" {
			continue
		}
		if err := write(zw, file); err != nil {
			return err
		}
	}
	return nil
}

/**
Write files to zip writer to safety close in loop.
*/
func write(zw *zip.Writer, file *zip.File) error {
	f, err := file.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	fw, err := zw.Create(file.Name)
	if err != nil {
		return err
	}

	if _, err := io.Copy(fw, f); err != nil {
		return err
	}
	return nil
}
