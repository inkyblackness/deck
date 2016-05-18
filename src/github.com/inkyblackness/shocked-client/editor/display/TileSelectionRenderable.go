package display

import (
	"fmt"
	"os"

	mgl "github.com/go-gl/mathgl/mgl32"

	editormodel "github.com/inkyblackness/shocked-client/editor/model"
	"github.com/inkyblackness/shocked-client/opengl"
)

var mapTileSelectionVertexShaderSource = `
  attribute vec3 vertexPosition;

  uniform mat4 modelMatrix;
  uniform mat4 viewMatrix;
  uniform mat4 projectionMatrix;

  void main(void) {
    gl_Position = projectionMatrix * viewMatrix * modelMatrix * vec4(vertexPosition, 1.0);
  }
`

var mapTileSelectionFragmentShaderSource = `
  #ifdef GL_ES
    precision mediump float;
  #endif

  void main(void) {
    gl_FragColor = vec4(1.0, 1.0, 1.0, 0.75);
  }
`

// TileSelectionCallback is a function receiving the coordinate of a selected tile.
type TileSelectionCallback func(coord editormodel.TileCoordinate)

// TileSelectionQuery is a query function to receive all selected tile coordinates.
type TileSelectionQuery func(TileSelectionCallback)

// TileSelectionRenderable is a renderable for textures.
type TileSelectionRenderable struct {
	gl opengl.OpenGl

	program                 uint32
	vertexArrayObject       uint32
	vertexPositionBuffer    uint32
	vertexPositionAttrib    int32
	modelMatrixUniform      int32
	viewMatrixUniform       int32
	projectionMatrixUniform int32

	query TileSelectionQuery
}

// NewTileSelectionRenderable returns a new instance of a renderable for tile selections
func NewTileSelectionRenderable(gl opengl.OpenGl, query TileSelectionQuery) *TileSelectionRenderable {
	vertexShader, err1 := opengl.CompileNewShader(gl, opengl.VERTEX_SHADER, mapTileSelectionVertexShaderSource)
	defer gl.DeleteShader(vertexShader)
	fragmentShader, err2 := opengl.CompileNewShader(gl, opengl.FRAGMENT_SHADER, mapTileSelectionFragmentShaderSource)
	defer gl.DeleteShader(fragmentShader)
	program, _ := opengl.LinkNewProgram(gl, vertexShader, fragmentShader)

	if err1 != nil {
		fmt.Fprintf(os.Stderr, "Failed to compile shader 1:\n", err1)
	}
	if err2 != nil {
		fmt.Fprintf(os.Stderr, "Failed to compile shader 2:\n", err2)
	}

	renderable := &TileSelectionRenderable{
		gl:                      gl,
		program:                 program,
		vertexArrayObject:       gl.GenVertexArrays(1)[0],
		vertexPositionBuffer:    gl.GenBuffers(1)[0],
		vertexPositionAttrib:    gl.GetAttribLocation(program, "vertexPosition"),
		modelMatrixUniform:      gl.GetUniformLocation(program, "modelMatrix"),
		viewMatrixUniform:       gl.GetUniformLocation(program, "viewMatrix"),
		projectionMatrixUniform: gl.GetUniformLocation(program, "projectionMatrix"),
		query: query}

	renderable.withShader(func() {
		limit := float32(32.0)
		var vertices = []float32{
			0.0, 0.0, 0.0,
			limit, 0.0, 0.0,
			limit, limit, 0.0,

			limit, limit, 0.0,
			0.0, limit, 0.0,
			0.0, 0.0, 0.0}

		gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexPositionBuffer)
		gl.BufferData(opengl.ARRAY_BUFFER, len(vertices)*4, vertices, opengl.STATIC_DRAW)
	})

	return renderable
}

// Dispose releases any internal resources
func (renderable *TileSelectionRenderable) Dispose() {
	renderable.gl.DeleteProgram(renderable.program)
	renderable.gl.DeleteBuffers([]uint32{renderable.vertexPositionBuffer})
	renderable.gl.DeleteVertexArrays([]uint32{renderable.vertexArrayObject})
}

// Render renders
func (renderable *TileSelectionRenderable) Render(context *RenderContext) {
	gl := renderable.gl

	renderable.withShader(func() {
		renderable.setMatrix32(renderable.viewMatrixUniform, context.ViewMatrix())
		renderable.setMatrix32(renderable.projectionMatrixUniform, context.ProjectionMatrix())

		gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexPositionBuffer)
		gl.VertexAttribOffset(uint32(renderable.vertexPositionAttrib), 3, opengl.FLOAT, false, 0, 0)

		renderable.query(func(coord editormodel.TileCoordinate) {
			x, y := coord.XY()
			modelMatrix := mgl.Translate3D(float32(x)*32.0, float32(63-y)*32.0, 0.0)
			renderable.setMatrix32(renderable.modelMatrixUniform, &modelMatrix)
			gl.DrawArrays(opengl.TRIANGLES, 0, 6)
		})
	})
}

func (renderable *TileSelectionRenderable) withShader(task func()) {
	gl := renderable.gl

	gl.UseProgram(renderable.program)
	gl.BindVertexArray(renderable.vertexArrayObject)
	gl.EnableVertexAttribArray(uint32(renderable.vertexPositionAttrib))

	defer func() {
		gl.EnableVertexAttribArray(0)
		gl.BindVertexArray(0)
		gl.UseProgram(0)
	}()

	task()
}

func (renderable *TileSelectionRenderable) setMatrix32(uniform int32, matrix *mgl.Mat4) {
	matrixArray := ([16]float32)(*matrix)
	renderable.gl.UniformMatrix4fv(uniform, false, &matrixArray)
}
