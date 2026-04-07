基于 HTML 子集的 XML 格式描述飞书文档内容。

# 一、标准 HTML 标签
p, h1-h9, ul, ol, li, table, thead, tbody, tr, th, td, blockquote, pre, code, hr, img, b, em, u, del, a, br, span 语义不变

# 二、扩展标签速查表
## 块级标签
|标签|说明|关键属性|
|-|-|-|
| `<title>` | 文档标题（每篇唯一）| `align` |
| `<checkbox>` | 待办项| `done="true"\|"false"` |

## 容器标签
|标签|说明|关键属性|
|-|-|-|
| `<callout>` | 高亮提示框，子块仅支持文本、标题、列表、待办、引用 | `emoji`(默认 bulb), `background-color`, `border-color`, `text-color` |
| `<grid>` + `<column>` | 分栏布局，各列 width-ratio 之和为 1 | `width-ratio` |
| `<whiteboard>` | 嵌入画板 | `type`: `mermaid` \| `plantuml` \| `blank` |
| `<pre>` | （代码块，内含 `code`）| `lang`, `caption` |
| `<figure>` | 视图容器 | `view-type` |

## 行内组件
| 标签 | 说明 | 关键属性 |
|-------------------------|-|-|
| `<cite type="user">` | @人 | `<cite type="user" user-id="张三"></cite>` |
| `<cite type="doc">` | @文档 | `<cite type="doc" doc-id="需求文档"></cite>` |
| `<latex>` | 行内公式 | `<latex>E = mc^2</latex>` |
| `<img>` | 图片（可独立成块或内联） | `<img width="800" height="600" caption="说明" name="图.png"/>` |
| `<source>` | 文件附件（可独立成块或内联） | `<source name="报告.pdf"/>` |
| `<a type="url-preview">` | 预览卡片 | `<a type="url-preview" href="...">标题</a>` |
| `<button>` | 操作按钮 | `background-color`,`action` |
| `<time>` | 提醒 | |

## 文本块通用属性
- `align` — `"left"`|`"center"`|`"right"`（适用于 p / h1-h9 / li / checkbox）
- 有序列表项用 `seq="auto"` 自动编号

# 三、补充规则

## 富文本样式嵌套顺序
- 行内样式标签必须按以下固定顺序嵌套（外 → 内），关闭顺序严格反转：`<a> → <b> → <em> → <del> → <u> → <code> → <span> → 文本内容`

## 列表分组
- 连续同类型列表项自动合并为一个 `<ul>` 或 `<ol>`
- 嵌套子列表放在 `<li>` 内部


## 表格扩展
标准 HTML table 结构不变，扩展点：
- `<colgroup>` / `<col>` 定义列宽，紧跟 `<table>` 之后：`<col span="2" width="100"/>`
- `<th>` / `<td>` 增加 `background-color` 和 `vertical-align`（top | middle | bottom）
- 有表头时第一行在 `<thead>` 用 `<th>`，其余在 `<tbody>` 用 `<td>`
- 合并单元格仅起始格输出 `colspan` / `rowspan`，被合并的格不出现

# 四、美化系统
- 颜色优先使用命名色，也可写 `rgb(r,g,b)` / `rgba(r,g,b,a)`。

| 场景 | 可选命名色 |
|-|-|
| 文字颜色（span text-color）| gray, red, orange, yellow, green, blue, purple |
| 文字背景（span background-color）| light-{gray,red,orange,yellow,green,blue,purple}, medium-gray, gray, red, orange, yellow, green, blue, purple |
| 高亮框边框（callout border-color）| gray, red, orange, yellow, green, blue, purple |
| 高亮框填充（callout background-color）| light-{...}, medium-{gray,red,orange,yellow,green,blue,purple}, gray |
- 常用 emoji： 💡(默认)✅❌⚠️📝❓❗👍❤️📌🏁⭐

# 五、**重要规则**
## 标签内部的文本内容针对特殊含义符号必须进行转义(标签本身不需要)：
| 原始字符 | 转义后 |
|-|-|
| `<` | `&lt;` |
| `>` | `&gt;` |
| `&` | `&amp;` |
| `\n`（换行符） | `<br/>` |