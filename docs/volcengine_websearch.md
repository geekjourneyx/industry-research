# Volcengine Ark Web Search API Notes

This file is a provider API source note for `researcher/internal/provider/volcengine`.

Implementation alignment:

- Default endpoint: `https://ark.cn-beijing.volces.com/api/v3/responses`.
- Secret: `ARK_API_KEY`.
- Provider type: `model_answer_search`.
- Public CLI command: `researcher answer volcengine "<query>" --json`.
- Model answer text is not final evidence. Citations and `web_search_call` output are leads that must be verified or cross-checked before supporting high-confidence claims.
- Do not send unsupported `caching` parameters.

---

Web Search（联网内容插件）是一款基础联网搜索工具，能通过 Responses API 为您的大模型获取实时的公开网络信息（如新闻、商品、天气等）。使用此工具可以解决数据时效性、知识盲区、信息同步等核心问题，并且您无需自行开发搜索引擎或维护数据资源。

---


<span id="v7Us2ZbDGr"></span>
# **核心功能**

* **多轮自动搜索**：针对复杂问题，**模型会自动判断是否需要补充搜索**，您无需手动触发。
* **图文输入兼容**：支持 VLM 模型接收图文混合输入，模型结合图片识别结果判断是否发起关联搜索。
* **多工具混合调用**：可与自定义函数、MCP 等工具协同使用，模型会自动判断调用优先级与必要性。
* **响应模式灵活**：支持同步和流式两种响应模式。流式响应可实时返回思考、搜索、回答的全过程。


---


<span id="f5d73093"></span>
# 注意事项

* 命名冲突：自定义函数名称应避免使用`web_search`，否则模型将按内置优先级判断调用逻辑。
* 权限要求：需具备火山方舟平台基础访问权限，默认支持账号维度 5 QPS，有更高并发需求推荐接入[联网搜索](https://www.volcengine.com/product/SearchInfinity)产品。
* 计费说明：按联网内容插件实际使用次数计费，**是否触发搜索由模型判断**，一轮用户查询可能触发多个关键词搜索，可通过合理设置`max_keyword`参数限制，避免过多关键词导致调用次数增加，建议根据场景设置 1~10 个。关于`max_keyword`参数的具体使用方法和示例，请参见[创建模型响应](/docs/82379/1569618)。具体收费标准详见[联网内容插件产品计费](/docs/82379/1338550)。
* 用量查询：通过响应参数`usage.tool_usage`（总次数）和`usage.tool_usage_details`（明细）查看插件使用情况。关于这两个参数的详细信息，请参见[The response object](/docs/82379/1783703)。
* 错误信息：暂不支持`caching`参数，使用该参数会返回`400`错误。


---


<span id="b44cae6e"></span>
# **快速开始**
<span id="21052a40"></span>
## 开通联网内容插件

1. 在方舟控制台打开[服务组件库](https://console.volcengine.com/ark/region:ark+cn-beijing/components)页面。
2. 在**服务组件库**页面上找到联网内容插件，并在**操作**列中点击**开通。** 

:::warning
是否触发搜索调用由模型判断，一轮搜索调用（若模型判定需要）可能会发起多个关键词同时搜索，会多次使用联网内容插件，您可以通过`max_keyword`参数来限制一轮搜索最大的关键词数量，进一步控制调用频次与成本。
:::
<span id="dc8effde"></span>
## 示例代码
:::tip
方舟平台的新用户？获取 API Key 及 开通模型等准备工作，请参见 [快速入门](/docs/82379/1399008)。

:::
```mixin-react
return (<Tabs>
<Tabs.TabPane title="cURL" key="ADqJArkJoo"><RenderMd content={`\`\`\`Bash
curl --location 'https://ark.cn-beijing.volces.com/api/v3/responses' \\
--header "Authorization: Bearer $ARK_API_KEY" \\
--header 'Content-Type: application/json' \\
--data '{
    "model": "doubao-seed-2-0-lite-260215",
    "stream": true,
    "tools": [
        {
            "type": "web_search",
            "max_keyword": 2
        }
    ],
    "input": [
        {
            "role": "user",
            "content": [
                {
                    "type": "input_text",
                    "text": "今天有什么热点新闻？"
                }
            ]
        }
    ]
}'
\`\`\`

`}></RenderMd></Tabs.TabPane>
<Tabs.TabPane title="Python SDK" key="A2f4gTS5UE"><RenderMd content={`\`\`\`Python
import os
from volcenginesdkarkruntime import Ark

# To obtain the api key from your environment variables, see: https://www.volcengine.com/docs/82379/1820161
api_key = os.getenv('ARK_API_KEY')

client = Ark(
    base_url='https://ark.cn-beijing.volces.com/api/v3',
    api_key=api_key,
)

tools = [{
    "type": "web_search",
    "max_keyword": 2,  
}]

# Create a conversation request
response = client.responses.create(
    model="doubao-seed-2-0-lite-260215",
    input=[{"role": "user", "content": "今天有什么热点新闻？"}],
    tools=tools,
)

print(response)
\`\`\`

`}></RenderMd></Tabs.TabPane>
<Tabs.TabPane title="Java SDK" key="eIIVawNar2"><RenderMd content={`\`\`\`Java
package com.ark.example;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.JsonNode;
import com.volcengine.ark.runtime.model.responses.item.*;
import com.volcengine.ark.runtime.model.responses.request.*;
import com.volcengine.ark.runtime.model.responses.response.ResponseObject;
import com.volcengine.ark.runtime.model.responses.constant.ResponsesConstants;
import com.volcengine.ark.runtime.model.responses.content.InputContentItemText;
import com.volcengine.ark.runtime.model.responses.tool.ResponsesTool;
import com.volcengine.ark.runtime.model.responses.tool.ToolFunction;
import com.volcengine.ark.runtime.model.responses.tool.ToolWebSearch;
import com.volcengine.ark.runtime.service.ArkService;

import java.util.Arrays;
import java.util.List;

public class demo {
    public static ObjectMapper om = new ObjectMapper();

    public demo() throws JsonProcessingException {
    }

    public static List<ResponsesTool> buildTools() {
        ToolWebSearch t = ToolWebSearch.builder()
                .build();
        System.out.println(Arrays.asList(t));
        return Arrays.asList(t);
    }

    public static void main(String[] args) throws JsonProcessingException {
        String apiKey = System.getenv("ARK_API_KEY");


        // Create an ArkService instance
        ArkService arkService = ArkService.builder().apiKey(apiKey).baseUrl("https://ark.cn-beijing.volces.com/api/v3").build();
        CreateResponsesRequest req = CreateResponsesRequest.builder()
                .model("doubao-seed-2-0-lite-260215")
                .input(ResponsesInput.builder().addListItem(
                        ItemEasyMessage.builder().role(ResponsesConstants.MESSAGE_ROLE_USER).content(
                                MessageContent.builder()
                                        .addListItem(InputContentItemText.builder().text("今天有什么热点新闻？").build())
                                        .build()
                        ).build()
                ).build())
                .tools(buildTools())
                .build();
        ResponseObject resp = arkService.createResponse(req);
        System.out.println(resp);

        arkService.shutdownExecutor();
    }
}
\`\`\`

`}></RenderMd></Tabs.TabPane>
<Tabs.TabPane title="Go SDK" key="QGNOQZ0cwT"><RenderMd content={`\`\`\`Go
package main

import (
        "context"
        "fmt"
        "github.com/volcengine/volcengine-go-sdk/service/arkruntime"
        "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model/responses"
        "io"
        "os"
)

/**
 * Authentication
 * If you authorize your endpoint using an API key, you can set your api key to environment variable "ARK_API_KEY"
 * client := arkruntime.NewClientWithApiKey(os.Getenv("ARK_API_KEY"))
 * Note: If you use an API key, this API key will not be refreshed.
 * To prevent the API from expiring and failing after some time, choose an API key with no expiration date.
 */

func main() {
        stream()
}

func stream() {
        fmt.Println("stream")
        client := arkruntime.NewClientWithApiKey(os.Getenv("ARK_API_KEY"))
        ctx := context.Background()
        maxToolCalls := int64(1)

        inputMessage := &responses.ItemInputMessage{
                Role: responses.MessageRole_user,
                Content: []*responses.ContentItem{
                        {
                                Union: &responses.ContentItem_Text{
                                        Text: &responses.ContentItemText{
                                                Type: responses.ContentItemType_input_text,
                                                Text: "今天有什么热点新闻？",
                                        },
                                },
                        },
                },
        }
        createResponsesReq := &responses.ResponsesRequest{
                Model: "doubao-seed-2-0-lite-260215",
                Input: &responses.ResponsesInput{
                        Union: &responses.ResponsesInput_ListValue{
                                ListValue: &responses.InputItemList{ListValue: []*responses.InputItem{{
                                        Union: &responses.InputItem_InputMessage{
                                                InputMessage: inputMessage,
                                        },
                                }}},
                        },
                },
                Tools: []*responses.ResponsesTool{
                        {
                                Union: &responses.ResponsesTool_ToolWebSearch{
                                        ToolWebSearch: &responses.ToolWebSearch{
                                                Type: responses.ToolType_web_search,
                                        },
                                },
                        },
                },
                MaxToolCalls: &maxToolCalls,
        }

        resp, err := client.CreateResponsesStream(ctx, createResponsesReq)
        if err != nil {
                fmt.Printf("stream error: %v\\n", err)
                return
        }

        for {
                event, err := resp.Recv()
                if err == io.EOF {
                        break
                }
                if err != nil {
                        fmt.Printf("stream error: %v\\n", err)
                        return
                }
                handleEvent(event)
        }

        fmt.Println()
}

func handleEvent(event *responses.Event) {
        switch event.GetEventType() {
        case responses.EventType_response_reasoning_summary_text_delta.String():
                print(event.GetReasoningText().GetDelta())
        case responses.EventType_response_reasoning_summary_text_done.String(): // aggregated reasoning text
                fmt.Printf("\\naggregated reasoning text: %s\\n", event.GetReasoningText().GetText())
        case responses.EventType_response_output_text_delta.String():
                print(event.GetText().GetDelta())
        case responses.EventType_response_output_text_done.String(): // aggregated output text
                fmt.Printf("\\naggregated output text: %s\\n", event.GetText().GetText())
        case responses.EventType_response_output_item_done.String():
                if event.GetItem().GetItem().GetFunctionMcpApprovalRequest() != nil {
                        fmt.Printf("\\nmcp call: %s; arguments %s\\n", event.GetItem().GetItem().GetFunctionMcpApprovalRequest().GetName(), event.GetItem().GetItem().GetFunctionMcpApprovalRequest().GetArguments())
                }
        case responses.EventType_response_mcp_call_arguments_delta.String():
                fmt.Printf("\\nmcp call arguments: %s\\n", event.GetResponseMcpCallArgumentsDelta().GetDelta())
        case responses.EventType_response_mcp_call_completed.String():
                fmt.Printf("\\nmcp call completed.\\n")
        case responses.EventType_response_mcp_list_tools_in_progress.String():
                fmt.Printf("\\nlisting mcp tools.\\n")
        case responses.EventType_response_mcp_list_tools_completed.String():
                fmt.Printf("\\nDone listing mcp tools.\\n")
        default:
                return
        }
}
\`\`\`

`}></RenderMd></Tabs.TabPane></Tabs>);
```


---


<span id="19b51cc2"></span>
# **参数说明**
详情请参见[创建 Responses 模型请求](https://www.volcengine.com/docs/82379/1569618)。

---


<span id="fab8ba6d"></span>
# **支持的模型**
参见[工具调用能力](/docs/82379/1330310#15a31773)。
> 关于 thinking 模式的设置，可以参考[开启/关闭深度思考](/docs/82379/1449737#fa3f44fa)。


---


<span id="55502143"></span>
# 模型输出
使用 Web Search 工具的模型回答，一般包含以下部分：

* `web_search_call`：包含搜索调用的 ID 以及执行的操作`web_search_call.action`。
* `message`：
   * `message.content[0].text`：文本结果。
   * `message.content[0].annotations`：引用网址的注释。


---


<span id="569e71b9"></span>
# 常用配置
本节介绍常见的使用场景及其字段设置方法，以帮助您更高效地运用联网内容插件。完整的参考代码见本节末尾。
<span id="b3a6340e"></span>
## 开启流式调用
在 response 中设置`stream=True`，即可开启流式调用，使响应以流式方式返回，从而更快地获取部分结果，同时可实时查看模型判断是否调用搜索的思考过程。
示例：
```Python
response = client.responses.create(
    model="doubao-seed-2-0-lite-260215",
    input=[  # 输入内容，包含系统提示和用户问题
        ...
    ],
    tools=[  # 使用工具及参数
        ...
    ],
    stream=True,  # 启用流式响应（实时返回结果，而非等待全部完成）
)
```

<span id="4a053e74"></span>
## 设置搜索来源
默认情况下，Web Search 工具会通过`search_engine`搜索全网内容。您也可以通过在 tools 中设置`sources`字段，来添加额外的内容源以优化搜索结果。是否触发搜索及使用哪些来源搜索，由模型判断。可用的附加搜索源包括：

* `"douyin"`：抖音百科
* `"moji"`：墨迹天气
* `"toutiao"`：头条图文

示例：
```Python
tools=[
    {
        "type": "web_search",  # 配置工具类型为联网内容插件
        "sources": ["douyin", "moji", "toutiao"],  # 附加搜索来源（抖音百科、墨迹天气、头条图文等平台）
    }
],
```

<span id="68c643c8"></span>
## 指定用户地理位置
您可以在 tools 中设置`user_location`字段，并提供用户的国家、地区和城市信息，以优化与地理位置相关的搜索结果。
示例：
```Python
tools=[
    {
        "type": "web_search",  # 配置工具类型为联网内容插件 
        "user_location": {  # 指定用户地理位置（用于优化搜索结果）
            "type": "approximate",  # 大致位置
            "country": "中国",
            "region": "浙江",
            "city": "杭州"
        }
    }
],
```

<span id="57cc853a"></span>
## 设置搜索限制
为平衡搜索效果与成本，您可以通过以下参数精细控制搜索行为，避免资源浪费或性能损耗：

* `tools.`**`max_keyword`**
   * 作用：限制单轮搜索中可使用的最大关键词数量。
   * 取值范围：`1`～`50`。
   * 示例：如果模型原本计划搜索三个关键词（例如“大模型最新进展”、“2025 年科技创新”），但将`max_keyword`设置为`1`，那么模型将仅使用第一个关键词进行搜索。关于`max_keyword`参数的具体使用方法和示例，请参见[创建模型响应](/docs/82379/1569618)。
* `tools.`**`limit`**
   * 作用：限制单轮搜索操作返回的最大结果条数。
   * 取值范围：`1`～`50`。
   * 默认值：`10`。
   * 说明：此参数会影响返回内容的规模和请求性能。单次搜索最多可返回 20 条结果，但单轮可能有多次搜索，默认召回 10 条。
* **`max_tool_calls`**
   * 作用：限制在一次完整的模型响应中可以执行工具调用的最大轮次。
   * 取值范围：`1` ～`10`。
   * 默认值：`3`。

示例：
```Python
tools = [{
    "type": "web_search",
    "max_keyword": 2, 
    "limit": 10, 
}],
max_tool_calls = 3,
```

<span id="MNcGnHGEfF"></span>
## 查询搜索用量
您可以在 API 返回的参数中查看联网内容插件的使用情况：

* `usage.tool_usage`：显示插件的总调用次数。
* `usage.tool_usage_details`：显示每个搜索源（例如`search_engine`、`toutiao`）的详细调用次数。

关于这两个参数的详细信息，请参见[The response object](/docs/82379/1783703)。
<span id="5638ef0e"></span>
## 设置系统提示词
系统提示词的设置对搜索请求的影响较大，建议进行优化以提升搜索的准确性与效率。以下两种系统提示词模板供您参考。
<span id="74d42962"></span>
### 模板一
```Python
# 定义系统提示词
system_prompt = """
你是AI个人助手，负责解答用户的各种问题。你的主要职责是：
1. **信息准确性守护者**：确保提供的信息准确无误。
2. **搜索成本优化师**：在信息准确性和搜索成本之间找到最佳平衡。
# 任务说明
## 1. 联网意图判断
当用户提出的问题涉及以下情况时，需使用 `web_search` 进行联网搜索：
- **时效性**：问题需要最新或实时的信息。
- **知识盲区**：问题超出当前知识范围，无法准确解答。
- **信息不足**：现有知识库无法提供完整或详细的解答。
## 2. 联网后回答
- 在回答中，优先使用已搜索到的资料。
- 回复结构应清晰，使用序号、分段等方式帮助用户理解。
## 3. 引用已搜索资料
- 当使用联网搜索的资料时，在正文中明确引用来源，引用格式为：  
`[1]  (URL地址)`。
## 4. 总结与参考资料
- 在回复的最后，列出所有已参考的资料。格式为：  
1. [资料标题](URL地址1)
2. [资料标题](URL地址2)
"""
```

<span id="dc10226a"></span>
### 模板二
```Python
# 定义系统提示词
system_prompt = """
# 角色
你是AI个人助手，负责解答用户的各种问题。你的主要职责是：
1. **信息准确性守护者**：确保提供的信息准确无误。
2. **回答更生动活泼**：请在模型的回复中多使用适当的 emoji 标签 🌟😊🎉
# 任务说明
## 1. 联网意图判断
当用户提出的问题涉及以下情况时，需使用 `web_search` 进行联网搜索：
- **时效性**：问题需要最新或实时的信息。
- **知识盲区**：问题超出当前知识范围，无法准确解答。
- **信息不足**：现有知识库无法提供完整或详细的解答。
## 2. 联网后回答
- 在回答中，优先使用已搜索到的资料。
- 回复结构应清晰，使用序号、分段等方式帮助用户理解。
## 3. 引用已搜索资料
- 当使用联网搜索的资料时，在正文中明确引用来源，引用格式为：  
`[1]  (URL地址)`。
## 4. 总结与参考资料
- 在回复的最后，列出所有已参考的资料。格式为：  
1. [资料标题](URL地址1)
2. [资料标题](URL地址2)
"""
```


---


<span id="e8909299"></span>
# 示例代码
<span id="b87a8642"></span>
## 示例 1：常用功能
请求：
```Python
import os
from volcenginesdkarkruntime import Ark

# To obtain the api key from your environment variables, see: https://www.volcengine.com/docs/82379/1820161
api_key = os.getenv('ARK_API_KEY')

client = Ark(
    base_url='https://ark.cn-beijing.volces.com/api/v3',
    api_key=api_key,
)

tools = [{
    "type": "web_search",
    "max_keyword": 3,  
    "sources": ["douyin", "moji", "toutiao"],
}]

# Send request with streaming mode enabled
response = client.responses.create(
    model="doubao-seed-1-6-250615",
    input=[{"role": "user", "content": "搜索一下大模型领域最近有什么热门的科技新闻？火山方舟最近发布了什么新模型"}],
    tools=tools,
    stream=True
)

# Store the final event (for extracting usage info)
final_event = None

for event in response:
    print(event)  # Print all events
    # Capture the completion event
    if event.type == "response.completed":
        final_event = event


# Extract statistics of tool use
if final_event:
    tool_usage = final_event.response.usage.tool_usage
    tool_usage_details = final_event.response.usage.tool_usage_details
    print("\n===== 工具使用统计 =====")
    print(f"tool_usage: {tool_usage}")
    print(f"tool_usage_details: {tool_usage_details}")
```

返回：
```Python
...

===== 工具使用统计 =====
tool_usage: ToolUsage(web_search=2, mcp=None)
tool_usage_details: ToolUsageDetails(web_search={'search_engine': 2, 'toutiao': 2}, mcp=None)
```

<span id="dd0bb90d"></span>
## 示例 2：边想边搜
以下代码通过 OpenAI SDK 调用火山方舟 Web Search 工具，实现 “AI 思考 \- 联网搜索 \- 答案生成” 全链路自动化。**触发边想边搜的必需条件由系统提示词定义**：当用户问题满足**时效性（如近 3 年数据）、知识盲区（如小众领域信息）、信息不足（如细节缺失）**  三者之一时，自动触发工具补数据。通过**流式响应**实时输出思考、搜索、回答过程，保障信息可追溯、决策可感知。
```Python
import os
from openai import OpenAI
from datetime import datetime

def realize_think_while_search():

    # 1. 初始化 OpenAI 客户端
    client = OpenAI(
        base_url="https://ark.cn-beijing.volces.com/api/v3", 
        api_key=os.getenv("ARK_API_KEY")
    )

    # 2. 定义系统提示词（核心：规范“何时搜”“怎么搜”“怎么展示思考”）
    system_prompt = """
    你是AI个人助手，需实现“边想边搜边答”，核心规则如下：
    一、思考与搜索判断（必须实时输出思考过程）：
    1. 若问题涉及“时效性（如近3年数据）、知识盲区（如具体企业薪资）、信息不足”，必须调用web_search；
    2. 思考时需说明“是否需要搜索”“为什么搜”“搜索关键词是什么”。

    二、回答规则：
    1. 优先使用搜索到的资料，引用格式为`[1] (URL地址)`；
    2. 结构清晰（用序号、分段），多使用简单易懂的表述；
    3. 结尾需列出所有参考资料（格式：1. [资料标题](URL)）。
    """

    # 3. 构造 API 请求（触发思考-搜索-回答联动）
    response = client.responses.create(
        model="doubao-seed-1-6-250615",  
        input=[
            # 系统提示词（指导 AI 行为）
            {"role": "system", "content": [{"type": "input_text", "text": system_prompt}]},
            # 用户问题（可替换为任意需边想边搜的问题）
            {"role": "user", "content": [{"type": "input_text", "text": "世界500强企业在国内所在的城市，近三年的平均工资是多少？"}]}
        ],
        tools=[
            # 配置 Web Search 工具参数
            {
                "type": "web_search",
                "limit": 10,  # 最多返回 10 条搜索结果
                "sources": ["toutiao", "douyin", "moji"],  # 优先从头条内容、抖音百科，墨迹天气搜索
                "user_location": {  # 优化地域相关搜索结果（如国内城市）
                    "type": "approximate",
                    "country": "中国",
                    "region": "浙江",
                    "city": "杭州"
                }
            }
        ],
        stream=True,  # 启用流式响应（核心：实时获取思考、搜索、回答片段）
    )

    # 4. 处理流式响应（实时展示“思考-搜索-回答”过程）
    # 状态变量：避免重复打印标题
    thinking_started = False  # AI 思考过程是否已开始打印
    answering_started = False  # AI 回答是否已开始打印

    print("=== 边想边搜启动 ===")
    for chunk in response:  # 遍历每一个实时返回的片段（chunk）
        chunk_type = getattr(chunk, "type", "")  # 获取片段类型（思考/搜索/回答）

        # ① 处理 AI 思考过程（实时打印“为什么搜、搜什么”）
        if chunk_type == "response.reasoning_summary_text.delta":
            if not thinking_started:
                print(f"\n🤔 AI思考中 [{datetime.now().strftime('%H:%M:%S')}]:")
                thinking_started = True
            # 打印思考内容（delta 为实时增量文本）
            print(getattr(chunk, "delta", ""), end="", flush=True)

        # ② 处理搜索状态（开始/完成提示）
        elif "web_search_call" in chunk_type:
            if "in_progress" in chunk_type:
                print(f"\n\n🔍 开始搜索 [{datetime.now().strftime('%H:%M:%S')}]")
            elif "completed" in chunk_type:
                print(f"\n✅ 搜索完成 [{datetime.now().strftime('%H:%M:%S')}]")

        # ③ 处理搜索关键词（展示 AI 实际搜索的内容）
        elif (chunk_type == "response.output_item.done" 
              and hasattr(chunk, "item") 
              and str(getattr(chunk.item, "id", "")).startswith("ws_")):  # ws_为搜索结果标识
            if hasattr(chunk.item.action, "query"):
                search_keyword = chunk.item.action.query
                print(f"\n📝 本次搜索关键词：{search_keyword}")

        # ④ 处理最终回答（实时整合搜索结果并输出）
        elif chunk_type == "response.output_text.delta":
            if not answering_started:
                print(f"\n\n💬 AI回答 [{datetime.now().strftime('%H:%M:%S')}]:")
                print("-" * 50)
                answering_started = True
            # 打印回答内容（实时增量输出）
            print(getattr(chunk, "delta", ""), end="", flush=True)

    # 5. 流程结束
    print(f"\n\n=== 边想边搜完成 [{datetime.now().strftime('%Y-%m-%d %H:%M:%S')}] ===")

# 运行函数
if __name__ == "__main__":
    realize_think_while_search()
```


