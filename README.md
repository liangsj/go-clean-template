![Go Clean Template](docs/img/logo.svg)

# Go Clean template
# Go 整洁模板

Clean Architecture template for Golang services
golang服务的整洁架构图模板
[![Go Report Card](https://goreportcard.com/badge/github.com/evrone/go-clean-template)](https://goreportcard.com/report/github.com/evrone/go-clean-template)
[![License](https://img.shields.io/github/license/evrone/go-clean-template.svg)](https://github.com/evrone/go-clean-template/blob/master/LICENSE)
[![Release](https://img.shields.io/github/v/release/evrone/go-clean-template.svg)](https://github.com/evrone/go-clean-template/releases/)
[![codecov](https://codecov.io/gh/evrone/go-clean-template/branch/master/graph/badge.svg?token=XE3E0X3EVQ)](https://codecov.io/gh/evrone/go-clean-template)

## Overview
## 综述
The purpose of the template is to show:
模板的目的为了展示：
- 如何组织项目，以防止项目演化成难以维护的代码
- 在哪里处理业务逻辑，以维护代码的独立、清晰和可扩展
- 当微服务增长时，不要失去控制

使用了  Robert Martin 的原则

[Go 整洁模板](https://evrone.com/go-clean-template?utm_source=github&utm_campaign=go-clean-template) 由 创建和提供支持[Evrone](https://evrone.com/?utm_source=github&utm_campaign=go-clean-template).

- how to organize a project and prevent it from turning into spaghetti code
- where to store business logic so that it remains independent, clean, and extensible
- how not to lose control when a microservice grows

Using the principles of Robert Martin (aka Uncle Bob).

[Go-clean-template](https://evrone.com/go-clean-template?utm_source=github&utm_campaign=go-clean-template) is created & supported by [Evrone](https://evrone.com/?utm_source=github&utm_campaign=go-clean-template).

## 内容
- [快速开始](#快速开始)
- [工程架构](#工程架构)
- [依赖注入](#依赖注入)
- [整洁架构](#整洁架构)
## Content
- [Quick start](#quick-start)
- [Project structure](#project-structure)
- [Dependency Injection](#dependency-injection)
- [Clean Architecture](#clean-architecture)


## 快速开始
本地开发
```sh
# Postgres, RabbitMQ
$ make compose-up
# Run app with migrations
$ make run
```
集成测试（可以在CI中运行）
```sh
# DB, app + migrations, integration tests
$ make compose-up-integration-test
```
## Quick start
Local development:
```sh
# Postgres, RabbitMQ
$ make compose-up
# Run app with migrations
$ make run
```

Integration tests (can be run in CI):
```sh
# DB, app + migrations, integration tests
$ make compose-up-integration-test
```
## 工程架构
### `cmd/app/main.go`
配置和日志能力初始化。主要的功能在 `internal/app/app.go` 
## Project structure
### `cmd/app/main.go`
Configuration and logger initialization. Then the main function "continues" in
`internal/app/app.go`.

### `conf`
配置，首先读取 `config.yml`中的内容，如果环境变量里面有符合的变量，将其覆盖yml中的配置
配置的结构在 `config.go`中
`env-required:true` 标签强制你指定值（在yml文件或者环境变量中）

对于配置，我们选择[cleanenv](https://github.com/ilyakaznacheev/cleanenv) 库
cleanenv没有很多的starts数，但是它简单并能满足我们的需求
从yaml配置文件中读取配置违背了12种原则。但是比起从环境变量种读取所有的配置，它确实是更加方便的。
这种模式假设配置默认值在yaml文件种，对安全敏感的配置在环境变量种
### `config`
Configuration. First, `config.yml` is read, then environment variables overwrite the yaml config if they match.
The config structure is in the `config.go`.
The `env-required: true` tag obliges you to specify a value (either in yaml, or in environment variables).

For configuration, we chose the [cleanenv](https://github.com/ilyakaznacheev/cleanenv) library.
It does not have many stars on GitHub, but is simple and meets all the requirements.

Reading the config from yaml contradicts the ideology of 12 factors, but in practice, it is more convenient than
reading the entire config from ENV.
It is assumed that default values are in yaml, and security-sensitive variables are defined in ENV.

### `docs`
Swagger 文档。由  [swag](https://github.com/swaggo/swag) 库自动生成
你不需要自己修改任何内容
### `docs`
Swagger documentation. Auto-generated by [swag](https://github.com/swaggo/swag) library.
You don't need to correct anything by yourself.

### 集成测试
集成测试
它会在应用容器旁启动独立的容器
使用[go-hit](https://github.com/Eun/go-hit) 会很方便的测试Rest API
### `integration-test`
Integration tests.
They are launched as a separate container, next to the application container.
It is convenient to test the Rest API using [go-hit](https://github.com/Eun/go-hit).

## `internal/app`
这里只有通常只有一个 _Run_ 函数在 `app.go` 文件种。它是 _main_ 函数的延续

主要的对象在这里生成
依赖注入通过 "New ..." 构造 (阅读依赖注入)，这个技术允许使用[依赖注入](#依赖注入)的原则进行分层，使得业务逻辑独立于其他层。

接下来，我们启动服务器并阻塞等待_select_ 中的信号正常完成。

如果 `app.go`的规模增长了，你可以将它拆分为多个文件.

对于大量的依赖，可以使用[wire](https://github.com/google/wire)

`migrate.go` 文件用于是数据库自动构建
它显示的包扣了 _migrate_ 标签
```sh
$ go run -tags migrate ./cmd/app
```

### `internal/app`
There is always one _Run_ function in the `app.go` file, which "continues" the _main_ function.

This is where all the main objects are created.
Dependency injection occurs through the "New ..." constructors (see Dependency Injection).
This technique allows us to layer the application using the [Dependency Injection](#dependency-injection) principle.
This makes the business logic independent from other layers.

Next, we start the server and wait for signals in _select_ for graceful completion.
If `app.go` starts to grow, you can split it into multiple files.

For a large number of injections, [wire](https://github.com/google/wire) can be used.

The `migrate.go` file is used for database auto migrations.
It is included if an argument with the _migrate_ tag is specified.
For example:

```sh
$ go run -tags migrate ./cmd/app
```

### `internal/controller`
服务器处理层（MVC 控制层），这个模块展示两个服务
- RPC 
- REST http 

服务的路由用同样的风格进行编写
- handler 按照应用领域进行分组（又共同的基础）
- 对于每一个分组，创建自己的路由结构体和请求path路径
- 业务逻辑的结构被注入到路由器结构中，它将被处理程序调用

### `internal/controller`
Server handler layer (MVC controllers). The template shows 2 servers:
- RPC (RabbitMQ as transport)
- REST http (Gin framework)

Server routers are written in the same style:
- Handlers are grouped by area of application (by a common basis)
- For each group, its own router structure is created, the methods of which process paths
- The structure of the business logic is injected into the router structure, which will be called by the handlers

#### `internal/controller/http`


#### `internal/controller/http`
Simple REST versioning.
For v2, we will need to add the `http/v2` folder with the same content.
And in the file `internal/app` add the line:
```go
handler := gin.New()
v1.NewRouter(handler, t)
v2.NewRouter(handler, t)
```

Instead of Gin, you can use any other http framework or even the standard `net/http` library.

In `v1/router.go` and above the handler methods, there are comments for generating swagger documentation using [swag](https://github.com/swaggo/swag).

### `internal/entity`
Entities of business logic (models) can be used in any layer.
There can also be methods, for example, for validation.

### `internal/usecase`
Business logic.
- Methods are grouped by area of application (on a common basis)
- Each group has its own structure
- One file - one structure

Repositories, webapi, rpc, and other business logic structures are injected into business logic structures
(see [Dependency Injection](#dependency-injection)).

#### `internal/usecase/repo`
A repository is an abstract storage (database) that business logic works with.

#### `internal/usecase/webapi`
It is an abstract web API that business logic works with.
For example, it could be another microservice that business logic accesses via the REST API.
The package name changes depending on the purpose.

### `pkg/rabbitmq`
RabbitMQ RPC pattern:
- There is no routing inside RabbitMQ
- Exchange fanout is used, to which 1 exclusive queue is bound, this is the most productive config
- Reconnect on the loss of connection

## Dependency Injection
In order to remove the dependence of business logic on external packages, dependency injection is used.

For example, through the New constructor, we inject the dependency into the structure of the business logic.
This makes the business logic independent (and portable).
We can override the implementation of the interface without making changes to the `usecase` package.

```go
package usecase

import (
    // Nothing!
)

type Repository interface {
    Get()
}

type UseCase struct {
    repo Repository
}

func New(r Repository) *UseCase{
    return &UseCase{
        repo: r,
    }
}

func (uc *UseCase) Do()  {
    uc.repo.Get()
}
```

It will also allow us to do auto-generation of mocks (for example with [mockery](https://github.com/vektra/mockery)) and easily write unit tests.

> We are not tied to specific implementations in order to always be able to change one component to another.
> If the new component implements the interface, nothing needs to be changed in the business logic.

## Clean Architecture
### Key idea
Programmers realize the optimal architecture for an application after most of the code has been written.

> A good architecture allows decisions to be delayed to as late as possible.

### The main principle
Dependency Inversion (the same one from SOLID) is the principle of dependency inversion.
The direction of dependencies goes from the outer layer to the inner layer.
Due to this, business logic and entities remain independent from other parts of the system.

So, the application is divided into 2 layers, internal and external:
1. **Business logic** (Go standard library).
2. **Tools** (databases, servers, message brokers, any other packages and frameworks).

![Clean Architecture](docs/img/layers-1.png)

**The inner layer** with business logic should be clean. It should:
- Not have package imports from the outer layer.
- Use only the capabilities of the standard library.
- Make calls to the outer layer through the interface (!).

The business logic doesn't know anything about Postgres or a specific web API.
Business logic has an interface for working with an _abstract_ database or _abstract_ web API.

**The outer layer** has other limitations:
- All components of this layer are unaware of each other's existence. How to call another from one tool? Not directly, only through the inner layer of business logic.
- All calls to the inner layer are made through the interface (!).
- Data is transferred in a format that is convenient for business logic (`internal/entity`).

For example, you need to access the database from HTTP (controller).
Both HTTP and database are in the outer layer, which means they know nothing about each other.
The communication between them is carried out through `usecase` (business logic):

```
    HTTP > usecase
           usecase > repository (Postgres)
           usecase < repository (Postgres)
    HTTP < usecase
```
The symbols > and < show the intersection of layer boundaries through Interfaces.
The same is shown in the picture:

![Example](docs/img/example-http-db.png)

Or more complex business logic:
```
    HTTP > usecase
           usecase > repository
           usecase < repository
           usecase > webapi
           usecase < webapi
           usecase > RPC
           usecase < RPC
           usecase > repository
           usecase < repository
    HTTP < usecase
```

### Layers
![Example](docs/img/layers-2.png)

### Clean Architecture Terminology
- **Entities** are structures that business logic operates on.
  They are located in the `internal/entity` folder.
  In MVC terms, entities are models.
  
- **Use Cases** is business logic located in `internal/usecase`.

The layer with which business logic directly interacts is usually called the _infrastructure_ layer.
These can be repositories `internal/usecase/repo`, external webapi `internal/usecase/webapi`, any pkg, and other microservices.
In the template, the _infrastructure_ packages are located inside `internal/usecase`.

You can choose how to call the entry points as you wish. The options are:
- controller (in our case)
- delivery
- transport
- gateways
- entrypoints
- primary
- input

### Additional layers
The classic version of [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) was designed for building large monolithic applications and has 4 layers.

In the original version, the outer layer is divided into two more, which also have an inversion of dependencies
to each other (directed inward) and communicate through interfaces.

The inner layer is also divided into two (with separation of interfaces), in the case of complex logic.

_______________________________

Complex tools can be divided into additional layers.
However, you should add layers only if really necessary.

### Alternative approaches
In addition to Clean architecture, _Onion architecture_ and _Hexagonal_ (_Ports and adapters_) are similar to it.
Both are based on the principle of Dependency Inversion.
_Ports and adapters_ are very close to _Clean Architecture_, the differences are mainly in terminology.


## Similar projects
- [https://github.com/bxcodec/go-clean-arch](https://github.com/bxcodec/go-clean-arch)
- [https://github.com/zhashkevych/courses-backend](https://github.com/zhashkevych/courses-backend)

## Useful links
- [The Clean Architecture article](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Twelve factors](https://12factor.net/ru/)
