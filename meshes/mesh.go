package meshes

import (
	"errors"
	"fmt"

	"github.com/bloeys/assimp-go/asig"
	"github.com/bloeys/gglm/gglm"
	"github.com/bloeys/nmage/asserts"
	"github.com/bloeys/nmage/buffers"
)

type Mesh struct {
	Name   string
	BufObj *buffers.BufferObject
}

func NewMesh(name, modelPath string, postProcessFlags asig.PostProcess) (*Mesh, error) {

	scene, release, err := asig.ImportFile(modelPath, asig.PostProcessTriangulate|postProcessFlags)
	if err != nil {
		return nil, errors.New("Failed to load model. Err: " + err.Error())
	}
	defer release()

	if len(scene.Meshes) == 0 {
		return nil, errors.New("No meshes found in file: " + modelPath)
	}

	mesh := &Mesh{Name: name, BufObj: buffers.NewBufferObject()}

	sceneMesh := scene.Meshes[0]
	mesh.BufObj.GenBuffer(flattenVec3(sceneMesh.Normals), buffers.BufUsageStatic, buffers.BufTypeNormal, buffers.DataTypeVec3)
	mesh.BufObj.GenBuffer(flattenVec3(sceneMesh.Vertices), buffers.BufUsageStatic, buffers.BufTypeVertPos, buffers.DataTypeVec3)
	mesh.BufObj.GenBufferUint32(flattenFaces(sceneMesh.Faces), buffers.BufUsageStatic, buffers.BufTypeIndex, buffers.DataTypeUint32)

	if len(sceneMesh.ColorSets) > 0 {
		mesh.BufObj.GenBuffer(flattenVec4(sceneMesh.ColorSets[0]), buffers.BufUsageStatic, buffers.BufTypeColor, buffers.DataTypeVec4)
	}

	return mesh, nil
}

func flattenVec3(vec3s []gglm.Vec3) []float32 {

	floats := make([]float32, len(vec3s)*3)
	for i := 0; i < len(vec3s); i++ {
		floats[i*3+0] = vec3s[i].X()
		floats[i*3+1] = vec3s[i].Y()
		floats[i*3+2] = vec3s[i].Z()
	}

	return floats
}

func flattenVec4(vec4s []gglm.Vec4) []float32 {

	floats := make([]float32, len(vec4s)*4)
	for i := 0; i < len(vec4s); i++ {
		floats[i*4+0] = vec4s[i].X()
		floats[i*4+1] = vec4s[i].Y()
		floats[i*4+2] = vec4s[i].Z()
		floats[i*4+3] = vec4s[i].W()
	}

	return floats
}

func flattenFaces(faces []asig.Face) []uint32 {

	asserts.True(len(faces[0].Indices) == 3, fmt.Sprintf("Face doesn't have 3 indices. Index count: %v\n", len(faces[0].Indices)))

	uints := make([]uint32, len(faces)*3)
	for i := 0; i < len(faces); i++ {
		uints[i*3+0] = uint32(faces[i].Indices[0])
		uints[i*3+1] = uint32(faces[i].Indices[1])
		uints[i*3+2] = uint32(faces[i].Indices[2])
	}

	return uints
}
