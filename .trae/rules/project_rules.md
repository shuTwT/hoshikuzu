1.项目使用框架版本:gofiber v2.52.9,ent orm v0.14.5

2.包名应为小写。例如，使用 `apihandler` 而不是 `api_handler` 或 `apiHandler`

3.避免使用 `util`, `common`, `base` 等过于宽泛的名称。如果需要，请创建更具描述性的包，例如 `strutil`。

4.后端在引用一个包时，需要检查该包是否存在于 import 语句中，即检查该包是否被导入，需要先导入才能添加使用

4.ui文件夹下为前端项目。vue 组件采用大驼峰命名，ts 文件采用小驼峰命名。

5.在 `pkg/domain/model/error.go` 中定义业务错误，方便统一处理。`services` 层应该处理和包装来自底层的错误，并返回给 `handlers` 层。

6.所有数据库操作都应通过 `ent` ORM 进行。 对于需要多个数据库操作的业务逻辑，使用事务来保证数据一致性。

7.使用 `slog` 或 `zap` 等库进行结构化日志记录。日志应包含时间戳、日志级别、消息以及相关的上下文信息（如 request ID, user ID）

8.接口路由命名规范：
- 所有路由都应采用 RESTful 风格。
- 路由路径应以 `/api/v1` 开头，其中 `v1` 是 API 版本号。
- 每个路由都应使用 HTTP 方法（GET, POST, PUT, DELETE 等）来表示操作类型。
- 一般情况下，路由结尾含义为:
  - `/list`：获取资源列表 GET
  - `/page`：获取资源分页列表 GET
  - `/create`：创建新资源 POST
  - `/update/:id`：更新指定 ID 的资源 PUT
  - `/query/:id`：查询指定 ID 的资源 GET
  - `/delete/:id`：删除指定 ID 的资源 DELETE
- 在未经说明的情况下，`/list`的路由无需被生成

9.在完成修改`ent/schema/`目录下的文件后，需要运行`go generate ./ent`命令来生成`ent`相关的代码，通过判断是否执行成功来确认schema文件编写是否正确，如果失败，宣告本次任务执行失败，中断任务，让用户自行处理schema中的错误。

10.后端运行在`localhost:13000`，前端运行在`localhost:5732`

11.前端项目在`ui`文件夹下，采用vue3 + typescript + naive-ui 框架。

12.前端 api 接口地址以`${BASE_URL}/v1`开头，其中`${BASE_URL}`为后端运行地址，一般情况下为`localhost:13000/api`。

13.后端handler中c.QueryParser()传入的参数为指针类型，该传入参数的变量类型不能被定义为指针，需要用`&`符号来获取指针地址。例如：
```
var pageQuery model.PageQuery
if err := c.QueryParser(&pageQuery); err != nil {
	return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
}
```

14.一般情况下前端和后端都区分一级模块和二级模块，在后端中体现为`/api/v1/{一级模块}/{二级模块}`开头的路由结构,`handlers/{一级模块}/{二级模块}/handler.go`,`services/{一级模块}/{二级模块}/service.go`的文件夹结构。在前端中体现为`/src/views/{一级模块}/{二级模块}/index.vue`,`/src/api/{一级模块}/{二级模块}.ts`,`/src/router/{一级模块}/{二级模块}.ts`的文件夹结构。

15.每一个`handler`和`service`一级模块和二级模块下都存放有一个`README.md`文件，用于说明该模块的功能和使用方法,如果没有，你需要根据文件夹名字来猜测该模块的功能和使用方法。

16.每当用户要求添加一个新的一级模块或二级模块时，都需要在`/api/v1/{一级模块}/{二级模块}`开头的路由结构中添加一个新的路由，同时在`handlers/{一级模块}/{二级模块}/handler.go`,`services/{一级模块}/{二级模块}/service.go`的文件夹结构中添加一个新的文件，同时在`/src/views/{一级模块}/{二级模块}/index.vue`,`/src/api/{一级模块}/{二级模块}.ts`,`/src/router/{一级模块}/{二级模块}.ts`的文件夹结构中添加一个新的文件。

17.在操作项目模块时，需要先将项目模块结构列成一个表格，表格中包含一级模块、二级模块、路由、handler、service、vue组件、api接口、路由文件等信息。

18. `n-date-picker`组件只能接受`number`类型的时间戳作为值，不能接受`string`类型的时间字符串。

19.使用`make proto`命令来为protobuf文件生成代码