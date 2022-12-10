package jira

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
)

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
