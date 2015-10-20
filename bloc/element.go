package bloc

import (
	"regexp"
	"strconv"
	"errors"
	)

type LayoutType int
const (
	LayoutAbsolute LayoutType = iota
	LayoutLeft
	LayoutRight
	LayoutDown
	LayoutUp
	)

type Unit int
const (
	Pixels Unit = iota
	ParentWidth
	ParentHeight
	Pixels
	GridWidth
	GridHeight
	)

var valueRegExp

type Value struct {
	Scalar float32
	Unit Unit
}

func matchWithNames(s string) map[string]string {
	match  := valueRegExp.FindStringSubmatch(s)
	result := make(map[string]string)
  for i, name := range myExp.SubexpNames() {
    result[name] = match[i]
  }
  return result
}

func (t *Value) UnmarshalYAML(data []byte) error {
	matches := matchWithNames(string(data))
	switch matches["unit"] {
	case "px":
		t.Unit = Pixels
	case "W":
		t.Unit = ParentWidth
	case "H":
		t.Unit = ParentHeight
	case "gw":
		t.Unit = GridWidth
	case "gh":
		t.Unit = GridHeight
	default:
		return errors.New("Unmarshal Value: Unit type not recognized.")
	}

	v, err := strconv.ParseFloat(matches["value"], 32)
	if err {
		return errors.New("Unmarshal Value: Scalar is not valid floating point number.")
	} else {
		t.Value = v
	}

	return nil
}

func (t *Value) MarshalYAML() ([]byte, error) {
	result := strconv.FormatFloat(float64(t.Value), 'f', -1, 32)
	switch v.Unit {
	case Pixels:
		append(result, "px")
	case ParentWidth:
		append(result, "W")
	case ParentHeight:
		append(result, "H")
	case GridWidth:
		append(result, "gw")
	case GridHeight:
		append(result, "gh")
	default:
		return errors.New("Marshal Value: Unit enum value not valid.")
	}

	return []byte(result), nil
}

func (v *Value) ToPixels(parent *Element) float32 {
	switch v.Unit {
		case Pixels:
			return v.Scalar
		case GridWidth:
			row := int(v.Scalar)
			l := len(parent.GridRows)
			switch {
			case l == 0:
				return 0
			case row >= l:
				return parent.GridRows[l-1].ToPixels(parent.Parent)
			default:
				return parent.GridRows[row].ToPixels(parent.parent)
			}
		case GridHeight:
			col := int(v.Scalar)
			l := len(parent.GridColumns)
			switch {
			case l == 0:
				return 0
			case col >= l:
				return parent.GridColumns[l-1].ToPixels(parent.Parent)
			default:
				return parent.GridColumns[row].ToPixels(parent.parent)
			}
		case ParentWidth:
			return parent.Size[0] * v.Scalar
		case ParentHeight:
			return parent.Size[1] * v.Scalar
	}
}

type Element struct {
	*Bloc
	Layout LayoutType							"yaml:layout,		omitempty"
	Anchor [2]value 							"yaml:anchor,		omitempty" 
	Dimensions [2]value 					"yaml:size,			omitempty"
	Location [2]Value 						"yaml:position,	omitempty"
	GridRows []Value 							"yaml:rows,			omitempty"
	GridColumns []Value 					"yaml:columns,	omitempty"
	Flow bool											"yaml:flow,			omitempty" // TODO Implement flow
	positioningOffset float32
}

func NewElement() *Element {
	result := new(Element)
	result.Bloc = NewBloc()
	result.Layout = LayoutAbsolute
	result.Position[0] = Value{0, ParentWidth}
	result.Position[1] = Value{0, ParentHeight}
	result.Anchor[0] = Value{.5, Width}
	result.Anchor[1] = Value{.5, Height}
	result.SizeField[0] = Value{0, ParentWidth}
	result.SizeField[1] = Value{0, ParentHeight}
	return result
}


func (e *Element) ComputeAbsoluteLayout() {
	parent := interface{}(e.Parent).(*Element)
	e.Size[0] = e.Dimensions[0].ToPixels(parent)
  e.Size[1] = e.Dimensions[1].ToPixels(parent)
  e.Position[0] = e.Location[0].ToPixels(parent) + e.Anchor[0].ToPixels(parent)
  e.Position[1] = e.Location[1].ToPixels(parent) + e.Anchor[1].ToPixels(parent)
}

func (e *Element) ComputeLeftLayout() {
	parent := interface{}(e.Parent).(*Element)
	e.Size[0] = e.Dimensions[0].ToPixels(parent)
  e.Size[1] = e.Dimensions[1].ToPixels(parent)
  e.Position[0] = parent.positioningOffset
  e.Position[1] = e.Location[1].ToPixels(parent) + e.Anchor[1].ToPixels(parent)
  parent.positioningOffset += e.Size[0]
}

func (e *Element) ComputeRightLayout() {
	parent := interface{}(e.Parent).(*Element)
	e.Size[0] = e.Dimensions[0].ToPixels(parent)
  e.Size[1] = e.Dimensions[1].ToPixels(parent)
  e.Position[0] = parent.positioningOffset - e.Size[0]
  e.Position[1] = e.Location[1].ToPixels(parent) + e.Anchor[1].ToPixels(parent)
  parent.positioningOffset -= e.Size[0]
}

func (e *Element) ComputeDownLayout() {
	parent := interface{}(e.Parent).(*Element)
	e.Size[0] = e.Dimensions[0].ToPixels(parent)
  e.Size[1] = e.Dimensions[1].ToPixels(parent)
  e.Position[0] = e.Location[0].ToPixels(parent) + e.Anchor[0].ToPixels(parent)
  e.Position[1] = parent.positioningOffset
  parent.positioningOffset += e.Size[1]
}

func (e *Element) ComputeUpLayout() {
	parent := interface{}(e.Parent).(*Element)
	e.Size[0] = e.Dimensions[0].ToPixels(parent)
  e.Size[1] = e.Dimensions[1].ToPixels(parent)
  e.Position[0] = e.Location[0].ToPixels(parent) + e.Anchor[0].ToPixels(parent)
  e.Position[1] = parent.positioningOffset - e.Size[1]
  parent.positioningOffset -= e.Size[1]
}

func (e *Element) ComputeLayout() {
	switch e.Layout {
		case LayoutLeft:
			e.positioningOffset = e.Position[0]
		case LayoutRight:
			e.positioningOffset = e.Position[0] + e.Size[0]
		case LayoutDown:
			e.positioningOffset = e.Position[1]
		case LayoutUp:
			e.positioningOffset = e.Position[1] + e.Size[1]
	}

	for child := range(e.children) {
		if c, ok = (interface{})(child).(*Element); ok {
			switch e.Layout {
			case LayoutAbsolute:
				c.ComputeAbsoluteLayout()
			case LayoutLeft:
				c.ComputeLeftLayout()
			case LayoutRight:
				c.ComputeRightLayout()
			case LayoutDown:
				c.ComputeDownLayout()
			case LayoutUp:
				c.ComputeUpLayout()
			}
			c.ComputeLayout()
		}
	}
}

func init() {
	valueRegExp := regexp.MustCompile(`(?P<value>[+-]?\d+(\.\d+)?(?P<unit>px|w|h|W|H|gw|gh)`)
}