# Week 04 工程化实践

## 作业题目

- 按照自己的构想，写一个项目满足基本的目录结构和工程化实践。
- 代码需要包含对数据层、业务层、API 注册，以及 main 函数对于服务的注册和启动、信号处理。
- 使用 Wire 构建依赖。可以使用自己熟悉的框架。

## 作业代码 - Catalog API

- 一个简单的对 Product 资源进行 CRUD 操作的 API.
- 主要目的是熟悉课程中教授的目录结构和工程化实践。

## 如何运行

- 运行一个 MySQL 实例，当中创建一个 `testdb` 数据库。
- 修改 `catalog/configs/config.yaml` 文件使用相应的 MySQL 实例信息。

```yaml
data:
  database:
    driver: mysql
    source: root:password@tcp(127.0.0.1:3306)/testdb?parseTime=True
```

- 运行：

```bash
cd {PATH_TO_THE_PROJECT}/geek_university_go_course/week_04/catalog
make run
```

- 在本地测试 HTTP 服务使用 `:8080` 端口。