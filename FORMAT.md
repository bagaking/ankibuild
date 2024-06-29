## TOML 格式配置说明

TOML（Tom's Obvious, Minimal Language）是一种简洁明了的配置文件格式，适用于定义结构化数据。以下是如何在 AnkiBuild 项目中使用 TOML 格式配置文件的详细说明。

### 基本结构

一个典型的 `.apkg.toml` 配置文件包含以下几个部分：

1. **标题和全局设置**：
    - `title`：表示卡片包的标题。
    - `tags`：全局标签列表，适用于所有卡片。
    - `content_fmt`：内容格式，可以是 `markdown` 或 `plaintext`。

2. **问答卡片**：
    - `question`：问题文本。
    - `answer`：答案文本。
    - `tags`：特定卡片的标签列表。
    - `runtime`：运行时信息（可选），包括卡片和笔记的 ID 和 GUID。

### 示例配置

以下是一个示例 `.apkg.toml` 配置文件：

```
title = "TheExampleBarn"
tags = ["physics"]
content_fmt = "markdown"

[[q_a]]
question = "question 1"
answer = "ans 1"
tags = ["algorithms", "complexity"]

[[q_a]]
question = "are multi-line ans supported"
answer = """of course ~

- you can try this
- by yourself
"""
tags = ["algorithms"]
```

### 详细说明

- **全局设置**：
    - `title`：定义卡片包的标题。
    - `tags`：定义一个全局标签列表，这些标签将应用于所有卡片。
    - `content_fmt`：定义内容格式，可以是 `markdown` 或 `plaintext`。

- **问答卡片**：
    - `[[q_a]]`：表示一个问答卡片的开始。
    - `question`：定义卡片的正面内容，即问题。
    - `answer`：定义卡片的背面内容，即答案。支持多行文本和 Markdown 格式。
    - `tags`：定义特定卡片的标签列表，这些标签将覆盖全局标签。
    - `runtime`：可选部分，用于记录运行时信息，包括卡片和笔记的 ID 和 GUID。

### 运行时配置示例

以下是一个包含运行时信息的 `.apkg.toml` 配置文件示例：

```
title = "RuntimeExampleBarn"
tags = ["complexity"]
content_fmt = "markdown"
runtime = true

[[q_a]]
  question = "this is a runtime example"
  answer = "# runtime ans\n\n  - 1\n  "
  tags = ["algorithms", "complexity"]
  [q_a.runtime]
    cid = 1705923679008
    nid = 1705923679007
    guid = "cHuNXOn9hM"
```

在这个示例中，`runtime` 设置为 `true`，表示启用运行时信息记录。每个问答卡片的 `runtime` 部分记录了卡片和笔记的 ID 和 GUID。

通过这种方式，TOML 格式配置文件可以方便地定义和管理 Anki 问答卡片包的内容和设置。