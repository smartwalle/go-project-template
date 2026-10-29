// Package handler 是业务接口的 HTTP 适配层。
//
// 这里只做协议适配：绑定请求、调用用例、输出结果。
// 业务错误的响应由 jsoncodec 统一输出，handler 不自己写状态码。
//
// 接口全部用 POST，参数放在 JSON 请求体；路径是"资源 + 用例名"，
// 用例名与 Service 的方法名一致：列表 /users/list、创建 /users/create、
// 查询 /users/info、更新 /users/update。
package handler
