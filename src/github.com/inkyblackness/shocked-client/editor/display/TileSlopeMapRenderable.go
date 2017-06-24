package display

import (
	"fmt"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/shocked-model"

	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/opengl"
)

var mapTileSlopeVertexShaderSource = `
#version 150
precision mediump float;

in vec3 vertexPosition;

uniform mat4 modelMatrix;
uniform mat4 viewMatrix;
uniform mat4 projectionMatrix;

out float hue;

void main(void) {
   gl_Position = projectionMatrix * viewMatrix * modelMatrix * vec4(vertexPosition.xy, 0.0, 1.0);
   hue = vertexPosition.z;
}
`

var mapTileSlopeFragmentShaderSource = `
#version 150
precision mediump float;

in float hue;

out vec4 fragColor;

vec3 hsv2rgb(vec3 c) {
   vec4 K = vec4(1.0, 2.0 / 3.0, 1.0 / 3.0, 3.0);
   vec3 p = abs(fract(c.xxx + K.xyz) * 6.0 - K.www);

   return c.z * mix(K.xxx, clamp(p - K.xxx, 0.0, 1.0), c.y);
}

void main(void) {
   fragColor = vec4(hsv2rgb(vec3(hue, 1.0, 0.8)), 0.5);
}
`

// TileSlopeMapRenderable is a renderable for the tile slopes.
type TileSlopeMapRenderable struct {
	context *graphics.RenderContext

	program                 uint32
	vao                     *opengl.VertexArrayObject
	vertexPositionBuffer    uint32
	vertexPositionAttrib    int32
	modelMatrixUniform      opengl.Matrix4Uniform
	viewMatrixUniform       opengl.Matrix4Uniform
	projectionMatrixUniform opengl.Matrix4Uniform

	tiles [][]*model.TileProperties
}

// NewTileSlopeMapRenderable returns a new instance of a renderable for tile slopes.
func NewTileSlopeMapRenderable(context *graphics.RenderContext) *TileSlopeMapRenderable {
	gl := context.OpenGl()
	program, programErr := opengl.LinkNewStandardProgram(gl, mapTileSlopeVertexShaderSource, mapTileSlopeFragmentShaderSource)

	if programErr != nil {
		panic(fmt.Errorf("TileSlopeMapRenderable shader failed: %v", programErr))
	}
	renderable := &TileSlopeMapRenderable{
		context:                 context,
		program:                 program,
		vao:                     opengl.NewVertexArrayObject(gl, program),
		vertexPositionBuffer:    gl.GenBuffers(1)[0],
		vertexPositionAttrib:    gl.GetAttribLocation(program, "vertexPosition"),
		modelMatrixUniform:      opengl.Matrix4Uniform(gl.GetUniformLocation(program, "modelMatrix")),
		viewMatrixUniform:       opengl.Matrix4Uniform(gl.GetUniformLocation(program, "viewMatrix")),
		projectionMatrixUniform: opengl.Matrix4Uniform(gl.GetUniformLocation(program, "projectionMatrix")),

		tiles: make([][]*model.TileProperties, int(tilesPerMapSide))}

	for i := 0; i < len(renderable.tiles); i++ {
		renderable.tiles[i] = make([]*model.TileProperties, int(tilesPerMapSide))
	}

	renderable.vao.OnShader(func() {
		dotHalf := float32(0.05)
		dotBase := float32(0.5) - (dotHalf * 2.0)
		floorHue := float32(0.3)
		ceilingHue := float32(0.0)

		top := dotBase + dotHalf
		topEnd := dotBase - dotHalf
		left := -dotBase - dotHalf
		leftEnd := -dotBase + dotHalf

		right := top
		rightEnd := topEnd
		bottom := left
		bottomEnd := leftEnd

		vertices := []float32{
			left, top, floorHue, leftEnd, top, floorHue, left, topEnd, floorHue,
			leftEnd, top, ceilingHue, leftEnd, topEnd, ceilingHue, left, topEnd, ceilingHue,

			right, top, floorHue, right, topEnd, floorHue, rightEnd, top, floorHue,
			rightEnd, top, ceilingHue, right, topEnd, ceilingHue, rightEnd, topEnd, ceilingHue,

			right, bottom, floorHue, rightEnd, bottom, floorHue, right, bottomEnd, floorHue,
			right, bottomEnd, ceilingHue, rightEnd, bottom, ceilingHue, rightEnd, bottomEnd, ceilingHue,

			left, bottomEnd, floorHue, leftEnd, bottom, floorHue, left, bottom, floorHue,
			left, bottomEnd, ceilingHue, leftEnd, bottomEnd, ceilingHue, leftEnd, bottom, ceilingHue}

		gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexPositionBuffer)
		gl.BufferData(opengl.ARRAY_BUFFER, len(vertices)*4, vertices, opengl.STATIC_DRAW)
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)

	})

	renderable.vao.WithSetter(func(gl opengl.OpenGl) {
		gl.EnableVertexAttribArray(uint32(renderable.vertexPositionAttrib))
		gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexPositionBuffer)
		gl.VertexAttribOffset(uint32(renderable.vertexPositionAttrib), 3, opengl.FLOAT, false, 0, 0)
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	})

	return renderable
}

// Dispose releases any internal resources
func (renderable *TileSlopeMapRenderable) Dispose() {
	gl := renderable.context.OpenGl()
	gl.DeleteProgram(renderable.program)
	gl.DeleteBuffers([]uint32{renderable.vertexPositionBuffer})
	renderable.vao.Dispose()
}

// SetTile sets the properties for the specified tile coordinate.
func (renderable *TileSlopeMapRenderable) SetTile(x, y int, properties *model.TileProperties) {
	renderable.tiles[y][x] = properties
}

// Clear resets all tiles.
func (renderable *TileSlopeMapRenderable) Clear() {
	for _, row := range renderable.tiles {
		for index := 0; index < len(row); index++ {
			row[index] = nil
		}
	}
}

var slopeTicksByType = map[model.TileType][]int{
	model.SlopeSouthToNorth: []int{0, 1},
	model.SlopeWestToEast:   []int{1, 2},
	model.SlopeNorthToSouth: []int{2, 3},
	model.SlopeEastToWest:   []int{3, 0},

	model.ValleySouthEastToNorthWest: []int{3, 0, 1},
	model.ValleySouthWestToNorthEast: []int{0, 1, 2},
	model.ValleyNorthWestToSouthEast: []int{1, 2, 3},
	model.ValleyNorthEastToSouthWest: []int{2, 3, 0},

	model.RidgeSouthEastToNorthWest: []int{0},
	model.RidgeSouthWestToNorthEast: []int{1},
	model.RidgeNorthWestToSouthEast: []int{2},
	model.RidgeNorthEastToSouthWest: []int{3}}

var invertedTileTypes = map[model.TileType]model.TileType{
	model.SlopeSouthToNorth: model.SlopeNorthToSouth,
	model.SlopeWestToEast:   model.SlopeEastToWest,
	model.SlopeNorthToSouth: model.SlopeSouthToNorth,
	model.SlopeEastToWest:   model.SlopeWestToEast,

	model.ValleySouthEastToNorthWest: model.RidgeNorthWestToSouthEast,
	model.ValleySouthWestToNorthEast: model.RidgeNorthEastToSouthWest,
	model.ValleyNorthWestToSouthEast: model.RidgeSouthEastToNorthWest,
	model.ValleyNorthEastToSouthWest: model.RidgeSouthWestToNorthEast,

	model.RidgeSouthEastToNorthWest: model.ValleyNorthWestToSouthEast,
	model.RidgeSouthWestToNorthEast: model.ValleyNorthEastToSouthWest,
	model.RidgeNorthWestToSouthEast: model.ValleySouthEastToNorthWest,
	model.RidgeNorthEastToSouthWest: model.ValleySouthWestToNorthEast}

// Render renders
func (renderable *TileSlopeMapRenderable) Render() {
	gl := renderable.context.OpenGl()

	floorStarts := []int32{0, 6, 12, 18}
	ceilingStarts := []int32{3, 9, 15, 21}

	renderable.vao.OnShader(func() {
		renderable.viewMatrixUniform.Set(gl, renderable.context.ViewMatrix())
		renderable.projectionMatrixUniform.Set(gl, renderable.context.ProjectionMatrix())

		for y, row := range renderable.tiles {
			for x, tile := range row {
				if tile != nil && (*tile.SlopeHeight > 0) {
					modelMatrix := mgl.Ident4().
						Mul4(mgl.Translate3D((float32(x)+0.5)*fineCoordinatesPerTileSide, (float32(y)+0.5)*fineCoordinatesPerTileSide, 0.0)).
						Mul4(mgl.Scale3D(fineCoordinatesPerTileSide, fineCoordinatesPerTileSide, 1.0))
					floorTicks := []int{}
					ceilingTicks := []int{}

					if *tile.SlopeControl != model.SlopeFloorFlat {
						floorTicks = slopeTicksByType[*tile.Type]
					}
					if *tile.SlopeControl == model.SlopeCeilingMirrored {
						ceilingTicks = slopeTicksByType[*tile.Type]
					} else if *tile.SlopeControl != model.SlopeCeilingFlat {
						if invertedType, inversible := invertedTileTypes[*tile.Type]; inversible {
							ceilingTicks = slopeTicksByType[invertedType]
						}
					}

					renderable.modelMatrixUniform.Set(gl, &modelMatrix)
					for _, index := range floorTicks {
						gl.DrawArrays(opengl.TRIANGLES, floorStarts[index], 3)
					}
					for _, index := range ceilingTicks {
						gl.DrawArrays(opengl.TRIANGLES, ceilingStarts[index], 3)
					}
				}
			}
		}
	})
}
