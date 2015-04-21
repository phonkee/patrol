package models

import "github.com/phonkee/patrol/rest/metadata"

func UpdateProjectNameMetadata(f *metadata.Field) {
	f.SetLabel("Name").SetHelpText("Project name")
}

func UpdateProjectPlatformMetadata(f *metadata.Field) {
	f.SetLabel("Platform").SetHelpText("Platform on which is project running")
}
