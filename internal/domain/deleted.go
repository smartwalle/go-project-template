package domain

// DeleteStatus 是逻辑删除状态。
//
// 表里用 deleted 列记录删除标记：0-正常 1-删除。
// 查询时带上 deleted = 0，就是只取没被删除的数据。
type DeleteStatus int8

const (
	// DeleteStatusNormal 是正常状态，未删除。
	DeleteStatusNormal DeleteStatus = 0

	// DeleteStatusDeleted 是已删除状态。
	DeleteStatusDeleted DeleteStatus = 1
)
