package graphics

import (
	"fmt"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/shocked-client/opengl"
)

var bitmapTextureVertexShaderSource = `
#version 150
precision mediump float;

attribute vec2 vertexPosition;
attribute vec2 uvPosition;

uniform mat4 modelMatrix;
uniform mat4 viewMatrix;
uniform mat4 projectionMatrix;

varying vec2 uv;

void main(void) {
   gl_Position = projectionMatrix * viewMatrix * modelMatrix * vec4(vertexPosition, 0.0, 1.0);

   uv = uvPosition;
}
`

var bitmapTextureFragmentShaderSource = `
#version 150
precision mediump float;

uniform sampler2D palette;
uniform sampler2D bitmap;

varying vec2 uv;

void main(void) {
   vec4 pixel = texture2D(bitmap, uv);

   if (pixel.a > 0.0) {
      vec4 color = texture2D(palette, vec2(pixel.a, 0.5));

      gl_FragColor = color;
   } else {
      discard;
   }
}
`

// BitmapTextureRenderer renders bitmapped textures based on a palette.
type BitmapTextureRenderer struct {
	renderContext *RenderContext

	program                 uint32
	vao                     *opengl.VertexArrayObject
	vertexPositionBuffer    uint32
	vertexPositionAttrib    int32
	uvPositionAttrib        int32
	modelMatrixUniform      opengl.Matrix4Uniform
	viewMatrixUniform       opengl.Matrix4Uniform
	projectionMatrixUniform opengl.Matrix4Uniform

	paletteUniform int32
	bitmapUniform  int32

	paletteTexture Texture
}

// NewBitmapTextureRenderer returns a new instance of a texture renderer for bitmaps.
func NewBitmapTextureRenderer(renderContext *RenderContext, paletteTexture Texture) *BitmapTextureRenderer {
	gl := renderContext.OpenGl()
	program, programErr := opengl.LinkNewStandardProgram(gl, bitmapTextureVertexShaderSource, bitmapTextureFragmentShaderSource)

	if programErr != nil {
		panic(fmt.Errorf("BitmapTextureRenderer shader failed: %v", programErr))
	}
	renderer := &BitmapTextureRenderer{
		renderContext: renderContext,
		program:       program,

		vao:                     opengl.NewVertexArrayObject(gl, program),
		vertexPositionBuffer:    gl.GenBuffers(1)[0],
		vertexPositionAttrib:    gl.GetAttribLocation(program, "vertexPosition"),
		uvPositionAttrib:        gl.GetAttribLocation(program, "uvPosition"),
		modelMatrixUniform:      opengl.Matrix4Uniform(gl.GetUniformLocation(program, "modelMatrix")),
		viewMatrixUniform:       opengl.Matrix4Uniform(gl.GetUniformLocation(program, "viewMatrix")),
		projectionMatrixUniform: opengl.Matrix4Uniform(gl.GetUniformLocation(program, "projectionMatrix")),
		paletteTexture:          paletteTexture,
		paletteUniform:          gl.GetUniformLocation(program, "palette"),
		bitmapUniform:           gl.GetUniformLocation(program, "bitmap")}

	renderer.vao.WithSetter(func(gl opengl.OpenGl) {
		floatSize := int(4)
		stride := int32(4 * floatSize)
		gl.EnableVertexAttribArray(uint32(renderer.vertexPositionAttrib))
		gl.EnableVertexAttribArray(uint32(renderer.uvPositionAttrib))
		gl.BindBuffer(opengl.ARRAY_BUFFER, renderer.vertexPositionBuffer)
		gl.VertexAttribOffset(uint32(renderer.vertexPositionAttrib), 2, opengl.FLOAT, false, stride, 0*floatSize)
		gl.VertexAttribOffset(uint32(renderer.uvPositionAttrib), 2, opengl.FLOAT, false, stride, 2*floatSize)
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	})

	return renderer
}

// Dispose clears any resources.
func (renderer *BitmapTextureRenderer) Dispose() {
	gl := renderer.renderContext.OpenGl()

	renderer.vao.Dispose()
	gl.DeleteBuffers([]uint32{renderer.vertexPositionBuffer})
	gl.DeleteProgram(renderer.program)
}

// Render implements the TextureRenderer interface.
func (renderer *BitmapTextureRenderer) Render(modelMatrix *mgl.Mat4, texture Texture, textureRect Rectangle) {
	gl := renderer.renderContext.OpenGl()

	{
		baseRect := RectByCoord(0, 0, 1.0, 1.0)
		var vertices = []float32{
			baseRect.Left(), baseRect.Top(), textureRect.Left(), textureRect.Top(),
			baseRect.Left(), baseRect.Bottom(), textureRect.Left(), textureRect.Bottom(),
			baseRect.Right(), baseRect.Top(), textureRect.Right(), textureRect.Top(),

			baseRect.Right(), baseRect.Top(), textureRect.Right(), textureRect.Top(),
			baseRect.Left(), baseRect.Bottom(), textureRect.Left(), textureRect.Bottom(),
			baseRect.Right(), baseRect.Bottom(), textureRect.Right(), textureRect.Bottom()}
		gl.BindBuffer(opengl.ARRAY_BUFFER, renderer.vertexPositionBuffer)
		gl.BufferData(opengl.ARRAY_BUFFER, len(vertices)*4, vertices, opengl.STATIC_DRAW)
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	}

	renderer.vao.OnShader(func() {
		renderer.modelMatrixUniform.Set(gl, modelMatrix)
		renderer.viewMatrixUniform.Set(gl, renderer.renderContext.ViewMatrix())
		renderer.projectionMatrixUniform.Set(gl, renderer.renderContext.ProjectionMatrix())

		textureUnit := int32(0)
		gl.ActiveTexture(opengl.TEXTURE0 + uint32(textureUnit))
		gl.BindTexture(opengl.TEXTURE_2D, renderer.paletteTexture.Handle())
		gl.Uniform1i(renderer.paletteUniform, textureUnit)

		textureUnit = 1
		gl.ActiveTexture(opengl.TEXTURE0 + uint32(textureUnit))
		gl.Uniform1i(renderer.bitmapUniform, textureUnit)
		gl.BindTexture(opengl.TEXTURE_2D, texture.Handle())

		gl.DrawArrays(opengl.TRIANGLES, 0, 6)
	})
}
