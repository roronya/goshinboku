package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/roronya/goshinboku/types"
	"io"
	"log"
	"os"
)

/**
xmindファイルの中にあるcontent.jsonを書き換えて元のxmindファイルに戻す。
以下の手順で処理を行う。
1. xmindファイルの実態はzipであるので、zipを解答し圧縮されたファイルの一覧を得る
2. content.jsonを編集する
3. 編集したcontent.jsonと残りのファイルで改めてzipに圧縮する
4. 元のxmindファイルを削除し、新しく作ったzipを元の名前にrenameする
*/
func main() {
	// 1. xmindファイルを解答しファイルの一覧を得る
	srcReader, err := zip.OpenReader("./sample.xmind")
	if err != nil {
		log.Fatal(err)
	}
	defer srcReader.Close()

	// 2. content.jsonを編集する
	contentJsonFile, err := findContentJsonFile(srcReader.File)
	if err != nil {
		log.Fatal(err)
	}
	contentJsonReader, err := contentJsonFile.Open()
	if err != nil {
		log.Fatal(err)
	}
	dec := json.NewDecoder(contentJsonReader)
	var c types.Contents
	for {
		if err := dec.Decode(&c); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
	}
	// TODO: content.jsonの編集
	c[0].RootTopic.Children.Attached[0].Title = "modified"
	c[0].RootTopic.Children.Attached[0].TitleUnedited = false

	j, err := json.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Printf(string(j))

	// 3. 編集したcontent.jsonと残りのファイルで改めてzipに圧縮する
	dst, err := os.Create("./new.xmind")
	if err != nil {
		log.Fatal(err)
	}
	dstWriter := zip.NewWriter(dst)

	contentJsonWriter, err := dstWriter.Create("content.json")
	if _, err := contentJsonWriter.Write(j); err != nil {
		log.Fatal(err)
	}

	srcFiles := srcReader.File
	for _, srcFile := range srcFiles {
		if srcFile.Name == "content.json" {
			continue
		}

		srcFileReader, err := srcFile.Open()
		if err != nil {
			log.Fatal(err)
		}

		dstFileWriter, err := dstWriter.Create(srcFile.Name)
		if err != nil {
			log.Fatal(err)
		}

		if _, err := io.Copy(dstFileWriter, srcFileReader); err != nil {
			log.Fatal(err)
		}
		srcFileReader.Close()
	}
	srcReader.Close()
	dstWriter.Close()

	/**
	// 4. 元のxmindファイルを削除し、新しく作ったzipを元の名前にrenameする
	if err := os.Remove("./sample.xmind"); err != nil {
		log.Fatal(err)
	}
	if err := os.Rename("./new.xmind", "./sample.xmind"); err != nil {
		log.Fatal(err)
	}
	*/
}

func findContentJsonFile(files []*zip.File) (*zip.File, error) {
	for _, f := range files {
		if f.Name == "content.json" {
			return f, nil
		}
	}
	return nil, fmt.Errorf("cannot find content.json")
}
