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

func TestAttached_ParseTitle(t *testing.T) {
	a := Attached{Title: "title\nissueId: JIRA-0000"}
	a.ParseTitle()
	if a.IssueId != "JIRA-0000" {
		t.Fatalf("want JIRA-0000, but got %s", a.IssueId)
	}

	a = Attached{Title: "title"}
	a.ParseTitle()
	if a.IssueId != "" {
		t.Fatalf("want empty string, but got %s", a.IssueId)
	}

}

func TestAttached_SetMakerAs(t *testing.T) {
	a := Attached{Markers: []Marker{}}
	a.SetMarkerAsProgress()
	if len(a.Markers) != 1 {
		t.Fatalf("want len(a.Markers) = 1, but got %d", len(a.Markers))
	}
	if m := a.Markers[0].MarkerId; m != "tag-yellow" {
		t.Fatalf("want tag-yellow, but got %s", m)
	}
	a.SetMakerAsDone()
	if len(a.Markers) != 1 {
		t.Fatalf("want len(a.Markers) = 1, but got %d", len(a.Markers))
	}
	if m := a.Markers[0].MarkerId; m != "tag-green" {
		t.Fatalf("want tag-green, but got %s", m)
	}
	a.SetMakerAsTodo()
	if len(a.Markers) != 0 {
		t.Fatalf("want len(a.Markers) = 0, but got %d", len(a.Markers))
	}
}

func TestAttached_SetIssueId(t *testing.T) {
	a := Attached{Title: "title\nissueId: JIRA-0000"}
	err := a.SetIssueId("JIRA-0001")
	if err == nil {
		t.Fatalf("want error, but got nil")
	}
	a = Attached{Title: "title"}
	a.SetIssueId("JIRA-0000")
	want := "title\nissueId: JIRA-0000"
	if a.Title != want {
		t.Fatalf("want %s, got %s", want, a.Title)
	}

}
