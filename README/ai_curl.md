# <font color="red">语音识别模块</font>
## Whisper_v1
暂不使用，cnn库冲突，后续调明白会使用，因为它是SOTA
### 整段识别
```cmd
curl -X POST http://120.79.25.184:7205/transcribe ^
  -F "file=@test.wav"
```
### 流式识别
```cmd
curl -N -X POST http://120.79.25.184:7205/transcribe_stream ^
  -F "file=@test.wav"
```

## SenseVoice_v1
* `/asr`（非流式）
* `/asr/stream`（流式）

```cmd
curl -X POST http://120.79.25.184:7205/asr ^
     -F "file=@test.wav"
```
```cmd
curl -N -X POST http://120.79.25.184/asr/stream ^
     -F "file=@test.wav"
```


# <font color="red">LLM模块</font>
* `/health`
* `/api/chat`（非流式）
* `/api/chat/stream`（流式）

##  健康检查  `/health`

```bash
curl http://127.0.0.1:7207/health
```

## 非流式对话 `/api/chat`

### 最简单对话

```bash
curl -X POST http://127.0.0.1:7207/api/chat ^
  -H "Content-Type: application/json" ^
  -d "{\"messages\":[{\"role\":\"user\",\"content\":\"你好，给我一句鸡汤\"}]}"
```

### 带 system 提示

```bash
curl -X POST http://127.0.0.1:7207/api/chat ^
  -H "Content-Type: application/json" ^
  -d "{\"messages\":[{\"role\":\"system\",\"content\":\"你是一个精简的助手\"},{\"role\":\"user\",\"content\":\"请用一句话总结人工智能\"}]}"
```

### 指定模型

```bash
curl -X POST http://127.0.0.1:7207/api/chat ^
  -H "Content-Type: application/json" ^
  -d "{\"model\":\"Qwen/Qwen3-Coder-480B-A35B-Instruct\",\"messages\":[{\"role\":\"user\",\"content\":\"Python 适合初学者吗？\"}]}"
```
### 保存到 JSON 再请求：

  ```json
  {
    "messages": [
      {"role": "system", "content": "你是一个简洁的助手"},
      {"role": "user", "content": "给我一句关于未来的格言"}
    ]
  }
  ```

  然后：

  ```bash
  curl -X POST http://127.0.0.1:7207/api/chat -H "Content-Type: application/json" -d @body.json
  ```

## 流式对话 `/api/chat/stream`

### 简单对话

```bash
curl -N -X POST http://127.0.0.1:7207/api/chat/stream ^
  -H "Content-Type: application/json" ^
  -d "{\"messages\":[{\"role\":\"user\",\"content\":\"请用一句话总结量子计算\"}]}"
```

### 带 system 提示

```bash
curl -N -X POST http://127.0.0.1:7207/api/chat/stream ^
  -H "Content-Type: application/json" ^
  -d "{\"messages\":[{\"role\":\"system\",\"content\":\"你是一位数学家\"},{\"role\":\"user\",\"content\":\"解释勾股定理\"}]}"
```

# <font color="red">TTS模块</font>
模型权重和源码不再下载，重构源代码的 webui.py 为 Flask API 接口，直接放在源码的根目录下，按照根目录的 README.md 配置环境。uv 安装一下 flask 模块然后运行 flask.py 即可。
```bash
uv pip install flask
python3 flask.py
```
## flaskapi_v1
可以正常基于说话人声音和文本进行语音合成，目前还有些小问题。

## flaskapi_v2
可以正常基于说话人声音和文本进行语音合成，并且可以返回合成后的语音文件。也可以根据参考音频生成对应风格的语音。

### 使用默认风格音频

```cmd
curl -X POST http://120.79.25.184:7206/synthesize ^
  -F "text=你好，这是默认风格的测试" ^
  -F "style=style2" ^
  --output output.wav
```

### 上传用户参考音频

```cmd
curl -X POST http://120.79.25.184:7206/synthesize ^
  -F "text=你好，这是用我自己的录音做参考的测试" ^
  -F "prompt_audio=@test.wav" ^
  --output output.wav
```


* 这里的 `test.wav` 必须是你 **本地 cmd 工作目录**下的文件；


### 带调参的请求

```cmd
curl -X POST http://120.79.25.184:7206/synthesize ^
  -F "text=你好，这是调节参数后的TTS测试" ^
  -F "style=style1" ^
  -F "temperature=0.6" ^
  -F "top_p=0.9" ^
  --output output.wav
```

### 健康检查

```cmd
curl http://120.79.25.184:7206/health
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
