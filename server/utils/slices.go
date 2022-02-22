package utils

func SubStr(a []string, b []string) (subA []string) {
	temp := map[string]struct{}{}
	for _, val := range b {
		if _, ok := temp[val]; !ok {
			temp[val] = struct{}{}
		}
	}
	for _, val := range a {
		if _, ok := temp[val]; !ok {
			subA = append(subA, val)
			temp[val] = struct{}{}
		}
	}
	return
}

func SubInt(a []int, b []int) (subA []int) {
	temp := map[int]struct{}{}
	for _, val := range b {
		if _, ok := temp[val]; !ok {
			temp[val] = struct{}{}
		}
	}
	for _, val := range a {
		if _, ok := temp[val]; !ok {
			subA = append(subA, val)
			temp[val] = struct{}{}
		}
	}
	return
}

func SubFloat64(a []float64, b []float64) (subA []float64) {
	temp := map[float64]struct{}{}
	for _, val := range b {
		if _, ok := temp[val]; !ok {
			temp[val] = struct{}{}
		}
	}
	for _, val := range a {
		if _, ok := temp[val]; !ok {
			subA = append(subA, val)
			temp[val] = struct{}{}
		}
	}
	return
}
