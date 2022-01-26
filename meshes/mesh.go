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

	dataSize := len(sceneMesh.Vertices)*3 + len(sceneMesh.Normals)*3
	layoutToUse := []buffers.Element{{ElementType: buffers.DataTypeVec3}, {ElementType: buffers.DataTypeVec3}}
	if len(sceneMesh.ColorSets) > 0 {
		layoutToUse = append(layoutToUse, buffers.Element{ElementType: buffers.DataTypeVec4})
		dataSize += len(sceneMesh.ColorSets) * 4
	}

	mesh.Buf.SetLayout(layoutToUse...)
	positions := flattenVec3(sceneMesh.Vertices)
	normals := flattenVec3(sceneMesh.Normals)
	colors := []float32{}
	if len(sceneMesh.ColorSets) > 0 {
		colors = flattenVec4(sceneMesh.ColorSets[0])
	}

	var values []float32
	if len(colors) > 0 {
		values = interleave(
			arrInfo{values: positions, valsPerComp: 3},
			arrInfo{values: normals, valsPerComp: 3},
			arrInfo{values: colors, valsPerComp: 4},
		)
	} else {
		values = interleave(
			arrInfo{values: positions, valsPerComp: 3},
			arrInfo{values: normals, valsPerComp: 3},
		)
	}

	mesh.Buf.SetData(values)
	mesh.Buf.SetIndexBufData(flattenFaces(sceneMesh.Faces))
	return mesh, nil
}

type arrInfo struct {
	values      []float32
	valsPerComp int
}

func interleave(arrs ...arrInfo) []float32 {

	if len(arrs) == 0 || len(arrs[0].values) == 0 {
		panic("No input to interleave or arrays are empty")
	}

	size := 0
	for i := 0; i < len(arrs); i++ {
		size += len(arrs[i].values)
	}

	out := make([]float32, 0, size)
	for posInArr := 0; posInArr < len(arrs[0].values)/arrs[0].valsPerComp; posInArr++ {
		for arrToUse := 0; arrToUse < len(arrs); arrToUse++ {
			for compToAdd := 0; compToAdd < arrs[arrToUse].valsPerComp; compToAdd++ {

				out = append(out, arrs[arrToUse].values[posInArr*arrs[arrToUse].valsPerComp+compToAdd])
			}
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
