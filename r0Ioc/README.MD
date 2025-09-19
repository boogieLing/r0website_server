# R0Ioc

简易的依赖注入容器，基于反射，于运行时执行依赖注入。

## HOW TO USE
将“依赖注入”分解为“容器准备”与“注入至注册项”两件事。
- 先自行初始化需要用的到组件，如DAO、Config，Cache等，并将组件注入到R0Ioc。

- 准备需要被注入的注册项，在本项目中为<code>R0Route</code>，也就是各个路由
所使用到的Controller，R0Ioc还提供了简单的对结构体初始化的方法

- 对于被注入的注册项的定义处，在需要注入的组件处添加<code>\`R0Ioc:"true"\`</code>

- 初始化注册项，并使用<code>Injection</code>方法执行注入。

- 在需要的地方使用<code>R0Route</code>中的方法即可。

此外，对于组件注入容器，有两个方法<code>RegisterComponents</code>与<code>RegisterComponentSingle</code>，
前者批量注入组件，后者注入单个组件但允许提供R0Ioc退出时执行的“动作”。

最终，用户只需要

- 将组件准备好，形成实例，调用<code>RegisterComponents</code>或<code>RegisterComponentSingle</code>

- 在定义结构体时在需要注入的字段后加上<code>\`R0Ioc:"true"\`</code>的标签，

- 需要注册的结构体添加到<code>R0Route</code>

对比常用的基于运行时的依赖注入库（facebookgo/inject、go.uber.org/dig）

R0Ioc最大的优势是用户不再需要显式地提供一个初始化的函数，注册项的构造过程完全无感。
也不再需要顾及注册项实例化时的依赖问题。

只需要在定义结构体时顺带添加一个标签，
并且多个结构体能使用同一个组件（可在debug模式下验证地址），避免了多副本的内存消耗。


-----
## 原理

R0Ioc主要有三个部分

1 - R0IocItem，注入到容器的每个组件实体的上层包装，保存组件的名称以及**注入时的实例类型（非指针）**。

注意，**组件名称也能表示类型**。
```golang
type R0IocItem struct {
	Name     string
	Type     string
	Instance interface{}
	ExitFunc func(item *R0IocItem)
}
```
如一个<code>basicDao := *dao.BasicDaoMongo</code>的组件进行注入，
那么它的Name是BasicDaoMongo，而Type则是dao.BasicDaoMongo。

<code>R0IocItem.ExitFunc</code>会在退出R0Ioc容器时被执行，可以在此执行db连接的关闭等操作。

2 - R0Ioc，容器本体，本质上是一个Map，并且容器初始化后不会再被修改，其元素皆是R0IocItem

3 - R0Route，需要被注入的注册项，可以自定义名称和成员，在本项目中主要是路由的Controller。

```go
var R0Route = &struct {
	AdminArticleController *admin.ArticleController
	AdminUserController    *admin.UserController
	BaseArticleController  *base.ArticleController
	BaseUserController     *base.UserController
}{}
```
### RegisterComponents 注册组件

此方法用于将各个组件转为<code>R0IocItem</code>后注入到<code>R0Ioc</code>

它接受一个<code>...interface{}</code>，为每个组件调用<code>RegisterComponentsSingle</code>。

对于每个组件，检查它是否一个指针并获取上述所需要的信息，将组件本体引用到<code>R0IocItem.Instance</code>。

### Injection 注入器

此方法检查注册项，并将注册项所需的组件进行注入。

它接受一个<code>...interface{}</code>，每个注册项获取其字段和值信息后，调用<code>dfsInjection</code>进行递归注册。

递归注册时，检查注册项是否一个指针，否则不能进行注入。遍历注册项的字段，并获取对应的类型和名称，

一旦匹配到字段类型名称在容器中存在，并且其Tag中标识了<code>\`R0Ioc:"true"\`</code>，直接将组件的实例设置到该字段上，
此时该成员的地址与组件地址应当是相符的。

如果没有匹配到对应的组件名，但标识了<code>\`R0Ioc:"true"\`</code>，就以该字段类型新建一个指针对象（在内存中开辟了空间），并为其分配一个构造函数，将其转为<code>interface{}</code>之后设置到注册项的该字段上。
此时该注册项的成员成功“占领”了内存，所以将此成员作为注册项继续递归地进行注入。

----
## 存在的问题

- 为了方便修改，R0Ioc只接受指针类型的组件注入和注册项，以副作用的方式执行设置。

- 对于注册项<code>R0Route</code>，其在注入之前必须初始化（获得内存），本项目中的结构较为简单，所以能提供一个简易的“一键初始化”。
但面对复杂的注册项，使用者还是有不小的负担（需要自行初始化）

- 对于很多corner case没有考虑和规避，比如在组件的字段上也添加<code>\`R0Ioc:"true"\`</code>，实际上需要终止递归。

- 容器内组件与注册项需要组件的匹配规则可能需要斟酌，目前的方案面对<code>*dao.BasicDaoMongo</code>与<code>*otherDao.BasicDaoMongo</code>同时注入的情况就会出现问题。
