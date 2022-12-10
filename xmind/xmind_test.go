package xmind

import (
	"sort"
	"testing"
)

func TestRootTopic_ParseTitle(t *testing.T) {
	r := RootTopic{
		Title: "title\nproject:project\nepic:epic\ncomponent:component",
	}
	r.ParseTitle()
	if r.Project != "project" || r.Epic != "epic" || r.Component != "component" {
		t.Fatalf("want project=project, epic=epic and component=compoent, but r project=%s, epic=%s and component=%s", r.Project, r.Epic, r.Component)
	}
}

func TestRootTopic_FindLeaves(t *testing.T) {
	r := RootTopic{
		Children: Children{
			Attached: []Attached{
				{
					Title: "intermediate0",
					Children: Children{
						Attached: []Attached{
							{Title: "leaf0"},
						},
					},
				},
				{
					Title: "intermediate1",
					Children: Children{
						Attached: []Attached{
							{Title: "leaf1"},
						},
					},
				},
				{
					Title: "leaf2",
				},
			},
		},
	}

	leaves := r.FindLeaves()
	if len(leaves) != 3 {
		t.Fatalf("want len(leaves)=3, but got %d", len(leaves))
	}
	var titles []string
	for _, l := range leaves {
		titles = append(titles, l.Title)
	}
	sort.Strings(titles)
	wants := []string{"leaf0", "leaf1", "leaf2"}
	for i := 0; i < len(leaves); i++ {
		if wants[i] != titles[i] {
			t.Fatalf("want %s, but got %s", wants[i], titles[i])
		}
	}
}
