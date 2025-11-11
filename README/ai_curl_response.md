这里详细地介绍API接口的调用示例以及对应的返回结果
# <font color="red">语音识别模块</font>
## Whisper_v1
暂不使用，cnn库冲突，后续调明白会使用，因为它是SOTA

## SenseVoice_v1
* `/asr`（非流式）


```bash
curl -X POST http://127.0.0.1:7205/asr ^
     -F "file=@test.wav"
```
```bash
{"text":"\u4f60\u597d\uff0c\u8fd9\u662f\u9ed8\u8ba4\u98ce\u683c\u7684\u6d4b\u8bd5\u3002"}
```

* `/asr/stream`（流式）
```bash
curl -N -X POST http://127.0.0.1:7205/asr/stream ^
     -F "file=@test.wav"
```
```
{"partial_text": "你好，这是默认风格的测试。"}
```

# <font color="red">LLM模块</font>
* `/health`
* `/api/chat`（非流式）
* `/api/chat/stream`（流式）

##  健康检查  `/health`

```bash
curl http://127.0.0.1:7207/health
```
```
{"model":"Qwen/Qwen3-Coder-480B-A35B-Instruct","status":"ok"}
```
## 非流式对话 `/api/chat`

### 最简单对话

```bash
curl -X POST http://127.0.0.1:7207/api/chat ^
  -H "Content-Type: application/json" ^
  -d "{\"messages\":[{\"role\":\"user\",\"content\":\"你好，给我一句鸡汤\"}]}"
```
```
{"content":"\u4f60\u597d\uff01\u8fd9\u91cc\u6709\u4e00\u53e5\u6e29\u6696\u7684\u9e21\u6c64\u9001\u7ed9\u4f60\uff1a\n\n**\"\u6bcf\u4e00\u4e2a\u5e73\u51e1\u7684\u65e5\u5b50\uff0c\u90fd\u85cf\u7740\u4e0d\u5e73\u51e1\u7684\u53ef\u80fd\u3002\"**\n\n\u613f\u4f60\u4eca\u5929\u8fc7\u5f97\u5145\u5b9e\u7f8e\u597d\uff01\u2728","model":"Qwen/Qwen3-Coder-480B-A35B-Instruct","raw":{"choices":[{"delta":{"content":"","function_calls":null,"reasoning_content":"","role":null,"tool_calls":null},"finish_reason":"stop","index":0,"logprobs":null,"message":{"annotations":null,"audio":null,"content":"\u4f60\u597d\uff01\u8fd9\u91cc\u6709\u4e00\u53e5\u6e29\u6696\u7684\u9e21\u6c64\u9001\u7ed9\u4f60\uff1a\n\n**\"\u6bcf\u4e00\u4e2a\u5e73\u51e1\u7684\u65e5\u5b50\uff0c\u90fd\u85cf\u7740\u4e0d\u5e73\u51e1\u7684\u53ef\u80fd\u3002\"**\n\n\u613f\u4f60\u4eca\u5929\u8fc7\u5f97\u5145\u5b9e\u7f8e\u597d\uff01\u2728","function_call":null,"function_calls":null,"reasoning_content":"","refusal":null,"role":"assistant","tool_calls":null}}],"created":1758968053760,"id":"chatcmpl-47c483ec-ff8d-4fd1-ad1a-86dc36805546","model":"Qwen/Qwen3-Coder-480B-A35B-Instruct","object":"chat.completion","service_tier":null,"system_fingerprint":"","usage":{"completion_tokens":34,"completion_tokens_details":null,"prompt_tokens":13,"prompt_tokens_details":null,"total_tokens":47}}}
```
### 带 system 提示

```bash
curl -X POST http://127.0.0.1:7207/api/chat ^
  -H "Content-Type: application/json" ^
  -d "{\"messages\":[{\"role\":\"system\",\"content\":\"你是一个精简的助手\"},{\"role\":\"user\",\"content\":\"请用一句话总结人工智能\"}]}"
```
```
{"content":"\u4eba\u5de5\u667a\u80fd\u662f\u8ba9\u673a\u5668\u6a21\u62df\u4eba\u7c7b\u667a\u80fd\u884c\u4e3a\u7684\u6280\u672f\u3002","model":"Qwen/Qwen3-Coder-480B-A35B-Instruct","raw":{"choices":[{"delta":{"content":"","function_calls":null,"reasoning_content":"","role":null,"tool_calls":null},"finish_reason":"stop","index":0,"logprobs":null,"message":{"annotations":null,"audio":null,"content":"\u4eba\u5de5\u667a\u80fd\u662f\u8ba9\u673a\u5668\u6a21\u62df\u4eba\u7c7b\u667a\u80fd\u884c\u4e3a\u7684\u6280\u672f\u3002","function_call":null,"function_calls":null,"reasoning_content":"","refusal":null,"role":"assistant","tool_calls":null}}],"created":1758968084004,"id":"chatcmpl-7f92ebbd-c078-4c52-b1bd-40fdfdabf65e","model":"Qwen/Qwen3-Coder-480B-A35B-Instruct","object":"chat.completion","service_tier":null,"system_fingerprint":"","usage":{"completion_tokens":10,"completion_tokens_details":null,"prompt_tokens":24,"prompt_tokens_details":null,"total_tokens":34}}}
```
### 指定模型

```bash
curl -X POST http://127.0.0.1:7207/api/chat ^
  -H "Content-Type: application/json" ^
  -d "{\"model\":\"Qwen/Qwen3-Coder-480B-A35B-Instruct\",\"messages\":[{\"role\":\"user\",\"content\":\"Python 适合初学者吗？\"}]}"
```
```
{"content":"\u662f\u7684\uff0cPython \u975e\u5e38\u9002\u5408\u521d\u5b66\u8005\uff01\u4ee5\u4e0b\u662f\u4e3b\u8981\u539f\u56e0\uff1a\n\n## \ud83c\udfaf **\u8bed\u6cd5\u7b80\u6d01\u6613\u61c2**\n```python\n# Python \u4ee3\u7801\u793a\u4f8b\nprint(\"Hello, World!\")\nname = input(\"\u4f60\u7684\u540d\u5b57\uff1a\")\nprint(f\"\u4f60\u597d\uff0c{name}\uff01\")\n```\n\n## \ud83c\udf1f **Python \u7684\u4f18\u52bf**\n\n### 1. **\u8bed\u6cd5\u63a5\u8fd1\u81ea\u7136\u8bed\u8a00**\n- \u4ee3\u7801\u53ef\u8bfb\u6027\u5f3a\n- \u5b66\u4e60\u66f2\u7ebf\u5e73\u7f13\n- \u51cf\u5c11\u521d\u5b66\u8005\u7684\u8ba4\u77e5\u8d1f\u62c5\n\n### 2. **\u4e30\u5bcc\u7684\u5b66\u4e60\u8d44\u6e90**\n- \u5927\u91cf\u514d\u8d39\u6559\u7a0b\u548c\u6587\u6863\n- \u6d3b\u8dc3\u7684\u793e\u533a\u652f\u6301\n- \u4e30\u5bcc\u7684\u7b2c\u4e09\u65b9\u5e93\n\n### 3. **\u5e94\u7528\u9886\u57df\u5e7f\u6cdb**\n- Web \u5f00\u53d1\n- \u6570\u636e\u5206\u6790\n- \u4eba\u5de5\u667a\u80fd\n- \u81ea\u52a8\u5316\u811a\u672c\n- \u6e38\u620f\u5f00\u53d1\n\n### 4. **\u5f00\u53d1\u73af\u5883\u53cb\u597d**\n- \u89e3\u91ca\u578b\u8bed\u8a00\uff0c\u5373\u65f6\u53cd\u9988\n- \u5f3a\u5927\u7684 IDE \u652f\u6301\n- \u8be6\u7ec6\u7684\u9519\u8bef\u63d0\u793a\n\n## \ud83d\ude80 **\u5b66\u4e60\u5efa\u8bae**\n\n1. **\u4ece\u57fa\u7840\u5f00\u59cb**\uff1a\u53d8\u91cf\u3001\u5faa\u73af\u3001\u51fd\u6570\n2. **\u52a8\u624b\u5b9e\u8df5**\uff1a\u591a\u5199\u4ee3\u7801\uff0c\u591a\u505a\u7ec3\u4e60\n3. **\u5faa\u5e8f\u6e10\u8fdb**\uff1a\u4ece\u7b80\u5355\u9879\u76ee\u5f00\u59cb\n4. **\u5584\u7528\u8d44\u6e90**\uff1a\u5229\u7528\u5728\u7ebf\u6559\u7a0b\u548c\u793e\u533a\n\nPython \u7684\u8bbe\u8ba1\u7406\u5ff5\u5c31\u662f\"\u7b80\u6d01\u4f18\u96c5\"\uff0c\u8fd9\u8ba9\u5b83\u6210\u4e3a\u7f16\u7a0b\u5165\u95e8\u7684\u7edd\u4f73\u9009\u62e9\uff01","model":"Qwen/Qwen3-Coder-480B-A35B-Instruct","raw":{"choices":[{"delta":{"content":"","function_calls":null,"reasoning_content":"","role":null,"tool_calls":null},"finish_reason":"stop","index":0,"logprobs":null,"message":{"annotations":null,"audio":null,"content":"\u662f\u7684\uff0cPython \u975e\u5e38\u9002\u5408\u521d\u5b66\u8005\uff01\u4ee5\u4e0b\u662f\u4e3b\u8981\u539f\u56e0\uff1a\n\n## \ud83c\udfaf **\u8bed\u6cd5\u7b80\u6d01\u6613\u61c2**\n```python\n# Python \u4ee3\u7801\u793a\u4f8b\nprint(\"Hello, World!\")\nname = input(\"\u4f60\u7684\u540d\u5b57\uff1a\")\nprint(f\"\u4f60\u597d\uff0c{name}\uff01\")\n```\n\n## \ud83c\udf1f **Python \u7684\u4f18\u52bf**\n\n### 1. **\u8bed\u6cd5\u63a5\u8fd1\u81ea\u7136\u8bed\u8a00**\n- \u4ee3\u7801\u53ef\u8bfb\u6027\u5f3a\n- \u5b66\u4e60\u66f2\u7ebf\u5e73\u7f13\n- \u51cf\u5c11\u521d\u5b66\u8005\u7684\u8ba4\u77e5\u8d1f\u62c5\n\n### 2. **\u4e30\u5bcc\u7684\u5b66\u4e60\u8d44\u6e90**\n- \u5927\u91cf\u514d\u8d39\u6559\u7a0b\u548c\u6587\u6863\n- \u6d3b\u8dc3\u7684\u793e\u533a\u652f\u6301\n- \u4e30\u5bcc\u7684\u7b2c\u4e09\u65b9\u5e93\n\n### 3. **\u5e94\u7528\u9886\u57df\u5e7f\u6cdb**\n- Web \u5f00\u53d1\n- \u6570\u636e\u5206\u6790\n- \u4eba\u5de5\u667a\u80fd\n- \u81ea\u52a8\u5316\u811a\u672c\n- \u6e38\u620f\u5f00\u53d1\n\n### 4. **\u5f00\u53d1\u73af\u5883\u53cb\u597d**\n- \u89e3\u91ca\u578b\u8bed\u8a00\uff0c\u5373\u65f6\u53cd\u9988\n- \u5f3a\u5927\u7684 IDE \u652f\u6301\n- \u8be6\u7ec6\u7684\u9519\u8bef\u63d0\u793a\n\n## \ud83d\ude80 **\u5b66\u4e60\u5efa\u8bae**\n\n1. **\u4ece\u57fa\u7840\u5f00\u59cb**\uff1a\u53d8\u91cf\u3001\u5faa\u73af\u3001\u51fd\u6570\n2. **\u52a8\u624b\u5b9e\u8df5**\uff1a\u591a\u5199\u4ee3\u7801\uff0c\u591a\u505a\u7ec3\u4e60\n3. **\u5faa\u5e8f\u6e10\u8fdb**\uff1a\u4ece\u7b80\u5355\u9879\u76ee\u5f00\u59cb\n4. **\u5584\u7528\u8d44\u6e90**\uff1a\u5229\u7528\u5728\u7ebf\u6559\u7a0b\u548c\u793e\u533a\n\nPython \u7684\u8bbe\u8ba1\u7406\u5ff5\u5c31\u662f\"\u7b80\u6d01\u4f18\u96c5\"\uff0c\u8fd9\u8ba9\u5b83\u6210\u4e3a\u7f16\u7a0b\u5165\u95e8\u7684\u7edd\u4f73\u9009\u62e9\uff01","function_call":null,"function_calls":null,"reasoning_content":"","refusal":null,"role":"assistant","tool_calls":null}}],"created":1758968114624,"id":"chatcmpl-1a64fca7-28f3-429b-b35b-bfd4487ca272","model":"Qwen/Qwen3-Coder-480B-A35B-Instruct","object":"chat.completion","service_tier":null,"system_fingerprint":"","usage":{"completion_tokens":288,"completion_tokens_details":null,"prompt_tokens":15,"prompt_tokens_details":null,"total_tokens":303}}}
```

## 流式对话 `/api/chat/stream`

### 简单对话

```bash
curl -N -X POST http://127.0.0.1:7207/api/chat/stream ^
  -H "Content-Type: application/json" ^
  -d "{\"messages\":[{\"role\":\"user\",\"content\":\"请用一句话总结量子计算\"}]}"
```
```bash
{"delta": "量子"}
{"delta": "计算是"}
{"delta": "利用量子力学"}
{"delta": "原理，"}
{"delta": "通过量子比特的"}
{"delta": "叠加和纠缠特性"}
{"delta": "，实现并行"}
{"delta": "计算和"}
{"delta": "指数级加速的"}
{"delta": "新型计算范式"}
{"delta": "。"}
{"event": "done"}
```
### 带 system 提示

```bash
curl -N -X POST http://127.0.0.1:7207/api/chat/stream ^
  -H "Content-Type: application/json" ^
  -d "{\"messages\":[{\"role\":\"system\",\"content\":\"你是一位数学家\"},{\"role\":\"user\",\"content\":\"解释勾股定理\"}]}"
```
```bash
{"delta": "勾"}
{"delta": "股定理是"}
{"delta": "几何"}
#(太长了，这里省略中间的)
{"event": "done"}
```

# <font color="red">TTS模块</font>
模型权重和源码不再下载，重构源代码的 webui.py 为 Flask API 接口，直接放在源码的根目录下，按照根目录的 README.md 配置环境。uv 安装一下 flask 模块然后运行 flask.py 即可。
```bash
uv pip install flask
python3 flask_v2.py
```
## flaskapi_v1
可以正常基于说话人声音和文本进行语音合成，目前还有些小问题。

## flaskapi_v2
可以正常基于说话人声音和文本进行语音合成，并且可以返回合成后的语音文件。也可以根据参考音频生成对应风格的语音。

### 使用默认风格音频

```bash
curl -X POST http://127.0.0.1:7206/synthesize ^
  -F "text=你好，这是默认风格的测试" ^
  -F "style=style2" ^
  --output output.wav
```
```bash
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  182k  100  182k  100   293  28289     44  0:00:06  0:00:06 --:--:-- 52952
```
### 上传用户参考音频

```cmd
curl -X POST http://127.0.0.1:7206/synthesize ^
  -F "text=你好，这是用我自己的录音做参考的测试" ^
  -F "prompt_audio=@test.wav" ^
  --output output.wav
```
```bash
与使用默认风格音频的响应相同
```

* 这里的 `test.wav` 必须是你 **本地 cmd 工作目录**下的文件；


### 带调参的请求

```bash
curl -X POST http://127.0.0.1:7206/synthesize ^
  -F "text=你好，这是调节参数后的TTS测试" ^
  -F "style=style1" ^
  -F "temperature=0.6" ^
  -F "top_p=0.9" ^
  --output output.wav
```
```bash
与使用默认风格音频的响应相同
```

### 健康检查

```bash
curl http://127.0.0.1:7206/health
```
```bash
{"status":"ok"}
```

# <font color="red"> Embedding 模块</font>

虽然你主要是通过 `core.py` 调用，但也可以单独测试：

## 单条 GET

```bash
curl "http://127.0.0.1:7202/Qwen3-Embedding-4B/hello%20world"
```

## 批量 POST

```bash
curl -X POST http://127.0.0.1:7202/Qwen3-Embedding-4B ^
  -H "Content-Type: application/json" ^
  -d "{\"texts\":[\"hello world\",\"python语言\"]}"
```

## 相似度

```bash
curl -X POST http://127.0.0.1:7202/similarity ^
  -H "Content-Type: application/json" ^
  -d "{\"a\":\"python programming\",\"b\":\"learn python\"}"
```

## Rerank

```bash
curl -X POST http://127.0.0.1:7202/rerank ^
  -H "Content-Type: application/json" ^
  -d "{\"query\":\"python\",\"candidates\":[\"python tutorial\",\"java guide\",\"python flask\"]}"
```
