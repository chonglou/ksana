package ksana

type Migration struct {
	serial string
	up     []string
	down   []string
}
