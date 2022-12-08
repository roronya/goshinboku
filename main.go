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
	zr, err := zip.OpenReader("./sample.xmind")
	if err != nil {
		log.Fatal(err)
	}
	defer zr.Close()

	// 2. content.jsonを編集する
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
	if err := save(zr.File, j); err != nil {
		log.Fatal(err)
	}

	/**
	// 4. 元のxmindファイルを削除し、新しく作ったzipを元の名前にrenameする
	zr.Close() // removeする前にcloseしておく
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

/**
元のxmindのFileと新しく作ったcontent.jsonから、新しいxmindファイルを作成する
*/
func save(files []*zip.File, c []byte) error {
	z, err := os.Create("./new.xmind") // TODO: /tmpにユニークな名前で一時ファイルを作る
	if err != nil {
		return err
	}

	zw := zip.NewWriter(z)
	defer zw.Close()

	fw, err := zw.Create("content.json")
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
