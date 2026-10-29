// Package pagination 提供分页请求参数和分页计算。
//
// 各模块的请求参数内嵌 Page，就能直接拿到分页字段和分页计算：
// PageNumber、PageSize 是分页参数，Normalize 负责归一化，
// Limit 和 Offset 负责算条数上限和偏移量，OutOfRange 负责判断页码是否越界。
package pagination

const (
	// DefaultPageNumber 是默认页码，从 1 开始。
	DefaultPageNumber = 1

	// DefaultPageSize 是默认每页条数。
	DefaultPageSize = 20

	// MaxPageSize 是每页条数上限，超过时按上限取值。
	MaxPageSize = 100
)

// Page 是分页请求参数。
type Page struct {
	// PageNumber 是页码，从 1 开始。
	PageNumber int `json:"page"`

	// PageSize 是每页条数。
	PageSize int `json:"page_size"`
}

// Normalize 归一化分页参数：页码小于 1 按 1 处理，每页条数小于 1 按默认值处理，
// 超过上限时按上限处理。归一化之后直接用 PageNumber 和 PageSize。
func (p *Page) Normalize() {
	if p.PageNumber <= 0 {
		p.PageNumber = DefaultPageNumber
	}
	if p.PageSize <= 0 {
		p.PageSize = DefaultPageSize
	}
	if p.PageSize > MaxPageSize {
		p.PageSize = MaxPageSize
	}
}

// Limit 是当前页的条数上限。
func (p *Page) Limit() int64 {
	return int64(p.PageSize)
}

// Offset 是当前页相对第一页的偏移量。
func (p *Page) Offset() int64 {
	return int64((p.PageNumber - 1) * p.PageSize)
}

// OutOfRange 判断页码是否超出范围：总页数向上取整，页码超过总页数就是超出。
func (p *Page) OutOfRange(total int64) bool {
	var totalPage = (total + int64(p.PageSize) - 1) / int64(p.PageSize)
	return int64(p.PageNumber) > totalPage
}
