package sheet_generator

import (
    "fmt"
    "os"
    "encoding/json"
)

var val = 128
var neg_sign = uint8(val)

type generate_conditions struct{
    Value int8
    Mask [3]int8
}

func (cond *generate_conditions) CheckCondition(val int8) int {
    var conditions [3]int8
    
    conditions[0] = int8((uint8(val - cond.Value) & neg_sign) >> 7) // 1 if val < params  else 0
    conditions[2] = int8((uint8(cond.Value - val) & neg_sign) >> 7) // 1 if val > params  else 0
    conditions[1] = int8(1) + (conditions[0] - conditions[2]) * (conditions[2] - conditions[0]) // 1 if val == params else 0
    
    res := 1
    
    for i := int8(0); i < int8(len(conditions)); i++ {
	res *= 1 + int((conditions[i] - cond.Mask[i]) * (cond.Mask[i] - conditions[i]))
    }
    
    return res
}

func check_all_condition(values []int8, conditions map[int8]generate_conditions) int {
    res := 1
    for i:= int8(0); i < int8(len(values)) && res == 1; i++ {
	cond, exists := conditions[i]
	if exists {
	    res *= cond.CheckCondition(values[i])
	}
    }
    return res
}

type sheet_generator_params struct {
    Conditions []map[int8]generate_conditions
    RepeatCnt int
}

func read_conditions() ([]map[int8]generate_conditions, int) {
    file_path := "generate_conditions.json"
    file_bytes, err := os.ReadFile(file_path)
    
    if err != nil {
	fmt.Println("error occurred: ", err)
	return nil, 0
    }
    
    var params sheet_generator_params
    
    json.Unmarshal(file_bytes, &params)
    
    fmt.Println(params)
    
    return params.Conditions, params.RepeatCnt
}