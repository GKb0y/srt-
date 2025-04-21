import whisper
from pydub import AudioSegment
import datetime

def audio_to_srt(audio_path, model_name="base", output_srt="output.srt"):
    """
    将音频文件转换为 SRT 字幕文件
    :param audio_path: 输入音频文件路径（支持 mp3, wav 等）
    :param model_name: Whisper 模型名称（可选：tiny, base, small, medium, large）
    :param output_srt: 输出 SRT 文件路径
    """
    # 加载 Whisper 模型
    model = whisper.load_model(model_name)
    
    # 识别音频并获取带时间戳的结果
    result = model.transcribe(audio_path)
    
    # 生成 SRT 内容
    srt_content = ""
    for i, segment in enumerate(result["segments"]):
        start_time = datetime.timedelta(seconds=segment["start"])
        end_time = datetime.timedelta(seconds=segment["end"])
        
        # 格式化时间戳为 SRT 标准格式 (HH:MM:SS,mmm)
        start_str = format_timedelta(start_time)
        end_str = format_timedelta(end_time)
        
        # 构建 SRT 条目
        srt_content += f"{i+1}\n{start_str} --> {end_str}\n{segment['text'].strip()}\n\n"
    
    # 写入文件
    with open(output_srt, "w", encoding="utf-8") as f:
        f.write(srt_content)
    print(f"SRT 文件已生成：{output_srt}")

def format_timedelta(td: datetime.timedelta) -> str:
    """将时间戳格式化为 SRT 标准时间格式 (HH:MM:SS,mmm)"""
    hours, remainder = divmod(td.seconds, 3600)
    minutes, seconds = divmod(remainder, 60)
    milliseconds = td.microseconds // 1000
    return f"{hours:02}:{minutes:02}:{seconds:02},{milliseconds:03}"

if __name__ == "__main__":
    # 示例：将 output.wav 转换为字幕，使用 base 模型
    audio_to_srt("output.wav", model_name="base", output_srt="output.srt")