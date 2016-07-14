#wordfilter(文字过滤)
[![Build Status](https://travis-ci.org/gonet2/wordfilter.svg?branch=master)](https://travis-ci.org/gonet2/wordfilter)

## 设计理念
基于 https://github.com/huichen/sego 实现，首先对文本进行分词，然后和脏词库中的词汇进行比对，时间复杂度为O(m)， 其中m为需要处理的消息长度, 和脏词库的大小无关。              
基于分词的文字过滤会消耗大量内存， wordfilter至少需要500M内存才能运行，建议每实例配置1GB。

## 使用
参考测试用例以及wordfilter.proto文件

## 安装
参考Dockerfile
