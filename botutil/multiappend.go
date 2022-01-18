package botutil

func Multiappend(slices ...[]byte) (res []byte) {
	for _, slice := range slices {
		res = append(res, slice...)
	}
	return
}
