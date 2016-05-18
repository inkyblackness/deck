package core

import (
	"sort"

	"github.com/inkyblackness/shocked-core/release"
)

type Workspace struct {
	source            release.Release
	projectsContainer release.ReleaseContainer

	projects map[string]*Project
}

// NewWorkspace takes a Release as a basis for existing resources and returns
// a new workspace instance. With this instance, projects from given projects container
// can be worked with.
func NewWorkspace(source release.Release, projects release.ReleaseContainer) *Workspace {
	ws := &Workspace{
		source:            source,
		projectsContainer: projects,
		projects:          make(map[string]*Project)}

	return ws
}

func (ws *Workspace) ProjectNames() []string {
	names := ws.projectsContainer.Names()

	sort.Strings(names)

	return names
}

func (ws *Workspace) Project(name string) (project *Project, err error) {
	project, existing := ws.projects[name]

	if !existing {
		rel, relErr := ws.projectsContainer.Get(name)
		if relErr == nil {
			project, err = NewProject(name, ws.source, rel)
			if err == nil {
				ws.projects[name] = project
			}
		} else {
			err = relErr
		}
	}

	return
}

func (ws *Workspace) NewProject(name string) (project *Project, err error) {
	rel, err := ws.projectsContainer.New(name)

	if err == nil {
		project, err = NewProject(name, ws.source, rel)
		if err == nil {
			ws.projects[name] = project
		}
	}

	return
}
