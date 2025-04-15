package opengl

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/glitter"
	"github.com/veandco/go-sdl2/sdl"
	"strings"
)

func InitializeGL(handle *sdl.Window) (string, sdl.GLContext) {
	// Create OpenGL context
	glContext, err := handle.GLCreateContext()
	if err != nil {
		core.Fatalf(ModuleName, "failed to create OpenGL context: %v\n", err)
	}

	sdl.GLSetSwapInterval(0)

	// Initialize OpenGL
	if err := gl.Init(); err != nil {
		core.Fatalf(glitter.ModuleName, "failed to initialize OpenGL: %v\n", err)
	}

	// Get and print extensions
	//numExtensions := int32(0)
	//gl.GetIntegerv(gl.NUM_EXTENSIONS, &numExtensions)
	//
	//for i := int32(0); i < numExtensions; i++ {
	//	extension := gl.GoStr(gl.GetStringi(gl.EXTENSIONS, uint32(i)))
	//	if strings.Contains(extension, "geometry") {
	//		fmt.Println("found geometry-related extension:", extension)
	//	}
	//}

	// Get OpenGL version
	return gl.GoStr(gl.GetString(gl.VERSION)), glContext
}

func CreateVBO(vertices []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	if len(vertices) > 0 {
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	return vbo
}

func CompileShader(src string, shaderType uint32) uint32 {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(src + "\x00")
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	// Check for compilation errors
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		core.Fatalf(glitter.ModuleName, "failed to compile shader: %v\n", log)
	}
	return shader
}

func LinkPrograms(shaderIDs ...uint32) uint32 {
	program := gl.CreateProgram()
	for _, shader := range shaderIDs {
		gl.AttachShader(program, shader)
	}
	gl.LinkProgram(program)

	// Check for linking errors
	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		core.Fatalf(glitter.ModuleName, "failed to link program: %v\n", log)
	}
	return program
}
