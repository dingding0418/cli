
# docs +update（更新飞书云文档）

> **前置条件：** 先阅读 [`../lark-shared/SKILL.md`](../../lark-shared/SKILL.md) 了解认证、全局参数和安全规则。

通过六种 XML 指令精确更新飞书云文档。支持字符串级别和 block 级别的操作。

## 命令

```bash
# 字符串替换
lark-cli docs +update --doc "<doc_id>" --command str_replace \
  --pattern "旧文本" --content "<p>新文本</p>"

# 字符串删除
lark-cli docs +update --doc "<doc_id>" --command str_delete \
  --pattern "要删除的文本"

# Block 插入（在指定 block 之后追加）
lark-cli docs +update --doc "<doc_id>" --command block_insert \
  --block-id "blkcnXXXX" --content "<h2>新章节</h2><p>新内容</p>"

# Block 替换
lark-cli docs +update --doc "<doc_id>" --command block_replace \
  --block-id "blkcnXXXX" --content "<p>替换后的内容</p>"

# Block 删除
lark-cli docs +update --doc "<doc_id>" --command block_delete \
  --block-id "blkcnXXXX"

# 全文覆盖（create）
lark-cli docs +update --doc "<doc_id>" --command create \
  --content "<title>全新文档</title><p>新内容</p>"

# 使用 Markdown 格式
lark-cli docs +update --doc "<doc_id>" --command str_replace \
  --doc-format markdown --pattern "旧内容" --content "新内容"
```

## 返回值

```json
{
  "ok": true,
  "identity": "user",
  "data": {
    "document": { "revision_id": 13 },
    "result": "success",
    "updated_blocks_count": 3,
    "warnings": []
  }
}
```

- **`result`**：`success` | `partial_success` | `failed`
- **`updated_blocks_count`**：实际更新的 block 数量
- **`warnings`**：警告信息列表

## 参数

| 参数 | 必填 | 说明 |
|------|------|------|
| `--doc` | 是 | 文档 URL 或 token |
| `--command` | 是 | 操作指令（见下方六种指令） |
| `--doc-format` | 否 | 内容格式：`xml`（默认）\| `markdown` |
| `--content` | 视指令 | 写入内容 |
| `--pattern` | 视指令 | 匹配文本（str_replace / str_delete） |
| `--block-id` | 视指令 | 目标 block ID（block_* 操作） |
| `--revision-id` | 否 | 基准版本号，-1 = 最新（默认 `-1`） |


# 六种指令详解

## str_replace — 全文文本替换

查找文档中匹配 `--pattern` 的文本，替换为 `--content`。**replacement 支持富文本标签。**

**必需**：`--pattern`、`--content`

```bash
# 简单文本替换
lark-cli docs +update --doc "<doc_id>" --command str_replace \
  --pattern "张三" --content "李四"

# 将普通文本替换为富文本（加粗 + 链接）
lark-cli docs +update --doc "<doc_id>" --command str_replace \
  --pattern "旧链接" --content '<b>新链接</b> <a href="https://example.com">点击查看</a>'
```


## str_delete — 全文文本删除

查找并删除匹配 `--pattern` 的文本。

**必需**：`--pattern`

```bash
lark-cli docs +update --doc "<doc_id>" --command str_delete \
  --pattern "废弃的内容"
```

## block_insert — 在指定 block 之后插入

在 `--block-id` 指定的 block 之后追加新内容。

**必需**：`--block-id`、`--content`

```bash
lark-cli docs +update --doc "<doc_id>" --command block_insert \
  --block-id "blkcnXXXX" \
  --content '<h2>新章节</h2><ul><li>要点 1</li><li>要点 2</li></ul>'
```

## block_replace — 替换指定 block

将 `--block-id` 指定的 block 内容替换为 `--content`。

**必需**：`--block-id`、`--content`

```bash
lark-cli docs +update --doc "<doc_id>" --command block_replace \
  --block-id "blkcnXXXX" \
  --content '<p>替换后的段落内容</p>'
```

> **注意**：同一个 block_id 在一次请求中只能被替换一次。

## block_delete — 删除指定 block

删除 `--block-id` 指定的 block。支持逗号分隔的多个 ID。

**必需**：`--block-id`


## create — 全文覆盖

清空文档，用 `--content` 完全重写。

**必需**：`--content`

```bash
lark-cli docs +update --doc "<doc_id>" --command create \
  --content '<title>全新文档</title><h1>概述</h1><p>新的内容</p>'
```

⚠️ 会清空文档后重写，可能丢失图片、评论等。仅在需要完全重建文档时使用。


# 典型工作流

## 精确 block 级更新

1. **获取文档内容和 block ID**：
   ```bash
   lark-cli docs +fetch --doc "<doc_id>" --detail with-ids
   ```

2. **定位目标 block**：从返回的 XML 中找到要修改的 block 及其 `id` 属性

3. **执行更新**：
   ```bash
   # 替换特定 block
   lark-cli docs +update --doc "<doc_id>" --command block_replace \
     --block-id "blkcnXXXX" --content "<p>新内容</p>"

   # 在某 block 后插入
   lark-cli docs +update --doc "<doc_id>" --command block_insert \
     --block-id "blkcnXXXX" --content "<h2>追加的章节</h2>"
   ```

## 全量编辑工作流

1. **获取完整文档**（含样式和 block ID）：
   ```bash
   lark-cli docs +fetch --doc "<doc_id>" --detail full
   ```

2. **多步修改**：
   ```bash
   # 先替换标题
   lark-cli docs +update --doc "<doc_id>" --command str_replace \
     --pattern "旧内容" --content "新内容"

   # 再插入新段落
   lark-cli docs +update --doc "<doc_id>" --command block_insert \
     --block-id "blkcnXXXX" --content "<p>补充内容</p>"
   ```

## 简单文本替换

不需要 block ID，直接匹配替换：

```bash
lark-cli docs +update --doc "<doc_id>" --command str_replace \
  --pattern "v1.0" --content "v2.0"
```

# 最佳实践

- **精确操作优于全文覆盖**：使用 `block_replace`/`block_insert` 精确修改，避免 `create` 全文覆盖
- **保护不可重建的内容**：图片、画板、电子表格等以 token 形式存储，替换时避开这些 block
- **str_replace 的 replacement 支持富文本**：可以用行内标签 `<b>`、`<a>` 等替换普通文本为富文本
- **同一 block 只能被 replace 一次**：多次修改同一 block 请合并为一次 block_replace
- **block_delete 支持批量**：用逗号分隔多个 block_id 一次删除

## 参考

- [lark-doc-fetch](lark-doc-fetch.md) — 获取文档
- [lark-doc-create](lark-doc-create.md) — 创建文档（含完整 XML 语法参考）
- [lark-doc-media-insert](lark-doc-media-insert.md) — 插入图片/文件到文档
- [lark-shared](../../lark-shared/SKILL.md) — 认证和全局参数
