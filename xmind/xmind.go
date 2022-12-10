package xmind

import "strings"

// struct of content.json in xmind file

type Contents []Content

type Content struct {
	Id               string                   `json:"id"`
	Class            string                   `json:"class"`
	Title            string                   `json:"title"`
	RootTopic        RootTopic                `json:"rootTopic"`
	Extensions       []map[string]interface{} // extensionsとthemeは触らないのでinterface{]で適当に処理する
	Theme            map[string]interface{}
	TopicPositioning string `json:"topicPositioning"`
	CoreVersion      string `json:"coreVersion"`
}

type RootTopic struct {
	Id             string   `json:"id"`
	Class          string   `json:"class"`
	Title          string   `json:"title"`
	StructureClass string   `json:"structureClass"`
	TitleUnedited  bool     `json:"titleUnedited"`
	Children       Children `json:"children"`

	// 下記はこのツールのために拡張したプロパティでcontent.jsonには含まれない
	Project   string
	Component string
	Epic      string
}

// Titleは改行と:で構成されて、JIRAチケットを作るためのメタデータを持っている
// e.g. Title\nproject:project\ncomponent:component\nepic:epic
func (r *RootTopic) ParseTitle() {
	data := strings.Split(r.Title, "\n")
	for _, d := range data {
		s := strings.Split(d, ":")
		switch s[0] {
		case "project":
			r.Project = strings.TrimSpace(s[1])
		case "component":
			r.Component = strings.TrimSpace(s[1])
		case "epic":
			r.Epic = strings.TrimSpace(s[1])
		}
	}
}

type Children struct {
	Attached []Attached `json:"attached"`
}

type Attached struct {
	Id            string   `json:"id"`
	Title         string   `json:"title"`
	TitleUnedited bool     `json:"titleUnedited,omitempty"` // Decodeするときに、nullや空リストだったときにJSONから消したいキーはomitemptyをつける
	Children      Children `json:"children,omitempty"`
	Markers       []Marker `json:"makers,omitempty"`
}

type Marker struct {
	MarkerId string `json:"makerId"`
}
