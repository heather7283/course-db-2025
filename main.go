package main

import (
	"runtime"
	"unsafe"
	_ "embed"

	"github.com/AllenDang/cimgui-go/backend"
	"github.com/AllenDang/cimgui-go/backend/glfwbackend"
	"github.com/AllenDang/cimgui-go/imgui"
	_ "github.com/AllenDang/cimgui-go/impl/glfw"
)

//go:embed assets/JetBrainsMonoNLNerdFont-Regular.ttf
var font []byte

var currentBackend backend.Backend[glfwbackend.GLFWWindowFlags]

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := dbOpen("db.sqlite3"); err != nil {
		panic(err)
	}

	currentBackend, _ = backend.CreateBackend(glfwbackend.NewGLFWBackend())
	currentBackend.SetAfterCreateContextHook(func() {
		fontDataPtr := uintptr(unsafe.Pointer(&font[0]))
		fontDataLen := int32(len(font))
		f := imgui.CurrentIO().Fonts().AddFontFromMemoryTTF(fontDataPtr, fontDataLen)
		imgui.CurrentIO().SetFontDefault(f)

		imgui.CurrentIO().SetIniFilename("/doesnotexist")
	})

	currentBackend.CreateWindow("Олимпиада", 1200, 900)

	initUI()
	currentBackend.Run(runUI)
}

