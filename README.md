# dora-logs

Mọi log luôn có session = uid + time.unix (UTC).

Tao session khi start,restart app và session hết hạn khi quá 30 min từ lần bắn log cuối cùng hoặc app ở trạng thái stop.

##Log
###Log article  
1. Log active bắn khi app start, restart.
2. Log deactivate bắn khi app ở trạng thái stop.
3. Log action bắn khi click event, article.

###Log video 
unsupported.