import whisper
from pydub import AudioSegment
import datetime
import argparse
import time
import sys
import subprocess
import os

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
def get_wav(video_path,out_wav):
    # 静默执行命令（不显示任何输出）
    result =subprocess.run(
    ["ffmpeg", "-i",video_path,"-vn","-acodec","pcm_s16le","-ar","44100","-ac","2",out_wav],
    stdout=subprocess.DEVNULL,
    stderr=subprocess.DEVNULL
    )
    if result.returncode==0:
        print("成功提取到音频文件:"+out_wav)
    else:
        print("提取失败")
        sys.exit()

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='语音识别转字幕工具')
    parser.add_argument('input', help='输入文件')
    parser.add_argument('output', help='输出文件')
    args = parser.parse_args()

    # 记录开始时间
    start_time = time.time()
    
    #提取音频
    out_wav=str(start_time)+".wav"
    get_wav(args.input,out_wav)

    # 示例：将 output.wav 转换为字幕，使用 base 模型
    print("开始转换音频为srt字幕")
    audio_to_srt("output.wav", model_name="base", output_srt=args.output)
    end_time = time.time()

    os.remove(out_wav)
    print(f"执行时间: {end_time - start_time:.4f} 秒")
