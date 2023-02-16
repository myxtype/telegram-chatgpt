# telegram-chatgpt
A simple and easy telegram bot for using chatgpt

# config

```toml
[Bot]
Token = "" # token of the telegram robot

[ChatGPT]
ApiKey = "" # OpenAI AccessToken
MaxTokens = 1024 # maximum characters per reply
Temperature = 0.9 # default

[Session]
TokensLimit = 1024 # the maximum number of characters to save the session, if exceeded, it will be deleted from the old to the new
Exp = 10 # session expiration time

[Limiter]
Tokens = 3 # speed limiter times
Interval = 1 # minutes
```