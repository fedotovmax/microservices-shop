package utils

func SliceToSliceInterface[T, I any](in []T) []I {
	out := make([]I, len(in))
	for i := range in {
		out[i] = any(in[i]).(I)
	}
	return out
}
