package display

import (
	"fmt"
	"os"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/shocked-model"

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

// TileGridMapRenderable is a renderable for textures.
type TileGridMapRenderable struct {
	gl opengl.OpenGl

	program                 uint32
	vertexArrayObject       uint32
	vertexPositionBuffer    uint32
	vertexPositionAttrib    int32
	viewMatrixUniform       int32
	projectionMatrixUniform int32

	tiles [][]*model.TileProperties
}

// NewTileGridMapRenderable returns a new instance of a renderable for tile grid maps
func NewTileGridMapRenderable(gl opengl.OpenGl) *TileGridMapRenderable {
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
		gl:                      gl,
		program:                 program,
		vertexArrayObject:       gl.GenVertexArrays(1)[0],
		vertexPositionBuffer:    gl.GenBuffers(1)[0],
		vertexPositionAttrib:    gl.GetAttribLocation(program, "vertexPosition"),
		viewMatrixUniform:       gl.GetUniformLocation(program, "viewMatrix"),
		projectionMatrixUniform: gl.GetUniformLocation(program, "projectionMatrix"),
		tiles: make([][]*model.TileProperties, 64)}

	for i := 0; i < 64; i++ {
		renderable.tiles[i] = make([]*model.TileProperties, 64)
	}

	return renderable
}

// Dispose releases any internal resources
func (renderable *TileGridMapRenderable) Dispose() {
	renderable.gl.DeleteProgram(renderable.program)
	renderable.gl.DeleteBuffers([]uint32{renderable.vertexPositionBuffer})
	renderable.gl.DeleteVertexArrays([]uint32{renderable.vertexArrayObject})
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
func (renderable *TileGridMapRenderable) Render(context *RenderContext) {
	gl := renderable.gl

	renderable.withShader(func() {
		renderable.setMatrix32(renderable.viewMatrixUniform, context.ViewMatrix())
		renderable.setMatrix32(renderable.projectionMatrixUniform, context.ProjectionMatrix())

		gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexPositionBuffer)
		gl.VertexAttribOffset(uint32(renderable.vertexPositionAttrib), 3, opengl.FLOAT, false, 0, 0)

		for y, row := range renderable.tiles {
			for x, tile := range row {
				if tile != nil {
					left := float32(x) * 32.0
					right := left + 32.0
					top := float32(y) * 32.0
					bottom := top + 32.0

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
	})
}

func (renderable *TileGridMapRenderable) withShader(task func()) {
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

func (renderable *TileGridMapRenderable) setMatrix32(uniform int32, matrix *mgl.Mat4) {
	matrixArray := ([16]float32)(*matrix)
	renderable.gl.UniformMatrix4fv(uniform, false, &matrixArray)
}
