package display

import (
	"fmt"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/shocked-model"

	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/opengl"
)

var mapTileColoringVertexShaderSource = `
#version 150
precision mediump float;

in vec4 vertexColor;
in vec2 vertexPosition;

uniform mat4 modelMatrix;
uniform mat4 viewMatrix;
uniform mat4 projectionMatrix;

out vec4 color;

void main(void) {
   gl_Position = projectionMatrix * viewMatrix * modelMatrix * vec4(vertexPosition.xy, 0.0, 1.0);
   color = vertexColor;
}
`

var mapTileColoringFragmentShaderSource = `
#version 150
precision mediump float;

in vec4 color;

out vec4 fragColor;

void main(void) {
   fragColor = color;
}
`

// TilePropertiesQuery is a function to return the properties of requested tile, or nil if not available.
type TilePropertiesQuery func(x, y int) *model.TileProperties

// ColorQuery returns the 4 colors for the specified tile.
// Returned order: bottom-left, top-left, top-right, bottom-right.
type ColorQuery func(x, y int, properties *model.TileProperties, query TilePropertiesQuery) [4]graphics.Color

func darknessOfSingleTile(properties *model.TileProperties, resolver func(properties *model.TileProperties) int) graphics.Color {
	darkValue := 0
	if (properties != nil) && (properties.RealWorld != nil) {
		darkValue = resolver(properties)
	}
	return graphics.RGBA(0.0, 0.0, 0.0, float32(darkValue)/15)
}

// FloorShadow returns the colors for the floor darkness.
func FloorShadow(x, y int, properties *model.TileProperties, query TilePropertiesQuery) [4]graphics.Color {
	darkValue := func(prop *model.TileProperties) int { return *prop.RealWorld.FloorShadow }
	return [4]graphics.Color{
		darknessOfSingleTile(properties, darkValue),
		darknessOfSingleTile(query(x, y+1), darkValue),
		darknessOfSingleTile(query(x+1, y+1), darkValue),
		darknessOfSingleTile(query(x+1, y), darkValue)}
}

// CeilingShadow returns the colors for the ceiling darkness.
func CeilingShadow(x, y int, properties *model.TileProperties, query TilePropertiesQuery) [4]graphics.Color {
	darkValue := func(prop *model.TileProperties) int { return *prop.RealWorld.CeilingShadow }
	return [4]graphics.Color{
		darknessOfSingleTile(properties, darkValue),
		darknessOfSingleTile(query(x, y+1), darkValue),
		darknessOfSingleTile(query(x+1, y+1), darkValue),
		darknessOfSingleTile(query(x+1, y), darkValue)}
}

// TileColorMapRenderable is a renderable for the tile colorings.
type TileColorMapRenderable struct {
	context *graphics.RenderContext

	program                 uint32
	vao                     *opengl.VertexArrayObject
	vertexPositionBuffer    uint32
	vertexPositionAttrib    int32
	vertexColorBuffer       uint32
	vertexColorAttrib       int32
	modelMatrixUniform      opengl.Matrix4Uniform
	viewMatrixUniform       opengl.Matrix4Uniform
	projectionMatrixUniform opengl.Matrix4Uniform

	tiles [][]*model.TileProperties

	colorQuery ColorQuery
}

// NewTileColorMapRenderable returns a new instance of a renderable for tile colorings.
func NewTileColorMapRenderable(context *graphics.RenderContext) *TileColorMapRenderable {
	gl := context.OpenGl()
	program, programErr := opengl.LinkNewStandardProgram(gl, mapTileColoringVertexShaderSource, mapTileColoringFragmentShaderSource)

	if programErr != nil {
		panic(fmt.Errorf("TileColorMapRenderable shader failed: %v", programErr))
	}
	renderable := &TileColorMapRenderable{
		context:                 context,
		program:                 program,
		vao:                     opengl.NewVertexArrayObject(gl, program),
		vertexPositionBuffer:    gl.GenBuffers(1)[0],
		vertexPositionAttrib:    gl.GetAttribLocation(program, "vertexPosition"),
		vertexColorBuffer:       gl.GenBuffers(1)[0],
		vertexColorAttrib:       gl.GetAttribLocation(program, "vertexColor"),
		modelMatrixUniform:      opengl.Matrix4Uniform(gl.GetUniformLocation(program, "modelMatrix")),
		viewMatrixUniform:       opengl.Matrix4Uniform(gl.GetUniformLocation(program, "viewMatrix")),
		projectionMatrixUniform: opengl.Matrix4Uniform(gl.GetUniformLocation(program, "projectionMatrix")),

		tiles: make([][]*model.TileProperties, int(tilesPerMapSide)),

		colorQuery: nil}

	for i := 0; i < len(renderable.tiles); i++ {
		renderable.tiles[i] = make([]*model.TileProperties, int(tilesPerMapSide))
	}

	{
		top := float32(fineCoordinatesPerTileSide)
		left := float32(0.0)
		right := float32(fineCoordinatesPerTileSide)
		bottom := float32(0.0)

		vertices := []float32{
			left, bottom,
			left, top,
			right, top,

			right, top,
			right, bottom,
			left, bottom}

		gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexPositionBuffer)
		gl.BufferData(opengl.ARRAY_BUFFER, len(vertices)*4, vertices, opengl.STATIC_DRAW)
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	}

	renderable.vao.WithSetter(func(gl opengl.OpenGl) {
		gl.EnableVertexAttribArray(uint32(renderable.vertexPositionAttrib))
		gl.EnableVertexAttribArray(uint32(renderable.vertexColorAttrib))
		gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexPositionBuffer)
		gl.VertexAttribOffset(uint32(renderable.vertexPositionAttrib), 2, opengl.FLOAT, false, 0, 0)
		gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexColorBuffer)
		gl.VertexAttribOffset(uint32(renderable.vertexColorAttrib), 4, opengl.FLOAT, false, 0, 0)
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	})

	return renderable
}

// Dispose releases any internal resources
func (renderable *TileColorMapRenderable) Dispose() {
	gl := renderable.context.OpenGl()
	gl.DeleteProgram(renderable.program)
	gl.DeleteBuffers([]uint32{renderable.vertexPositionBuffer})
	renderable.vao.Dispose()
}

// SetColorQuery sets the query function for the coloring. nil disables coloring.
func (renderable *TileColorMapRenderable) SetColorQuery(colorQuery ColorQuery) {
	renderable.colorQuery = colorQuery
}

// SetTile sets the properties for the specified tile coordinate.
func (renderable *TileColorMapRenderable) SetTile(x, y int, properties *model.TileProperties) {
	renderable.tiles[y][x] = properties
}

func (renderable *TileColorMapRenderable) tile(x, y int) (properties *model.TileProperties) {
	if (y >= 0) && (y < len(renderable.tiles)) {
		row := renderable.tiles[y]
		if (x >= 0) && (x < len(renderable.tiles)) {
			properties = row[x]
		}
	}
	return
}

// Clear resets all tiles.
func (renderable *TileColorMapRenderable) Clear() {
	for _, row := range renderable.tiles {
		for index := 0; index < len(row); index++ {
			row[index] = nil
		}
	}
}

// Render renders
func (renderable *TileColorMapRenderable) Render() {
	gl := renderable.context.OpenGl()

	if renderable.colorQuery != nil {
		renderable.vao.OnShader(func() {
			colors := make([]float32, 24)

			renderable.viewMatrixUniform.Set(gl, renderable.context.ViewMatrix())
			renderable.projectionMatrixUniform.Set(gl, renderable.context.ProjectionMatrix())

			gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexColorBuffer)
			for y, row := range renderable.tiles {
				for x, tile := range row {
					if tile != nil {
						modelMatrix := mgl.Ident4().
							Mul4(mgl.Translate3D((float32(x))*fineCoordinatesPerTileSide, (float32(y))*fineCoordinatesPerTileSide, 0.0))
						renderable.modelMatrixUniform.Set(gl, &modelMatrix)

						tileColors := renderable.colorQuery(x, y, tile, renderable.tile)
						copy(colors[0:4], tileColors[0].AsVector()[:])
						copy(colors[4:8], tileColors[1].AsVector()[:])
						copy(colors[8:12], tileColors[2].AsVector()[:])
						copy(colors[12:16], tileColors[2].AsVector()[:])
						copy(colors[16:20], tileColors[3].AsVector()[:])
						copy(colors[20:24], tileColors[0].AsVector()[:])

						gl.BufferData(opengl.ARRAY_BUFFER, len(colors)*4, colors, opengl.STATIC_DRAW)
						gl.DrawArrays(opengl.TRIANGLES, 0, 6)
					}
				}
			}
			gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
		})
	}
}
