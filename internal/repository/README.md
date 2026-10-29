# internal/repository

数据访问的实现：每个业务模块一个文件，实现领域包里定义的 Repository 契约。

```go
// UserRepository 是用户仓储的实现。
//
// 错误原样返回，不做包装、不做转换。
type UserRepository struct {
	dbs.Repository[userdomain.UserEntity]
}

// NewUserRepository 创建仓储。db 由应用组装层统一持有和关闭。
func NewUserRepository(db dbs.Database) *UserRepository {
	return &UserRepository{dbs.NewRepository[userdomain.UserEntity](db)}
}
```

- 嵌入 dbs 的通用 Repository，就有了 Find、FindOne、FindList、Create、Update、Delete、Transaction
- 数据库连接由应用组装层统一持有和关闭，仓储不自己管连接
- 仓储里只放会被复用的查询；一次性的在 service 里用 `SelectBuilder` 拼
