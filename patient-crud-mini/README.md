# Patient CRUD Mini

## 演示

[Patient CRUD Mini 代码演示.html](<Patient CRUD Mini 代码演示.html>) 由 Agent 创建，用于模拟展示项目的代码和请求流程，方便演示与学习。下载后用浏览器打开即可。

## 启动

需要 Go、Docker。以下命令使用 PowerShell，在本目录运行：

```powershell
docker run -d --name patient-mini-db -e MYSQL_ROOT_PASSWORD=demo123 -e MYSQL_DATABASE=patient_mini -p 3306:3306 mysql:8
$env:MYSQL_DSN = 'root:demo123@tcp(127.0.0.1:3306)/patient_mini?charset=utf8mb4&parseTime=True&loc=Local'
go run .
```

首次启动 MySQL 需等待初始化完成。服务监听 `localhost:8080`，启动时自动创建 `patients` 表。MySQL 连接信息通过 `MYSQL_DSN` 配置。

## 验证

另开一个 PowerShell 窗口，使用虚构数据：

```powershell
curl.exe -X POST http://localhost:8080/patients -H "Content-Type: application/json" -d '{\"name\":\"Alice\",\"phone\":\"13800000001\",\"admission_number\":\"ZY001\"}'
curl.exe http://localhost:8080/patients
curl.exe -X PUT http://localhost:8080/patients/1 -H "Content-Type: application/json" -d '{\"name\":\"Bob\",\"phone\":\"13800000002\",\"admission_number\":\"ZY002\"}'
curl.exe -i -X DELETE http://localhost:8080/patients/1
```

修改和删除时，将 `1` 换成新增返回的 ID。PUT 需要提交全部三个字段。
