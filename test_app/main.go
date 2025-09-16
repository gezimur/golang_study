package main

import (
    "fmt"
    "os"
    "image/jpeg"
    "image/color"
    "sheet_generator"
)

func save_image(sheet *sheet_generator.ColorfullSheet){
    file, err := os.Create("output/saved_image.jpg")
    if err != nil {
        fmt.Printf("Error creating file: %v\n", err)
        return
    }
    defer file.Close()
    
    err = jpeg.Encode(file, sheet, nil)
    if err != nil {
        fmt.Printf("Error encoding image: %v\n", err)
        return
    }
    fmt.Println("Image saved successfully!")
}

func main(){
    base_colors := []color.RGBA{
	color.RGBA{0, 0, 0, 255}, 
	color.RGBA{255, 0, 0, 255}, 
	color.RGBA{0, 255, 0, 255}, 
	color.RGBA{0, 0, 0, 255}, 
	color.RGBA{0, 0, 0, 255}, 
	color.RGBA{0, 0, 0, 255}, 
	color.RGBA{0, 0, 0, 255}, 
	color.RGBA{0, 0, 0, 255}, 
    }
    
    object_maker := func(id uint) *sheet_generator.ColorfullObject {
	obj := sheet_generator.MakeColorfullObject(
	    base_colors[id & 7],
	    true,
	)
	return &obj
    }
    
    var sheet sheet_generator.ColorfullSheet
    sheet.SetSize(500, 500)

    res_sheet := sheet_generator.FillSheet(&sheet, object_maker)

    save_image(res_sheet)
}
