package pathfinder

func Sign(num float64) float64 {
	if num > 0 {
		return 1
	} else if num < 0 {
		return -1
	}
	return 0
}
