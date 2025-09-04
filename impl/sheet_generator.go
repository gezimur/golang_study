package sheet_generator

import (
    "fmt"
    "math/rand"
    "image"
    "image/color"
)

type SheetObject interface {
    IsValid() bool

    View() color.RGBA
    IsCrossable() bool
}

type Sheet [T SheetObject] interface {
    Width() uint
    SetWidth(uint)
    
    Height() uint
    SetHeight(uint)
    
    SetSize(uint, uint)
    
    ObjectAt(image.Point) T
    SetObjectAt(image.Point, T)
}

func find_empty_points[T_Object SheetObject, T_Sheet Sheet[T_Object]](sheet T_Sheet, current_point image.Point) []image.Point {
    turns := [8][2]int{
	{-1, -1},
	{-1,  0},
	{-1,  1},
	{ 0, -1},
	{ 0,  1},
	{ 1, -1},
	{ 1,  0},
	{ 1,  1},
    }
    
    empty_points := make([]image.Point, 0)
    for i:= 0; i < len(turns); i++ {
        next_point := image.Point {
	    X : current_point.X + turns[i][0],
	    Y : current_point.Y + turns[i][1],
	}
	if int(sheet.Width()) > next_point.X && 
	    int(sheet.Height()) > next_point.Y && 
	    0 <= next_point.X && 
	    0 <= next_point.Y && 
	    !sheet.ObjectAt(next_point).IsValid() {
	    
	    empty_points = append(empty_points, next_point)
	}
    }
    return empty_points
}

func has_point(point_slice []image.Point, point image.Point) bool {
    for i := 0; i < len(point_slice); i++ {
	if point_slice[i] == point {
	    return true
	}
    }
    return false
}

func remove_copy(point_slice []image.Point) []image.Point{
    res := make([]image.Point, 0)
    for i := 0; i < len(point_slice); i++ {
	if !has_point(res, point_slice[i]) {
	    res = append(res, point_slice[i])
	}
    }
    return res
}

func fill_step_bfs[T_Object SheetObject, T_Sheet Sheet[T_Object]](sheet T_Sheet, point_to_visit []image.Point, object_maker func()T_Object) []image.Point {
    next_point_to_visit := make([]image.Point, 0)
    for _, point := range point_to_visit {
	empty_points := find_empty_points(sheet, point)
	for _, next_point := range empty_points {
	    sheet.SetObjectAt(next_point, object_maker())
	}
	next_point_to_visit = append(next_point_to_visit, empty_points...)
    }
    return remove_copy(next_point_to_visit)
}

func fill_step_dfs[T_Object SheetObject, T_Sheet Sheet[T_Object]](sheet T_Sheet, traveled_path []image.Point, object_maker func()T_Object) []image.Point {
    // todo make optimize
    last_point_id := len(traveled_path) - 1
    
    for i := last_point_id; i >= 0; i-- {
	empty_points := find_empty_points(sheet, traveled_path[i])

	if len(empty_points) > 0 {
	    random_point := empty_points[rand.Intn(len(empty_points))]
	    sheet.SetObjectAt(random_point, object_maker())
	    return remove_copy(append(traveled_path[: i + 1], random_point)) 
	}
    }
    
    return []image.Point{}
}

type sheet_filler_struct [T_Object SheetObject, T_Sheet Sheet[T_Object]] struct {
    arg_points []image.Point
    object_maker func()T_Object
    fill_step func(T_Sheet, []image.Point, func()T_Object) []image.Point
}

func fill[T_Object SheetObject, T_Sheet Sheet[T_Object]](filler_ptr *sheet_filler_struct[T_Object, T_Sheet], sheet T_Sheet) bool{
    filler_ptr.arg_points = filler_ptr.fill_step(sheet, filler_ptr.arg_points, filler_ptr.object_maker)
    return len(filler_ptr.arg_points) > 0
}

func FillSheet [T_Object SheetObject, T_Sheet Sheet[T_Object]] (sheet T_Sheet, object_maker func(uint)T_Object) T_Sheet{
    points_cnt := 8
    sheet_fillers := make([]*sheet_filler_struct[T_Object, T_Sheet], points_cnt)
    
    for i, _ := range sheet_fillers {
	random_point := image.Point{rand.Intn(int(sheet.Width())), rand.Intn(int(sheet.Height()))}
	special_object_maker := func () T_Object {
	    return object_maker(uint(i & 7))
	}
	var special_fill_step func(T_Sheet, []image.Point, func()T_Object) []image.Point
	if rand.Intn(5) > 2 {
	    special_fill_step = fill_step_dfs
	} else {
	    special_fill_step = fill_step_bfs
	}
	
	fmt.Printf("start point: %v\n", random_point)
	
	sheet.SetObjectAt(random_point, special_object_maker())
	sheet_fillers[i] = &sheet_filler_struct[T_Object, T_Sheet]{
	    arg_points : []image.Point{random_point},
	    object_maker : special_object_maker,
	    fill_step : special_fill_step,
	}
    }
    
    has_not_filled := true
    for ; has_not_filled; {
	has_not_filled = false
	for _, filler_ptr := range sheet_fillers {
	    do_smth := fill(filler_ptr, sheet)
	    has_not_filled = has_not_filled || do_smth
	}
    }
    
    return sheet
}
