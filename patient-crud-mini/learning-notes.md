# Patient CRUD Mini Learning Notes

启动命令和 curl 示例见 `README.md`。

## 四个相关技术

| 技术 | 在项目中的作用 |
| --- | --- |
| Go | 定义数据结构和函数，执行程序逻辑 |
| Gin | 接收 HTTP 请求，匹配路由，返回响应 |
| GORM | 把 Go 结构体和数据库表关联，执行数据库操作 |
| MySQL | 保存 `patients` 表中的患者数据 |

```text
curl → Gin 路由 → 处理函数 → GORM → MySQL
curl ← HTTP 响应 ← 处理函数 ← 查询或写入结果
```

## 服务启动

1. `os.Getenv("MYSQL_DSN")` 读取数据库连接信息；没有配置时，`log.Fatal` 输出错误并结束程序。
2. `gorm.Open(mysql.Open(dsn), ...)` 使用 MySQL 驱动连接数据库。连接失败时退出。
3. `db.AutoMigrate(&Patient{})` 根据模型创建或调整表结构。
4. `gin.Default()` 创建 Gin 引擎，随后注册四条路由。
5. `r.Run(":8080")` 监听 8080 端口，等待 HTTP 请求。

MySQL 数据库 `patient_mini` 需要先存在；`AutoMigrate` 处理其中的表。

## 模型和输入

`Patient` 对应 `patients` 表。`ID` 是主键，其余字段是姓名、手机号、住院号。字段后的标签有两种用途：

- `json:"admission_number"` 指定 JSON 字段名。
- `gorm:"size:50;not null"` 指定数据库列的长度和非空约束。

`PatientInput` 只包含客户端需要提交的三个字段。新增时，程序用它构造 `Patient`，ID 由数据库生成。

## 一次新增请求

客户端发送 `POST /patients`：

```json
{"name":"Alice","phone":"13800000001","admission_number":"ZY001"}
```

1. Gin 匹配 `r.POST("/patients", ...)`，调用对应函数。
2. 函数收到 `c *gin.Context`，通过它读取本次请求、写入本次响应。
3. `readInput(c)` 调用 `c.ShouldBindJSON(&input)`，按 JSON 标签把请求体解析到 `PatientInput`。
4. `TrimSpace` 去掉字段首尾空白；代码检查三个字段非空，并按字符数限制长度。校验失败返回 400，函数提前结束。
5. 校验通过后，程序构造 `Patient`，调用 `db.Create(&patient)` 写入 MySQL。
6. GORM 把数据库生成的 ID 填回 `patient`；Gin 用 `c.JSON` 返回 201 和患者记录。写入失败则返回 500。

`&input` 和 `&patient` 是变量的地址。把指针传给函数后，函数可以修改原变量：JSON 绑定写入 `input`，GORM 写回 `patient.ID`。

## 其余三个接口

| 请求 | 代码流程 | 成功响应 |
| --- | --- | --- |
| `GET /patients` | `db.Find(&patients)` 查询全部患者 | 200，JSON 数组；没有记录时为 `[]` |
| `PUT /patients/:id` | `readID` 解析 ID，校验请求体；`First` 查找，再修改字段并用 `Save` 写回 | 200，更新后的记录 |
| `DELETE /patients/:id` | `readID` 解析 ID；`Delete` 删除，检查 `Error` 和 `RowsAffected` | 204，没有响应体 |

PUT 需要提交全部三个字段。ID 不是正整数时返回 400；PUT 查不到患者、DELETE 未删除任何记录时返回 404；数据库操作失败时返回 500。

## 语法要点

- `:=` 声明并初始化变量，例如 `dsn := os.Getenv("MYSQL_DSN")`。
- `if err != nil` 检查操作是否失败；失败后返回响应或结束程序。
- `input, ok := readInput(c)` 同时接收输入值和是否通过校验；`!ok` 时 `return`，后面的写库代码不会执行。
- `[]Patient{}` 创建空切片，作为查询结果的容器。
- `c.Param("id")` 取得路径里的字符串；`strconv.ParseUint` 将其转换为无符号整数。
- `result.Error` 表示删除操作是否出错；`result.RowsAffected` 表示删除了多少行。

## 自查

1. 新增请求中的 JSON 从哪里进入 Go 变量？
2. 患者 ID 在哪一步生成，又如何出现在响应中？
3. 修改患者时，哪些代码只改变内存中的变量，哪一行把变化写入数据库？
4. 删除不存在的 ID 时，为什么还需要检查 `RowsAffected`？
5. 如果数据库连接失败，服务能否继续启动？

## 回答

1. Gin 将请求交给 POST 处理函数后，`readInput` 调用 `c.ShouldBindJSON(&input)`，把请求体按 JSON 标签解析到 `PatientInput` 变量中。
2. `db.Create(&patient)` 插入记录时，MySQL 生成自增 ID，GORM 将 ID 写回 `patient.ID`；随后 `c.JSON(http.StatusCreated, patient)` 把它包含在响应中。
3. `patient.Name, patient.Phone, patient.AdmissionNumber = ...` 只修改 Go 内存中的变量；`db.Save(&patient)` 将修改写回数据库。
4. 删除不存在的记录通常不会产生数据库错误。`RowsAffected == 0` 表示没有删除任何行，程序据此返回 404。
5. 不能。`gorm.Open` 返回错误时，`log.Fatalf` 输出错误并结束进程，后续路由注册和服务监听都不会执行。
