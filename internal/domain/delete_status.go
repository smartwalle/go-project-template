package domain

// DeleteStatus 表示数据的逻辑删除状态。
//
// 正常状态表示数据未被删除，删除状态表示数据已被逻辑删除。
type DeleteStatus int8

const (
	// DeleteStatusNormal 表示数据正常，未被删除。
	DeleteStatusNormal DeleteStatus = 0

	// DeleteStatusDeleted 表示数据已被逻辑删除。
	DeleteStatusDeleted DeleteStatus = 1
)
