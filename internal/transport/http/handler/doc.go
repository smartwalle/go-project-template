// Package handler 是业务接口的 HTTP 适配层。
//
// 这里只做协议适配：绑定请求、调用用例、输出结果。
// 业务错误的响应由 jsoncodec 统一输出，handler 不自己写状态码。
//
// 业务接口统一使用 POST；健康检查等基础设施接口按语义选择 HTTP 方法，
// 参数放在 JSON 请求体。路径按资源和接口语义确定，不要求路径动作名称
// 与 Service 方法名完全一致。例如：列表 /users/list 对应 List，
// 创建 /users/create 对应 Create，查询 /users/info 对应 Get，
// 更新 /users/update 对应 Update。
package handler
