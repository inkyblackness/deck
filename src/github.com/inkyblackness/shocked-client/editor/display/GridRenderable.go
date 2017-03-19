package display

import (
	"fmt"

	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/opengl"
)

var gridVertexShaderSource = `
#version 150
precision mediump float;

in vec3 vertexPosition;

uniform mat4 viewMatrix;
uniform mat4 projectionMatrix;

out vec4 gridColor;
out vec3 originalPosition;

void main(void) {
   gl_Position = projectionMatrix * viewMatrix * vec4(vertexPosition, 1.0);

   gridColor = vec4(0.0, 0.1, 0.0, 0.6);
   originalPosition = vertexPosition;
}
`

var gridFragmentShaderSource = `
#version 150
precision mediump float;

in vec4 gridColor;
in vec3 originalPosition;

out vec4 fragColor;

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
   float alphaX = nearGrid(256.0, originalPosition.x);
   float alphaY = nearGrid(256.0, originalPosition.y);
   bool beyondX = (originalPosition.x / 256.0) >= 64.0 || (originalPosition.x < 0.0);
   bool beyondY = (originalPosition.y / 256.0) >= 64.0 || (originalPosition.y < 0.0);
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

   fragColor = vec4(gridColor.rgb, gridColor.a * alpha);
}
`

// GridRenderable renders a grid with transparent holes.
type GridRenderable struct {
	context *graphics.RenderContext

	program                 uint32
	vao                     *opengl.VertexArrayObject
	vertexPositionBuffer    uint32
	vertexPositionAttrib    int32
	viewMatrixUniform       opengl.Matrix4Uniform
	projectionMatrixUniform opengl.Matrix4Uniform
}

// NewGridRenderable returns a new instance of GridRenderable.
func NewGridRenderable(context *graphics.RenderContext) *GridRenderable {
	gl := context.OpenGl()
	program, programErr := opengl.LinkNewStandardProgram(gl, gridVertexShaderSource, gridFragmentShaderSource)

	if programErr != nil {
		panic(fmt.Errorf("GridRenderable shader failed: %v", programErr))
	}
	renderable := &GridRenderable{
		context:                 context,
		program:                 program,
		vao:                     opengl.NewVertexArrayObject(gl, program),
		vertexPositionBuffer:    gl.GenBuffers(1)[0],
		vertexPositionAttrib:    gl.GetAttribLocation(program, "vertexPosition"),
		viewMatrixUniform:       opengl.Matrix4Uniform(gl.GetUniformLocation(program, "viewMatrix")),
		projectionMatrixUniform: opengl.Matrix4Uniform(gl.GetUniformLocation(program, "projectionMatrix"))}

	{
		gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexPositionBuffer)
		half := fineCoordinatesPerTileSide / float32(2.0)
		limit := float32(fineCoordinatesPerTileSide*tilesPerMapSide + half)
		var vertices = []float32{
			-half, -half, 0.0,
			limit, -half, 0.0,
			limit, limit, 0.0,

			limit, limit, 0.0,
			-half, limit, 0.0,
			-half, -half, 0.0}
		gl.BufferData(opengl.ARRAY_BUFFER, len(vertices)*4, vertices, opengl.STATIC_DRAW)
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	}
	renderable.vao.WithSetter(func(gl opengl.OpenGl) {
		gl.EnableVertexAttribArray(uint32(renderable.vertexPositionAttrib))
		gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexPositionBuffer)
		gl.VertexAttribOffset(uint32(renderable.vertexPositionAttrib), 3, opengl.FLOAT, false, 0, 0)
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	})

	return renderable
}

// Render renders
func (renderable *GridRenderable) Render() {
	gl := renderable.context.OpenGl()

	renderable.vao.OnShader(func() {
		renderable.viewMatrixUniform.Set(gl, renderable.context.ViewMatrix())
		renderable.projectionMatrixUniform.Set(gl, renderable.context.ProjectionMatrix())

		gl.DrawArrays(opengl.TRIANGLES, 0, 6)
	})
}
