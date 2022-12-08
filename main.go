package main

import (
	"archive/zip"
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
	// TODO

	// 3. 編集したcontent.jsonと残りのファイルで改めてzipに圧縮する
	dst, err := os.Create("./new.xmind")
	if err != nil {
		log.Fatal(err)
	}
	dstWriter := zip.NewWriter(dst)

	srcFiles := srcReader.File
	for _, srcFile := range srcFiles {
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

	// 4. 元のxmindファイルを削除し、新しく作ったzipを元の名前にrenameする
	if err := os.Remove("./sample.xmind"); err != nil {
		log.Fatal(err)
	}
	if err := os.Rename("./new.xmind", "./sample.xmind"); err != nil {
		log.Fatal(err)
	}
}
