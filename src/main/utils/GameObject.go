package utils

import (
	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type GameObject struct {
	mesh      *Mesh
	size      float32
	transform mgl32.Mat4
	objLoc    int32
}

func NewObject(mesh *Mesh, size float32, transform_loc int32) GameObject {
	obj := GameObject{mesh, 0, mgl32.Ident4(), transform_loc}
	obj.Update(mgl32.Vec3{0.0, 0.0, 0.0}, size)
	return obj
}

func (object *GameObject) Draw() {
	gl.UniformMatrix4fv(object.objLoc, 1, false, &object.transform[0])
	object.mesh.Draw()
}

func (object *GameObject) Update(pos mgl32.Vec3, size float32) {
	object.size = size
	object.transform = mgl32.Translate3D(pos.X(), pos.Y(), 0).Mul4(mgl32.Scale3D(size, size, 1))
}
