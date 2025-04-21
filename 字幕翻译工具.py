import sys
import os
import hashlib
import hmac
import time
import json
import requests
from zhconv import convert
import srt
import argparse
from datetime import datetime
from tqdm import tqdm

class TencentTranslator:
    def __init__(self, secret_id, secret_key):
        self.secret_id = "SecretId"
        self.secret_key = "SecretKey"
        self.endpoint = "tmt.tencentcloudapi.com"
        self.version = "2018-03-21"
        self.region = "ap-beijing"
        self.session = requests.Session()
        self.session.headers.update({'Content-Type': 'application/json; charset=utf-8'})

    def _sign(self, params, payload):
        service = 'tmt'
        timestamp = int(time.time())
        date = datetime.utcfromtimestamp(timestamp).strftime('%Y-%m-%d')
        
        canonical_headers = f"content-type:application/json\nhost:{self.endpoint}\n"
        signed_headers = "content-type;host"
        hashed_payload = hashlib.sha256(payload.encode('utf-8')).hexdigest()
        
        canonical_request = (
            "POST\n"
            "/\n"
            "\n"
            f"{canonical_headers}\n"
            f"{signed_headers}\n"
            f"{hashed_payload}"
        )

        key_date = hmac.new(f"TC3{self.secret_key}".encode('utf-8'), date.encode('utf-8'), 'sha256').digest()
        key_service = hmac.new(key_date, service.encode('utf-8'), 'sha256').digest()
        key_signing = hmac.new(key_service, b"tc3_request", 'sha256').digest()
        
        string_to_sign = (
            "TC3-HMAC-SHA256\n"
            f"{timestamp}\n"
            f"{date}/{service}/tc3_request\n"
            f"{hashlib.sha256(canonical_request.encode('utf-8')).hexdigest()}"
        ).encode('utf-8')
        
        signature = hmac.new(key_signing, string_to_sign, 'sha256').hexdigest()
        
        return f'TC3-HMAC-SHA256 Credential={self.secret_id}/{date}/{service}/tc3_request, SignedHeaders={signed_headers}, Signature={signature}'

    def _call_api(self, params, retry=3):
        for attempt in range(retry):
            try:
                payload_data = {
                    "SourceText": params['SourceText'],
                    "Source": params['Source'],
                    "Target": params['Target'],
                    "ProjectId": params['ProjectId']
                }
                payload = json.dumps(payload_data, ensure_ascii=False)
                
                params.update({
                    'Region': self.region,
                    'Timestamp': int(time.time()),
                    'Version': self.version
                })
                
                signature = self._sign(params, payload)
                
                headers = {
                    'Authorization': signature,
                    'Content-Type': 'application/json',
                    'Host': self.endpoint,
                    'X-TC-Action': params['Action'],
                    'X-TC-Timestamp': str(params['Timestamp']),
                    'X-TC-Version': self.version,
                    'X-TC-Region': self.region
                }
                
                response = self.session.post(
                    f'https://{self.endpoint}',
                    headers=headers,
                    data=payload.encode('utf-8')
                )
                
                result = response.json()
                if 'Response' in result and 'Error' not in result['Response']:
                    return result['Response']
                tqdm.write(f'API错误: {result.get("Response", {}).get("Error", {}).get("Message")}')
                time.sleep(2 ** attempt)
            except Exception as e:
                tqdm.write(f'请求失败: {str(e)}')
                time.sleep(3)
        return {}

    def detect_language(self, text):
        params = {
            'Action': 'LanguageDetect',
            'Text': text,
            'ProjectId': 0,
            'Version': self.version
        }
        return self._call_api(params).get('Lang', 'auto')

    def translate(self, text, source='auto', target='zh'):
        params = {
            'Action': 'TextTranslate',
            'SourceText': text,
            'Source': source,
            'Target': target,
            'ProjectId': 0,
            'Version': self.version
        }
        return self._call_api(params).get('TargetText', text)

def process_subtitle(content, translator):
    try:
        if any(0x4E00 <= ord(c) <= 0x9FFF for c in content.content):
            simplified = convert(content.content, 'zh-cn')
            return srt.Subtitle(
                index=content.index,
                start=content.start,
                end=content.end,
                content=simplified
            )
        
        lang = translator.detect_language(content.content)
        if lang in ['zh', 'zh-TW']:
            return content
        
        translated = translator.translate(content.content, source=lang, target='zh')
        return srt.Subtitle(
            index=content.index,
            start=content.start,
            end=content.end,
            content=translated
        )
    except Exception as e:
        tqdm.write(f'字幕块 {content.index} 处理失败: {str(e)}')
        return content

def main():
    parser = argparse.ArgumentParser(description='腾讯云字幕转换工具')
    parser.add_argument('input', help='输入SRT文件路径')
    parser.add_argument('-o', '--output', help='输出文件路径')
    parser.add_argument('--batch', type=int, default=5, help='批量处理行数')
    secret_id = ""
    secret_key = ""
    
    args = parser.parse_args()
    
    if not args.output:
        base = os.path.splitext(args.input)[0]
        args.output = f'{base}_tc_translated.srt'

    translator = TencentTranslator(secret_id, secret_key)
    
    try:
        with open(args.input, 'r', encoding='utf-8-sig') as f:
            subs = list(srt.parse(f.read()))
        
        batches = [subs[i:i+args.batch] for i in range(0, len(subs), args.batch)]
        converted = []
        
        # 添加进度条
        with tqdm(
            total=len(batches),
            desc="🚀 翻译进度",
            unit="batch",
            bar_format="{l_bar}{bar}| {n_fmt}/{total_fmt} [{elapsed}<{remaining}]",
            dynamic_ncols=True
        ) as pbar:
            for batch_idx, batch in enumerate(batches):
                start_time = time.time()
                combined = '\n'.join([sub.content for sub in batch])
                
                try:
                    translated = translator.translate(combined)
                except Exception as e:
                    tqdm.write(f"批处理 {batch_idx} 失败: {str(e)}")
                    translated = combined  # 保留原文
                
                translated_lines = translated.split('\n')
                
                for i, sub in enumerate(batch):
                    new_content = translated_lines[i] if i < len(translated_lines) else sub.content
                    converted.append(srt.Subtitle(
                        index=sub.index,
                        start=sub.start,
                        end=sub.end,
                        content=new_content
                    ))
                
                # 更新进度信息
                pbar.update(1)
                processed = min((batch_idx+1)*args.batch, len(subs))
                pbar.set_postfix({
                    "进度": f"{processed}/{len(subs)}",
                    "批大小": args.batch,
                    "速度": f"{args.batch/(time.time()-start_time):.1f}行/秒"
                })
        
        with open(args.output, 'w', encoding='utf-8') as f:
            f.write(srt.compose(converted))
            
        print(f'\n✅ 转换完成！文件已保存至: {os.path.abspath(args.output)}')
    except Exception as e:
        tqdm.write(f'\n❌ 处理失败: {str(e)}')
        sys.exit(1)

if __name__ == '__main__':
    main()
