package utils

import (
	"github.com/go-gl/gl/v4.2-core/gl"
	"math"
)

const (
	intSize = 4
)

type Mesh struct {
	vbo, vao uint32
	count    int32
}

func NewMesh(data []float32) Mesh {
	mesh := Mesh{0, 0, int32(len(data) / 9)}

	gl.GenVertexArrays(1, &mesh.vao)
	gl.BindVertexArray(mesh.vao)

	gl.GenBuffers(1, &mesh.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*intSize, gl.Ptr(data), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 9*intSize, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 9*intSize, gl.PtrOffset(4*intSize))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointer(2, 3, gl.FLOAT, false, 9*intSize, gl.PtrOffset(6*intSize))
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

func IdentCircle(segment int) []float32 {
	data := make([]float32, segment*3*9)

	for i := 0; i < segment; i++ {
		var angle float64
		//CENTER
		data[i*3*9+0] = 0
		data[i*3*9+1] = 0
		data[i*3*9+2] = 0
		data[i*3*9+3] = 1

		data[i*3*9+4] = 0
		data[i*3*9+5] = 0
		data[i*3*9+6] = 1

		data[i*3*9+7] = 0.5
		data[i*3*9+8] = 0.5
		//FIRST
		angle = 2.0 * math.Pi / float64(segment) * float64(i)
		data[i*3*9+9+0] = float32(math.Cos(angle))
		data[i*3*9+9+1] = float32(math.Sin(angle))
		data[i*3*9+9+2] = 0
		data[i*3*9+9+3] = 1

		data[i*3*9+9+4] = data[i*3*9+9+1]
		data[i*3*9+9+5] = data[i*3*9+9+0]
		data[i*3*9+9+6] = 0

		data[i*3*9+9+7] = data[i*3*9+9+0]/2 + 0.5
		data[i*3*9+9+8] = data[i*3*9+9+1]/2 + 0.5
		//SECOND
		angle = 2.0 * math.Pi / float64(segment) * float64(i+1)
		data[i*3*9+18+0] = float32(math.Cos(angle))
		data[i*3*9+18+1] = float32(math.Sin(angle))
		data[i*3*9+18+2] = 0
		data[i*3*9+18+3] = 1

		data[i*3*9+18+4] = data[i*3*9+18+0]
		data[i*3*9+18+5] = data[i*3*9+18+1]
		data[i*3*9+18+6] = 0

		data[i*3*9+18+7] = data[i*3*9+18+0]/2 + 0.5
		data[i*3*9+18+8] = data[i*3*9+18+1]/2 + 0.5
	}

	return data
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
