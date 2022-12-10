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
xmindファイルの中にあるcontent.jsonを書き換えて元のxmindファイルに戻す。
以下の手順で処理を行う。
1. xmindファイルの実態はzipであるので、zipを解答し圧縮されたファイルの一覧を得る
2. JIRAチケットを作成し、その結果でcontent.jsonを上書きする
3. 編集したcontent.jsonと残りのファイルで改めてzipに圧縮する
4. 元のxmindファイルを削除し、新しく作ったzipを元の名前にrenameする
*/
func main() {
	flag.BoolVar(&dryrun, "dryrun", false, "skip creating jira tickets if dryrun option is true")
	flag.Parse()

	// 1. xmindファイルを解答しファイルの一覧を得る
	zr, err := zip.OpenReader("./sample.xmind")
	if err != nil {
		log.Fatal(err)
	}
	defer zr.Close()

	// 2. JIRAチケットを作成し、その結果でcontent.jsonを上書きする
	f, err := xmind.FindContentJsonFile(zr.File)
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

	// content.jsonのrootオブジェクトは配列だが、要素は1つなのでc[0]で取得する
	r := c[0].RootTopic
	r.ParseTitle()
	// projectとepicが設定されてないとチケットを作れないという仕様にする
	if r.Project == "" || r.Epic == "" {
		log.Fatal("RootTopic must be set project and epic")
	}

	leaves := r.FindLeaves()

	if dryrun == true {
		fmt.Println("skipped creating jira tickets because dryrun option is true")
		fmt.Println("below tickets will be create")
		fmt.Printf("project:%s\ncomponent:%s\nepic:%s\n", r.Project, r.Component, r.Epic)
		for i, leaf := range leaves {
			fmt.Printf("ticket %d: %s", i, leaf.Title)
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
			url, err := jira.IssueCreate(client, r.Project, r.Component, r.Epic, "Task", leaves[i].Title, "", "", "")
			if err != nil {
				log.Printf("got an error creating a ticket titled \"%s\". error is below:\n%s", leaves[i].Title, err)
				continue
			}
			fmt.Printf("created a ticket titled \"%s\": %s\n", leaves[i].Title, url)
			// チケットのURLをXMindの葉に書き足す
			leaves[i].Title = fmt.Sprintf("%s\nurl: %s\n", leaves[i].Title, url)
		}
	}

	// 構造体からjsonに戻す
	j, err := json.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}

	// 3. 編集したcontent.jsonと残りのファイルで改めてzipに圧縮する
	z, err := os.CreateTemp("", "new.xmind")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(z.Name())

	if err := save(zr.File, j, z); err != nil {
		log.Fatal(err)
	}

	// 4. 元のxmindファイルを削除し、新しく作ったzipを元の名前にrenameする
	zr.Close() // removeする前にcloseしておく
	if err := os.Remove("./sample.xmind"); err != nil {
		log.Fatal(err)
	}
	if err := os.Rename(z.Name(), "./sample.xmind"); err != nil {
		log.Fatal(err)
	}
}

/**
元のxmindのFileと新しく作ったcontent.jsonを新しいzipファイルに保存する
*/
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
zipWriterにfileを書き込む
安全にfileをcloseするために、ループの本体を関数に書き出す
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
