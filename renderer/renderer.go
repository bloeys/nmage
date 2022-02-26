package renderer

import (
	"github.com/bloeys/gglm/gglm"
	"github.com/bloeys/nmage/materials"
	"github.com/bloeys/nmage/meshes"
)

type Render interface {
	Draw(mesh *meshes.Mesh, trMat *gglm.TrMat, mat *materials.Material)
	FrameEnd()
}
