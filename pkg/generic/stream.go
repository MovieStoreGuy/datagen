package generic

func StreamMap[In any, Out any](fn func(in In) Out, ins ...In) []Out {
	out := make([]Out, len(ins))
	for i, in := range ins {
		out[i] = fn(in)
	}
	return out
}

func StreamFilter[In any](predicate func(in In) bool, values ...In) []In {
	out := make([]In, 0, len(values))
	for _, v := range values {
		if predicate(v) {
			out = append(out, v)
		}
	}
	return out
}
