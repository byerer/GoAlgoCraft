package dsu

type DSU struct {
	parent []int
}

func (d *DSU) Init(n int) {
	d.parent = make([]int, n)
	for i := range d.parent {
		d.parent[i] = i
	}
}

func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		// path compression 路径压缩
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) Union(x, y int) {
	d.parent[d.Find(x)] = d.Find(y)
}
