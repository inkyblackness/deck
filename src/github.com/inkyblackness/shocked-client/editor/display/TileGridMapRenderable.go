package display

import (
	"fmt"
	"os"

	"github.com/inkyblackness/shocked-model"

	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/opengl"
)

var mapTileGridVertexShaderSource = `
  attribute vec3 vertexPosition;

  uniform mat4 viewMatrix;
  uniform mat4 projectionMatrix;

  varying float height;

  void main(void) {
    gl_Position = projectionMatrix * viewMatrix * vec4(vertexPosition.xy, 0.0, 1.0);
    height = vertexPosition.z;
  }
`

var mapTileGridFragmentShaderSource = `
  #ifdef GL_ES
    precision mediump float;
  #endif

  varying float height;

  void main(void) {
    gl_FragColor = vec4(0.0, 0.8, 0.0, height);
  }
`

// TileGridMapRenderable is a renderable for the tile grid.
type TileGridMapRenderable struct {
	context *graphics.RenderContext

	program                 uint32
	vao                     *opengl.VertexArrayObject
	vertexPositionBuffer    uint32
	vertexPositionAttrib    int32
	viewMatrixUniform       opengl.Matrix4Uniform
	projectionMatrixUniform opengl.Matrix4Uniform

	tiles [][]*model.TileProperties
}

// NewTileGridMapRenderable returns a new instance of a renderable for tile grids.
func NewTileGridMapRenderable(context *graphics.RenderContext) *TileGridMapRenderable {
	gl := context.OpenGl()
	vertexShader, err1 := opengl.CompileNewShader(gl, opengl.VERTEX_SHADER, mapTileGridVertexShaderSource)
	defer gl.DeleteShader(vertexShader)
	fragmentShader, err2 := opengl.CompileNewShader(gl, opengl.FRAGMENT_SHADER, mapTileGridFragmentShaderSource)
	defer gl.DeleteShader(fragmentShader)
	program, _ := opengl.LinkNewProgram(gl, vertexShader, fragmentShader)

	if err1 != nil {
		fmt.Fprintf(os.Stderr, "Failed to compile shader 1:\n", err1)
	}
	if err2 != nil {
		fmt.Fprintf(os.Stderr, "Failed to compile shader 2:\n", err2)
	}

	renderable := &TileGridMapRenderable{
		context:                 context,
		program:                 program,
		vao:                     opengl.NewVertexArrayObject(gl, program),
		vertexPositionBuffer:    gl.GenBuffers(1)[0],
		vertexPositionAttrib:    gl.GetAttribLocation(program, "vertexPosition"),
		viewMatrixUniform:       opengl.Matrix4Uniform(gl.GetUniformLocation(program, "viewMatrix")),
		projectionMatrixUniform: opengl.Matrix4Uniform(gl.GetUniformLocation(program, "projectionMatrix")),

		tiles: make([][]*model.TileProperties, int(tilesPerMapSide))}

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
func (renderable *TileGridMapRenderable) Dispose() {
	gl := renderable.context.OpenGl()
	gl.DeleteProgram(renderable.program)
	gl.DeleteBuffers([]uint32{renderable.vertexPositionBuffer})
	renderable.vao.Dispose()
}

// SetTile sets the properties for the specified tile coordinate.
func (renderable *TileGridMapRenderable) SetTile(x, y int, properties *model.TileProperties) {
	renderable.tiles[y][x] = properties
}

// Clear resets all tiles.
func (renderable *TileGridMapRenderable) Clear() {
	for _, row := range renderable.tiles {
		for index := 0; index < len(row); index++ {
			row[index] = nil
		}
	}
}

// Render renders
func (renderable *TileGridMapRenderable) Render() {
	gl := renderable.context.OpenGl()

	renderable.vao.OnShader(func() {
		renderable.viewMatrixUniform.Set(gl, renderable.context.ViewMatrix())
		renderable.projectionMatrixUniform.Set(gl, renderable.context.ProjectionMatrix())

		gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexPositionBuffer)
		for y, row := range renderable.tiles {
			for x, tile := range row {
				if tile != nil {
					left := float32(x) * fineCoordinatesPerTileSide
					right := left + fineCoordinatesPerTileSide
					bottom := float32(y) * fineCoordinatesPerTileSide
					top := bottom + fineCoordinatesPerTileSide

					vertices := make([]float32, 0, 6*2*3)

					if tile.CalculatedWallHeights.North > 0 {
						vertices = append(vertices, left, top, tile.CalculatedWallHeights.North, right, top, tile.CalculatedWallHeights.North)
					}
					if tile.CalculatedWallHeights.South > 0 {
						vertices = append(vertices, left, bottom, tile.CalculatedWallHeights.South, right, bottom, tile.CalculatedWallHeights.South)
					}
					if tile.CalculatedWallHeights.West > 0 {
						vertices = append(vertices, left, top, tile.CalculatedWallHeights.West, left, bottom, tile.CalculatedWallHeights.West)
					}
					if tile.CalculatedWallHeights.East > 0 {
						vertices = append(vertices, right, top, tile.CalculatedWallHeights.East, right, bottom, tile.CalculatedWallHeights.East)
					}
					if *tile.Type == model.DiagonalOpenNorthEast || *tile.Type == model.DiagonalOpenSouthWest {
						vertices = append(vertices, left, top, 1.0, right, bottom, 1.0)
					}
					if *tile.Type == model.DiagonalOpenNorthWest || *tile.Type == model.DiagonalOpenSouthEast {
						vertices = append(vertices, left, bottom, 1.0, right, top, 1.0)
					}

					if len(vertices) > 0 {
						gl.BufferData(opengl.ARRAY_BUFFER, len(vertices)*4, vertices, opengl.STATIC_DRAW)
						gl.DrawArrays(opengl.LINES, 0, int32(len(vertices)/3))
					}
				}
			}
		}
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	})
}
