package browser

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
)

// WebGl is a wrapper for the WebGL related operations.
type WebGl struct {
	gl *webgl.Context

	buffers           ObjectMapper
	programs          ObjectMapper
	shaders           ObjectMapper
	textures          ObjectMapper
	uniforms          ObjectMapper
	uniformsByProgram map[uint32][]uint32
}

// NewWebGl returns a new instance of WebGl, wrapping the provided context.
func NewWebGl(gl *webgl.Context) *WebGl {
	result := &WebGl{
		gl:                gl,
		buffers:           NewObjectMapper(),
		programs:          NewObjectMapper(),
		shaders:           NewObjectMapper(),
		textures:          NewObjectMapper(),
		uniforms:          NewObjectMapper(),
		uniformsByProgram: make(map[uint32][]uint32)}

	return result
}

// ActiveTexture implements the opengl.OpenGl interface.
func (web *WebGl) ActiveTexture(texture uint32) {
	web.gl.ActiveTexture(int(texture))
}

// AttachShader implements the opengl.OpenGl interface.
func (web *WebGl) AttachShader(program uint32, shader uint32) {
	objShader := web.shaders.get(shader)
	objProgram := web.programs.get(program)

	web.gl.AttachShader(objProgram, objShader)
}

// BindAttribLocation implements the opengl.OpenGl interface.
func (web *WebGl) BindAttribLocation(program uint32, index uint32, name string) {
	web.gl.BindAttribLocation(web.programs.get(program), int(index), name)
}

// BindBuffer implements the opengl.OpenGl interface.
func (web *WebGl) BindBuffer(target uint32, buffer uint32) {
	web.gl.BindBuffer(int(target), web.buffers.get(buffer))
}

// BindTexture implements the opengl.OpenGl interface.
func (web *WebGl) BindTexture(target uint32, texture uint32) {
	web.gl.BindTexture(int(target), web.textures.get(texture))
}

// BindVertexArray implements the opengl.OpenGl interface.
func (web *WebGl) BindVertexArray(array uint32) {
	// not supported in WebGL, can be ignored
}

// BlendFunc implements the OpenGl interface.
func (web *WebGl) BlendFunc(sfactor uint32, dfactor uint32) {
	web.gl.BlendFunc(int(sfactor), int(dfactor))
}

// BufferData implements the opengl.OpenGl interface.
func (web *WebGl) BufferData(target uint32, size int, data interface{}, usage uint32) {
	web.gl.BufferData(int(target), data, int(usage))
}

// Clear implements the opengl.OpenGl interface.
func (web *WebGl) Clear(mask uint32) {
	web.gl.Clear(int(mask))
}

// ClearColor implements the opengl.OpenGl interface.
func (web *WebGl) ClearColor(red float32, green float32, blue float32, alpha float32) {
	web.gl.ClearColor(red, green, blue, alpha)
}

// CompileShader implements the opengl.OpenGl interface.
func (web *WebGl) CompileShader(shader uint32) {
	web.gl.CompileShader(web.shaders.get(shader))
}

// CreateProgram implements the opengl.OpenGl interface.
func (web *WebGl) CreateProgram() uint32 {
	key := web.programs.put(web.gl.CreateProgram())
	web.uniformsByProgram[key] = make([]uint32, 0)

	return key
}

// CreateShader implements the opengl.OpenGl interface.
func (web *WebGl) CreateShader(shaderType uint32) uint32 {
	return web.shaders.put(web.gl.CreateShader(int(shaderType)))
}

// DeleteBuffers implements the opengl.OpenGl interface.
func (web *WebGl) DeleteBuffers(buffers []uint32) {
	for _, buffer := range buffers {
		web.gl.DeleteBuffer(web.buffers.del(buffer))
	}
}

// DeleteProgram implements the opengl.OpenGl interface.
func (web *WebGl) DeleteProgram(program uint32) {
	web.gl.DeleteProgram(web.programs.del(program))
	for _, value := range web.uniformsByProgram[program] {
		web.uniforms.del(value)
	}
	delete(web.uniformsByProgram, program)
}

// DeleteShader implements the opengl.OpenGl interface.
func (web *WebGl) DeleteShader(shader uint32) {
	web.gl.DeleteShader(web.shaders.del(shader))
}

// DeleteTextures implements the opengl.OpenGl interface.
func (web *WebGl) DeleteTextures(textures []uint32) {
	for _, texture := range textures {
		web.gl.DeleteTexture(web.textures.del(texture))
	}
}

// DeleteVertexArrays implements the opengl.OpenGl interface.
func (web *WebGl) DeleteVertexArrays(arrays []uint32) {
	// Not supported in WebGL, can be ignored
}

// DrawArrays implements the opengl.OpenGl interface.
func (web *WebGl) DrawArrays(mode uint32, first int32, count int32) {
	web.gl.DrawArrays(int(mode), int(first), int(count))
}

// Enable implements the opengl.OpenGl interface.
func (web *WebGl) Enable(cap uint32) {
	web.gl.Enable(int(cap))
}

// EnableVertexAttribArray implements the opengl.OpenGl interface.
func (web *WebGl) EnableVertexAttribArray(index uint32) {
	web.gl.EnableVertexAttribArray(int(index))
}

// GenerateMipmap implements the opengl.OpenGl interface.
func (web *WebGl) GenerateMipmap(target uint32) {
	web.gl.GenerateMipmap(int(target))
}

// GenBuffers implements the opengl.OpenGl interface.
func (web *WebGl) GenBuffers(n int32) []uint32 {
	ids := make([]uint32, n)

	for i := int32(0); i < n; i++ {
		ids[i] = web.buffers.put(web.gl.CreateBuffer())
	}

	return ids
}

// GenTextures implements the opengl.OpenGl interface.
func (web *WebGl) GenTextures(n int32) []uint32 {
	ids := make([]uint32, n)

	for i := int32(0); i < n; i++ {
		ids[i] = web.textures.put(web.gl.CreateTexture())
	}

	return ids
}

// GenVertexArrays implements the opengl.OpenGl interface.
func (web *WebGl) GenVertexArrays(n int32) []uint32 {
	// Not supported in WebGL, can be ignored; Creating dummy IDs
	ids := []uint32{}
	for i := int32(0); i < n; i++ {
		ids = append(ids, uint32(i+1))
	}
	return ids
}

// GetAttribLocation implements the opengl.OpenGl interface.
func (web *WebGl) GetAttribLocation(program uint32, name string) int32 {
	return int32(web.gl.GetAttribLocation(web.programs.get(program), name))
}

// GetError implements the opengl.OpenGl interface.
func (web *WebGl) GetError() uint32 {
	return uint32(web.gl.GetError())
}

// GetShaderInfoLog implements the opengl.OpenGl interface.
func (web *WebGl) GetShaderInfoLog(shader uint32) string {
	return web.gl.GetShaderInfoLog(web.shaders.get(shader))
}

func paramToInt(value *js.Object) int32 {
	result := int32(value.Int())

	if value.String() == "true" {
		result = 1
	}

	return result
}

// GetShaderParameter implements the opengl.OpenGl interface.
func (web *WebGl) GetShaderParameter(shader uint32, param uint32) int32 {
	value := web.gl.GetShaderParameter(web.shaders.get(shader), int(param))

	return paramToInt(value)
}

// GetProgramInfoLog implements the opengl.OpenGl interface.
func (web *WebGl) GetProgramInfoLog(program uint32) string {
	return web.gl.GetProgramInfoLog(web.programs.get(program))
}

// GetProgramParameter implements the opengl.OpenGl interface.
func (web *WebGl) GetProgramParameter(program uint32, param uint32) int32 {
	// Call function directly since the wrapping function does not cover properly convert strings.
	value := web.gl.Call("getProgramParameter", web.programs.get(program), int(param))

	return paramToInt(value)
}

// GetUniformLocation implements the opengl.OpenGl interface.
func (web *WebGl) GetUniformLocation(program uint32, name string) int32 {
	uniform := web.gl.GetUniformLocation(web.programs.get(program), name)
	key := web.uniforms.put(uniform)

	web.uniformsByProgram[program] = append(web.uniformsByProgram[program], key)

	return int32(key)
}

// LinkProgram implements the opengl.OpenGl interface.
func (web *WebGl) LinkProgram(program uint32) {
	web.gl.LinkProgram(web.programs.get(program))
}

// ReadPixels implements the opengl.OpenGl interface.
func (web *WebGl) ReadPixels(x int32, y int32, width int32, height int32, format uint32, pixelType uint32, pixels interface{}) {
	//web.gl.ReadPixels(int(x), int(y), int(width), int(height), int(format), int(pixelType), pixels)
	//web.gl.Call("readPixels", x, y, width, height, int(format), int(pixelType), pixels)
}

// ShaderSource implements the opengl.OpenGl interface.
func (web *WebGl) ShaderSource(shader uint32, source string) {
	web.gl.ShaderSource(web.shaders.get(shader), source)
}

// TexImage2D implements the opengl.OpenGl interface.
func (web *WebGl) TexImage2D(target uint32, level int32, internalFormat uint32, width int32, height int32,
	border int32, format uint32, xtype uint32, pixels interface{}) {
	web.gl.Call("texImage2D", int(target), int(level), int(internalFormat), int(width), int(height),
		int(border), int(format), int(xtype), pixels)
}

// TexParameteri implements the opengl.OpenGl interface.
func (web *WebGl) TexParameteri(target uint32, pname uint32, param int32) {
	web.gl.TexParameteri(int(target), int(pname), int(param))
}

// Uniform1i implements the opengl.OpenGl interface.
func (web *WebGl) Uniform1i(location int32, value int32) {
	web.gl.Uniform1i(web.uniforms.get(uint32(location)), int(value))
}

// Uniform4fv implements the opengl.OpenGl interface.
func (web *WebGl) Uniform4fv(location int32, value *[4]float32) {
	web.gl.Uniform4f(web.uniforms.get(uint32(location)), value[0], value[1], value[2], value[3])
}

// UniformMatrix4fv implements the opengl.OpenGl interface.
func (web *WebGl) UniformMatrix4fv(location int32, transpose bool, value *[16]float32) {
	web.gl.UniformMatrix4fv(web.uniforms.get(uint32(location)), transpose, (*value)[:])
	//web.gl.Call("uniformMatrix4fv", web.uniforms.get(uint32(location)), transpose, value)
}

// UseProgram implements the opengl.OpenGl interface.
func (web *WebGl) UseProgram(program uint32) {
	web.gl.UseProgram(web.programs.get(program))
}

// VertexAttribOffset implements the opengl.OpenGl interface.
func (web *WebGl) VertexAttribOffset(index uint32, size int32, attribType uint32, normalized bool, stride int32, offset int) {
	web.gl.VertexAttribPointer(int(index), int(size), int(attribType), normalized, int(stride), offset)
}

// Viewport implements the opengl.OpenGl interface.
func (web *WebGl) Viewport(x int32, y int32, width int32, height int32) {
	web.gl.Viewport(int(x), int(y), int(width), int(height))
}
