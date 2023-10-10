package array

// Sort 按照 Id 进行排序
type Sort []interface{}

func (a Sort) Len() int      { return len(a) }
func (a Sort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a Sort) Less(i, j int) bool {
	item1, ok1 := a[i].(struct{ Id int })
	item2, ok2 := a[j].(struct{ Id int })
	if !ok1 || !ok2 {
		panic("SortArray contains non-structs or structs without an Id field")
	}
	return item1.Id < item2.Id
}
