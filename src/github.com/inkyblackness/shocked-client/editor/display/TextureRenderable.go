package display

import (
	"fmt"
	"os"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/opengl"
)

var textureVertexShaderSource = `
  attribute vec3 vertexPosition;

  uniform mat4 modelMatrix;
  uniform mat4 viewMatrix;
  uniform mat4 projectionMatrix;

  varying vec2 uv;

  void main(void) {
    gl_Position = projectionMatrix * viewMatrix * modelMatrix * vec4(vertexPosition, 1.0);

    uv = vertexPosition.xy;
  }
`

var textureFragmentShaderSource = `
  #ifdef GL_ES
    precision mediump float;
  #endif

  uniform sampler2D palette;
  uniform sampler2D bitmap;

  varying vec2 uv;

  void main(void) {
    vec4 pixel = texture2D(bitmap, uv);
    vec4 color = texture2D(palette, vec2(pixel.a, 0.5));

    gl_FragColor = color;
  }
`

// TextureRenderable is a renderable for textures.
type TextureRenderable struct {
	gl opengl.OpenGl

	modelMatrix mgl.Mat4

	program                 uint32
	vertexArrayObject       uint32
	vertexPositionBuffer    uint32
	vertexPositionAttrib    int32
	modelMatrixUniform      int32
	viewMatrixUniform       int32
	projectionMatrixUniform int32

	paletteUniform int32
	bitmapUniform  int32

	paletteTexture graphics.Texture
	bitmapTexture  graphics.Texture
}

// NewTextureRenderable returns a new instance of a texture renderable
func NewTextureRenderable(gl opengl.OpenGl, positionX, positionY float32, displaySize float32,
	paletteTexture graphics.Texture, bitmapTexture graphics.Texture) *TextureRenderable {
	vertexShader, err1 := opengl.CompileNewShader(gl, opengl.VERTEX_SHADER, textureVertexShaderSource)
	defer gl.DeleteShader(vertexShader)
	fragmentShader, err2 := opengl.CompileNewShader(gl, opengl.FRAGMENT_SHADER, textureFragmentShaderSource)
	defer gl.DeleteShader(fragmentShader)
	program, _ := opengl.LinkNewProgram(gl, vertexShader, fragmentShader)

	if err1 != nil {
		fmt.Fprintf(os.Stderr, "Failed to compile shader 1:\n", err1)
	}
	if err2 != nil {
		fmt.Fprintf(os.Stderr, "Failed to compile shader 2:\n", err2)
	}

	renderable := &TextureRenderable{
		gl:      gl,
		program: program,
		modelMatrix: mgl.Ident4().
			Mul4(mgl.Translate3D(positionX, positionY, 0.0)).
			Mul4(mgl.Scale3D(displaySize, displaySize, 1.0)),

		vertexArrayObject:       gl.GenVertexArrays(1)[0],
		vertexPositionBuffer:    gl.GenBuffers(1)[0],
		vertexPositionAttrib:    gl.GetAttribLocation(program, "vertexPosition"),
		modelMatrixUniform:      gl.GetUniformLocation(program, "modelMatrix"),
		viewMatrixUniform:       gl.GetUniformLocation(program, "viewMatrix"),
		projectionMatrixUniform: gl.GetUniformLocation(program, "projectionMatrix"),
		paletteTexture:          paletteTexture,
		paletteUniform:          gl.GetUniformLocation(program, "palette"),
		bitmapTexture:           bitmapTexture,
		bitmapUniform:           gl.GetUniformLocation(program, "bitmap")}

	renderable.withShader(func() {
		gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexPositionBuffer)
		limit := float32(1.0)
		var vertices = []float32{
			0.0, 0.0, 0.0,
			limit, 0.0, 0.0,
			limit, limit, 0.0,

			limit, limit, 0.0,
			0.0, limit, 0.0,
			0.0, 0.0, 0.0}
		gl.BufferData(opengl.ARRAY_BUFFER, len(vertices)*4, vertices, opengl.STATIC_DRAW)
	})

	return renderable
}

// Render renders
func (renderable *TextureRenderable) Render(context *RenderContext) {
	gl := renderable.gl

	renderable.withShader(func() {
		renderable.setMatrix(renderable.modelMatrixUniform, &renderable.modelMatrix)
		renderable.setMatrix(renderable.viewMatrixUniform, context.ViewMatrix())
		renderable.setMatrix(renderable.projectionMatrixUniform, context.ProjectionMatrix())

		gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexPositionBuffer)
		gl.VertexAttribOffset(uint32(renderable.vertexPositionAttrib), 3, opengl.FLOAT, false, 0, 0)

		textureUnit := int32(0)
		gl.ActiveTexture(opengl.TEXTURE0 + uint32(textureUnit))
		gl.BindTexture(opengl.TEXTURE_2D, renderable.paletteTexture.Handle())
		gl.Uniform1i(renderable.paletteUniform, textureUnit)

		textureUnit = 1
		gl.ActiveTexture(opengl.TEXTURE0 + uint32(textureUnit))
		gl.BindTexture(opengl.TEXTURE_2D, renderable.bitmapTexture.Handle())
		gl.Uniform1i(renderable.bitmapUniform, textureUnit)

		gl.DrawArrays(opengl.TRIANGLES, 0, 6)

		gl.BindTexture(opengl.TEXTURE_2D, 0)
	})
}

func (renderable *TextureRenderable) withShader(task func()) {
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

func (renderable *TextureRenderable) setMatrix(uniform int32, matrix *mgl.Mat4) {
	matrixArray := ([16]float32)(*matrix)
	renderable.gl.UniformMatrix4fv(uniform, false, &matrixArray)
}
