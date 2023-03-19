package host

// 分页迭代器
type Pager interface {
	Next() *PagerResult
}

func NewPageResult() *PagerResult {
	return &PagerResult{
		Data: NewHostSet(),
	}
}

type PagerResult struct {
	Data    *HostSet
	Err     error
	HasNext bool
}
