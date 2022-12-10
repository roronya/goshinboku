package main

import (
	"github.com/roronya/goshinboku/types"
	"testing"
)

func TestGetMetaData(t *testing.T) {
	r := types.RootTopic{
		Id:             "id",
		Class:          "class",
		Title:          "title\nproject:project\nepic:epic\ncomponent:component",
		StructureClass: "structureClass",
		TitleUnedited:  true,
	}
	got := getMetaData(r)
	if got.Project != "project" || got.Epic != "epic" || got.Component != "component" {
		t.Fatalf("want project=project, epic=epic and component=compoent, but got project=%s, epic=%s and component=%s", got.Project, got.Epic, got.Component)
	}
}
