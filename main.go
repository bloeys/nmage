package main

import (
	"fmt"

	"github.com/bloeys/assimp-go/asig"
	"github.com/bloeys/gglm/gglm"
	"github.com/bloeys/nmage/assets"
	"github.com/bloeys/nmage/camera"
	"github.com/bloeys/nmage/engine"
	"github.com/bloeys/nmage/input"
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
// Entities and components
// Integrate physx
// Create VAO struct independent from VBO to support multi-VBO use cases (e.g. instancing)
// Renderer batching
// Scene graph
// Separate engine loop from rendering loop? or leave it to the user?
// Abstract keys enum away from sdl
// Proper Asset loading
// Frustum culling
// Material system editor with fields automatically extracted from the shader

var (
	window *engine.Window

	cam *camera.Camera

	simpleMat *materials.Material
	cubeMesh  *meshes.Mesh

	modelMat = gglm.NewTrMatId()

	lightPos1   = gglm.NewVec3(2, 2, 0)
	lightColor1 = gglm.NewVec3(1, 1, 1)
)

type OurGame struct {
	Win       *engine.Window
	ImGUIInfo nmageimgui.ImguiInfo
}

func main() {

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

	modelMat.Mul(translationMat.Mul(rotMat.Mul(scaleMat)))
	simpleMat.SetUnifMat4("modelMat", &modelMat.Mat4)

	// Camera
	winWidth, winHeight := g.Win.SDLWin.GetSize()
	cam = camera.NewPerspective(
		gglm.NewVec3(0, 0, -10),
		gglm.NewVec3(0, 0, -9),
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

func (g *OurGame) Update() {

	if input.IsQuitClicked() || input.KeyClicked(sdl.K_ESCAPE) {
		engine.Quit()
	}

	//Camera movement
	var camSpeed float32 = 15
	if input.KeyDown(sdl.K_w) {
		cam.Pos.AddY(camSpeed * timing.DT())
		updateViewMat()
	}
	if input.KeyDown(sdl.K_s) {
		cam.Pos.AddY(-camSpeed * timing.DT())
		updateViewMat()
	}
	if input.KeyDown(sdl.K_d) {
		cam.Pos.AddX(camSpeed * timing.DT())
		updateViewMat()
	}
	if input.KeyDown(sdl.K_a) {
		cam.Pos.AddX(-camSpeed * timing.DT())
		updateViewMat()
	}

	if input.GetMouseWheelYNorm() > 0 {
		cam.Pos.AddZ(1)
		updateViewMat()
	} else if input.GetMouseWheelYNorm() < 0 {
		cam.Pos.AddZ(-1)
		updateViewMat()
	}

	//Rotating cubes
	if input.KeyDown(sdl.K_SPACE) {
		modelMat.Rotate(10*timing.DT()*gglm.Deg2Rad, gglm.NewVec3(1, 1, 1).Normalize())
		simpleMat.SetUnifMat4("modelMat", &modelMat.Mat4)
	}

	imgui.DragFloat3("Cam Pos", &cam.Pos.Data)
}

func (g *OurGame) Render() {

	tempModelMat := modelMat.Clone()

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
	target := cam.Pos.Clone()
	target.AddZ(1)
	cam.Target = *target
	cam.Update()

	simpleMat.SetUnifMat4("viewMat", &cam.ViewMat)
}
