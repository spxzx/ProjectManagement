package code_gen

import "testing"

func TestGenStruct(t *testing.T) {
	GenStruct("pm_project_auth_node", "ProjectAuthNode")
}

func TestGenProtoMessage(t *testing.T) {
	GenProtoMessage("pm_project", "Project")
}
