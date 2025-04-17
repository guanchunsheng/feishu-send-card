# feishu-send-card
通过飞书群组机器人webhook，向群组中发送卡片的脚本。

1. 配置文件conf.json：
 - webhook和secret是群组机器人的配置
 - cardID和version是卡片的配置
2. 运行脚本
  go run sendcard

参考飞书官方文档
群组机器人api：https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot
卡片发送api https://open.feishu.cn/document/feishu-cards/quick-start/send-message-cards-with-custom-bot
