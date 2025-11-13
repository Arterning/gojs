# GoJS - JavaScript 运行时

一个用 Go 语言实现的 JavaScript 运行时，支持事件循环、Promise、定时器和 Node.js 风格的模块系统。

## 特性

✅ **事件循环** - 完整的宏任务和微任务队列实现
✅ **Promise 支持** - 完整的 Promise API (Promise.all, Promise.race, Promise.allSettled 等)
✅ **Async/Await** - 支持异步函数
✅ **定时器** - setTimeout, setInterval, setImmediate, clearTimeout, clearInterval
✅ **微任务** - queueMicrotask 支持
✅ **Console API** - console.log, console.error, console.warn 等
✅ **Node.js 模块** - fs (文件系统) 和 path (路径处理)
✅ **CommonJS** - require() 模块加载系统
✅ **REPL** - 交互式命令行
✅ **ES 语法** - 支持 ES5.1+ 主流语法

## 安装

### 编译项目

```bash
go build -o gojs.exe
```

## 使用方法

### 运行 JS 文件

```bash
gojs test.js
```

### 启动 REPL

```bash
gojs
```

在 REPL 中可以使用以下命令：
- `.help` - 显示帮助信息
- `.exit` - 退出 REPL
- `.clear` - 清屏

### 查看帮助

```bash
gojs --help
```

### 查看版本

```bash
gojs --version
```

## 示例

### 基础示例

```javascript
// hello.js
console.log("Hello, GoJS!");

const add = (a, b) => a + b;
console.log("5 + 3 =", add(5, 3));
```

运行：
```bash
gojs hello.js
```

### Promise 示例

```javascript
// promise.js
const promise = new Promise((resolve, reject) => {
    setTimeout(() => {
        resolve("操作成功!");
    }, 1000);
});

promise.then((result) => {
    console.log(result);
});

console.log("等待 Promise...");
```

### Async/Await 示例

```javascript
// async.js
async function fetchData() {
    console.log("开始获取数据...");
    const result = await Promise.resolve("数据获取成功!");
    console.log(result);
}

fetchData();
```

### 文件系统示例

```javascript
// fs-example.js
const fs = require('fs');
const path = require('path');

// 写入文件
fs.writeFileSync('hello.txt', 'Hello from GoJS!');
console.log("文件已写入");

// 读取文件
const content = fs.readFileSync('hello.txt', 'utf8');
console.log("文件内容:", content);

// 检查文件是否存在
const exists = fs.existsSync('hello.txt');
console.log("文件存在:", exists);

// 路径操作
const fullPath = path.join('foo', 'bar', 'baz.js');
console.log("完整路径:", fullPath);

// 清理
fs.unlinkSync('hello.txt');
console.log("文件已删除");
```

### 事件循环示例

```javascript
// eventloop.js
console.log("1. 同步代码");

setTimeout(() => {
    console.log("4. 宏任务 (setTimeout)");
}, 0);

Promise.resolve().then(() => {
    console.log("3. 微任务 (Promise)");
});

queueMicrotask(() => {
    console.log("3.5. 微任务 (queueMicrotask)");
});

console.log("2. 同步代码");

// 输出顺序：
// 1. 同步代码
// 2. 同步代码
// 3. 微任务 (Promise)
// 3.5. 微任务 (queueMicrotask)
// 4. 宏任务 (setTimeout)
```

## 项目结构

```
gojs/
├── main.go              # CLI 入口
├── runtime/             # 运行时核心
│   ├── runtime.go       # 运行时主逻辑
│   ├── eventloop.go     # 事件循环实现
│   └── promise.go       # Promise 实现
├── modules/             # 内置模块
│   ├── console.go       # Console API
│   ├── fs.go            # 文件系统模块
│   ├── path.go          # 路径处理模块
│   └── require.go       # 模块加载系统
├── repl/                # REPL 实现
│   └── repl.go
└── README.md
```

## API 参考

### 全局函数

- `setTimeout(callback, delay, ...args)` - 延迟执行
- `setInterval(callback, delay, ...args)` - 定时重复执行
- `setImmediate(callback, ...args)` - 立即执行（下一个事件循环）
- `clearTimeout(id)` - 取消 setTimeout
- `clearInterval(id)` - 取消 setInterval
- `queueMicrotask(callback)` - 队列微任务

### Console API

- `console.log(...args)` - 输出日志
- `console.info(...args)` - 输出信息
- `console.warn(...args)` - 输出警告
- `console.error(...args)` - 输出错误
- `console.debug(...args)` - 输出调试信息
- `console.dir(obj)` - 输出对象
- `console.assert(condition, ...args)` - 断言
- `console.clear()` - 清屏
- `console.time(label)` - 开始计时
- `console.timeEnd(label)` - 结束计时

### fs 模块

- `fs.readFileSync(path, encoding)` - 同步读取文件
- `fs.writeFileSync(path, data)` - 同步写入文件
- `fs.existsSync(path)` - 检查文件/目录是否存在
- `fs.mkdirSync(path, options)` - 创建目录
- `fs.readdirSync(path)` - 读取目录
- `fs.unlinkSync(path)` - 删除文件
- `fs.statSync(path)` - 获取文件信息

### path 模块

- `path.join(...paths)` - 连接路径
- `path.resolve(...paths)` - 解析绝对路径
- `path.basename(path, ext)` - 获取文件名
- `path.dirname(path)` - 获取目录名
- `path.extname(path)` - 获取扩展名
- `path.parse(path)` - 解析路径
- `path.format(pathObject)` - 格式化路径
- `path.isAbsolute(path)` - 判断是否绝对路径
- `path.normalize(path)` - 规范化路径
- `path.relative(from, to)` - 计算相对路径

### Promise API

- `new Promise(executor)`
- `Promise.resolve(value)`
- `Promise.reject(reason)`
- `Promise.all(promises)`
- `Promise.race(promises)`
- `Promise.allSettled(promises)`
- `promise.then(onFulfilled, onRejected)`
- `promise.catch(onRejected)`
- `promise.finally(onFinally)`

## 技术栈

- **Go** - 主要编程语言
- **goja** - JavaScript 引擎（词法分析、语法分析、执行）
- **自实现** - 事件循环、Promise、定时器、模块系统

## 限制

- 不支持浏览器 API（DOM, fetch, XMLHttpRequest 等）
- 不支持部分 ES6+ 新特性（取决于 goja 支持情况）
- fs 模块仅支持同步操作
- 性能可能不如 Node.js

## 开发计划

- [ ] 支持 ES6 模块 (import/export)
- [ ] 添加更多 Node.js 内置模块
- [ ] 异步 I/O 支持
- [ ] 性能优化
- [ ] 更完善的错误处理

## 许可证

MIT

## 贡献

欢迎提交 Issue 和 Pull Request！
