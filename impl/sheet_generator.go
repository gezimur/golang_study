package sheet_generator

import (
    "image"
    "math/rand"
)

var turns = [8]image.Point{
	{-1, -1}, {-1,  0}, {-1,  1},
	{ 0, -1},           { 0,  1},
	{ 1, -1}, { 1,  0}, { 1,  1},
    }

func calc_color(values []int8, conditions []map[int8]generate_conditions) int {
    res := 0
    for i := 0; i < len(conditions) && res == 0; i++ {
	res = i * check_all_condition(values, conditions[i])
    }
    return res
}

func sum_point_neighbor(color_map *[][]int, point image.Point, conditions []map[int8]generate_conditions) int {
    available_rect := image.Rect(0,0,len(*color_map),len((*color_map)[0]))

    sum := make([]int8, 8)
    
    for _, turn := range turns {
	next_point := point.Add(turn)
	
	if !next_point.In(available_rect){
	    continue
	}
	
	val := (*color_map)[next_point.X][next_point.Y]
	sum[val]++
    }
    
    return calc_color(sum, conditions)
}

func make_color_sheet(sheet *ColorfullSheet, max_val uint) [][]int {
    color_sheet := make([][]int, int(sheet.Width()))
    for x := 0; x < int(sheet.Width()); x++ {
	color_sheet[x] = make([]int, int(sheet.Height()))
	
	if max_val == 0 {
	    continue
	}
	
	for y := 0; y < len(color_sheet[x]); y++ {
	    color_sheet[x][y] = rand.Intn(int(max_val))
	}
    }
    return color_sheet
}

func FillSheet (sheet *ColorfullSheet, object_maker func(uint)*ColorfullObject) *ColorfullSheet{
    color_sheet := make_color_sheet(sheet, uint(7))
    swap_color_sheet := make_color_sheet(sheet, 0)
    
    conditions, repeat_cnt := read_conditions()
    
    for i := 0; i < repeat_cnt; i++ {
	for x := 0; x < len(color_sheet); x++ {
	    for y := 0; y < len(color_sheet[x]); y++ {
		swap_color_sheet[x][y] = sum_point_neighbor(&color_sheet, image.Point{x, y}, conditions)
	    }
	}
	color_sheet, swap_color_sheet = swap_color_sheet, color_sheet
    }
    
    for x := 0; x < len(color_sheet); x++ {
	for y := 0; y < len(color_sheet[x]); y++ {
	    sheet.SetObjectAt(image.Point{x, y}, object_maker(uint(color_sheet[x][y])))
	}
    }
    
    return sheet
}
