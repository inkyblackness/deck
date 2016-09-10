package display

import (
	"fmt"
	"os"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/opengl"
)

var textVertexShaderSource = `
  attribute vec3 vertexPosition;
  attribute vec3 uvPosition;

  uniform mat4 modelMatrix;
  uniform mat4 viewMatrix;
  uniform mat4 projectionMatrix;

  varying vec2 uv;

  void main(void) {
    gl_Position = projectionMatrix * viewMatrix * modelMatrix * vec4(vertexPosition, 1.0);

    uv = uvPosition.xy;
  }
`

var textFragmentShaderSource = `
  #ifdef GL_ES
    precision mediump float;
  #endif

  uniform sampler2D palette;
  uniform sampler2D bitmap;

  varying vec2 uv;

  void main(void) {
    vec4 pixel = texture2D(bitmap, uv);
    vec4 color = texture2D(palette, vec2(pixel.a, 0.5));

    if (pixel.a > 0.0) {
      gl_FragColor = color;
    } else {
      gl_FragColor = vec4(0.0, 0.0, 0.0, 0.0);
    }
  }
`

// TextRenderable is a renderable for texts.
type TextRenderable struct {
	gl opengl.OpenGl

	modelMatrix mgl.Mat4

	program                 uint32
	vertexArrayObject       uint32
	vertexPositionBuffer    uint32
	vertexPositionAttrib    int32
	uvPositionBuffer        uint32
	uvPositionAttrib        int32
	modelMatrixUniform      int32
	viewMatrixUniform       int32
	projectionMatrixUniform int32

	paletteUniform int32
	bitmapUniform  int32

	paletteTexture graphics.Texture
	bitmapTexture  graphics.Texture

	textRenderer graphics.TextRenderer
	text         string
}

// NewTextRenderable returns a new instance of a text renderable
func NewTextRenderable(gl opengl.OpenGl, positionX, positionY float32, displaySize float32,
	paletteTexture graphics.Texture) *TextRenderable {
	vertexShader, err1 := opengl.CompileNewShader(gl, opengl.VERTEX_SHADER, textVertexShaderSource)
	defer gl.DeleteShader(vertexShader)
	fragmentShader, err2 := opengl.CompileNewShader(gl, opengl.FRAGMENT_SHADER, textFragmentShaderSource)
	defer gl.DeleteShader(fragmentShader)
	program, _ := opengl.LinkNewProgram(gl, vertexShader, fragmentShader)

	if err1 != nil {
		fmt.Fprintf(os.Stderr, "Failed to compile shader 1:\n", err1)
	}
	if err2 != nil {
		fmt.Fprintf(os.Stderr, "Failed to compile shader 2:\n", err2)
	}

	renderable := &TextRenderable{
		gl:      gl,
		program: program,
		modelMatrix: mgl.Ident4().
			Mul4(mgl.Translate3D(positionX, positionY, 0.0)).
			Mul4(mgl.Scale3D(displaySize, displaySize, 1.0)),

		vertexArrayObject:       gl.GenVertexArrays(1)[0],
		vertexPositionBuffer:    gl.GenBuffers(1)[0],
		vertexPositionAttrib:    gl.GetAttribLocation(program, "vertexPosition"),
		uvPositionBuffer:        gl.GenBuffers(1)[0],
		uvPositionAttrib:        gl.GetAttribLocation(program, "uvPosition"),
		modelMatrixUniform:      gl.GetUniformLocation(program, "modelMatrix"),
		viewMatrixUniform:       gl.GetUniformLocation(program, "viewMatrix"),
		projectionMatrixUniform: gl.GetUniformLocation(program, "projectionMatrix"),
		paletteTexture:          paletteTexture,
		paletteUniform:          gl.GetUniformLocation(program, "palette"),
		bitmapUniform:           gl.GetUniformLocation(program, "bitmap"),
		text:                    ""}

	return renderable
}

// SetTextRenderer sets the renderer to use. This update the bitmap.
func (renderable *TextRenderable) SetTextRenderer(renderer graphics.TextRenderer) {
	renderable.textRenderer = renderer
	renderable.updateBitmap()
}

// SetText sets the text to display. This updates the bitmap.
func (renderable *TextRenderable) SetText(text string) {
	if renderable.text != text {
		renderable.text = text
		renderable.updateBitmap()
	}
}

func (renderable *TextRenderable) updateBitmap() {
	if renderable.bitmapTexture != nil {
		renderable.bitmapTexture.Dispose()
		renderable.bitmapTexture = nil
	}
	if renderable.textRenderer != nil {
		textBitmap := renderable.textRenderer.Render(renderable.text)
		renderable.bitmapTexture = graphics.NewBitmapTexture(renderable.gl, textBitmap.Width, textBitmap.Height, textBitmap.Pixels)

		renderable.withShader(func() {
			gl := renderable.gl
			width := float32(textBitmap.Width)
			height := float32(textBitmap.Height)
			var vertices = []float32{
				0.0, 0.0, 0.0,
				width, 0.0, 0.0,
				width, height, 0.0,

				width, height, 0.0,
				0.0, height, 0.0,
				0.0, 0.0, 0.0}
			gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexPositionBuffer)
			gl.BufferData(opengl.ARRAY_BUFFER, len(vertices)*4, vertices, opengl.STATIC_DRAW)

			limit := float32(1.0)
			var uv = []float32{
				0.0, 0.0, 0.0,
				limit, 0.0, 0.0,
				limit, limit, 0.0,

				limit, limit, 0.0,
				0.0, limit, 0.0,
				0.0, 0.0, 0.0}
			gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.uvPositionBuffer)
			gl.BufferData(opengl.ARRAY_BUFFER, len(uv)*4, uv, opengl.STATIC_DRAW)
		})
	}
}

// Render renders
func (renderable *TextRenderable) Render(context *RenderContext) {
	gl := renderable.gl

	if renderable.bitmapTexture != nil {
		renderable.withShader(func() {
			renderable.setMatrix(renderable.modelMatrixUniform, &renderable.modelMatrix)
			renderable.setMatrix(renderable.viewMatrixUniform, context.ViewMatrix())
			renderable.setMatrix(renderable.projectionMatrixUniform, context.ProjectionMatrix())

			gl.EnableVertexAttribArray(uint32(renderable.vertexPositionAttrib))
			gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexPositionBuffer)
			gl.VertexAttribOffset(uint32(renderable.vertexPositionAttrib), 3, opengl.FLOAT, false, 0, 0)
			gl.EnableVertexAttribArray(uint32(renderable.uvPositionAttrib))
			gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.uvPositionBuffer)
			gl.VertexAttribOffset(uint32(renderable.uvPositionAttrib), 3, opengl.FLOAT, false, 0, 0)

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
}

func (renderable *TextRenderable) withShader(task func()) {
	gl := renderable.gl

	gl.UseProgram(renderable.program)
	gl.BindVertexArray(renderable.vertexArrayObject)

	defer func() {
		gl.EnableVertexAttribArray(0)
		gl.BindVertexArray(0)
		gl.UseProgram(0)
	}()

	task()
}

func (renderable *TextRenderable) setMatrix(uniform int32, matrix *mgl.Mat4) {
	matrixArray := ([16]float32)(*matrix)
	renderable.gl.UniformMatrix4fv(uniform, false, &matrixArray)
}
