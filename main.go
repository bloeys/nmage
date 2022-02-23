package main

import (
	"fmt"
	"runtime"

	"github.com/bloeys/assimp-go/asig"
	"github.com/bloeys/gglm/gglm"
	"github.com/bloeys/nmage/assets"
	"github.com/bloeys/nmage/engine"
	"github.com/bloeys/nmage/input"
	"github.com/bloeys/nmage/logging"
	"github.com/bloeys/nmage/materials"
	"github.com/bloeys/nmage/meshes"
	"github.com/bloeys/nmage/timing"
	nmageimgui "github.com/bloeys/nmage/ui/imgui"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/inkyblackness/imgui-go/v4"
	"github.com/veandco/go-sdl2/sdl"
)

//TODO: Tasks:
//Proper rendering setup
//Entities and components
//Camera class
//Audio
//Flesh out the material system

//Low Priority:
// Abstract keys enum away from sdl
// Abstract UI
// Proper Asset loading

var (
	isRunning bool = true
	window    *engine.Window

	simpleMat *materials.Material
	cubeMesh  *meshes.Mesh

	modelMat = gglm.NewTrMatId()
	projMat  = &gglm.Mat4{}

	camPos     = gglm.NewVec3(0, 0, -10)
	camForward = gglm.NewVec3(0, 0, 1)
)

type OurGame struct {
	Win       *engine.Window
	ImGUIInfo nmageimgui.ImguiInfo
	Quitting  bool
}

func (g *OurGame) Init() {

	//Create materials
	simpleMat = materials.NewMaterial("Simple Mat", "./res/shaders/simple")

	//Load meshes
	var err error
	cubeMesh, err = meshes.NewMesh("Cube", "./res/models/tex-cube.fbx", asig.PostProcess(0))
	if err != nil {
		logging.ErrLog.Fatalln("Failed to load cube mesh. Err: ", err)
	}

	//Load textures
	tex, err := assets.LoadPNGTexture("./res/textures/Low poly planet.png")
	if err != nil {
		logging.ErrLog.Fatalln("Failed to load texture. Err: ", err)
	}
	cubeMesh.AddTexture(tex)

	//Set mesh textures on material
	for _, v := range cubeMesh.TextureIDs {
		simpleMat.AddTextureID(v)
	}

	//Enable vertex attributes
	simpleMat.SetAttribute(cubeMesh.Buf)

	//Movement, scale and rotation
	translationMat := gglm.NewTranslationMat(gglm.NewVec3(0, 0, 0))
	scaleMat := gglm.NewScaleMat(gglm.NewVec3(0.25, 0.25, 0.25))
	rotMat := gglm.NewRotMat(gglm.NewQuatEuler(gglm.NewVec3(0, 0, 0).AsRad()))

	modelMat.Mul(translationMat.Mul(rotMat.Mul(scaleMat)))
	simpleMat.SetUnifMat4("modelMat", &modelMat.Mat4)

	//Moves objects into the cameras view
	updateViewMat()

	//Perspective/Depth
	projMat := gglm.Perspective(45*gglm.Deg2Rad, float32(1280)/float32(720), 0.1, 500)
	simpleMat.SetUnifMat4("projMat", projMat)

	//Lights
	simpleMat.SetUnifVec3("lightPos1", &lightPos1)
	simpleMat.SetUnifVec3("lightColor1", &lightColor1)
}

func (g *OurGame) FrameStart() {
}

func (g *OurGame) Update() {

	if input.IsQuitClicked() {
		g.Quitting = true
	}

	winWidth, winHeight := g.Win.SDLWin.GetSize()
	projMat = gglm.Perspective(45*gglm.Deg2Rad, float32(winWidth)/float32(winHeight), 0.1, 20)
	simpleMat.SetUnifMat4("projMat", projMat)

	//Camera movement
	var camSpeed float32 = 15
	if input.KeyDown(sdl.K_w) {
		camPos.Data[1] += camSpeed * timing.DT()
		updateViewMat()
	}
	if input.KeyDown(sdl.K_s) {
		camPos.Data[1] -= camSpeed * timing.DT()
		updateViewMat()
	}
	if input.KeyDown(sdl.K_d) {
		camPos.Data[0] += camSpeed * timing.DT()
		updateViewMat()
	}
	if input.KeyDown(sdl.K_a) {
		camPos.Data[0] -= camSpeed * timing.DT()
		updateViewMat()
	}

	if input.GetMouseWheelYNorm() > 0 {
		camPos.Data[2] += 1
		updateViewMat()
	} else if input.GetMouseWheelYNorm() < 0 {
		camPos.Data[2] -= 1
		updateViewMat()
	}

	//Rotating cubes
	if input.KeyDown(sdl.K_SPACE) {
		modelMat.Rotate(10*timing.DT()*gglm.Deg2Rad, gglm.NewVec3(1, 1, 1).Normalize())
		simpleMat.SetUnifMat4("modelMat", &modelMat.Mat4)
	}

	imgui.DragFloat3("Cam Pos", &camPos.Data)
}

var dtAccum float32 = 0
var lastElapsedTime uint64 = 0
var framesSinceLastFPSUpdate uint = 0

func (g *OurGame) Render() {

	simpleMat.Bind()
	cubeMesh.Buf.Bind()
	tempModelMat := modelMat.Clone()

	rowSize := 100
	for y := 0; y < rowSize; y++ {
		for x := 0; x < rowSize; x++ {
			simpleMat.SetUnifMat4("modelMat", &tempModelMat.Translate(gglm.NewVec3(-1, 0, 0)).Mat4)
			gl.DrawElements(gl.TRIANGLES, cubeMesh.Buf.IndexBufCount, gl.UNSIGNED_INT, gl.PtrOffset(0))
		}
		simpleMat.SetUnifMat4("modelMat", &tempModelMat.Translate(gglm.NewVec3(float32(rowSize), -1, 0)).Mat4)
	}
	simpleMat.SetUnifMat4("modelMat", &modelMat.Mat4)

	dtAccum += timing.DT()
	framesSinceLastFPSUpdate++
	if timing.ElapsedTime() > lastElapsedTime {

		avgDT := dtAccum / float32(framesSinceLastFPSUpdate)
		g.GetWindow().SDLWin.SetTitle("nMage (" + fmt.Sprint(1/avgDT, " fps)"))

		dtAccum = 0
		framesSinceLastFPSUpdate = 0
	}

	lastElapsedTime = timing.ElapsedTime()
}

func (g *OurGame) FrameEnd() {

}

func (g *OurGame) ShouldRun() bool {
	return !g.Quitting
}

func (g *OurGame) GetWindow() *engine.Window {
	return g.Win
}

func (g *OurGame) GetImGUI() nmageimgui.ImguiInfo {
	return g.ImGUIInfo
}

func (g *OurGame) Deinit() {
	g.Win.Destroy()
}

func main() {

	runtime.LockOSThread()

	//Init engine
	err := engine.Init()
	if err != nil {
		logging.ErrLog.Fatalln("Failed to init nMage. Err:", err)
	}

	//Create window
	window, err = engine.CreateOpenGLWindowCentered("nMage", 1280, 720, engine.WindowFlags_RESIZABLE)
	if err != nil {
		logging.ErrLog.Fatalln("Failed to create window. Err: ", err)
	}
	defer window.Destroy()

	engine.SetVSync(false)

	game := &OurGame{
		Win:       window,
		ImGUIInfo: nmageimgui.NewImGUI(),
	}

	engine.Run(game)
	return
}

func updateViewMat() {
	targetPos := camPos.Clone().Add(camForward)
	viewMat := gglm.LookAt(camPos, targetPos, gglm.NewVec3(0, 1, 0))
	simpleMat.SetUnifMat4("viewMat", &viewMat.Mat4)
}

var time uint64 = 0
var name string = ""

var ambientColor gglm.Vec3 = *gglm.NewVec3(1, 1, 1)
var ambientColorStrength float32 = 0.1

var lightPos1 gglm.Vec3 = *gglm.NewVec3(2, 2, 0)
var lightColor1 gglm.Vec3 = *gglm.NewVec3(1, 1, 1)
