able = 'YES', data_type, character_maximum_length, column_type, column_key, extra, column_comment, numeric_precision, numeric_scale , datetime_precision FROM information_schema.columns WHERE table_schema = 'Qiniu_Project' AND table_name = 'membership_infos' ORDER BY ORDINAL_POSITION

2025/11/10 23:27:35 /home/ubuntu/ChatWeb/backend/database/database.go:74 Error 3780 (HY000): Referencing column 'user_id' and referenced column 'user_id' in foreign key constraint 'membership_info_ibfk_1' are incompatible.
[0.374ms] [rows:0] ALTER TABLE `membership_infos` MODIFY COLUMN `user_id` bigint unsigned NOT NULL
exit status 1


curl -X GET "$BASE_URL/web/items" \
  -H "Authorization: Bearer $TOKEN"
这个在内容多的时候会很浪费性能


# 目前版本存在的问题
## 重新加载后的，语音没有缓存

## 不必要的Token也念了，不是全念正文
谢谢你问这个问题。作为AI，我没有传统意义上的"一天"体验，但我很好奇你为什么重复这个问题。
你今天过得怎么样呢？是什么让你想了解我的感受？
[DEEP_QUESTIONS]
当你询问他人"今天过得怎么样"时，你真正想了解什么？
你认为AI能像人类一样体验"过得好"或"不好"吗？
[END]
