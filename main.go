package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/roronya/goshinboku/xmind"
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
	if r.Project == "" || r.Epic == "" {
		log.Fatal("RootTopic must be set project and epic")
	}

	leaves := r.FindLeaves()

	user := os.Getenv("JIRA_USER")
	pass := os.Getenv("JIRA_PASSWORD")
	server := os.Getenv("JIRA_SERVER")
	client, err := NewClient(user, pass, server)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(leaves); i++ {
		url, err := IssueCreate(client, r.Project, r.Component, r.Epic, "Task", leaves[i].Title, "", "", "")
		if err != nil {
			log.Printf("got an error creating a ticket titled \"%s\". error is below:\n%s", leaves[i].Title, err)
			continue
		}
		fmt.Printf("created a ticket titled \"%s\": %s\n", leaves[i].Title, url)
		leaves[i].Title = fmt.Sprintf("%s\nurl: %s\n", leaves[i].Title, url)
	}

	j, err := json.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Printf(string(j))

	// 3. 編集したcontent.jsonと残りのファイルで改めてzipに圧縮する
	if err := save(zr.File, j); err != nil {
		log.Fatal(err)
	}

	// 4. 元のxmindファイルを削除し、新しく作ったzipを元の名前にrenameする
	zr.Close() // removeする前にcloseしておく
	if err := os.Remove("./sample.xmind"); err != nil {
		log.Fatal(err)
	}
	if err := os.Rename("./new.xmind", "./sample.xmind"); err != nil {
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

func NewClient(
	username string,
	password string,
	baseUrl string,
) (*jira.Client, error) {
	tp := jira.BasicAuthTransport{
		Username: username,
		Password: password,
	}
	client, err := jira.NewClient(tp.Client(), baseUrl)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// GetUser
// emailをもとにUserオブジェクトを探して返す
// emailに紐づくユーザーが見つからなかった場合はerrorを返す
func GetUser(
	client *jira.Client,
	email string,
) (*jira.User, error) {
	users, _, err := client.User.Find(email)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("couldn't find a such user<%s>", email)
	}
	return &users[0], nil
}

// IssueCreate
// client: NewClientによって作ったclient
// project: JIRAのプロジェクト。存在するJIRAのプロジェクトを必ず渡す。存在しないJIRAのプロジェクトを渡すとエラーになる。
// component: コンポーネント。設定しない場合は空文字を渡す。
// epic: エピック。設定しない場合は空文字を渡す。
// issueType: 現状はTaskのみ対応している。Task以外を渡すとエラーになる。
// summary: JIRAのチケットのタイトルになる
// assignee: アサインする人のアカウントIDを渡す。アカウントIDはGetUserによって取得できる。アサインしない場合は空文字を渡すと非アサイン状態になる。
// reporter: アサインする人のアカウントIDを渡す。アカウントIDはGetUserによって取得できる。アサインしない場合は空文字を渡すとチケットを作った人が報告者になる。
// description: 説明文
func IssueCreate(
	client *jira.Client,
	project string,
	component string,
	epic string,
	issueType string,
	summary string,
	assignee string,
	reporter string,
	description string,
) (url string, err error) {
	f := &jira.IssueFields{
		Project: jira.Project{
			Key: project,
		},
		Type: jira.IssueType{
			Name: issueType,
		},
		Summary:     summary,
		Description: description,
	}
	if component != "" {
		f.Components = []*jira.Component{{Name: component}}
	}
	// see: https://github.com/andygrunwald/go-jira/issues/307
	// FIXME: epicはcustomefieldとして作られていて、プロジェクトによって違う
	if epic != "" {
		f.Unknowns = map[string]interface{}{
			"customfield_10006": epic,
		}
	}

	if assignee != "" {
		f.Assignee = &jira.User{AccountID: assignee}
	}
	if reporter != "" {
		f.Reporter = &jira.User{AccountID: reporter}
	}
	i := jira.Issue{Fields: f}
	issue, _, err := client.Issue.Create(&i)
	if err != nil {
		return "", err
	}
	baseURL := client.GetBaseURL()
	return fmt.Sprintf("%sbrowse/%s", baseURL.String(), issue.Key), nil
}
