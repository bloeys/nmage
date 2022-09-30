package main

import (
	"fmt"

	"github.com/bloeys/assimp-go/asig"
	"github.com/bloeys/gglm/gglm"
	"github.com/bloeys/nmage/assets"
	"github.com/bloeys/nmage/camera"
	"github.com/bloeys/nmage/engine"
	"github.com/bloeys/nmage/entity"
	"github.com/bloeys/nmage/input"
	"github.com/bloeys/nmage/level"
	"github.com/bloeys/nmage/logging"
	"github.com/bloeys/nmage/materials"
	"github.com/bloeys/nmage/meshes"
	"github.com/bloeys/nmage/renderer/rend3dgl"
	"github.com/bloeys/nmage/timing"
	nmageimgui "github.com/bloeys/nmage/ui/imgui"
	"github.com/inkyblackness/imgui-go/v4"
	"github.com/veandco/go-sdl2/sdl"
)

// @Todo:
// Complete entity registry (e.g. HasEntity, GetEntity, Generational Indices etc...)
// Helper functions to update active entities

// Integrate physx
// Create VAO struct independent from VBO to support multi-VBO use cases (e.g. instancing)
// Renderer batching
// Scene graph
// Separate engine loop from rendering loop? or leave it to the user?
// Abstract keys enum away from sdl
// Proper Asset loading
// Frustum culling
// Material system editor with fields automatically extracted from the shader

const (
	camSpeed         = 15
	mouseSensitivity = 0.5
)

var (
	window *engine.Window

	pitch float32 = 0
	yaw   float32 = -90
	cam   *camera.Camera

	simpleMat *materials.Material
	cubeMesh  *meshes.Mesh

	cubeModelMat = gglm.NewTrMatId()

	lightPos1   = gglm.NewVec3(-2, 0, 2)
	lightColor1 = gglm.NewVec3(1, 1, 1)
)

type OurGame struct {
	Win       *engine.Window
	ImGUIInfo nmageimgui.ImguiInfo
}

type TransformComp struct {
	Pos   *gglm.Vec3
	Rot   *gglm.Quat
	Scale *gglm.Vec3
}

func (t TransformComp) Name() string {
	return "Transform Component"
}

func Test() {

	lvl := level.NewLevel("test level", 1000)
	e1 := lvl.Registry.NewEntity()

	trComp := entity.GetComp[*TransformComp](e1)
	fmt.Println("Got comp 1:", trComp)

	e1.Comps = append(e1.Comps, &TransformComp{
		Pos:   gglm.NewVec3(0, 0, 0),
		Rot:   gglm.NewQuatEulerXYZ(0, 0, 0),
		Scale: gglm.NewVec3(0, 0, 0),
	}, &TransformComp{
		Pos:   gglm.NewVec3(0, 0, 0),
		Rot:   gglm.NewQuatEulerXYZ(0, 0, 0),
		Scale: gglm.NewVec3(1, 1, 1),
	})

	trComp = entity.GetComp[*TransformComp](e1)
	fmt.Println("Got comp 2:", trComp)

	trComps := entity.GetAllCompOfType[*TransformComp](e1)
	fmt.Printf("Got comp 3: %+v, %+v\n", trComps[0], trComps[1])

	fmt.Printf("Entity: %+v\n", e1)
	fmt.Printf("Entity: %+v\n", lvl.Registry.NewEntity())
	fmt.Printf("Entity: %+v\n", lvl.Registry.NewEntity())
	fmt.Printf("Entity: %+v\n", lvl.Registry.NewEntity())
}

func main() {

	// Test()
	// return

	//Init engine
	err := engine.Init()
	if err != nil {
		logging.ErrLog.Fatalln("Failed to init nMage. Err:", err)
	}

	//Create window
	window, err = engine.CreateOpenGLWindowCentered("nMage", 1280, 720, engine.WindowFlags_RESIZABLE, rend3dgl.NewRend3DGL())
	if err != nil {
		logging.ErrLog.Fatalln("Failed to create window. Err: ", err)
	}
	defer window.Destroy()

	engine.SetVSync(false)

	game := &OurGame{
		Win:       window,
		ImGUIInfo: nmageimgui.NewImGUI(),
	}
	window.EventCallbacks = append(window.EventCallbacks, game.handleWindowEvents)

	engine.Run(game, window, game.ImGUIInfo)
}

func (g *OurGame) handleWindowEvents(e sdl.Event) {

	switch e := e.(type) {
	case *sdl.WindowEvent:
		if e.Event == sdl.WINDOWEVENT_SIZE_CHANGED {

			width := e.Data1
			height := e.Data2
			cam.AspectRatio = float32(width) / float32(height)
			cam.Update()

			simpleMat.SetUnifMat4("projMat", &cam.ProjMat)
		}
	}
}

func (g *OurGame) Init() {

	//Create materials
	simpleMat = materials.NewMaterial("Simple Mat", "./res/shaders/simple.glsl")

	//Load meshes
	var err error
	cubeMesh, err = meshes.NewMesh("Cube", "./res/models/tex-cube.fbx", asig.PostProcess(0))
	if err != nil {
		logging.ErrLog.Fatalln("Failed to load cube mesh. Err: ", err)
	}

	//Load textures
	tex, err := assets.LoadTexturePNG("./res/textures/Low poly planet.png", nil)
	if err != nil {
		logging.ErrLog.Fatalln("Failed to load texture. Err: ", err)
	}

	//Configure material
	simpleMat.DiffuseTex = tex.TexID

	//Movement, scale and rotation
	translationMat := gglm.NewTranslationMat(gglm.NewVec3(0, 0, 0))
	scaleMat := gglm.NewScaleMat(gglm.NewVec3(0.25, 0.25, 0.25))
	rotMat := gglm.NewRotMat(gglm.NewQuatEuler(gglm.NewVec3(0, 0, 0).AsRad()))

	cubeModelMat.Mul(translationMat.Mul(rotMat.Mul(scaleMat)))
	simpleMat.SetUnifMat4("modelMat", &cubeModelMat.Mat4)

	// Camera
	winWidth, winHeight := g.Win.SDLWin.GetSize()
	cam = camera.NewPerspective(
		gglm.NewVec3(0, 0, 10),
		gglm.NewVec3(0, 0, -1),
		gglm.NewVec3(0, 1, 0),
		0.1, 20,
		45*gglm.Deg2Rad,
		float32(winWidth)/float32(winHeight),
	)
	simpleMat.SetUnifMat4("projMat", &cam.ProjMat)

	updateViewMat()

	//Lights
	simpleMat.SetUnifVec3("lightPos1", lightPos1)
	simpleMat.SetUnifVec3("lightColor1", lightColor1)
}

// @TODO: Add this to gglm
// func vecRotByQuat(v *gglm.Vec3, q *gglm.Quat) *gglm.Vec3 {

// 	// Reference: https://gamedev.stackexchange.com/questions/28395/rotating-vector3-by-a-quaternion
// 	qVec := gglm.NewVec3(q.X(), q.Y(), q.Z())

// 	rotatedVec := qVec.Clone().Scale(2 * gglm.DotVec3(v, qVec))

// 	t1 := q.W()*q.W() - gglm.DotVec3(qVec, qVec)
// 	rotatedVec.Add(v.Clone().Scale(t1))

// 	rotatedVec.Add(gglm.Cross(qVec, v).Scale(2 * q.W()))
// 	return rotatedVec
// }

func (g *OurGame) Update() {

	if input.IsQuitClicked() || input.KeyClicked(sdl.K_ESCAPE) {
		engine.Quit()
	}

	g.updateCameraLookAround()
	g.updateCameraPos()

	//Rotating cubes
	if input.KeyDown(sdl.K_SPACE) {
		cubeModelMat.Rotate(10*timing.DT()*gglm.Deg2Rad, gglm.NewVec3(1, 1, 1).Normalize())
		simpleMat.SetUnifMat4("modelMat", &cubeModelMat.Mat4)
	}

	if imgui.DragFloat3("Cam Pos", &cam.Pos.Data) {
		updateViewMat()
	}
	if imgui.DragFloat3("Cam Forward", &cam.Forward.Data) {
		updateViewMat()
	}

	if input.KeyClicked(sdl.K_F4) {
		fmt.Printf("Pos: %s; Forward: %s; |Forward|: %f\n", cam.Pos.String(), cam.Forward.String(), cam.Forward.Mag())
	}
}

func (g *OurGame) updateCameraLookAround() {

	mouseX, mouseY := input.GetMouseMotion()
	if mouseX == 0 && mouseY == 0 {
		return
	}

	// Yaw
	yaw += float32(mouseX) * mouseSensitivity * timing.DT()

	// Pitch
	pitch += float32(-mouseY) * mouseSensitivity * timing.DT()
	if pitch > 89.0 {
		pitch = 89.0
	}

	if pitch < -89.0 {
		pitch = -89.0
	}

	// Update cam forward
	cam.UpdateRotation(pitch, yaw)

	updateViewMat()
}

func (g *OurGame) updateCameraPos() {

	update := false

	// Forward and backward
	if input.KeyDown(sdl.K_w) {
		cam.Pos.Add(cam.Forward.Clone().Scale(camSpeed * timing.DT()))
		update = true
	} else if input.KeyDown(sdl.K_s) {
		cam.Pos.Add(cam.Forward.Clone().Scale(-camSpeed * timing.DT()))
		update = true
	}

	// Left and right
	if input.KeyDown(sdl.K_d) {
		cam.Pos.Add(gglm.Cross(&cam.Forward, &cam.WorldUp).Normalize().Scale(camSpeed * timing.DT()))
		update = true
	} else if input.KeyDown(sdl.K_a) {
		cam.Pos.Add(gglm.Cross(&cam.Forward, &cam.WorldUp).Normalize().Scale(-camSpeed * timing.DT()))
		update = true
	}

	if update {
		updateViewMat()
	}
}

func (g *OurGame) Render() {

	tempModelMat := cubeModelMat.Clone()

	rowSize := 100
	for y := 0; y < rowSize; y++ {
		for x := 0; x < rowSize; x++ {
			tempModelMat.Translate(gglm.NewVec3(-1, 0, 0))
			window.Rend.Draw(cubeMesh, tempModelMat, simpleMat)
		}
		tempModelMat.Translate(gglm.NewVec3(float32(rowSize), -1, 0))
	}

	g.Win.SDLWin.SetTitle(fmt.Sprint("nMage (", timing.GetAvgFPS(), " fps)"))
}

func (g *OurGame) FrameEnd() {
}

func (g *OurGame) DeInit() {
	g.Win.Destroy()
}

func updateViewMat() {
	cam.Update()
	simpleMat.SetUnifMat4("viewMat", &cam.ViewMat)
}
