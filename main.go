package main

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"runtime"
	"strings"
)

const FLOAT_SIZE = 4 // size of float32/64 guaranteed to be 4 bytes

const vertexShaderSource = `
#version 330
layout (location = 0) in vec3 aPos;

void main()
{
gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);
}`

const fragmentShaderSource = `
#version 330
out vec4 FragColor;

void main()
{
FragColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);
}`

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(vertexShader)

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(fragmentShader)

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	return program, nil
}

func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
}

func framebufferSizeCallback(window *glfw.Window, width, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}

func main() {
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(800, 600, "Engine", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	window.SetFramebufferSizeCallback(framebufferSizeCallback)

	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.Viewport(0, 0, 800, 600)

	vertices := []float32{
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		0.0, 0.5, 0.0,
	}

	var vbo, vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	defer gl.DeleteVertexArrays(1, &vao)
	defer gl.DeleteBuffers(1, &vbo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*FLOAT_SIZE, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*FLOAT_SIZE, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	gl.BindVertexArray(0)

	program, err := newProgram(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		panic(err)
	}

	for !window.ShouldClose() {
		processInput(window)

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.UseProgram(program)
		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
