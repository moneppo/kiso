// bloc is the 2D scenegraph for kiso, abstracted from mode styling, NanoVG rendering,
// tap input management, tween animation system, and various layout services. It serves as
// a mixin definition for blocks in a screen.

package bloc

import (
	"github.com/kiso/math"
)

type Bloc struct {
	Visible bool
	Size math.Size
	Parent *Bloc

	children []*Bloc
	transform math.Mat3
	invTransform math.Mat3
	transformDirty bool
	position math.Vec2
	rotation float32
	scale math.Vec2
}

func NewBloc() *Bloc {
	result := new(Bloc)
	result.visible = true
	result.scale = math.Vec2{1,1}
	return result
}

func (b *Bloc) SetAffineValues(position math.Vec2, rotation float32, scale math.Vec2) {
	b.position = position
	b.rotation = rotation
	b.scale = scale

	b.transformDirty = true
}

func (b *Bloc) Transform() math.Mat3 {
	if (b.transformDirty) {
		translate := TranslateFromVec2(b.position)
		rotate := Rotate(b.rotation)
		scale := ScaleFromVec2(b.scale)

		b.transform = translate.Mul(scale.Mul(rotate))
		b.invTransform = rotate.Mul(scale.Mul(translate))
	}

	return b.transform
}

func (b *Bloc) InverseTransform() math.Mat3 {
	if (b.transformDirty) {
		translate := TranslateFromVec2(b.position)
		rotate := Rotate(b.rotation)
		scale := ScaleFromVec2(b.scale)

		b.transform = translate.Mul(scale.Mul(rotate))
		b.invTransform = rotate.Mul(scale.Mul(translate))
	}

	return b.invTransform
}

func (b *Bloc) Children() []*Bloc {
	var result []*Bloc
	copy(result, b.children)
	return result
}

func (b *Bloc) AddChild(c *Bloc) {
	b.children = append(b.children, c)
}

func (b *Bloc) AddChildAt(index int, c *Bloc) {
	b.children = append(b.children, nil)
	copy(b.children[index+1:], b.children[index:])
	b.children[index] = c
}

func (b *Bloc) RemoveChild(c *Bloc) {
	for i, x := range b.children {
		if (x == c) {
			b.RemoveChildAt(i)
			return
		}
	}
}

func (b *Bloc) HitTest(globalPoint math.Vec2) bool {
	p := b.GlobalToLocal(globalPoint)
	return p.X <= b.Size.W && p.Y <= b.Size.H 
}

func (b *Bloc) RemoveChildAt(index int) {
	b.children = append(b.children[:index], b.children[index+1:]...)
}

func (b *Bloc) LocalToGlobal(p math.Vec2) math.Vec2 {
	return p.Mul(b.Transform())
}

func (b *Bloc) GlobalToLocal(p math.Vec2) math.Vec2 {
	return p.Mul(b.InverseTransform())
}

func (b *Bloc) LocalToLocal(p math.Vec2, target *Bloc) math.Vec2 {
	pt := b.LocalToGlobal(p)
	return target.GlobalToLocal(pt)
}