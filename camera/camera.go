package camera

import (
	"github.com/bloeys/gglm/gglm"
)

type Type int32

const (
	Type_Unknown Type = iota
	Type_Perspective
	Type_Orthographic
)

type Camera struct {
	Type Type

	Pos    gglm.Vec3
	Target gglm.Vec3
	// Forward gglm.Vec3
	WorldUp gglm.Vec3

	NearClip float32
	FarClip  float32

	// Perspective data
	Fov         float32
	AspectRatio float32

	// Ortho data
	Left, Right, Top, Bottom float32

	// Matrices
	ViewMat gglm.Mat4
	ProjMat gglm.Mat4
}

// Update recalculates view and projection matrices
func (c *Camera) Update() {

	c.ViewMat = gglm.LookAt(&c.Pos, &c.Target, &c.WorldUp).Mat4

	if c.Type == Type_Perspective {
		c.ProjMat = *gglm.Perspective(c.Fov, c.AspectRatio, c.NearClip, c.FarClip)
	} else {
		c.ProjMat = gglm.Ortho(c.Left, c.Right, c.Top, c.Bottom, c.NearClip, c.FarClip).Mat4
	}
}

func (c *Camera) LookAt(targetPos, worldUp *gglm.Vec3) {
	c.Target = *targetPos
	c.WorldUp = *worldUp
	c.Update()
}

func NewPerspective(pos, targetPos, worldUp *gglm.Vec3, nearClip, farClip, fovRadians, aspectRatio float32) *Camera {

	cam := &Camera{
		Type: Type_Perspective,
		Pos:  *pos,
		// Forward: *gglm.NewVec3(0, 0, 1),
		Target:  *targetPos,
		WorldUp: *worldUp,

		NearClip: nearClip,
		FarClip:  farClip,

		Fov:         fovRadians,
		AspectRatio: aspectRatio,
	}
	cam.Update()

	return cam
}

func NewOrthographic(pos, targetPos, worldUp *gglm.Vec3, nearClip, farClip, left, right, top, bottom float32) *Camera {

	cam := &Camera{
		Type: Type_Orthographic,
		Pos:  *pos,
		// Forward: *gglm.NewVec3(0, 0, 0),
		Target:  *targetPos,
		WorldUp: *worldUp,

		NearClip: nearClip,
		FarClip:  farClip,

		Left:   left,
		Right:  right,
		Top:    top,
		Bottom: bottom,
	}
	cam.Update()

	return cam
}
