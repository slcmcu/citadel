package citadel

type engineMap map[string]*Docker

func (m engineMap) slice() []*Docker {
	var (
		i   int
		out = make([]*Docker, len(m))
	)

	for _, e := range m {
		out[i] = e
		i++
	}

	return out
}
