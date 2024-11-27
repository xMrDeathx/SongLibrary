package utils

func ListMaker[TFrom any, TTo any](list []TFrom, f func(TFrom) TTo) []TTo {
	vsm := make([]TTo, len(list))
	for i, v := range list {
		vsm[i] = f(v)
	}
	return vsm
}
