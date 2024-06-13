# Sensitive

敏感词查找，验证，过滤和替换

FindAll, Validate, Filter and Replace words. 

[![Go Reference](https://pkg.go.dev/badge/github.com/houseme/sensitive.svg)](https://pkg.go.dev/github.com/houseme/sensitive)
[![sensitive CI](https://github.com/houseme/sensitive/actions/workflows/go.yml/badge.svg)](https://github.com/houseme/sensitive/actions/workflows/go.yml)
[![License](https://img.shields.io/github/license/houseme/sensitive.svg?style=flat)](https://github.com/houseme/sensitive)
![GitHub go.mod Go version (branch)](https://img.shields.io/github/go-mod/go-version/houseme/sensitive/main)


## 特别说明

基于 [sensitive](https://github.com/importcjj/sensitive), 修复了一些问题，增加了一些功能。

#### 用法

```go
package main

import (
    "github.com/houseme/sensitive"
)

func main() {
    filter := sensitive.New()
    _ = filter.LoadWordDict("path/to/dict")
    // Do something
}
```

#### AddWord

添加敏感词

```go
filter.AddWord("垃圾")
```

#### Replace

把词语中的字符替换成指定的字符，这里的字符指的是 rune 字符，比如`*`就是`'*'`。

```go
filter.Replace("这篇文章真的好垃圾", '*')
// output => 这篇文章真的好**
```

#### Filter

直接移除词语

```go
filter.Filter("这篇文章真的好垃圾啊")
// output => 这篇文章真的好啊
```

#### FindIn

查找并返回第一个敏感词，如果没有则返回`false`

```go
filter.FindIn("这篇文章真的好垃圾")
// output => true, 垃圾
```

#### Validate

验证内容是否 ok，如果含有敏感词，则返回`false`和第一个敏感词。

```go
filter.Validate("这篇文章真的好垃圾")
// output => false, 垃圾
```

#### FindAll

查找内容中的全部敏感词，以数组返回。

```go
filter.FindAll("这篇文章真的好垃圾")
// output => [垃圾]
```

#### LoadNetWordDict

加载网络词库。

```go
filter.LoadNetWordDict("https://raw.githubusercontent.com/houseme/sensitive/main/dict/dict.txt")
```

#### UpdateNoisePattern

设置噪音模式，排除噪音字符。

```go
// failed
filter.FindIn("这篇文章真的好垃 x 圾")      // false
filter.UpdateNoisePattern(`x`)
// success
filter.FindIn("这篇文章真的好垃 x 圾")      // true, 垃圾
filter.Validate("这篇文章真的好垃 x 圾")    // False, 垃圾
```

## License

`Sensitive` is licensed under the [MIT License](LICENSE), 100% free and open-source, forever.