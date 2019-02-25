package jobbuilder



func New(input *JobCreateRequest) (*project, error) {

	if input.SCMClass == "" {
		input.SCMClass = "hudson.scm.NullSCM"
	}

	return &project{
		SCM: SCM{
			Class: input.SCMClass,
		},
		Description:      input.Description,
		KeepDependencies: input.KeepDependencies,
		Disabled:         input.Disable,
		Triggers:         Trigger{},
		Builders:         Builders{},
	}, nil

}
