package comps

import (
	"github.com/bloeys/gglm/gglm"
	"github.com/bloeys/nmage/entity"
)

type TransformComp struct {
	entity.BaseComp

	Pos   *gglm.Vec3
	Rot   *gglm.Quat
	Scale *gglm.Vec3
}

func (t *TransformComp) Name() string {
	return "Transform Component"
}
