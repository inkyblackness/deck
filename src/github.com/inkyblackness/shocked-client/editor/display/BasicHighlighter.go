package display

import (
	"fmt"
	"os"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/shocked-client/opengl"
)

var basicHighlighterVertexShaderSource = `
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
  #ifdef GL_ES
    precision mediump float;
  #endif

  varying vec4 color;

  void main(void) {
    gl_FragColor = color;
  }
`

// BasicHighlighter draws a simple highlighting of a rectangular area.
type BasicHighlighter struct {
	gl opengl.OpenGl

	program                 uint32
	vertexArrayObject       uint32
	vertexPositionBuffer    uint32
	vertexPositionAttrib    int32
	modelMatrixUniform      int32
	viewMatrixUniform       int32
	projectionMatrixUniform int32
	inColorUniform          int32
}

// NewBasicHighlighter returns a new instance of BasicHighlighter.
func NewBasicHighlighter(gl opengl.OpenGl, color [4]float32) *BasicHighlighter {
	vertexShader, err1 := opengl.CompileNewShader(gl, opengl.VERTEX_SHADER, basicHighlighterVertexShaderSource)
	defer gl.DeleteShader(vertexShader)
	fragmentShader, err2 := opengl.CompileNewShader(gl, opengl.FRAGMENT_SHADER, basicHighlighterFragmentShaderSource)
	defer gl.DeleteShader(fragmentShader)
	program, _ := opengl.LinkNewProgram(gl, vertexShader, fragmentShader)

	if err1 != nil {
		fmt.Fprintf(os.Stderr, "Failed to compile shader 1:\n", err1)
	}
	if err2 != nil {
		fmt.Fprintf(os.Stderr, "Failed to compile shader 2:\n", err2)
	}

	highlighter := &BasicHighlighter{
		gl:                      gl,
		program:                 program,
		vertexArrayObject:       gl.GenVertexArrays(1)[0],
		vertexPositionBuffer:    gl.GenBuffers(1)[0],
		vertexPositionAttrib:    gl.GetAttribLocation(program, "vertexPosition"),
		modelMatrixUniform:      gl.GetUniformLocation(program, "modelMatrix"),
		viewMatrixUniform:       gl.GetUniformLocation(program, "viewMatrix"),
		projectionMatrixUniform: gl.GetUniformLocation(program, "projectionMatrix"),
		inColorUniform:          gl.GetUniformLocation(program, "inColor")}

	highlighter.withShader(func() {
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

		gl.Uniform4fv(highlighter.inColorUniform, &color)
	})

	return highlighter
}

// Dispose releases all resources.
func (highlighter *BasicHighlighter) Dispose() {
	gl := highlighter.gl

	gl.DeleteBuffers([]uint32{highlighter.vertexPositionBuffer})
	gl.DeleteVertexArrays([]uint32{highlighter.vertexArrayObject})
	gl.DeleteShader(highlighter.program)
}

// Render renders the highlight.
func (highlighter *BasicHighlighter) Render(context *RenderContext, areas []Area) {
	gl := highlighter.gl

	highlighter.withShader(func() {
		highlighter.setMatrix(highlighter.viewMatrixUniform, context.ViewMatrix())
		highlighter.setMatrix(highlighter.projectionMatrixUniform, context.ProjectionMatrix())

		gl.BindBuffer(opengl.ARRAY_BUFFER, highlighter.vertexPositionBuffer)
		gl.VertexAttribOffset(uint32(highlighter.vertexPositionAttrib), 3, opengl.FLOAT, false, 0, 0)

		for _, area := range areas {
			x, y := area.Center()
			width, height := area.Size()
			modelMatrix := mgl.Ident4().
				Mul4(mgl.Translate3D(x, y, 0.0)).
				Mul4(mgl.Scale3D(width, height, 1.0))

			highlighter.setMatrix(highlighter.modelMatrixUniform, &modelMatrix)
			gl.DrawArrays(opengl.TRIANGLES, 0, 6)
		}
	})
}

func (highlighter *BasicHighlighter) withShader(task func()) {
	gl := highlighter.gl

	gl.UseProgram(highlighter.program)
	gl.BindVertexArray(highlighter.vertexArrayObject)
	gl.EnableVertexAttribArray(uint32(highlighter.vertexPositionAttrib))

	defer func() {
		gl.EnableVertexAttribArray(0)
		gl.BindVertexArray(0)
		gl.UseProgram(0)
	}()

	task()
}

func (highlighter *BasicHighlighter) setMatrix(uniform int32, matrix *mgl.Mat4) {
	matrixArray := ([16]float32)(*matrix)
	highlighter.gl.UniformMatrix4fv(uniform, false, &matrixArray)
}
