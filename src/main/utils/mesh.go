package utils

import (
	"github.com/go-gl/gl/v4.2-core/gl"
	//"github.com/go-gl/mathgl/mgl32"
)

const (
	intSize = 4
)

type Mesh struct {
	vbo, vao uint32
	count int32
}

func NewMesh(data []float32)  Mesh{
	mesh := Mesh{0, 0, int32(len(data) / 9)}

	gl.GenVertexArrays(1, &mesh.vao)
	gl.BindVertexArray(mesh.vao)

	gl.GenBuffers(1, &mesh.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(data) * intSize, gl.Ptr(data), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 9 * intSize, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 9 * intSize, gl.PtrOffset(4 * intSize))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointer(2, 3, gl.FLOAT, false, 9*intSize, gl.PtrOffset(6 * intSize))
	gl.EnableVertexAttribArray(2)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	return mesh
}

func (mesh *Mesh) Draw() {
	gl.BindVertexArray(mesh.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, mesh.count)
	gl.BindVertexArray(0)
}

func GenCube(size float32) []float32 {
	a := size / 2
	res := []float32{
		-a, -a, -a, 1, 0.0, 0.0, 0.0, -1.0, 0.0,
		a, -a, -a, 1, 1.0, 0.0, 0.0, -1.0, 0.0,
		-a, -a, a, 1, 0.0, 1.0, 0.0, -1.0, 0.0,
		a, -a, -a, 1, 1.0, 0.0, 0.0, -1.0, 0.0,
		a, -a, a, 1, 1.0, 1.0, 0.0, -1.0, 0.0,
		-a, -a, a, 1, 0.0, 1.0, 0.0, -1.0, 0.0,

		// Top
		-a, a, -a, 1, 0.0, 0.0, 0.0, 1.0, 0.0,
		-a, a, a, 1, 0.0, 1.0, 0.0, 1.0, 0.0,
		a, a, -a, 1, 1.0, 0.0, 0.0, 1.0, 0.0,
		a, a, -a, 1, 1.0, 0.0, 0.0, 1.0, 0.0,
		-a, a, a, 1, 0.0, 1.0, 0.0, 1.0, 0.0,
		a, a, a, 1, 1.0, 1.0, 0.0, 1.0, 0.0,

		// Front
		-a, -a, a, 1, 0.0, 0.0, 0.0, 0.0, 1.0,
		a, -a, a, 1, 1.0, 0.0, 0.0, 0.0, 1.0,
		-a, a, a, 1, 0.0, 1.0, 0.0, 0.0, 1.0,
		a, -a, a, 1, 1.0, 0.0, 0.0, 0.0, 1.0,
		a, a, a, 1, 1.0, 1.0, 0.0, 0.0, 1.0,
		-a, a, a, 1, 0.0, 1.0, 0.0, 0.0, 1.0,

		// Back
		-a, -a, -a, 1, 0.0, 0.0, 0.0, 0.0, -1.0,
		-a, a, -a, 1, 0.0, 1.0, 0.0, 0.0, -1.0,
		a, -a, -a, 1, 1.0, 0.0, 0.0, 0.0, -1.0,
		a, -a, -a, 1, 1.0, 0.0, 0.0, 0.0, -1.0,
		-a, a, -a, 1, 0.0, 1.0, 0.0, 0.0, -1.0,
		a, a, -a, 1, 1.0, 1.0, 0.0, 0.0, -1.0,

		// Left
		-a, -a, a, 1, 0.0, 1.0, -1.0, 0.0, 0.0,
		-a, a, -a, 1, 1.0, 0.0, -1.0, 0.0, 0.0,
		-a, -a, -a, 1, 0.0, 0.0, -1.0, 0.0, 0.0,
		-a, -a, a, 1, 0.0, 1.0, -1.0, 0.0, 0.0,
		-a, a, a, 1, 1.0, 1.0, -1.0, 0.0, 0.0,
		-a, a, -a, 1, 1.0, 0.0, -1.0, 0.0, 0.0,

		// Right
		a, -a, a, 1, 0.0, 1.0, 1.0, 0.0, 0.0,
		a, -a, -a, 1, 0.0, 0.0, 1.0, 0.0, 0.0,
		a, a, -a, 1, 1.0, 0.0, 1.0, 0.0, 0.0,
		a, -a, a, 1, 0.0, 1.0, 1.0, 0.0, 0.0,
		a, a, -a, 1, 1.0, 0.0, 1.0, 0.0, 0.0,
		a, a, a, 1, 1.0, 1.0, 1.0, 0.0, 0.0,
	}

	return res
}
