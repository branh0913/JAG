package jobbuilder

type Trigger struct {
}
type SCM struct {
	Class string `xml:"class,attr"`
	Value string `xml:",chardata"`
}

type project struct {
	Description                      string   `xml:"description"`
	KeepDependencies                 bool     `xml:"keepDependencies"`
	SCM                              SCM      `xml:"scm"`
	CanRoam                          bool     `xml:"canRoam"`
	Disabled                         bool     `xml:"disabled"`
	BlockBuildWhenDownstreamBuilding bool     `xml:"blockBuildWhenDownstreamBuilding"`
	BlockBuildWhenUpstreamBuilding   bool     `xml:"blockBuildWhenUpstreamBuilding"`
	Triggers                         Trigger  `xml:"triggers,omitempty"`
	ConcurrentBuild                  bool     `xml:"concurrentBuild"`
	Builders                         Builders `xml:"builders"`
}

type Builders struct {
	HudsonTasksShell HudsonTasksShell `xml:"hudson.tasks.Shell,omitempty"`
}

type HudsonTasksShell struct {
	Command string `xml:"command"`
}

type JobCreateRequest struct {
	CanRoam                          bool
	Disable                          bool
	BlockBuildWhenDownstreamBuilding bool
	BlockBuildWhenUpstreamBuilding   bool
	ConcurrentBuild                  bool
	Description                      string
	KeepDependencies                 bool
	SCMClass                         string
}

func New(input *JobCreateRequest) (*project, error) {

	return &project{
		SCM: SCM{
			Class: "hudson.scm.NullSCM",
		},
		Description:      input.Description,
		KeepDependencies: input.KeepDependencies,
		Disabled:         input.Disable,
		Triggers:         Trigger{},
		Builders:         Builders{},
	}, nil

}
