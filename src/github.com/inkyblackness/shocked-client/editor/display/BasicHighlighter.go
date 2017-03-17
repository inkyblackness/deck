package display

import (
	"fmt"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/opengl"
)

var basicHighlighterVertexShaderSource = `
#version 150
precision mediump float;

attribute vec3 vertexPosition;

uniform mat4 modelMatrix;
uniform mat4 viewMatrix;
uniform mat4 projectionMatrix;
uniform vec4 inColor;

varying vec4 color;

void main(void) {
	gl_Position = projectionMatrix * viewMatrix * modelMatrix * vec4(vertexPosition, 1.0);

	color = inColor;
}
`

var basicHighlighterFragmentShaderSource = `
#version 150
precision mediump float;

varying vec4 color;

void main(void) {
	gl_FragColor = color;
}
`

// BasicHighlighter draws a simple highlighting of a rectangular area.
type BasicHighlighter struct {
	context *graphics.RenderContext

	program                 uint32
	vao                     *opengl.VertexArrayObject
	vertexPositionBuffer    uint32
	vertexPositionAttrib    int32
	modelMatrixUniform      opengl.Matrix4Uniform
	viewMatrixUniform       opengl.Matrix4Uniform
	projectionMatrixUniform opengl.Matrix4Uniform
	inColorUniform          opengl.Vector4Uniform
}

// NewBasicHighlighter returns a new instance of BasicHighlighter.
func NewBasicHighlighter(context *graphics.RenderContext) *BasicHighlighter {
	gl := context.OpenGl()
	program, programErr := opengl.LinkNewStandardProgram(gl, basicHighlighterVertexShaderSource, basicHighlighterFragmentShaderSource)

	if programErr != nil {
		panic(fmt.Errorf("BasicHighlighter shader failed: %v", programErr))
	}
	highlighter := &BasicHighlighter{
		context: context,
		program: program,

		vao:                     opengl.NewVertexArrayObject(gl, program),
		vertexPositionBuffer:    gl.GenBuffers(1)[0],
		vertexPositionAttrib:    gl.GetAttribLocation(program, "vertexPosition"),
		modelMatrixUniform:      opengl.Matrix4Uniform(gl.GetUniformLocation(program, "modelMatrix")),
		viewMatrixUniform:       opengl.Matrix4Uniform(gl.GetUniformLocation(program, "viewMatrix")),
		projectionMatrixUniform: opengl.Matrix4Uniform(gl.GetUniformLocation(program, "projectionMatrix")),
		inColorUniform:          opengl.Vector4Uniform(gl.GetUniformLocation(program, "inColor"))}

	{
		gl.BindBuffer(opengl.ARRAY_BUFFER, highlighter.vertexPositionBuffer)
		half := float32(0.5)
		var vertices = []float32{
			-half, -half, 0.0,
			half, -half, 0.0,
			half, half, 0.0,

			half, half, 0.0,
			-half, half, 0.0,
			-half, -half, 0.0}
		gl.BufferData(opengl.ARRAY_BUFFER, len(vertices)*4, vertices, opengl.STATIC_DRAW)
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	}

	highlighter.vao.OnShader(func() {
		gl.EnableVertexAttribArray(uint32(highlighter.vertexPositionAttrib))
		gl.BindBuffer(opengl.ARRAY_BUFFER, highlighter.vertexPositionBuffer)
		gl.VertexAttribOffset(uint32(highlighter.vertexPositionAttrib), 3, opengl.FLOAT, false, 0, 0)
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	})

	return highlighter
}

// Dispose releases all resources.
func (highlighter *BasicHighlighter) Dispose() {
	gl := highlighter.context.OpenGl()

	highlighter.vao.Dispose()
	gl.DeleteBuffers([]uint32{highlighter.vertexPositionBuffer})
	gl.DeleteShader(highlighter.program)
}

// Render renders the highlights.
func (highlighter *BasicHighlighter) Render(areas []Area, color graphics.Color) {
	gl := highlighter.context.OpenGl()

	highlighter.vao.OnShader(func() {
		highlighter.viewMatrixUniform.Set(gl, highlighter.context.ViewMatrix())
		highlighter.projectionMatrixUniform.Set(gl, highlighter.context.ProjectionMatrix())
		highlighter.inColorUniform.Set(gl, color.AsVector())

		for _, area := range areas {
			x, y := area.Center()
			width, height := area.Size()
			modelMatrix := mgl.Ident4().
				Mul4(mgl.Translate3D(x, y, 0.0)).
				Mul4(mgl.Scale3D(width, height, 1.0))

			highlighter.modelMatrixUniform.Set(gl, &modelMatrix)

			gl.DrawArrays(opengl.TRIANGLES, 0, 6)
		}
	})
}
