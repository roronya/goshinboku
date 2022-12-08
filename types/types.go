package types

type Contents []Content

type Content struct {
	Id               string                   `json:"id""`
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
