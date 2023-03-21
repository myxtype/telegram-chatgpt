# telegram-chatgpt

使用ChatGPT或者NewBing接口来建立一个可限速控制的电报机器人。

# 配置

复制`config.example.toml`为`config.toml`，然后进行以下配置。

```toml
[Bot]
Token = "" # 电报机器人的Token
HelloText = "欢迎使用！\n1.@我、消息前加'/'或者直接回复我的消息，即可向我提问\n2.会话清除：发送/clear给我即可清除会话" # /start 的提示语句
SessionClearText = "会话已清除" # 清除会话时的提示
LimiterText = "限制每1分钟2次请求(剩余%v秒)🐢" # 超过请求数量的提示
ThinkingText = "我正在思考......" # 请求数据时占位语句

[ChatGPT]
Foreword = "你开始假装女仆，每次回答的结尾跟上：喵" # 给机器人加一个人设！～
ApiKey = "" # OpenAI的Token
MaxTokens = 1024 # AI最大回复字符数，一个中文占2 tokens
Temperature = 0.9 # 温度控制，越高随机性越强

[Session]
TokensLimit = 1024 # 保存会话的最大字符，utf8字符统计，超过这个会删除部分对话
Exp = 10 # 会话过期时间，填写分钟

[Limiter]
Tokens = 3 # 限速器在周期内可以回答的次数
Interval = 1 # 限速周期，分钟。这里填写1表示1分钟内可以发送Tokens次。
```

# 运行

执行`sh build.sh`来编译，或者`cd bot && go build`进行编译。

然后执行`./bot/bot`运行机器人，注意当前目录需要有`config.toml`配置文件。

# 说明

- @机器人、消息前加'/'或者回复机器人的消息，即可向机器人提问
- 会话清除：发送`/clear`给机器人即可清除会话

# 示例

https://t.me/botaigpt_bot