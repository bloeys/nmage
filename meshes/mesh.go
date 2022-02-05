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
	Name string
	Buf  buffers.Buffer
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

	mesh := &Mesh{Name: name}
	sceneMesh := scene.Meshes[0]
	mesh.Buf = buffers.NewBuffer()

	layoutToUse := []buffers.Element{{ElementType: buffers.DataTypeVec3}, {ElementType: buffers.DataTypeVec3}}
	if len(sceneMesh.ColorSets) > 0 {
		layoutToUse = append(layoutToUse, buffers.Element{ElementType: buffers.DataTypeVec4})
	}
	mesh.Buf.SetLayout(layoutToUse...)

	var values []float32
	if len(sceneMesh.ColorSets) > 0 && len(sceneMesh.ColorSets[0]) > 0 {
		values = interleave(
			arrToInterleave{V3s: sceneMesh.Vertices},
			arrToInterleave{V3s: sceneMesh.Normals},
			arrToInterleave{V4s: sceneMesh.ColorSets[0]},
		)
	} else {
		values = interleave(
			arrToInterleave{V3s: sceneMesh.Vertices},
			arrToInterleave{V3s: sceneMesh.Normals},
		)
	}

	mesh.Buf.SetData(values)
	mesh.Buf.SetIndexBufData(flattenFaces(sceneMesh.Faces))
	return mesh, nil
}

type arrToInterleave struct {
	V3s []gglm.Vec3
	V4s []gglm.Vec4
}

func (a *arrToInterleave) get(i int) []float32 {

	asserts.T(len(a.V3s) == 0 || len(a.V4s) == 0, "One array should be set in arrToInterleave, but both arrays are set")

	if len(a.V3s) > 0 {
		return a.V3s[i].Data[:]
	} else {
		return a.V4s[i].Data[:]
	}
}

func interleave(arrs ...arrToInterleave) []float32 {

	asserts.T(len(arrs) > 0, "No input sent to interleave")
	asserts.T(len(arrs[0].V3s) > 0 || len(arrs[0].V4s) > 0, "Interleave arrays are empty")

	elementCount := 0
	if len(arrs[0].V3s) > 0 {
		elementCount = len(arrs[0].V3s)
	} else {
		elementCount = len(arrs[0].V4s)
	}

	totalSize := 0
	for i := 0; i < len(arrs); i++ {

		asserts.T(len(arrs[i].V3s) == elementCount || len(arrs[i].V4s) == elementCount, "Mesh vertex data given to interleave is not the same length")

		if len(arrs[i].V3s) > 0 {
			totalSize += len(arrs[i].V3s) * 3
		} else {
			totalSize += len(arrs[i].V4s) * 4
		}
	}

	out := make([]float32, 0, totalSize)
	for i := 0; i < elementCount; i++ {
		for arrToUse := 0; arrToUse < len(arrs); arrToUse++ {
			out = append(out, arrs[arrToUse].get(i)...)
		}
	}

	return out
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

	asserts.T(len(faces[0].Indices) == 3, fmt.Sprintf("Face doesn't have 3 indices. Index count: %v\n", len(faces[0].Indices)))

	uints := make([]uint32, len(faces)*3)
	for i := 0; i < len(faces); i++ {
		uints[i*3+0] = uint32(faces[i].Indices[0])
		uints[i*3+1] = uint32(faces[i].Indices[1])
		uints[i*3+2] = uint32(faces[i].Indices[2])
	}

	return uints
}
