# telegram-chatgpt
使用ChatGPT接口来建立一个可限速控制的电报机器人。

# config
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

# cmd
执行`sh build.sh`来编译，然后执行`./bot/bot`运行机器人。