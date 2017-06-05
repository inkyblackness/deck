package display

import (
	"fmt"
	"math"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/shocked-model"

	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/opengl"
)

var mapTileVertexShaderSource = `
#version 150
precision mediump float;

in vec3 vertexPosition;

uniform mat4 modelMatrix;
uniform mat4 viewMatrix;
uniform mat4 projectionMatrix;
uniform mat4 uvMatrix;

out vec2 uv;

void main(void) {
	gl_Position = projectionMatrix * viewMatrix * modelMatrix * vec4(vertexPosition, 1.0);

	uv = (uvMatrix * vec4(vertexPosition, 1.0)).xy;
}
`

var mapTileFragmentShaderSource = `
#version 150
precision mediump float;

in vec2 uv;

uniform sampler2D palette;
uniform sampler2D bitmap;

out vec4 fragColor;

void main(void) {
	vec4 pixel = texture2D(bitmap, uv);

	fragColor = texture2D(palette, vec2(pixel.a, 0.5));
}
`

// TextureQuery is a getter function to retrieve the texture for the given
// level texture index.
type TextureQuery func(index int) *graphics.BitmapTexture

// TextureIndexQuery is a getter function to retrieve the texture index and rotations
// from given tile properties.
type TextureIndexQuery func(properties *model.RealWorldTileProperties) (textureIndex int, textureRotations int)

// FloorTexture returns the information for the floor.
func FloorTexture(properties *model.RealWorldTileProperties) (textureIndex int, textureRotations int) {
	return *properties.FloorTexture, *properties.FloorTextureRotations
}

// CeilingTexture returns the information for the ceiling.
func CeilingTexture(properties *model.RealWorldTileProperties) (textureIndex int, textureRotations int) {
	return *properties.CeilingTexture, *properties.CeilingTextureRotations
}

// WallTexture returns the information for the wall.
func WallTexture(properties *model.RealWorldTileProperties) (textureIndex int, textureRotations int) {
	return *properties.WallTexture, 0
}

// TileTextureMapRenderable is a renderable for textures.
type TileTextureMapRenderable struct {
	context *graphics.RenderContext

	program                 uint32
	vao                     *opengl.VertexArrayObject
	vertexPositionBuffer    uint32
	vertexPositionAttrib    int32
	modelMatrixUniform      opengl.Matrix4Uniform
	viewMatrixUniform       opengl.Matrix4Uniform
	projectionMatrixUniform opengl.Matrix4Uniform
	uvMatrixUniform         opengl.Matrix4Uniform

	paletteUniform int32
	bitmapUniform  int32

	paletteTexture    graphics.Texture
	textureIndexQuery TextureIndexQuery
	textureQuery      TextureQuery

	tiles        [][]*model.TileProperties
	lastTileType model.TileType
}

var uvRotations map[int]*mgl.Mat4

func init() {
	uvRotations = make(map[int]*mgl.Mat4)
	for i := 0; i < 4; i++ {
		matrix := mgl.Translate3D(0.5, 0.5, 0.0).
			Mul4(mgl.HomogRotate3DZ(float32(math.Pi * float32(i) / 2.0))).
			Mul4(mgl.Translate3D(-0.5, -0.5, 0.0)).
			Mul4(mgl.Scale3D(1.0, -1.0, 1.0))
		uvRotations[i] = &matrix
	}
}

// NewTileTextureMapRenderable returns a new instance of a renderable for tile map textures.
func NewTileTextureMapRenderable(context *graphics.RenderContext, paletteTexture graphics.Texture,
	textureQuery TextureQuery) *TileTextureMapRenderable {
	gl := context.OpenGl()
	program, programErr := opengl.LinkNewStandardProgram(gl, mapTileVertexShaderSource, mapTileFragmentShaderSource)

	if programErr != nil {
		panic(fmt.Errorf("TileTextureMapRenderable shader failed: %v", programErr))
	}
	renderable := &TileTextureMapRenderable{
		context: context,
		program: program,

		vao:                     opengl.NewVertexArrayObject(gl, program),
		vertexPositionBuffer:    gl.GenBuffers(1)[0],
		vertexPositionAttrib:    gl.GetAttribLocation(program, "vertexPosition"),
		modelMatrixUniform:      opengl.Matrix4Uniform(gl.GetUniformLocation(program, "modelMatrix")),
		viewMatrixUniform:       opengl.Matrix4Uniform(gl.GetUniformLocation(program, "viewMatrix")),
		projectionMatrixUniform: opengl.Matrix4Uniform(gl.GetUniformLocation(program, "projectionMatrix")),
		uvMatrixUniform:         opengl.Matrix4Uniform(gl.GetUniformLocation(program, "uvMatrix")),
		paletteUniform:          gl.GetUniformLocation(program, "palette"),
		bitmapUniform:           gl.GetUniformLocation(program, "bitmap"),
		paletteTexture:          paletteTexture,
		textureIndexQuery:       FloorTexture,
		textureQuery:            textureQuery,
		tiles:                   make([][]*model.TileProperties, int(tilesPerMapSide)),
		lastTileType:            model.Solid}

	for i := 0; i < len(renderable.tiles); i++ {
		renderable.tiles[i] = make([]*model.TileProperties, int(tilesPerMapSide))
	}
	renderable.vao.WithSetter(func(gl opengl.OpenGl) {
		gl.EnableVertexAttribArray(uint32(renderable.vertexPositionAttrib))
		gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexPositionBuffer)
		gl.VertexAttribOffset(uint32(renderable.vertexPositionAttrib), 3, opengl.FLOAT, false, 0, 0)
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	})

	return renderable
}

// Dispose releases any internal resources
func (renderable *TileTextureMapRenderable) Dispose() {
	gl := renderable.context.OpenGl()

	renderable.vao.Dispose()
	gl.DeleteProgram(renderable.program)
	gl.DeleteBuffers([]uint32{renderable.vertexPositionBuffer})
}

// SetTextureIndexQuery sets which texture shall be shown.
func (renderable *TileTextureMapRenderable) SetTextureIndexQuery(query TextureIndexQuery) {
	renderable.textureIndexQuery = query
}

// SetTile sets the properties for the specified tile coordinate.
func (renderable *TileTextureMapRenderable) SetTile(x, y int, properties *model.TileProperties) {
	renderable.tiles[y][x] = properties
}

// Clear resets all tiles.
func (renderable *TileTextureMapRenderable) Clear() {
	for _, row := range renderable.tiles {
		for index := 0; index < len(row); index++ {
			row[index] = nil
		}
	}
}

// Render renders
func (renderable *TileTextureMapRenderable) Render() {
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

		scaling := mgl.Scale3D(fineCoordinatesPerTileSide, fineCoordinatesPerTileSide, 1.0)
		for y, row := range renderable.tiles {
			for x, tile := range row {
				if tile != nil && *tile.Type != model.Solid && tile.RealWorld != nil {
					textureIndex, textureRotations := renderable.textureIndexQuery(tile.RealWorld)
					texture := renderable.textureQuery(textureIndex)
					if texture != nil {
						modelMatrix := mgl.Translate3D(float32(x)*fineCoordinatesPerTileSide, float32(y)*fineCoordinatesPerTileSide, 0.0).
							Mul4(scaling)

						uvMatrix := uvRotations[textureRotations]
						renderable.uvMatrixUniform.Set(gl, uvMatrix)
						renderable.modelMatrixUniform.Set(gl, &modelMatrix)
						verticeCount := renderable.ensureTileType(*tile.Type)
						gl.BindTexture(opengl.TEXTURE_2D, texture.Handle())
						gl.Uniform1i(renderable.bitmapUniform, textureUnit)

						gl.DrawArrays(opengl.TRIANGLES, 0, int32(verticeCount))
					}
				}
			}
		}

		gl.BindTexture(opengl.TEXTURE_2D, 0)
	})
}

func (renderable *TileTextureMapRenderable) ensureTileType(tileType model.TileType) (verticeCount int) {
	displayedType := model.TileType(model.Open)

	verticeCount = 6
	if tileType == model.DiagonalOpenNorthEast || tileType == model.DiagonalOpenNorthWest ||
		tileType == model.DiagonalOpenSouthEast || tileType == model.DiagonalOpenSouthWest {
		displayedType = tileType
		verticeCount = 3
	}
	if renderable.lastTileType != displayedType {
		gl := renderable.context.OpenGl()
		var vertices []float32
		limit := float32(1.0)

		if displayedType == model.DiagonalOpenNorthEast {
			vertices = []float32{
				0.0, limit, 0.0,
				limit, limit, 0.0,
				limit, 0.0, 0.0}
		} else if displayedType == model.DiagonalOpenNorthWest {
			vertices = []float32{
				0.0, limit, 0.0,
				limit, limit, 0.0,
				0.0, 0.0, 0.0}
		} else if displayedType == model.DiagonalOpenSouthEast {
			vertices = []float32{
				limit, limit, 0.0,
				limit, 0.0, 0.0,
				0.0, 0.0, 0.0}
		} else if displayedType == model.DiagonalOpenSouthWest {
			vertices = []float32{
				0.0, limit, 0.0,
				limit, 0.0, 0.0,
				0.0, 0.0, 0.0}
		} else if displayedType == model.Open {
			vertices = []float32{
				0.0, 0.0, 0.0,
				limit, 0.0, 0.0,
				limit, limit, 0.0,

				limit, limit, 0.0,
				0.0, limit, 0.0,
				0.0, 0.0, 0.0}
		}
		gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexPositionBuffer)
		gl.BufferData(opengl.ARRAY_BUFFER, len(vertices)*4, vertices, opengl.STATIC_DRAW)
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)

		renderable.lastTileType = displayedType
	}

	return
}
