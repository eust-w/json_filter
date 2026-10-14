<div align="center">
  <img src="docs/readme/hero.svg" width="100%" alt="JSON Filter workflow overview">
  <h1>JSON Filter</h1>
  <p><strong>在本机上传、逐步筛选并下载 JSON 数据。</strong></p>
  <p><a href="LICENSE"><img alt="License: Apache-2.0" src="https://img.shields.io/badge/License-Apache--2.0-d99536?style=flat-square"></a> <a href="docs/readme/validation.md"><img alt="Validation notes" src="https://img.shields.io/badge/Validation-repository%20checks-547ac4?style=flat-square"></a></p>
  <p><a href="#overview">项目亮点</a> · <a href="#quick-start">快速开始</a> · <a href="#workflow">工作流</a> · <a href="#validation">验证</a> · <a href="#documentation">文档</a></p>
  <p><a href="README.en.md">English</a> · <a href="README.md">简体中文</a></p>
</div>

---

<a id="overview"></a>
## 项目亮点

小型 Go Web 工具，用于 JSON 记录中的字符串精确匹配。监听 127.0.0.1:8080，并串行处理同一份本地筛选状态。

| 能力 | 内容 |
| :--- | :--- |
| **简单输入** | 顶层 results 对象数组。 |
| **递归匹配** | 递归查找对象和数组中的精确字符串。 |
| **明确输出** | 无匹配时结果变为空数组，同一行不重复追加。 |

### 本地页面

![json_filter local UI](docs/readme/demo.jpg)

<a id="quick-start"></a>
## 快速开始

```bash
git clone https://github.com/eust-w/json_filter.git
cd json_filter
go mod download
go run .
# Open http://127.0.0.1:8080/ and upload examples/sample.json.
```

<a id="workflow"></a>
## 工作流

```text
UPLOAD → FILTER → DOWNLOAD
```

```bash
go test -race ./...
go build .
```

<a id="validation"></a>
## 验证

```bash
go test -race ./...
go build .
```

运行环境与已执行结果见 [validation notes](docs/readme/validation.md).

| 项目 | 当前边界 |
| :--- | :--- |
| 输入 | 使用 {"results":[{...}]}；错误结构返回 HTTP 400。 |
| 下载 | 输出是筛选后的数组，不保留其他顶层字段。 |
| 用户 | 单用户本地工具，共享当前筛选状态，不是多租户服务。 |

<a id="documentation"></a>
## 文档

- [Sample input / 输入样例](examples/sample.json)
- [Regression tests / 回归测试](main_test.go)
- [Implementation / 实现](main.go)

## 许可证与来源

原创代码和新增演示资源采用 [Apache-2.0](LICENSE). 第三方代码、模型、数据和外部资源保留各自许可，详见 [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md).
