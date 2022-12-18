package physx

import (
	"errors"

	"github.com/bloeys/gglm/gglm"
	"github.com/bloeys/nmage/logging"
	"github.com/bloeys/physx-go/pgo"
)

type PhysX struct {
	Foundation pgo.Foundation
	Physics    pgo.Physics
	Scene      pgo.Scene
}

type PhysXCreationOptions struct {

	// Good defaults are length=1 (1m sizes), and speed=9.81 (speed of gravity)
	TypicalObjectLength float32

	// Good defaults are length=1 (1m sizes), and speed=9.81 (speed of gravity)
	TypicalObjectSpeed float32

	// If EnableVisualDebugger=true then all VisualDebuggerXYZ variables must be set
	EnableVisualDebugger bool
	VisualDebuggerHost   string

	// Default port is 5425
	VisualDebuggerPort                 int
	VisualDebuggerTimeoutMillis        int
	VisualDebuggerTransmitConstraints  bool
	VisualDebuggerTransmitContacts     bool
	VisualDebuggerTransmitSceneQueries bool

	SceneGravity *gglm.Vec3

	// Number of internal PhysX threads that do work.
	// If this is zero then all work is done on the thread that calls simulate
	SceneCPUDispatcherThreads uint32

	// Gets called when two objects collide
	SceneContactHandler func(cph pgo.ContactPairHeader)
}

func NewPhysx(options PhysXCreationOptions) (px *PhysX, err error) {

	if options.EnableVisualDebugger && !pgo.PvdSupported {
		return nil, errors.New("can not enable PhysX visual debugger (PVD) because physx-go is in release mode. Please build without the 'physx_release' tag")
	}

	logging.InfoLog.Printf("Initializing PhysX. Worker threads: %d. PVD supported: %v\n", options.SceneCPUDispatcherThreads, pgo.PvdSupported)

	// Setup foundation, pvd, and physics
	px = &PhysX{}
	px.Foundation = pgo.CreateFoundation()

	ts := pgo.NewTolerancesScale(options.TypicalObjectLength, options.TypicalObjectSpeed)
	if options.EnableVisualDebugger {

		pvdTr := pgo.DefaultPvdSocketTransportCreate(options.VisualDebuggerHost, options.VisualDebuggerPort, options.VisualDebuggerTimeoutMillis)
		pvd := pgo.CreatePvd(px.Foundation)
		logging.InfoLog.Println("Connecting to PVD...")
		if !pvd.Connect(pvdTr, pgo.PvdInstrumentationFlag_eALL) {
			return nil, errors.New("failed to connect to PhysX Visual Debugger (PVD). Is it running? Did you pass correct visual debugger host/port (default port is 5425)?")
		}
		logging.InfoLog.Println("Connected to PVD successfully")

		px.Physics = pgo.CreatePhysics(px.Foundation, ts, false, pvd)

	} else {
		px.Physics = pgo.CreatePhysics(px.Foundation, ts, false, nil)
	}

	// Setup scene
	sd := pgo.NewSceneDesc(ts)
	sd.SetGravity(pgo.NewVec3(options.SceneGravity.X(), options.SceneGravity.Y(), options.SceneGravity.Z()))

	defaultCpuDispatcher := pgo.DefaultCpuDispatcherCreate(options.SceneCPUDispatcherThreads, nil)
	sd.SetCpuDispatcher(defaultCpuDispatcher.ToCpuDispatcher())
	sd.SetOnContactCallback(options.SceneContactHandler)

	px.Scene = px.Physics.CreateScene(sd)
	if options.EnableVisualDebugger {

		scenePvdClient := px.Scene.GetScenePvdClient()
		scenePvdClient.SetScenePvdFlag(pgo.PvdSceneFlag_eTRANSMIT_CONSTRAINTS, options.VisualDebuggerTransmitConstraints)
		scenePvdClient.SetScenePvdFlag(pgo.PvdSceneFlag_eTRANSMIT_CONTACTS, options.VisualDebuggerTransmitContacts)
		scenePvdClient.SetScenePvdFlag(pgo.PvdSceneFlag_eTRANSMIT_SCENEQUERIES, options.VisualDebuggerTransmitSceneQueries)
	}

	return px, nil
}
