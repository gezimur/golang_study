package sheet_generator

import (
    "image"
    "math/rand"
)

func calc_color(color_map *map[int]int, max_val int) int {
    max_color_id := 0
    diff_colors := 0
    
    for i := 1; i < int(max_val); i++ {
	if (*color_map)[i] > 0 {
	    diff_colors += 1
	} else {
	    continue
	}
	
	if (max_color_id == 0) {
	    max_color_id = i
	} else if (*color_map)[i] > (*color_map)[max_color_id] {
	    max_color_id = i
	}
	
    }

    if diff_colors > 3 {
	return 0
    } else {
	return max_color_id
    }
}

func sum_point_neighbor(color_map *[][]int, point image.Point, max_val int) int {
    turns := [8]image.Point{
	{-1, -1}, {-1,  0}, {-1,  1},
	{ 0, -1},           { 0,  1},
	{ 1, -1}, { 1,  0}, { 1,  1},
    }
    
    available_rect := image.Rect(0,0,len(*color_map),len((*color_map)[0]))

    sum := make(map[int]int)
    for i := 0; i < int(max_val); i++ {
	sum[i] = 0
    }
    
    for _, turn := range turns {
	next_point := point.Add(turn)
	
	if !next_point.In(available_rect){
	    continue
	}
	
	val := (*color_map)[next_point.X][next_point.Y]
	sum[val]++
    }
    
    return calc_color(&sum, max_val)
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

func FillSheet (sheet *ColorfullSheet, object_maker func(uint)*ColorfullObject, max_val uint) *ColorfullSheet{
    color_sheet := make_color_sheet(sheet, max_val)
    swap_color_sheet := make_color_sheet(sheet, 0)
    
    for i := 0; i < 10; i++ {
	for x := 0; x < len(color_sheet); x++ {
	    for y := 0; y < len(color_sheet[x]); y++ {
		swap_color_sheet[x][y] = sum_point_neighbor(&color_sheet, image.Point{x, y}, 7)
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
