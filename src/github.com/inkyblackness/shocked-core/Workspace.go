package core

import (
	"sort"

	"github.com/inkyblackness/shocked-core/io"
	"github.com/inkyblackness/shocked-core/release"
)

// Workspace describes a container of projects.
type Workspace struct {
	autoSaveTimeoutMSec int

	source            release.Release
	projectsContainer release.ReleaseContainer

	projects map[string]*Project
}

// NewWorkspace takes a Release as a basis for existing resources and returns
// a new workspace instance. With this instance, projects from given projects container
// can be worked with.
func NewWorkspace(source release.Release, projects release.ReleaseContainer, autoSaveTimeoutMSec int) *Workspace {
	ws := &Workspace{
		autoSaveTimeoutMSec: autoSaveTimeoutMSec,

		source:            source,
		projectsContainer: projects,
		projects:          make(map[string]*Project)}

	return ws
}

// ProjectNames returns all currently known project identifiers.
func (ws *Workspace) ProjectNames() []string {
	names := ws.projectsContainer.Names()

	sort.Strings(names)

	return names
}

// Project returns a project matching the given identifier.
func (ws *Workspace) Project(name string) (project *Project, err error) {
	project, existing := ws.projects[name]

	if !existing {
		rel, relErr := ws.projectsContainer.Get(name)
		if relErr == nil {
			library := io.NewReleaseStoreLibrary(ws.source, rel, ws.autoSaveTimeoutMSec)
			project, err = NewProject(name, library)
			if err == nil {
				ws.projects[name] = project
			}
		} else {
			err = relErr
		}
	}

	return
}

// NewProject tries to create a new project in the workspace.
func (ws *Workspace) NewProject(name string) (project *Project, err error) {
	rel, err := ws.projectsContainer.New(name)

	if err == nil {
		library := io.NewReleaseStoreLibrary(ws.source, rel, ws.autoSaveTimeoutMSec)
		project, err = NewProject(name, library)
		if err == nil {
			ws.projects[name] = project
		}
	}

	return
}
