package display

import (
	"fmt"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/opengl"
)

var placedIconsVertexShaderSource = `
#version 150
precision mediump float;

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

var placedIconsFragmentShaderSource = `
#version 150
precision mediump float;

uniform sampler2D palette;
uniform sampler2D bitmap;

varying vec2 uv;

void main(void) {
	vec4 pixel = texture2D(bitmap, uv);
	vec4 color = texture2D(palette, vec2(pixel.a, 0.5));

	if (pixel.a > 0.0) {
		gl_FragColor = color;
	} else {
		discard;
	}
}
`

// PlacedIconsRenderable is a renderable for simple bitmaps.
type PlacedIconsRenderable struct {
	context *graphics.RenderContext

	program                 uint32
	vao                     *opengl.VertexArrayObject
	vertexPositionBuffer    uint32
	vertexPositionAttrib    int32
	uvPositionBuffer        uint32
	uvPositionAttrib        int32
	modelMatrixUniform      opengl.Matrix4Uniform
	viewMatrixUniform       opengl.Matrix4Uniform
	projectionMatrixUniform opengl.Matrix4Uniform

	paletteUniform int32
	bitmapUniform  int32

	paletteTexture graphics.Texture
}

// NewPlacedIconsRenderable returns a new instance of a simple bitmap renderable
func NewPlacedIconsRenderable(context *graphics.RenderContext, paletteTexture graphics.Texture) *PlacedIconsRenderable {
	gl := context.OpenGl()
	program, programErr := opengl.LinkNewStandardProgram(gl, placedIconsVertexShaderSource, placedIconsFragmentShaderSource)

	if programErr != nil {
		panic(fmt.Errorf("PlacedIconsRenderable shader failed: %v", programErr))
	}
	renderable := &PlacedIconsRenderable{
		context: context,
		program: program,

		vao:                     opengl.NewVertexArrayObject(gl, program),
		vertexPositionBuffer:    gl.GenBuffers(1)[0],
		vertexPositionAttrib:    gl.GetAttribLocation(program, "vertexPosition"),
		uvPositionBuffer:        gl.GenBuffers(1)[0],
		uvPositionAttrib:        gl.GetAttribLocation(program, "uvPosition"),
		modelMatrixUniform:      opengl.Matrix4Uniform(gl.GetUniformLocation(program, "modelMatrix")),
		viewMatrixUniform:       opengl.Matrix4Uniform(gl.GetUniformLocation(program, "viewMatrix")),
		projectionMatrixUniform: opengl.Matrix4Uniform(gl.GetUniformLocation(program, "projectionMatrix")),
		paletteTexture:          paletteTexture,
		paletteUniform:          gl.GetUniformLocation(program, "palette"),
		bitmapUniform:           gl.GetUniformLocation(program, "bitmap")}

	{
		half := float32(0.5)
		var vertices = []float32{
			-half, half, 0.0,
			half, half, 0.0,
			half, -half, 0.0,

			half, -half, 0.0,
			-half, -half, 0.0,
			-half, half, 0.0}
		gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexPositionBuffer)
		gl.BufferData(opengl.ARRAY_BUFFER, len(vertices)*4, vertices, opengl.STATIC_DRAW)
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	}
	renderable.vao.WithSetter(func(gl opengl.OpenGl) {
		gl.EnableVertexAttribArray(uint32(renderable.vertexPositionAttrib))
		gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexPositionBuffer)
		gl.VertexAttribOffset(uint32(renderable.vertexPositionAttrib), 3, opengl.FLOAT, false, 0, 0)
		gl.EnableVertexAttribArray(uint32(renderable.uvPositionAttrib))
		gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.uvPositionBuffer)
		gl.VertexAttribOffset(uint32(renderable.uvPositionAttrib), 3, opengl.FLOAT, false, 0, 0)
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	})

	return renderable
}

// Render renders the icons with their center at given position
func (renderable *PlacedIconsRenderable) Render(icons []PlacedIcon) {
	gl := renderable.context.OpenGl()

	renderable.vao.OnShader(func() {
		renderable.viewMatrixUniform.Set(gl, renderable.context.ViewMatrix())
		renderable.projectionMatrixUniform.Set(gl, renderable.context.ProjectionMatrix())

		textureUnit := int32(0)
		gl.ActiveTexture(opengl.TEXTURE0 + uint32(textureUnit))
		gl.BindTexture(opengl.TEXTURE_2D, renderable.paletteTexture.Handle())
		gl.Uniform1i(renderable.paletteUniform, textureUnit)

		textureUnit = 1
		gl.ActiveTexture(opengl.TEXTURE0 + uint32(textureUnit))
		gl.Uniform1i(renderable.bitmapUniform, textureUnit)
		for _, icon := range icons {
			texture := icon.Texture()

			if texture != nil {
				x, y := icon.Center()
				u, v := texture.UV()
				width, height := renderable.limitedSize(texture)
				modelMatrix := mgl.Ident4().
					Mul4(mgl.Translate3D(x, y, 0.0)).
					Mul4(mgl.Scale3D(width, height, 1.0))

				renderable.modelMatrixUniform.Set(gl, &modelMatrix)

				var uv = []float32{
					0.0, 0.0, 0.0,
					u, 0.0, 0.0,
					u, v, 0.0,

					u, v, 0.0,
					0.0, v, 0.0,
					0.0, 0.0, 0.0}
				gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.uvPositionBuffer)
				gl.BufferData(opengl.ARRAY_BUFFER, len(uv)*4, uv, opengl.STATIC_DRAW)
				gl.BindBuffer(opengl.ARRAY_BUFFER, 0)

				gl.BindTexture(opengl.TEXTURE_2D, texture.Handle())
				gl.DrawArrays(opengl.TRIANGLES, 0, 6)
			}
		}
		gl.BindTexture(opengl.TEXTURE_2D, 0)
	})
}

func (renderable *PlacedIconsRenderable) limitedSize(texture *graphics.BitmapTexture) (width, height float32) {
	width, height = texture.Size()
	referenceSize := float32(16.0)
	larger := width

	if larger < height {
		larger = height
	}
	if larger > referenceSize {
		ratio := referenceSize / larger
		width *= ratio
		height *= ratio
	}
	width *= (iconSize / referenceSize)
	height *= (iconSize / referenceSize)

	return
}
