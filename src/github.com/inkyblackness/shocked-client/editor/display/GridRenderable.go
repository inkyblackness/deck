package display

import (
	"fmt"
	"os"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/shocked-client/opengl"
)

var gridVertexShaderSource = `
  attribute vec3 vertexPosition;

  uniform mat4 viewMatrix;
  uniform mat4 projectionMatrix;

  varying vec4 color;
  varying vec3 originalPosition;

  void main(void) {
    gl_Position = projectionMatrix * viewMatrix * vec4(vertexPosition, 1.0);

    color = vec4(0.0, 0.1, 0.0, 0.6);
    originalPosition = vertexPosition;
  }
`

var gridFragmentShaderSource = `
  #ifdef GL_ES
    precision mediump float;
  #endif

  varying vec4 color;
  varying vec3 originalPosition;

  float modulo(float x, float y) {
    return x - y * floor(x/y);
  }

  float nearGrid(float stepSize, float value) {
    float remainder = modulo(value - (stepSize / 2.0), stepSize) * 2.0;

    if (remainder >= stepSize) {
      remainder = (stepSize * 2.0) - remainder;
    }

    return remainder / stepSize;
  }

  void main(void) {
    float alphaX = nearGrid(32.0, originalPosition.x);
    float alphaY = nearGrid(32.0, originalPosition.y);
    bool beyondX = (originalPosition.x / 32.0) >= 64.0;
    bool beyondY = (originalPosition.y / 32.0) >= 64.0;
    float alpha = 0.0;

    if (!beyondX && !beyondY) {
       alpha = max(alphaX, alphaY);
    } else if (beyondX && !beyondY) {
       alpha = alphaX;
    } else if (beyondY && !beyondX) {
       alpha = alphaY;
    } else {
       alpha = min(alphaX, alphaY);
    }

    alpha = pow(2.0, 10.0 * (alpha - 1.0));

    gl_FragColor = vec4(color.rgb, color.a * alpha);
  }
`

// GridRenderable renders a grid with transparent holes.
type GridRenderable struct {
	gl opengl.OpenGl

	program                 uint32
	vertexArrayObject       uint32
	vertexPositionBuffer    uint32
	vertexPositionAttrib    int32
	viewMatrixUniform       int32
	projectionMatrixUniform int32
}

// NewGridRenderable returns a new instance of GridRenderable.
func NewGridRenderable(gl opengl.OpenGl) *GridRenderable {
	vertexShader, err1 := opengl.CompileNewShader(gl, opengl.VERTEX_SHADER, gridVertexShaderSource)
	defer gl.DeleteShader(vertexShader)
	fragmentShader, err2 := opengl.CompileNewShader(gl, opengl.FRAGMENT_SHADER, gridFragmentShaderSource)
	defer gl.DeleteShader(fragmentShader)
	program, _ := opengl.LinkNewProgram(gl, vertexShader, fragmentShader)

	if err1 != nil {
		fmt.Fprintf(os.Stderr, "Failed to compile shader 1:\n", err1)
	}
	if err2 != nil {
		fmt.Fprintf(os.Stderr, "Failed to compile shader 2:\n", err2)
	}

	renderable := &GridRenderable{
		gl:                      gl,
		program:                 program,
		vertexArrayObject:       gl.GenVertexArrays(1)[0],
		vertexPositionBuffer:    gl.GenBuffers(1)[0],
		vertexPositionAttrib:    gl.GetAttribLocation(program, "vertexPosition"),
		viewMatrixUniform:       gl.GetUniformLocation(program, "viewMatrix"),
		projectionMatrixUniform: gl.GetUniformLocation(program, "projectionMatrix")}

	renderable.withShader(func() {
		gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexPositionBuffer)
		half := float32(16.0)
		limit := float32(32.0*64.0 + half)
		var vertices = []float32{
			-half, -half, 0.0,
			limit, -half, 0.0,
			limit, limit, 0.0,

			limit, limit, 0.0,
			-half, limit, 0.0,
			-half, -half, 0.0}
		gl.BufferData(opengl.ARRAY_BUFFER, len(vertices)*4, vertices, opengl.STATIC_DRAW)
	})

	return renderable
}

// Render renders
func (renderable *GridRenderable) Render(context *RenderContext) {
	gl := renderable.gl

	renderable.withShader(func() {
		renderable.setMatrix(renderable.viewMatrixUniform, context.ViewMatrix())
		renderable.setMatrix(renderable.projectionMatrixUniform, context.ProjectionMatrix())

		gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexPositionBuffer)
		gl.VertexAttribOffset(uint32(renderable.vertexPositionAttrib), 3, opengl.FLOAT, false, 0, 0)

		gl.DrawArrays(opengl.TRIANGLES, 0, 6)
	})
}

func (renderable *GridRenderable) withShader(task func()) {
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

func (renderable *GridRenderable) setMatrix(uniform int32, matrix *mgl32.Mat4) {
	matrixArray := ([16]float32)(*matrix)
	renderable.gl.UniformMatrix4fv(uniform, false, &matrixArray)
}
