package xmind

import (
	"testing"
)

func TestParseTitle(t *testing.T) {
	r := RootTopic{
		Id:             "id",
		Class:          "class",
		Title:          "title\nproject:project\nepic:epic\ncomponent:component",
		StructureClass: "structureClass",
		TitleUnedited:  true,
	}
	r.ParseTitle()
	if r.Project != "project" || r.Epic != "epic" || r.Component != "component" {
		t.Fatalf("want project=project, epic=epic and component=compoent, but r project=%s, epic=%s and component=%s", r.Project, r.Epic, r.Component)
	}
}
