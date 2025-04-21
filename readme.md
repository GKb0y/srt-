# 字幕翻译工具
自动将srt字幕中的非中文简体字幕转变为中文简体

## 功能概览

 ✓ 多语言转换   : 支持英/日/韩/繁 → 简中转换
  ✓ 智能编码检测 : 自动识别UTF-8/GBK等常见编码
  x 批量处理     : 支持目录批量转换
  ✓ 格式保留     : 完整保持原有时序和样式

## 支持语言

 • 输入语言:     auto：自动识别（识别为一种语言）    zh：简体中文    zh-TW：繁体中文    en：英语    ja：日语    ko：韩语    fr：法语    es：西班牙语    it：意大利语    de：德语    tr：土耳其语    ru：俄语    pt：葡萄牙语    vi：越南语    id：印尼语    th：泰语    ms：马来西亚语    ar：阿拉伯语        hi：印地语
  • 输出语言: 简体中文(zh-cn)

## 技术亮点

  ⦿ 基于腾讯云机器翻译API
  ⦿ 专业级SRT格式解析引擎
  ⦿ 自动错误恢复机制
  ⦿ 多线程加速处理

# python版本

## 工具安装

如果需要使用.py结尾的文件，请安装python环境，并学习一下如何运行python文件，谢谢。

安装python之后，请打开命令行。win+R输入cmd

之后安装所需要的python必要模块

```
pip install -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple #-i及之后的，可以不输入
```

## 工具使用

- 如果只需要将繁体转为简体，直接到工具目录之后使用繁转简.py即可

​	请打开命令行之后输入

```
python 繁转简.py input.srt output.srt #input.srt改为需要翻译的srt文件，output.srt改为需要输出的文件名
```

- 如果需要翻译其他语言，需要用到腾讯云的ak,sk(免费)。

  1.注册一个腾讯云账户。

  2.访问用户列表页面，https://console.cloud.tencent.com/cam

  3.创建一个子用户，并赋予机器翻译的权限。(请注意只需要机器翻译，不需要给多了，一切后果请自行承担)

  ![image-20250421094736994](C:\Users\gk\AppData\Roaming\Typora\typora-user-images\image-20250421094736994.png)

  随便输入一个用户名，修改访问方式为编程访问，修改用户权限。

  ![image-20250421094906522](C:\Users\gk\AppData\Roaming\Typora\typora-user-images\image-20250421094906522.png)

  点击上图最后一个箭头指向的笔，之后在下图的搜索框中，搜索机器翻译。选择全读写访问。之后点击确定。再点击创建用户，完成创建。

  ![image-20250421094838099](C:\Users\gk\AppData\Roaming\Typora\typora-user-images\image-20250421094838099.png)

  创建好之后，将下图的SecretId和SecretKey保存好，并且不要告诉他人。

  ![image-20250421095159666](C:\Users\gk\AppData\Roaming\Typora\typora-user-images\image-20250421095159666.png)

  4.配置工具

  打开字幕转换工具.py，在第153行和第154行的双引号中输入自己的SecretId和SecretKey。

  ![image-20250421095804402](C:\Users\gk\AppData\Roaming\Typora\typora-user-images\image-20250421095804402.png)

  5.弄好配置之后，请在文件所在目录打开cmd(点击文件夹空白处，按住shift点击鼠标右键，选择在终端打开)，或打开cmd之后cd到文件所在目录。

  ```
  python 字幕翻译工具.py input.srt output.srt ##input.srt改为需要翻译的srt文件，output.srt改为需要输出的文件名
  ```

  ![image-20250421100350282](C:\Users\gk\AppData\Roaming\Typora\typora-user-images\image-20250421100350282.png)

# go语言版本

1.如果需要自己编译，请下载源码编译

2.编译好的文件根据对应版本可以直接使用

- 参数配置

  首先，参照上文的腾讯云ak，sk获取步骤，获取SecretId和SecretKey。

  之后在工具目录下的key.txt输入SecretId，SecretKey。不需要引号

  ![image-20250421100809596](C:\Users\gk\AppData\Roaming\Typora\typora-user-images\image-20250421100809596.png)

- 工具使用

  直接双击，根据指引输入。

  ![image-20250421101322427](C:\Users\gk\AppData\Roaming\Typora\typora-user-images\image-20250421101322427.png)