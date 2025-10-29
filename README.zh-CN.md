# is

中文 | **[English](README.md)**

一个为 Go 语言设计的微型、快速、无依赖的验证/检测工具包。

本包专注于简单的谓词风格辅助函数（返回 bool），用于常见检查：

- 电子邮件、IP/MAC 地址、电话号码
- 数字与数字字符串、布尔值、大小写检查
- UUID/ULID、哈希、颜色、编码（HTML/URL/Base64）
- JWT 令牌、语义版本、标签
- 时间格式和时区
- 长度比较、数字/字符串/时间之间的比较
- 文件/目录存在性

模块路径：`go-slim.dev/is`

## 安装

```bash
go get go-slim.dev/is
```

## 快速开始

```go
package main

import (
    "fmt"
    is "go-slim.dev/is"
)

func main() {
    fmt.Println(is.Email("user@example.com")) // true
    fmt.Println(is.PhoneNumber("13800138000")) // true（中国手机号）

    // 数字与数字字符串
    fmt.Println(is.Number("123"))  // true（仅数字）
    fmt.Println(is.Number("12.3")) // false（包含小数）
    fmt.Println(is.Numeric("12.3")) // true（允许数字字符串）

    // 布尔值检测（字符串/整数）
    fmt.Println(is.Boolean("yes")) // true
    fmt.Println(is.Boolean(1))      // true

    // IP、MAC
    fmt.Println(is.IP("127.0.0.1")) // true
    fmt.Println(is.IPv6("::1"))     // true
    fmt.Println(is.MAC("01:23:45:67:89:ab")) // true

    // 颜色
    fmt.Println(is.Color("#fff"))        // true
    fmt.Println(is.RGBA("rgba(1,2,3,0.5)")) // true

    // 其他验证
    fmt.Println(is.JWT("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9")) // true
    fmt.Println(is.Semver("1.2.3")) // true

    // 时间和时区
    fmt.Println(is.Datetime("2023-01-02", "2006-01-02")) // true
    fmt.Println(is.Timezone("UTC"))                        // true

    // 长度辅助函数（字符串按符文长度；数字按位数）
    fmt.Println(is.Length("abc", 3, "="))            // true
    fmt.Println(is.LengthBetween([]int{1,2,3}, 1, 3))  // true

    // 比较和 Between/NotBetween
    fmt.Println(is.Compare(2, 1, ">")) // true
}
```

## API 亮点

- 字符串和字节
  - `Email(s)`, `Base64(s)`, `Base64URL(s)`, `URLEncoded(s)`, `HTMLEncoded(s)`
- 标识符和哈希
  - `UUID(s)`, `UUID3/4/5(s)`, `ULID(s)`, `MD5/SHA256/SHA384/SHA512(s)`, `JWT(s)`, `Semver(s)`, `Label(s)`
- 数字和布尔值
  - `Number(v)` 对字符串仅数字；对数字类型为数字
  - `Numeric(v)` 接受带小数的数字字符串和数字类型
  - `Boolean(v)` 接受 1/0、yes/no、on/off、true/false（不区分大小写）和数字 0/1
- IP、MAC 和电话
  - `IPv4(s)`, `IPv6(s)`, `IP(s)`, `MAC(s)`, `PhoneNumber(s)`, `E164(s)`
- 颜色
  - `HEXColor(s)`, `RGB(s)`, `RGBA(s)`, `HSL(s)`, `HSLA(s)`, `Color(s)`
- 时间
  - `Datetime(s, layout)`, `Timezone(name)`
- 长度和比较
  - `Length(v, n, op)`, `LengthBetween(v, min, max)`
  - `Compare(a, b, op)`, `GreaterThan/Equal/LessThan`, `Between/NotBetween`
- 文件系统
  - `File(path)`, `Dir(path)`

注意事项：

- `Number()` 与 `Numeric()` 故意设计为不同（仅数字 vs 包含小数的数字）。
- `URLEncoded()` 检测器根据当前正则表达式允许空格或有效的 %XX 序列。
- `Between()`/`NotBetween()` 守护语义由测试覆盖；对于直接关系检查推荐使用 `Compare()`。

## 示例

查看测试以获取实时示例：

- `is_test.go` 用于核心 API 覆盖
- `boundary_test.go` 用于边界情况

## 性能

基准测试（Apple M4 Pro, Go 1.24；仅供参考）：

```
BenchmarkCompare_Number-14   ~5.5 ns/op   0 B/op   0 allocs/op
BenchmarkIPv4-14             ~17 ns/op    0 B/op   0 allocs/op
BenchmarkURL-14              ~168 ns/op   192 B/op  2 allocs/op
```

本地运行：

```bash
go test -bench . -benchmem -run=^$
```

## 测试

```bash
go test -v ./...
```

## 版本控制与兼容性

- 推荐 Go 1.20+（用 Go 1.24 测试）
- 公共 API 旨在保持稳定；任何破坏性更改将增加主版本号

## 贡献

欢迎提交问题和 PR。请包括：

- 专注的更改、明确的理由
- 针对修复/功能的测试
- 性能相关的基准测试

## 许可证

MIT
