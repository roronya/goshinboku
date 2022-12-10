package xmind

import (
	"sort"
	"testing"
)

func TestRootTopic_ParseTitle(t *testing.T) {
	r := RootTopic{Title: "title\nproject:project\nepic:epic\ncomponent:component"}
	r.ParseTitle()
	if !(r.Project == "project" && r.Component == "component" && r.Epic == "epic") {
		t.Fatalf("want project=\"project\", component=\"component\" and epic=\"epic\", but r project=\"%s\", component=\"%s\" and epic=\"%s\"", r.Project, r.Epic, r.Component)
	}

	r = RootTopic{Title: "title"}
	r.ParseTitle()
	if !(r.Project == "" && r.Epic == "" && r.Component == "") {
		t.Fatalf("want project=\"\", component=\"\" and epic=\"\", but r project=\"%s\", component=\"%s\" and epic=\"%s\"", r.Project, r.Epic, r.Component)
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
