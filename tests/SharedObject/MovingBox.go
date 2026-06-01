package sharedobject

import (
	"autoworld/modules/napi"
	"fmt"
)

type MovingBox struct {
	napi.IObject
	napi.Pos
	napi.Box
	napi.Velo
	napi.Spr

	napi.Col
	napi.Info
}

func NewMovingBox(x, y int) *MovingBox {
	o := &MovingBox{}
	napi.Obj.NewObject(o, "moving-box", "pos box velo spr sce-main col")

	o.SetX(float32(x))
	o.SetY(float32(y))
	return o
}

func (this *MovingBox) Create() {
	this.AddSprite("normal", napi.Assert.GetSprite("character"))
	this.SetCurrentSprite("normal")

	this.AddTag("ball")
	this.OnCollisionTag("ball", this.Handler)

	this.SetVelocityX(1)
	this.SetVelocityY(1)

	this.SetScaleX(3)
	this.SetScaleY(3)
	this.SetImageSpeed(0.25)

	this.SetBoxH(float32(this.GetCurrentSprite().Height() * 3))
	this.SetBoxW(float32(this.GetCurrentSprite().Width() * 3))

	fmt.Print(this.BoxH(), this.BoxW())
}

func (this *MovingBox) StepUpdate() {
	_vx, _vy := this.VelocityX(), this.VelocityY()
	_w, _h := napi.Store.Int("game-width"), napi.Store.Int("game-height")
	_x, _y := this.X(), this.Y()

	if _x > float32(_w) || _x < 0 {
		_vx = -_vx
	}
	if _y > float32(_h) || _y < 0 {
		_vy = -_vy
	}
	this.SetVelocityX(_vx)
	this.SetVelocityY(_vy)
}

func (this *MovingBox) Destroy() {}

func (this *MovingBox) Handler(other napi.Object) {
	fmt.Print(".")
	this.SetVelocityX(-this.VelocityX())
	this.SetVelocityY(-this.VelocityY())
}
