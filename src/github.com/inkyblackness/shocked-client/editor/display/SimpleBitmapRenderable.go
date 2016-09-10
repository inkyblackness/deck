package display

import (
	"fmt"
	"os"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/opengl"
)

var simpleBitmapVertexShaderSource = `
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

var simpleBitmapFragmentShaderSource = `
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

// SimpleBitmapRenderable is a renderable for simple bitmaps.
type SimpleBitmapRenderable struct {
	gl opengl.OpenGl

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
}

// NewSimpleBitmapRenderable returns a new instance of a simple bitmap renderable
func NewSimpleBitmapRenderable(gl opengl.OpenGl, paletteTexture graphics.Texture) *SimpleBitmapRenderable {
	vertexShader, err1 := opengl.CompileNewShader(gl, opengl.VERTEX_SHADER, simpleBitmapVertexShaderSource)
	defer gl.DeleteShader(vertexShader)
	fragmentShader, err2 := opengl.CompileNewShader(gl, opengl.FRAGMENT_SHADER, simpleBitmapFragmentShaderSource)
	defer gl.DeleteShader(fragmentShader)
	program, _ := opengl.LinkNewProgram(gl, vertexShader, fragmentShader)

	if err1 != nil {
		fmt.Fprintf(os.Stderr, "Failed to compile shader 1:\n", err1)
	}
	if err2 != nil {
		fmt.Fprintf(os.Stderr, "Failed to compile shader 2:\n", err2)
	}

	renderable := &SimpleBitmapRenderable{
		gl:      gl,
		program: program,

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
		bitmapUniform:           gl.GetUniformLocation(program, "bitmap")}

	renderable.withShader(func() {
		gl := renderable.gl
		half := float32(0.5)
		var vertices = []float32{
			-half, -half, 0.0,
			half, -half, 0.0,
			half, half, 0.0,

			half, half, 0.0,
			-half, half, 0.0,
			-half, -half, 0.0}
		gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexPositionBuffer)
		gl.BufferData(opengl.ARRAY_BUFFER, len(vertices)*4, vertices, opengl.STATIC_DRAW)
	})

	return renderable
}

// Render renders the bitmap with the center at given position
func (renderable *SimpleBitmapRenderable) Render(context *RenderContext, icons []PlacedIcon) {
	gl := renderable.gl

	renderable.withShader(func() {
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
		gl.Uniform1i(renderable.bitmapUniform, textureUnit)
		for _, icon := range icons {
			x, y := icon.Center()
			u, v := icon.Icon().UV()
			width, height := renderable.limitedSize(icon)
			modelMatrix := mgl.Ident4().
				Mul4(mgl.Translate3D(x, y, 0.0)).
				Mul4(mgl.Scale3D(width, height, 1.0))

			renderable.setMatrix(renderable.modelMatrixUniform, &modelMatrix)

			var uv = []float32{
				0.0, 0.0, 0.0,
				u, 0.0, 0.0,
				u, v, 0.0,

				u, v, 0.0,
				0.0, v, 0.0,
				0.0, 0.0, 0.0}
			gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.uvPositionBuffer)
			gl.BufferData(opengl.ARRAY_BUFFER, len(uv)*4, uv, opengl.STATIC_DRAW)

			gl.BindTexture(opengl.TEXTURE_2D, icon.Icon().Handle())
			gl.DrawArrays(opengl.TRIANGLES, 0, 6)
		}
		gl.BindTexture(opengl.TEXTURE_2D, 0)
	})
}

func (renderable *SimpleBitmapRenderable) limitedSize(icon PlacedIcon) (width, height float32) {
	width, height = icon.Icon().Size()
	larger := width

	if larger < height {
		larger = height
	}
	if larger > 16.0 {
		ratio := 16.0 / larger
		width *= ratio
		height *= ratio
	}

	return
}

func (renderable *SimpleBitmapRenderable) withShader(task func()) {
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

func (renderable *SimpleBitmapRenderable) setMatrix(uniform int32, matrix *mgl.Mat4) {
	matrixArray := ([16]float32)(*matrix)
	renderable.gl.UniformMatrix4fv(uniform, false, &matrixArray)
}
