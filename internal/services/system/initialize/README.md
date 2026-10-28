# initialize

系统首次运行时的初始化服务。

- `Initialize`：在一个事务中检查 `system_initialized` 标记、创建内置角色（超级管理员/访客/普通用户）、创建超级管理员账户并写入初始化标记。任一步失败整体回滚。
- 同一进程内的并发请求通过互斥锁串行执行；系统已初始化时返回 `model.ErrSystemInitialized`。
