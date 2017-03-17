package graphics

import (
	"fmt"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/shocked-client/opengl"
)

var fillRectVertexShaderSource = `
#version 150
precision mediump float;

attribute vec2 vertexPosition;

uniform mat4 projectionMatrix;

void main(void) {
	gl_Position = projectionMatrix * vec4(vertexPosition, 0.0, 1.0);
}
`

var fillRectFragmentShaderSource = `
#version 150
precision mediump float;

uniform vec4 color;

void main(void) {
	gl_FragColor = color;
}
`

// RectangleRenderer renders rectangular shapes.
type RectangleRenderer struct {
	gl               opengl.OpenGl
	projectionMatrix *mgl.Mat4

	program                 uint32
	vao                     *opengl.VertexArrayObject
	vertexPositionBuffer    uint32
	vertexPositionAttrib    int32
	projectionMatrixUniform opengl.Matrix4Uniform
	colorUniform            opengl.Vector4Uniform
}

// NewRectangleRenderer returns a new instance of an RectangleRenderer type.
func NewRectangleRenderer(gl opengl.OpenGl, projectionMatrix *mgl.Mat4) *RectangleRenderer {
	program, programErr := opengl.LinkNewStandardProgram(gl, fillRectVertexShaderSource, fillRectFragmentShaderSource)

	if programErr != nil {
		panic(fmt.Errorf("BitmapTextureRenderer shader failed: %v", programErr))
	}
	renderer := &RectangleRenderer{
		gl:               gl,
		projectionMatrix: projectionMatrix,

		program:                 program,
		vao:                     opengl.NewVertexArrayObject(gl, program),
		vertexPositionBuffer:    gl.GenBuffers(1)[0],
		vertexPositionAttrib:    gl.GetAttribLocation(program, "vertexPosition"),
		projectionMatrixUniform: opengl.Matrix4Uniform(gl.GetUniformLocation(program, "projectionMatrix")),
		colorUniform:            opengl.Vector4Uniform(gl.GetUniformLocation(program, "color"))}

	renderer.vao.WithSetter(func(gl opengl.OpenGl) {
		gl.EnableVertexAttribArray(uint32(renderer.vertexPositionAttrib))
		gl.BindBuffer(opengl.ARRAY_BUFFER, renderer.vertexPositionBuffer)
		gl.VertexAttribOffset(uint32(renderer.vertexPositionAttrib), 2, opengl.FLOAT, false, 0, 0)
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	})

	return renderer
}

// Dispose clears any open resources.
func (renderer *RectangleRenderer) Dispose() {
	renderer.vao.Dispose()
	renderer.gl.DeleteProgram(renderer.program)
	renderer.gl.DeleteBuffers([]uint32{renderer.vertexPositionBuffer})
}

// Fill renders a rectangle filled with a solid color.
func (renderer *RectangleRenderer) Fill(left, top, right, bottom float32, fillColor Color) {
	gl := renderer.gl

	{
		var vertices = []float32{
			left, top,
			right, top,
			left, bottom,

			left, bottom,
			right, top,
			right, bottom}
		gl.BindBuffer(opengl.ARRAY_BUFFER, renderer.vertexPositionBuffer)
		gl.BufferData(opengl.ARRAY_BUFFER, len(vertices)*4, vertices, opengl.STATIC_DRAW)
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	}

	renderer.vao.OnShader(func() {
		renderer.projectionMatrixUniform.Set(gl, renderer.projectionMatrix)
		renderer.colorUniform.Set(gl, fillColor.AsVector())

		gl.DrawArrays(opengl.TRIANGLES, 0, 6)
	})
}
