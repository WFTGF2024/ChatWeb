# -*- coding: utf-8 -*-
"""
简单说明：
这是一个基于 Flask 的推理服务，用来调用 IndexTTS2 做语音合成（Text-to-Speech）。
支持：
1. multipart/form-data 提交文本
2. 可选上传参考音频（说话人/风格提示）
3. 不上传则使用默认参考音频
4. 返回合成后的 wav 文件
"""

import os
import time
import tempfile
from flask import Flask, request, jsonify, send_file
from indextts.infer_v2 import IndexTTS2

# ========================= 模型初始化相关配置 =========================
# 模型所在目录
MODEL_DIR = "./checkpoints"

# 启动前必须存在的模型文件列表
REQUIRED_FILES = [
    "bpe.model",
    "gpt.pth",
    "config.yaml",
    "s2mel.pth",
    "wav2vec2bert_stats.pt",
]

# 检查模型文件是否齐全
for f in REQUIRED_FILES:
    if not os.path.exists(os.path.join(MODEL_DIR, f)):
        raise FileNotFoundError(f"缺少模型文件: {f}，请检查目录 {MODEL_DIR}")

# 初始化 TTS 模型
# 这里的参数根据你的实际 IndexTTS2 实现来写，不同版本可能稍有差异
tts = IndexTTS2(
    model_dir=MODEL_DIR,
    cfg_path=os.path.join(MODEL_DIR, "config.yaml"),
    use_fp16=True,          # 使用半精度，节省显存
    use_deepspeed=False,    # 是否用 deepspeed
    use_cuda_kernel=False,  # 是否用自定义 cuda kernel
)

# 创建 Flask 应用
app = Flask(__name__)

# ========================= 默认参考音频（说话风格） =========================
# 这里配置几条默认的参考音频，当用户没有上传 prompt_audio 时会用这些
DEFAULT_REFS = {
    "style1": "./default/style1.wav",
    "style2": "./default/style2.wav",
    "style3": "./default/style3.wav",
}
# 默认使用的风格
DEFAULT_STYLE = "style1"


# ========================= 合成接口 =========================
@app.route("/synthesize", methods=["POST"])
def synthesize():
    """
    TTS 主接口：
    - 只接受 multipart/form-data
    - 必填字段：text
    - 可选字段：style、prompt_audio、温度采样相关参数、情感控制参数
    - 返回：audio/wav
    """
    temp_files = []  # 用来记录临时文件，最后要删掉

    # 1. 检查 Content-Type
    if not request.content_type or "multipart/form-data" not in request.content_type:
        return jsonify({"error": "只支持 multipart/form-data"}), 400

    # 2. 获取要合成的文本
    text = request.form.get("text")
    if not text:
        return jsonify({"error": "缺少必要参数: text"}), 400

    # 3. 处理说话风格 / 参考音频
    # 用户可以指定一个 style（其实就是我们事先放的参考音频）
    style = request.form.get("style", DEFAULT_STYLE)
    if style not in DEFAULT_REFS:
        # 如果用户传了一个不存在的 style，就回退到默认的
        style = DEFAULT_STYLE

    # prompt_audio_path 是最终要传给模型的说话人提示音频
    if "prompt_audio" in request.files:
        # 用户上传了参考音频，就保存到临时文件
        f = request.files["prompt_audio"]
        tmp = tempfile.NamedTemporaryFile(delete=False, suffix=".wav")
        f.save(tmp.name)
        prompt_audio_path = tmp.name
        temp_files.append(tmp.name)
    else:
        # 用户没传，就用预置的默认音频
        prompt_audio_path = DEFAULT_REFS[style]

    # 4. 情感控制相关参数（这里先都从表单里拿，实际用哪个看你模型实现）
    # emo_control_method:
    #   0: 不使用情感控制
    #   1: 使用情感音频
    #   2: 使用情感向量
    #   3: 使用情感文本
    emo_control = int(request.form.get("emo_control_method", 0))
    emo_audio = None       # 情感参考音频路径（如果有）
    emo_weight = float(request.form.get("emo_weight", 0.65))  # 情感权重
    emo_vec = None         # 情感向量占位
    emo_text = None        # 情感文本占位
    emo_random = False     # 是否随机情感

    # 5. 生成相关的采样 / 解码参数
    # 这些参数和具体的 IndexTTS2 推理接口对应
    kwargs = dict(
        do_sample=True,  # 是否采样
        temperature=float(request.form.get("temperature", 0.8)),  # 温度
        top_p=float(request.form.get("top_p", 0.8)),              # nucleus sampling
        top_k=int(request.form.get("top_k", 30)) or None,         # top-k，0 表示不用
        num_beams=int(request.form.get("num_beams", 3)),          # beam search 数
        repetition_penalty=float(request.form.get("repetition_penalty", 10.0)),
        length_penalty=float(request.form.get("length_penalty", 0.0)),
        max_mel_tokens=int(request.form.get("max_mel_tokens", 1500)),  # 最长 mel token
    )

    # 文本过长时会被切成多段，这里控制每段最大 token 数
    max_text_tokens_per_segment = int(
        request.form.get("max_text_tokens_per_segment", 120)
    )

    # 6. 准备输出目录和文件名
    os.makedirs("outputs", exist_ok=True)
    output_path = os.path.join("outputs", f"tts_{int(time.time())}.wav")

    # 7. 调用模型进行推理
    # 注意：下面这些参数名字要和 indextts.infer_v2.IndexTTS2 的 infer 方法对得上
    out = tts.infer(
        spk_audio_prompt=prompt_audio_path,         # 说话人/风格参考音频
        text=text,                                  # 待合成文本
        output_path=output_path,                    # 输出 wav 路径
        # ===== 以下是情感相关，可按需开启 =====
        emo_audio_prompt=emo_audio if emo_control == 1 else None,
        emo_alpha=emo_weight,
        emo_vector=emo_vec,
        use_emo_text=(emo_control == 3),
        emo_text=emo_text,
        use_random=emo_random,
        # ======
        verbose=False,
        max_text_tokens_per_segment=max_text_tokens_per_segment,
        **kwargs,
    )

    # 8. 清理所有临时文件
    for f in temp_files:
        try:
            os.remove(f)
        except Exception:
            # 删除失败就算了，不影响主流程
            pass

    # 9. 返回生成的 wav
    # out 通常就是 output_path，也就是我们传进去的路径
    return send_file(out, mimetype="audio/wav")


# ========================= 健康检查接口 =========================
@app.route("/health", methods=["GET"])
def health():
    """
    简单的健康检查接口，部署时可以让监控系统访问这个接口。
    """
    return jsonify({"status": "ok"})


# ========================= 程序入口 =========================
if __name__ == "__main__":
    # 0.0.0.0 表示外网可访问
    # 端口改成你自己想要的也可以
    app.run(host="0.0.0.0", port=7206)
