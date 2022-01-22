package engine

import "github.com/veandco/go-sdl2/sdl"

type WindowFlags int

const (
	WindowFlags_FULLSCREEN         WindowFlags = sdl.WINDOW_FULLSCREEN
	WindowFlags_OPENGL             WindowFlags = sdl.WINDOW_OPENGL
	WindowFlags_SHOWN              WindowFlags = sdl.WINDOW_SHOWN
	WindowFlags_HIDDEN             WindowFlags = sdl.WINDOW_HIDDEN
	WindowFlags_BORDERLESS         WindowFlags = sdl.WINDOW_BORDERLESS
	WindowFlags_RESIZABLE          WindowFlags = sdl.WINDOW_RESIZABLE
	WindowFlags_MINIMIZED          WindowFlags = sdl.WINDOW_MINIMIZED
	WindowFlags_MAXIMIZED          WindowFlags = sdl.WINDOW_MAXIMIZED
	WindowFlags_INPUT_GRABBED      WindowFlags = sdl.WINDOW_INPUT_GRABBED
	WindowFlags_INPUT_FOCUS        WindowFlags = sdl.WINDOW_INPUT_FOCUS
	WindowFlags_MOUSE_FOCUS        WindowFlags = sdl.WINDOW_MOUSE_FOCUS
	WindowFlags_FULLSCREEN_DESKTOP WindowFlags = sdl.WINDOW_FULLSCREEN_DESKTOP
	WindowFlags_FOREIGN            WindowFlags = sdl.WINDOW_FOREIGN
	WindowFlags_ALLOW_HIGHDPI      WindowFlags = sdl.WINDOW_ALLOW_HIGHDPI
	WindowFlags_MOUSE_CAPTURE      WindowFlags = sdl.WINDOW_MOUSE_CAPTURE
	WindowFlags_ALWAYS_ON_TOP      WindowFlags = sdl.WINDOW_ALWAYS_ON_TOP
	WindowFlags_SKIP_TASKBAR       WindowFlags = sdl.WINDOW_SKIP_TASKBAR
	WindowFlags_UTILITY            WindowFlags = sdl.WINDOW_UTILITY
	WindowFlags_TOOLTIP            WindowFlags = sdl.WINDOW_TOOLTIP
	WindowFlags_POPUP_MENU         WindowFlags = sdl.WINDOW_POPUP_MENU
	// WindowFlags_VULKAN             WindowFlags = sdl.WINDOW_VULKAN
)
