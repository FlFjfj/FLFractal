package utils

import (
	"github.com/go-gl/mathgl/mgl32"
)

type OrthographicCamera struct {
	WIDTH, HEIGHT, ZOOM float32
	position mgl32.Vec2
	projection mgl32.Mat4
}

func (cam *OrthographicCamera) Update() {
	cam.projection = mgl32.Ortho2D(
		cam.position.X() - cam.WIDTH / 2 * cam.ZOOM,
		cam.position.X() + cam.WIDTH / 2 * cam.ZOOM,
		cam.position.Y() - cam.HEIGHT / 2 * cam.ZOOM,
		cam.position.Y() + cam.HEIGHT / 2 * cam.ZOOM)
}

func (cam *OrthographicCamera) SetPosition(pos mgl32.Vec2) {
	cam.position = pos;
}

func (cam *OrthographicCamera) Combined() mgl32.Mat4 {
	return cam.projection
}

func (cam *OrthographicCamera) Translate(trans mgl32.Vec2) {
	cam.position = cam.position.Add(trans.Mul(cam.ZOOM))
}

func NewOrthographicCamera(WIDTH, HEIGHT float32) OrthographicCamera {
	var cam OrthographicCamera = OrthographicCamera{WIDTH, HEIGHT, 1.0, mgl32.Vec2{0,0}, mgl32.Ident4()}
	cam.Update();

	return cam
}
