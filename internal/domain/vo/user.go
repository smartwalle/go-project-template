package vo

// DescribeProfileRequest 获取当前登录用户信息请求参数
type DescribeProfileRequest struct {
}

func (r DescribeProfileRequest) Check() error {
	return nil
}

// DescribeProfileResponse 获取当前登录用户信息返回数据
type DescribeProfileResponse struct {
}
