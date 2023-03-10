# telegram-chatgpt

使用ChatGPT或者NewBing接口来建立一个可限速控制的电报机器人。

# 配置

复制`config.example.toml`为`config.toml`，然后进行以下配置。

```toml
[Bot]
Token = "" # 电报机器人的Token

[ChatGPT]
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