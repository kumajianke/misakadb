# AST 在 MiQL 中的使用说明

## 项目背景

当前项目在 `client/mql/MQ.py` 中尝试支持一种面向 CLI 的轻量调用方式，例如：

```python
mq.createDB("test")
mq.createDB("test", engine="tinydb")
```

CLI 读取到这类字符串后，会调用 `MiQL.shot(expr)`，由 `MiQL` 对表达式进行解析，再映射成对象方法调用，最后发送给服务端。

这类设计的目标不是执行任意 Python 代码，而是把用户输入的 `mq.xxx(...)` 语句解释为受控的 MiQL DSL。

## AST 是什么

`ast` 是 Python 标准库中的抽象语法树模块。它的作用是把源码字符串解析成“语法结构”，而不是直接执行代码。

例如：

```python
mq.createDB("test")
```

在语法层面会被拆成如下结构：

- `Call`：表示一次函数或方法调用
- `Attribute`：表示属性访问，这里是 `.createDB`
- `Name`：表示变量名，这里是 `mq`
- `Constant`：表示字面量，这里是 `"test"`

也就是说，`ast.parse()` 只负责“看懂这段代码长什么样”，真正要不要执行、允许执行什么，取决于后续的解释逻辑。

## 当前实现的执行流程

当前 `MiQL` 的关键逻辑位于 `client/mql/MQ.py`，大致流程如下：

### 1. CLI 读取用户输入

交互式终端拿到一段文本，例如：

```python
mq.createDB("test")
```

### 2. `shot(expr)` 解析表达式

`MiQL.shot(expr)` 在收到字符串时，会先执行：

```python
node = ast.parse(expr, mode="eval").body
```

这里的含义是：

- `mode="eval"`：只允许解析表达式，不允许解析整段脚本
- `.body`：取出表达式的根节点

如果语法本身不合法，例如括号不匹配，`ast.parse()` 会直接抛出异常。

### 3. `_eval_node(node)` 递归解释语法树

当前实现支持三类节点：

#### `ast.Name`

```python
if isinstance(node, ast.Name):
    if node.id != "mq":
        raise ValueError("MiQL 语句必须以 mq 开头")
    return self
```

含义：

- 只认可名字 `mq`
- 解析到 `mq` 时，返回当前 `MiQL` 对象本身

#### `ast.Attribute`

```python
if isinstance(node, ast.Attribute):
    obj = self._eval_node(node.value)
    return getattr(obj, node.attr)
```

含义：

- 先解析左侧对象
- 再取它的属性

例如 `mq.createDB` 会被解释成：

1. 先把 `mq` 解析成当前 `MiQL` 实例
2. 再通过 `getattr(self, "createDB")` 取到对应方法

#### `ast.Call`

```python
if isinstance(node, ast.Call):
    func = self._eval_node(node.func)
    args = [ast.literal_eval(arg) for arg in node.args]
    kwargs = {
        kw.arg: ast.literal_eval(kw.value)
        for kw in node.keywords
    }
    return func(*args, **kwargs)
```

含义：

- 先拿到要调用的方法
- 再把参数节点转换成 Python 字面量
- 最后执行这个方法调用

例如 `mq.createDB("test", engine="tinydb")` 会变成：

```python
self.createDB("test", engine="tinydb")
```

### 4. 自动发送

如果表达式执行结果还是一个 `MiQL` 对象，例如：

```python
mq.createDB("test")
```

那么 `shot(expr)` 会继续调用无参版 `shot()`，把当前构造好的 `mq.{...}` 发送给服务端。

## 为什么这里没有直接使用 `eval`

如果使用：

```python
eval(expr)
```

问题在于：

- 会直接执行 Python 表达式
- 很容易访问运行环境中的其他对象
- 一旦暴露内建函数、模块或外部对象，风险非常高

当前实现改为：

1. 先用 `ast.parse()` 只做语法解析
2. 再由 `_eval_node()` 手动控制允许的节点和行为

这是一种“受控解释执行”方式，比直接 `eval()` 安全得多。

## `ast.literal_eval` 的作用

当前参数解析使用了：

```python
ast.literal_eval(...)
```

这个函数只允许解析字面量，例如：

- 字符串
- 数字
- 布尔值
- `None`
- 列表
- 字典
- 元组

它不会去执行函数调用，也不会运行表达式，因此参数层面的安全性比 `eval()` 高很多。

例如：

```python
"test"
123
True
{"engine": "tinydb"}
```

这些都可以安全转换。

但类似下面这种不会被当作合法字面量执行：

```python
__import__("os").system("calc")
```

## 当前实现的安全边界

需要注意，当前实现不是“任意 Python 代码执行”，但仍然存在“属性暴露过宽”的风险。

风险点主要在这里：

```python
return getattr(obj, node.attr)
```

这意味着只要某个对象能被解析到，用户就可以继续取它的属性。

如果没有额外限制，理论上可能访问到：

```python
mq.__class__
mq._eval_node
```

如果后续又暴露了更多内部对象，还可能进一步越过 MiQL 设计边界。

## 风险判断

### 当前实现已经避免的风险

- 没有直接使用 `eval()` 或 `exec()`
- 没有开放任意变量名
- 参数使用 `ast.literal_eval()`，不会执行复杂表达式

### 当前实现仍然存在的风险

- 可以通过 `Attribute` 访问不应该开放的属性
- 可能调用内部方法或调试方法
- 如果未来类中新增敏感属性，风险会被放大

## 更推荐的做法

对于 CLI DSL，这里更适合使用“白名单方法调用”，而不是开放任意属性访问。

建议做法：

### 1. 只允许根对象为 `mq`

保留当前约束：

```python
if node.id != "mq":
    raise ValueError(...)
```

### 2. 只允许访问显式开放的方法名

例如：

```python
ALLOWED_METHODS = {"createDB", "shot"}
```

然后在解析属性时限制：

```python
if node.attr not in ALLOWED_METHODS:
    raise ValueError("不允许访问该方法")
```

### 3. 禁止访问私有属性

至少禁止以下划线开头的属性：

```python
if node.attr.startswith("_"):
    raise ValueError("禁止访问私有属性")
```

### 4. 将 DSL 和对象内部状态解耦

长期看，更稳妥的方式不是把 `mq.xxx(...)` 当成“Python 对象调用”，而是当成“受限 DSL”：

- 只识别少量固定命令
- 只接受固定参数
- 解析结果直接映射到协议数据结构

这类思路更接近成熟项目中的命令解释器，而不是通用 Python 表达式求值器。

## 为什么当前方案适合原型阶段

对于当前项目阶段，这种 AST 方案有几个优点：

- 开发成本低，能快速支持 `mq.xxx(...)` 的交互体验
- 可复用 `MiQL` 现有链式 API
- 比直接 `eval()` 明显更安全
- 后续可以逐步收紧白名单，而不需要推翻整个入口

因此，这种做法适合作为原型和过渡方案，但不建议在未加白名单限制的情况下直接作为长期稳定接口。

## 与成熟项目实践的对比

主流成熟项目在“用户输入解释执行”这一类场景里，通常会避免开放完整运行时对象能力。

### 开源项目常见做法

- 数据库 CLI 和 REPL 工具通常使用专用语法解析器，而不是直接映射到宿主语言对象
- 配置语言通常限制为声明式结构，而非任意可执行表达式
- 命令式框架会使用白名单命令注册表，而不是开放任意属性链

### 可借鉴的实践

- 明确区分“DSL 解析层”和“执行层”
- 对外仅暴露稳定命令集
- 对输入做语法和语义双重校验
- 错误信息面向用户，而不是暴露内部实现细节

## 下一步建议

### 高优先级

1. 为 `MiQL` 增加方法白名单
2. 禁止访问私有属性和内部状态
3. 对错误输出做统一包装，避免暴露内部函数名

### 中优先级

1. 为支持的 MiQL 命令写一份正式语法说明
2. 增加简单示例，例如 `createDB`、`useDB`、`createTable`
3. 为非法表达式补充测试样例

### 低优先级

1. 评估是否将 `mq.xxx(...)` 迁移为真正的 DSL 解析器
2. 统一方法命名风格，例如 `createDB` 与未来命令的命名规范

## 总结

在当前项目中，`ast` 的作用不是执行代码，而是把用户输入解析成语法树，再由 `MiQL` 手动解释并执行受控操作。

因此：

- `ast.parse()` 本身没有直接代码执行风险
- 当前方案比 `eval()` 安全很多
- 真实风险集中在“你允许 AST 访问哪些属性、调用哪些方法”

如果把访问范围收紧到白名单方法，这种方案会成为一个比较适合当前项目阶段的 MiQL CLI 过渡实现。
