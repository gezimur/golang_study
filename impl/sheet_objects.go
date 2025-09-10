package sheet_generator

import (
    "fmt"
    "image"
    "image/color"
)

//type SheetObject interface {
//    IsValid() bool
//
//    View() color.RGBA
//    IsCrossable() bool
//}
//
//type Sheet interface {
//    Width() uint
//    SetWidth(uint)
//    
//    Height() uint
//    SetHeight(uint)
//    
//    SetSize(uint, uint)
//    
//    ObjectAt(image.Point) T
//    SetObjectAt(image.Point, T)
//}

type ColorfullObject struct {
    color color.RGBA
    is_crossable bool
}

func MakeColorfullObject(a1 color.RGBA, a2 bool) ColorfullObject {
    return ColorfullObject{
	color : a1,
	is_crossable : a2,
    }
}

func (obj *ColorfullObject) IsValid() bool {
    return obj.color.A > 0
}

func (obj *ColorfullObject) View() color.RGBA {
    return obj.color
}

func (obj *ColorfullObject) IsCrossable() bool {
    return obj.is_crossable
}

type ColorfullSheet struct {
    sheet [][]*ColorfullObject
}

func (sheet *ColorfullSheet) Width() uint {
    return uint(len(sheet.sheet))
}

func (sheet *ColorfullSheet) SetWidth(width uint) {
    sheet.SetSize(width, 1)
}

func (sheet *ColorfullSheet) Height() uint {
    if sheet.Width() == 0 {
	return 0
    }
    return uint(len(sheet.sheet[0]))
}

func (sheet *ColorfullSheet) SetHeight(height uint) {
    sheet.SetSize(1, height)
}

func (sheet *ColorfullSheet) SetSize(width uint, height uint) {
    if sheet.Height() != 0 {
	panic("I don't know that to do!") // todo just do it
    }
    
    sheet.sheet = make([][]*ColorfullObject, width)
    for x := uint(0); x < width; x++ {
	sheet.sheet[x] = make([]*ColorfullObject, height)
	for y := uint(0); y < height; y++ {
	    sheet.sheet[x][y] = &ColorfullObject {
		color : color.RGBA{0,0,0,0},
		is_crossable : false,
	    }
	}
    }
}

func (sheet *ColorfullSheet) ObjectAt(pos image.Point) *ColorfullObject {
    if pos.X >= int(sheet.Width()) || pos.Y >= int(sheet.Height()) {
	panic(fmt.Sprintf("Unreachable position passed, (%d;%d)", pos.X, pos.Y))
    }
    
    return sheet.sheet[pos.X][pos.Y]
}

func (sheet *ColorfullSheet) SetObjectAt(pos image.Point, obj *ColorfullObject) {
    if pos.X >= int(sheet.Width()) || pos.Y >= int(sheet.Height()) {
	panic(fmt.Sprintf("Unreachable position passed, (%d;%d)", pos.X, pos.Y))
    }
    
    sheet.sheet[pos.X][pos.Y] = obj
}

func (sheet *ColorfullSheet) ColorModel() color.Model {
    return color.RGBAModel
}

func (sheet *ColorfullSheet) Bounds() image.Rectangle {
    return image.Rect(0, 0, int(sheet.Width()), int(sheet.Height()))
}

func (sheet *ColorfullSheet) At(x, y int) color.Color {
    return sheet.ObjectAt(image.Point{x, y}).View()
}
