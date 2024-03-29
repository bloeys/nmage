package meshes

import (
	"errors"
	"fmt"

	"github.com/bloeys/assimp-go/asig"
	"github.com/bloeys/gglm/gglm"
	"github.com/bloeys/nmage/assert"
	"github.com/bloeys/nmage/buffers"
)

type SubMesh struct {
	BaseVertex int32
	BaseIndex  uint32
	IndexCount int32
}

type Mesh struct {
	Name      string
	Buf       buffers.Buffer
	SubMeshes []SubMesh
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

	mesh := &Mesh{
		Name:      name,
		Buf:       buffers.NewBuffer(),
		SubMeshes: make([]SubMesh, 0, 1),
	}

	// Initial sizes assuming one submesh that has vertex pos+normals+texCoords, and 3 indices per face
	var vertexBufData []float32 = make([]float32, 0, len(scene.Meshes[0].Vertices)*3*3*2)
	var indexBufData []uint32 = make([]uint32, 0, len(scene.Meshes[0].Faces)*3)

	for i := 0; i < len(scene.Meshes); i++ {

		sceneMesh := scene.Meshes[i]

		if len(sceneMesh.TexCoords[0]) == 0 {
			sceneMesh.TexCoords[0] = make([]gglm.Vec3, len(sceneMesh.Vertices))
			println("Zeroing tex coords for submesh", i)
		}

		layoutToUse := []buffers.Element{{ElementType: buffers.DataTypeVec3}, {ElementType: buffers.DataTypeVec3}, {ElementType: buffers.DataTypeVec2}}
		if len(sceneMesh.ColorSets) > 0 && len(sceneMesh.ColorSets[0]) > 0 {
			layoutToUse = append(layoutToUse, buffers.Element{ElementType: buffers.DataTypeVec4})
		}

		if i == 0 {
			mesh.Buf.SetLayout(layoutToUse...)
		} else {

			// @NOTE: Require that all submeshes have the same vertex buffer layout
			firstSubmeshLayout := mesh.Buf.GetLayout()
			assert.T(len(firstSubmeshLayout) == len(layoutToUse), fmt.Sprintf("Vertex layout of submesh %d does not equal vertex layout of the first submesh. Original layout: %v; This layout: %v", i, firstSubmeshLayout, layoutToUse))

			for i := 0; i < len(firstSubmeshLayout); i++ {
				if firstSubmeshLayout[i].ElementType != layoutToUse[i].ElementType {
					panic(fmt.Sprintf("Vertex layout of submesh %d does not equal vertex layout of the first submesh. Original layout: %v; This layout: %v", i, firstSubmeshLayout, layoutToUse))
				}
			}
		}

		arrs := []arrToInterleave{{V3s: sceneMesh.Vertices}, {V3s: sceneMesh.Normals}, {V2s: v3sToV2s(sceneMesh.TexCoords[0])}}
		if len(sceneMesh.ColorSets) > 0 && len(sceneMesh.ColorSets[0]) > 0 {
			arrs = append(arrs, arrToInterleave{V4s: sceneMesh.ColorSets[0]})
		}

		indices := flattenFaces(sceneMesh.Faces)
		mesh.SubMeshes = append(mesh.SubMeshes, SubMesh{

			// Index of the vertex to start from (e.g. if index buffer says use vertex 5, and BaseVertex=3, the vertex used will be vertex 8)
			BaseVertex: int32(len(vertexBufData)*4) / mesh.Buf.Stride,
			// Which index (in the index buffer) to start from
			BaseIndex: uint32(len(indexBufData)),
			// How many indices in this submesh
			IndexCount: int32(len(indices)),
		})

		vertexBufData = append(vertexBufData, interleave(arrs...)...)
		indexBufData = append(indexBufData, indices...)
	}

	// fmt.Printf("!!! Vertex count: %d; Submeshes: %+v\n", len(vertexBufData)*4/int(mesh.Buf.Stride), mesh.SubMeshes)
	mesh.Buf.SetData(vertexBufData)
	mesh.Buf.SetIndexBufData(indexBufData)
	return mesh, nil
}

func v3sToV2s(v3s []gglm.Vec3) []gglm.Vec2 {

	v2s := make([]gglm.Vec2, len(v3s))
	for i := 0; i < len(v3s); i++ {
		v2s[i] = gglm.Vec2{
			Data: [2]float32{v3s[i].X(), v3s[i].Y()},
		}
	}

	return v2s
}

type arrToInterleave struct {
	V2s []gglm.Vec2
	V3s []gglm.Vec3
	V4s []gglm.Vec4
}

func (a *arrToInterleave) get(i int) []float32 {

	assert.T(len(a.V2s) == 0 || len(a.V3s) == 0, "One array should be set in arrToInterleave, but both arrays are set")
	assert.T(len(a.V2s) == 0 || len(a.V4s) == 0, "One array should be set in arrToInterleave, but both arrays are set")
	assert.T(len(a.V3s) == 0 || len(a.V4s) == 0, "One array should be set in arrToInterleave, but both arrays are set")

	if len(a.V2s) > 0 {
		return a.V2s[i].Data[:]
	} else if len(a.V3s) > 0 {
		return a.V3s[i].Data[:]
	} else {
		return a.V4s[i].Data[:]
	}
}

func interleave(arrs ...arrToInterleave) []float32 {

	assert.T(len(arrs) > 0, "No input sent to interleave")
	assert.T(len(arrs[0].V2s) > 0 || len(arrs[0].V3s) > 0 || len(arrs[0].V4s) > 0, "Interleave arrays are empty")

	elementCount := 0
	if len(arrs[0].V2s) > 0 {
		elementCount = len(arrs[0].V2s)
	} else if len(arrs[0].V3s) > 0 {
		elementCount = len(arrs[0].V3s)
	} else {
		elementCount = len(arrs[0].V4s)
	}

	//Calculate final size of the float buffer
	totalSize := 0
	for i := 0; i < len(arrs); i++ {

		assert.T(len(arrs[i].V2s) == elementCount || len(arrs[i].V3s) == elementCount || len(arrs[i].V4s) == elementCount, "Mesh vertex data given to interleave is not the same length")

		if len(arrs[i].V2s) > 0 {
			totalSize += len(arrs[i].V2s) * 2
		} else if len(arrs[i].V3s) > 0 {
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

func flattenFaces(faces []asig.Face) []uint32 {

	assert.T(len(faces[0].Indices) == 3, fmt.Sprintf("Face doesn't have 3 indices. Index count: %v\n", len(faces[0].Indices)))

	uints := make([]uint32, len(faces)*3)
	for i := 0; i < len(faces); i++ {
		uints[i*3+0] = uint32(faces[i].Indices[0])
		uints[i*3+1] = uint32(faces[i].Indices[1])
		uints[i*3+2] = uint32(faces[i].Indices[2])
	}

	return uints
}
