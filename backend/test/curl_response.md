(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ curl -X POST "$BASE_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\":\"$USERNAME\",
    \"password\":\"$PASSWORD\",
    \"full_name\":\"$FULL_NAME\",
    \"email\":\"$EMAIL\",
    \"phone_number\":\"$PHONE\",
    \"security_question1\":\"$SQ1\",
    \"security_answer1\":\"$SA1\",
    \"security_question2\":\"$SQ2\",
    \"security_answer2\":\"$SA2\"
  }"
{"message":"注册成功","success":true,"user_id":1}(FA) ubuntu@WFTGF2025:(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ TOKEN=$(curl -sS -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}" | jq -r .token) && echo "$TOKEN"
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjI4NzI3MjMsInVzZXJfaWQiOjF9.-kkDJEChejtDVUsolTkZ_JE3jy1TcCvrRtBR12CADb0
(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ curl -X GET "$BASE_URL/api/auth/me" \
  -H "Authorization: Bearer $TOKEN"
{"user_id":1,"username":"testuser","full_name":"Test User","email":"testuser@example.com","phone_number":"13900001234","created_at":"2025-11-10T22:51:48.721+08:00","updated_at":"2025-11-10T22:51:48.721+08:00"}(FA)(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ USER_ID=$(curl -sS -X GET "$BASE_URL/api/auth/me" \
  -H "Authorization: Bearer $TOKEN" | jq -r '.user_id // .data.user_id') && echo "$USER_ID"
1
(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ curl -X PUT "$BASE_URL/api/users/$USER_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"full_name\":\"${FULL_NAME}_updated\",
    \"email\":\"updated_${EMAIL}\",
    \"phone_number\":\"$PHONE\"
  }"
(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ curl -X POST "$BASE_URL/api/auth/verify-security" \POST "$BASE_URL/api/auth/verify-security" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"security_answer1\":\"$SA1\",\"security_answer2\":\"$SA2\"}"
{"reset_token":"02ba320129cbfa1da4b5795f1cf38a16","success":true}(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ RESET_RESET_TOKEN=$(curl -sS -X POST "$BASE_URL/api/auth/verify-security" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"security_answer1\":\"$SA1\",\"security_answer2\":\"$SA2\"}" | jq -r '.reset_token') && echo "$RESET_TOKEN"
24d8f04d485bc07691b0b3f522bd47c4
(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ curl -X POST "$BASE_URL/api/auth/reset-password" \
  -H "Content-Type: application/json" \
  -d "{\"reset_token\":\"$RESET_TOKEN\",\"new_password\":\"$NEW_PASSWORD\"}"
(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ TOKEN=$(curl -sS -X POST "$BASE_URL/api/auth/login" \sS -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$NEW_PASSWORD\"}" | jq -r .token) && echo "$TOKEN"
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjI4NzI4MTgsInVzZXJfaWQiOjF9.sV6V0aW-YmB1Brw4xFh4cnXnp3ZEpXToCxLGhPkI3go
(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ curl -X GET "$BASE_URL/api/membership" \
  -H "Authorization: Bearer $TOKEN"
(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ curl -X POST "$BASE_URL/api/membership" \ \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"user_id\": $USER_ID,
    \"start_date\": \"2025-01-01\",
    \"expire_date\": \"2026-01-01\",
    \"status\": \"active\"
  }"
{"membership_id":1,"message":"会员信息已创建","success":true}(FA) ubunt(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ ST "$BASGET "$BASE_URL/api/membership"   -H "Authorization: Bearer $TOKEN"
[{"membership_id":1,"user_id":1,"start_date":"2025-01-01T00:00:00+08:00","expire_date":"2026-01-01T00:00:00+08:00","status":"active"}](FA) ubu(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ curl -X GET "$BASE_URL/api/membership/$USER_ID" \
  -H "Authorization: Bearer $TOKEN"
{"membership_id":1,"user_id":1,"start_date":"2025-01-01T00:00:00+08:00","expire_date":"2026-01-01T00:00:00+08:00","status":"active"}(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ export MEMBERSHIP_ID=1
(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ curl -X PUT "$BASE_URL/api/membership/$MEMBERSHIP_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"expire_date":"2027-01-01","status":"active"}'
(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ curl -X POST "$BASE_URL/api/membership/orders" \-X POST "$BASE_URL/api/membership/orders" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"user_id\": $USER_ID,
    \"duration_months\": 12,
    \"amount\": 199.00,
    \"payment_method\": \"other\"
  }"
{"message":"订单已创建","order_id":1,"success":true}(FA) ubuntu@WFTGF20(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ curl -X GET "$BASE_URL/api/membership/orders/$USER_ID" \
  -H "Authorization: Bearer $TOKEN"
[{"order_id":1,"user_id":1,"purchase_date":"2025-11-10T22:54:35.474+08:00","duration_months":12,"amount":199,"payment_method":"other"}](FA) ub(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ curl -X GET "$BASE_URL/api/membership/orders/$USER_ID/latest" \
  -H "Authorization: Bearer $TOKEN"
{"order_id":1,"user_id":1,"purchase_date":"2025-11-10T22:54:35.474+08:00","duration_months":12,"amount":199,"payment_method":"other"}(FA) ubun(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ curl -X GET "$BASE_URL/api/membership/orders/$USER_ID/recent?n=3" \
  -H "Authorization: Bearer $TOKEN"
[{"order_id":1,"user_id":1,"purchase_date":"2025-11-10T22:54:35.474+08:00","duration_months":12,"amount":199,"payment_method":"other"}](FA) ub(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ curl -X POST "$BASE_URL/web/items" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"url":"https://www.njupt.edu.cn/","title":"Example Home","fetch":true}'
{"id":1,"user_id":1,"url":"https://www.njupt.edu.cn/","title":"Example Home","content":"旧版入口 English 电子邮件 首页 南邮概况 \u003e 学校简介 学校章程 南邮精神 校标校训 南邮校史 南邮校歌 现任领导 视频展播 校园实景漫... 校园景色 校区地图 内设机构 \u003e 党政群部门 教学机构 基层党的组... 科研机构 直属单位和... 独立学院 学科建设 科学研究 \u003e 自然科学研... 社会科学研... 高等教育研... 科技基础条... 学术刊物 招生就业 \u003e 本科招生 研究生招生 留学生招生 继续教育招... 就业信息网 信息公开 人才培养 \u003e 本科 研究生 留学生 继续教育 实验教学示... 教学质量监... 智慧校园 通知通告 关闭 NJUPT 查看更多 2024 06.25 信长星书记调研南京邮电大学思政课建设 6月13日，省委书记信长星到南京邮电大学调研思政课建设工作。他指出，要深入学习贯彻习近平总书记关于高校思想政治工作的重要论述和对学校思政课建设的重要指示精神，全面贯彻党的教育方针，落实立德树人根本任务，守正创新推动思政课建设内涵式发展，努力培养更多让党放心、爱国奉献、担当民族复兴重任的时代新人。 因邮电而生，随通信而长，由信息而强，南京邮电大学历史悠久、特色鲜明。信长星走进校史馆，深入了解学校历史沿革、建设发展、人才培养等情况。他指出，南邮在抗日烽火中淬炼而生，在时代大潮中发展壮大，具有鲜明的红色基因、深厚的家国情怀，培养了一大批优秀人才。要加强校史资料的挖掘、整理和研究，结合最新办学实践，因地制宜开展思政教育，激励广大师生发扬优良传统、赓续红色血脉。在与学校辅导员、学生代表交流时，信长星说，培养时代新人，既要提升理论素养，也要涵养人文精神。办好思政课，既要老师用心教，也要学生用心悟。要注重提炼挖掘各门课程的思政教育元素，把做人做事的基本道理，社会主义核心价值观的要求，强国建设、民族复兴的理想责任等融入课程教学之中，把思政小课堂和社会大课堂结合起来，推进课程思政与思政课程同向同行 2025 08.15 南京邮电大学位居2025软科世界大学学术排名全球第376名 8月15日，高等教育评价机构软科发布“2025软科世界大学学术排名”，展示了全球领先的1000所研究型大学，中国内地共有222所大学上榜。南京邮电大学位居全球排名第376名，相较2024年上升82名，位列国内77名。 软科世界大学学术排名于2003年首次发布，是世界范围内首个综合性的全球大学排名，已成为全球最具影响力和权威性的大学排名之一。该排名全部采用国际可比的客观指标和第三方数据，包括获诺贝尔奖和菲尔兹奖的校友和教师数、高被引科学家数、在Nature和Science上发表的论文数、被Web of Science科学引文索引（SCI）和社会科学引文索引（SSCI）收录的论文数、师均学术表现等。软科世界大学学术排名每年排名的全球大学超过2500所，发布最为领先的前1000所大学名单。 点击浏览：2025软科世界大学学术排名完整榜单https://www.shanghairanking.cn/rankings/arwu/2025（撰稿：唐澄澄 初审：胡纵宇 编辑：王存宏 审核：张丰） 2025 03.22 南邮连续五年荣获江苏省综合考核第一等次 3月21日，江苏省委、省政府召开全省2024年度高质量发展总结暨2025年工作推进会议，省委书记信长星出席会议并讲话。会议通报了全省2024年度高质量发展综合考核结果，南邮获评2024年度省属高校高质量发展综合考核第一等次。这也是学校自2020年以来，连续第五次荣膺第一等次。 2024年，学校坚持以习近平新时代中国特色社会主义思想为指导，深入学习贯彻习近平总书记对江苏工作重要讲话精神，坚持社会主义办学方向，落实立德树人根本任务，扎实推进国家“双一流”和江苏高水平大学建设，在全面从严治党和内涵建设方面取得了突出成效。学校扎实开展党纪学习教育，自觉接受省委巡视和审计；牵头重组并成功获批柔性电子全国重点实验室，参与重组并成功获批藏语智能全国重点实验室；新增2个博士学位授权点，6个学科进入ESI全球排名前1%，组建未来信息学科交叉中心，入选省首批学科交叉中心试点建设高校；国家自然科学基金项目连续三年突破百项，国家社科基金项目再创历史新高，引育欧洲科学院院士、国家级人才等7人；开展“AI+创新人才培养”行动，获批省首批集成电路学院、省首批卓越工程师学院，高质量完成本科教育教学审核评估考察。 2025 04.03 国家邮政局局长赵冲久一行来南邮调研 4月2日上午，交通运输部党组成员、国家邮政局党组书记、局长赵冲久来南邮调研。国家邮政局市场监管司司长林虎，人事司司长孙广明，江苏省邮政管理局党组书记、局长、省交通运输厅副厅长蒋波，学校党委书记郭宇锋，党委常委、副校长周亮参加调研和座谈。座谈会由校长叶美兰主持。座谈会现场 郭宇锋书记在致辞中表示，此次调研是对南邮服务邮政事业发展的全面检阅，更是对进一步融入国家邮政战略的激励与鞭策。学校将以此次调研为契机，充分发挥信息科技与邮政物流交叉融合的学科优势，聚焦行业“卡脖子”技术攻关，助力邮政领域核心技术突破；深化协同育人，打造“政产学研用”一体化人才培养高地，为行业高质量发展提供智力支撑；围绕“快递进村”“快递出海”等国家战略需求，以科技赋能行业绿色低碳转型和国际化布局。 座谈会上，现代邮政学院围绕学院概况、党建工作、师资队伍、学科专业、科研创新、人才培养等进行汇报。学院师生代表分别就邮政人才培养、邮政红色基因传承进行发言。 赵冲久局长充分肯定南邮在邮政快递领域为高等教育作出的努力与贡献。他指出，邮政快递业作为国民经济的重要支撑，2024年业务收入已突破1.4万亿，亟需以高质量教育科技人 2025 04.02 南邮4个项目入选教育部2025年度高校思想政治工作质量提升综合改革与精品建设项目 近日，教育部思想政治工作司发布了2025年度高校思想政治工作质量提升综合改革与精品建设项目遴选结果，南邮共有4个项目入选，入选项目数量位列江苏高校第二、省属高校第一，这是学校在教育部高校思想政治工作质量提升综合改革与精品建设项目上的重要突破，也是学校守正创新推动思想政治工作改革发展的重大成果。 教育部高校思想政治工作质量提升综合改革与精品建设项目是不断加强和改进新时代学校思想政治教育、建设高质量思想政治工作体系的重要举措。学校按照教育部、江苏省项目建设要求和标准，结合思政工作实际和育人特色，精心培育校级思想政治工作有关建设项目，切实推动思想政治工作队伍在育人实践中总结经验规律，在改革创新中增强工作质效，形成了一批体系健全、特色鲜明、育人成效明显的具有示范作用的典型做法和经验。 学校将坚持不懈用习近平新时代中国特色社会主义思想铸魂育人，全面落实全国教育大会精神，深入实施新时代立德树人工程，以建好教育部、江苏省、校级项目为抓手，不断推动“时代新人铸魂工程”提质升级，在构建固本铸魂的思想政治教育体系、培养堪当民族复兴重任的时代新人这一使命任务上提供南邮方案、贡献南邮智慧。（撰稿：王波 01 2025.11.04 南邮获“挑战杯”国赛特等奖7项 再捧“优胜杯” 2025.11.10 南邮MBA项目获评江苏省首个QS Stars五星认证 2025.11.07 南邮举办研究生“鼎山研习营” 专题学习党的二十届四中全会精神 2025.11.07 南邮两门课程出海泰国和印尼国家慕课平台 2025.11.06 南邮党委理论学习中心组召开党的二十届四中全会精神专题学习会 省委教育工委列席旁听工作第三组到会指导 2025.11.05 无锡市惠山区区长赵强一行来南邮调研交流 查看更多 【江苏卫视·江苏新时空】信长星调研高校思政课建设时指出 坚持把立德树人作为根本任务 守正创新推动思政课建设内涵式发展 【新华网客户端】2025年电子信息特色高校发展大会在南京邮电大学举行 【中国教育电视台·全国教育新闻联播】悄悄给饭卡打钱 被这些高校暖到了！ 【江苏教育报】奏响教育强省建设的南邮奋进曲 【光明日报客户端】南京邮电大学：于“热血铸荣光”中厚植青春信仰 【新华社客户端】我国科研人员在无机光伏材料薄膜化领域取得新进展 预告 报道 查看更多 话剧《华罗庚的最后讲演》 2025-11-10 19:00:00 安恒青春剧场 2025年“读经典·颂中华”朗诵比赛 2025-11-07 19:00:00 大学生活动中心安恒青春剧场 美育大讲堂邀您共赏“聆听音乐——为何音乐与音乐何为” 2025-11-10 15:30:00 运满满报告厅（图书馆四楼） 医工交叉院士专家报告会 2025-10-31 09:00:00 行政北楼中天报告厅 人工智能与科研管理专题培训会 2025-10-28 13:30:10 仙林校区中兴报告厅 跨越时空的对话——南京邮电大学红色校史剧《赤子》复演 2025-10-21 19:30:00 大学生活动中心安恒青春剧场 美育大讲堂邀您共赏“《牡丹亭》——从临川笔下到昆曲场上” 2025-10-22 14:00:00 运满满报告厅（图书馆四楼） 2025智能传播与艺术国际研讨会 2025-10-18 14:00:00 计算机学科楼报告厅 查看更多 南邮承办2025江苏高校国际产学研用合作交流周“能源与新材料”分论坛 南邮主办第十七届无线通信与信号处理国际会议 南邮联合主办“中国社会学自主知识体系与中国式现代化”学术座谈会 跨越时空的对话——南邮红色校史剧《赤子》复演 南邮举办第五十六届田径运动会 南邮主办第七届“强邮论坛”暨2025年空地协同低空物流技术学术会议 南邮举办“2025年智能传播与艺术国际研讨会” 南邮举办“热血铸荣光”——弘扬伟大抗战精神 赓续传承红色基因主题文艺展演 南邮举行2025年“邮启新程”主题晚会暨第三十届大学生科技节开幕式 南邮联合主办第18届IEEE毫米波与太赫兹技术联合会议 yyyy MM.dd title text 1 南邮教师获第二届全国电子信息类专业高校教师智慧教学案例竞赛一等奖和最佳创意奖 2 Nature Communications刊发汪联辉教授和王婷教授团队在柔性电化学传感领域的最新研究进展 3 Nature Energy刊发黄维院士辛颢教授铜锌锡硫硒薄膜太阳能最新研究进展 4 南邮教师获全国高校教师教学创新大赛新文科赛道国奖 查看更多 2025 10.29 南邮教师获第二届全国电子信息类专业高校教师智慧教学案例竞赛一等奖和最佳创意奖 2025 10.10 Nature Communications刊发汪联辉教授和王婷教授团队在柔性电化学传感领域的最新研究进展 2025 09.15 Nature Energy刊发黄维院士辛颢教授铜锌锡硫硒薄膜太阳能最新研究进展 2025 08.20 南邮教师获全国高校教师教学创新大赛新文科赛道国奖 查看更多 查看更多 优秀教师 优秀教育工作者 最美辅导员 最美大学生 查看更多 查看更多 微信 视频号 跨越六十载，八旬教授曹伟的求索与守望 美育大讲堂——昆曲《牡丹亭》里的文化传承与美学觉醒 这个系列怎么能少的了南邮？ 在南邮，一起赴“桂”！ 2025.10.31 同心向邮，双倍精彩！ 2025.10.30 青年科学家 | 扎根南邮，他让无线网络更智能！ 2025.10.25 全省体测第一高校的运动会来了！ 2025.10.24 南邮师生热议党的二十届四中全会 科学谋划“十五五”发展新篇章 7 博士后流动站7个 10 一级学科博士学位授权点10个 11 博士学位授权点11个 23 一级学科硕士学位授权点23个 14 硕士专业学位授权点（类别）14个 52 本科专业52个 6 6个学科进入ESI学科排名全球前1% 27 国家一流专业27个 14 14个专业通过国家工程教育专业认证 校历查询 规章制度 服务流程 办公信息 书记信箱 校长信箱 图书馆 艺术馆 南邮校报 VPN登录 实景漫游 校友会 基金会 相关链接 仙林校区地址：南京市仙林大学城文苑路9号 邮编：210023 三牌楼校区地址：南京市新模范马路66号 邮编：210003 锁金村校区地址：南京市龙蟠路177号 邮编：210042 联系电话:（86）-25-85866888 传真:（86）-25-85866999 邮箱:njupt@njupt.edu.cn 苏公网安备32011302320419号 |苏ICP备11073489号-1 版权所有：南京邮电大学","created_at":"2025-11-10T22:54:55.49(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ curl -X GET "$BASE_URL/web/items" \tWeb/backend/test$ curl -X GET "$BASE_URL/web/items" \
  -H "Authorization: Bearer $TOKEN"
{"items":[{"id":1,"user_id":1,"url":"https://www.njupt.edu.cn/","title":"Example Home","content":"旧版入口 English 电子邮件 首页 南邮概况 \u003e 学校简介 学校章程 南邮精神 校标校训 南邮校史 南邮校歌 现任领导 视频展播 校园实景漫... 校园景色 校区地图 内设机构 \u003e 党政群部门 教学机构 基层党的组... 科研机构 直属单位和... 独立学院 学科建设 科学研究 \u003e 自然科学研... 社会科学研... 高等教育研... 科技基础条... 学术刊物 招生就业 \u003e 本科招生 研究生招生 留学生招生 继续教育招... 就业信息网 信息公开 人才培养 \u003e 本科 研究生 留学生 继续教育 实验教学示... 教学质量监... 智慧校园 通知通告 关闭 NJUPT 查看更多 2024 06.25 信长星书记调研南京邮电大学思政课建设 6月13日，省委书记信长星到南京邮电大学调研思政课建设工作。他指出，要深入学习贯彻习近平总书记关于高校思想政治工作的重要论述和对学校思政课建设的重要指示精神，全面贯彻党的教育方针，落实立德树人根本任务，守正创新推动思政课建设内涵式发展，努力培养更多让党放心、爱国奉献、担当民族复兴重任的时代新人。 因邮电而生，随通信而长，由信息而强，南京邮电大学历史悠久、特色鲜明。信长星走进校史馆，深入了解学校历史沿革、建设发展、人才培养等情况。他指出，南邮在抗日烽火中淬炼而生，在时代大潮中发展壮大，具有鲜明的红色基因、深厚的家国情怀，培养了一大批优秀人才。要加强校史资料的挖掘、整理和研究，结合最新办学实践，因地制宜开展思政教育，激励广大师生发扬优良传统、赓续红色血脉。在与学校辅导员、学生代表交流时，信长星说，培养时代新人，既要提升理论素养，也要涵养人文精神。办好思政课，既要老师用心教，也要学生用心悟。要注重提炼挖掘各门课程的思政教育元素，把做人做事的基本道理，社会主义核心价值观的要求，强国建设、民族复兴的理想责任等融入课程教学之中，把思政小课堂和社会大课堂结合起来，推进课程思政与思政课程同向同行 2025 08.15 南京邮电大学位居2025软科世界大学学术排名全球第376名 8月15日，高等教育评价机构软科发布“2025软科世界大学学术排名”，展示了全球领先的1000所研究型大学，中国内地共有222所大学上榜。南京邮电大学位居全球排名第376名，相较2024年上升82名，位列国内77名。 软科世界大学学术排名于2003年首次发布，是世界范围内首个综合性的全球大学排名，已成为全球最具影响力和权威性的大学排名之一。该排名全部采用国际可比的客观指标和第三方数据，包括获诺贝尔奖和菲尔兹奖的校友和教师数、高被引科学家数、在Nature和Science上发表的论文数、被Web of Science科学引文索引（SCI）和社会科学引文索引（SSCI）收录的论文数、师均学术表现等。软科世界大学学术排名每年排名的全球大学超过2500所，发布最为领先的前1000所大学名单。 点击浏览：2025软科世界大学学术排名完整榜单https://www.shanghairanking.cn/rankings/arwu/2025（撰稿：唐澄澄 初审：胡纵宇 编辑：王存宏 审核：张丰） 2025 03.22 南邮连续五年荣获江苏省综合考核第一等次 3月21日，江苏省委、省政府召开全省2024年度高质量发展总结暨2025年工作推进会议，省委书记信长星出席会议并讲话。会议通报了全省2024年度高质量发展综合考核结果，南邮获评2024年度省属高校高质量发展综合考核第一等次。这也是学校自2020年以来，连续第五次荣膺第一等次。 2024年，学校坚持以习近平新时代中国特色社会主义思想为指导，深入学习贯彻习近平总书记对江苏工作重要讲话精神，坚持社会主义办学方向，落实立德树人根本任务，扎实推进国家“双一流”和江苏高水平大学建设，在全面从严治党和内涵建设方面取得了突出成效。学校扎实开展党纪学习教育，自觉接受省委巡视和审计；牵头重组并成功获批柔性电子全国重点实验室，参与重组并成功获批藏语智能全国重点实验室；新增2个博士学位授权点，6个学科进入ESI全球排名前1%，组建未来信息学科交叉中心，入选省首批学科交叉中心试点建设高校；国家自然科学基金项目连续三年突破百项，国家社科基金项目再创历史新高，引育欧洲科学院院士、国家级人才等7人；开展“AI+创新人才培养”行动，获批省首批集成电路学院、省首批卓越工程师学院，高质量完成本科教育教学审核评估考察。 2025 04.03 国家邮政局局长赵冲久一行来南邮调研 4月2日上午，交通运输部党组成员、国家邮政局党组书记、局长赵冲久来南邮调研。国家邮政局市场监管司司长林虎，人事司司长孙广明，江苏省邮政管理局党组书记、局长、省交通运输厅副厅长蒋波，学校党委书记郭宇锋，党委常委、副校长周亮参加调研和座谈。座谈会由校长叶美兰主持。座谈会现场 郭宇锋书记在致辞中表示，此次调研是对南邮服务邮政事业发展的全面检阅，更是对进一步融入国家邮政战略的激励与鞭策。学校将以此次调研为契机，充分发挥信息科技与邮政物流交叉融合的学科优势，聚焦行业“卡脖子”技术攻关，助力邮政领域核心技术突破；深化协同育人，打造“政产学研用”一体化人才培养高地，为行业高质量发展提供智力支撑；围绕“快递进村”“快递出海”等国家战略需求，以科技赋能行业绿色低碳转型和国际化布局。 座谈会上，现代邮政学院围绕学院概况、党建工作、师资队伍、学科专业、科研创新、人才培养等进行汇报。学院师生代表分别就邮政人才培养、邮政红色基因传承进行发言。 赵冲久局长充分肯定南邮在邮政快递领域为高等教育作出的努力与贡献。他指出，邮政快递业作为国民经济的重要支撑，2024年业务收入已突破1.4万亿，亟需以高质量教育科技人 2025 04.02 南邮4个项目入选教育部2025年度高校思想政治工作质量提升综合改革与精品建设项目 近日，教育部思想政治工作司发布了2025年度高校思想政治工作质量提升综合改革与精品建设项目遴选结果，南邮共有4个项目入选，入选项目数量位列江苏高校第二、省属高校第一，这是学校在教育部高校思想政治工作质量提升综合改革与精品建设项目上的重要突破，也是学校守正创新推动思想政治工作改革发展的重大成果。 教育部高校思想政治工作质量提升综合改革与精品建设项目是不断加强和改进新时代学校思想政治教育、建设高质量思想政治工作体系的重要举措。学校按照教育部、江苏省项目建设要求和标准，结合思政工作实际和育人特色，精心培育校级思想政治工作有关建设项目，切实推动思想政治工作队伍在育人实践中总结经验规律，在改革创新中增强工作质效，形成了一批体系健全、特色鲜明、育人成效明显的具有示范作用的典型做法和经验。 学校将坚持不懈用习近平新时代中国特色社会主义思想铸魂育人，全面落实全国教育大会精神，深入实施新时代立德树人工程，以建好教育部、江苏省、校级项目为抓手，不断推动“时代新人铸魂工程”提质升级，在构建固本铸魂的思想政治教育体系、培养堪当民族复兴重任的时代新人这一使命任务上提供南邮方案、贡献南邮智慧。（撰稿：王波 01 2025.11.04 南邮获“挑战杯”国赛特等奖7项 再捧“优胜杯” 2025.11.10 南邮MBA项目获评江苏省首个QS Stars五星认证 2025.11.07 南邮举办研究生“鼎山研习营” 专题学习党的二十届四中全会精神 2025.11.07 南邮两门课程出海泰国和印尼国家慕课平台 2025.11.06 南邮党委理论学习中心组召开党的二十届四中全会精神专题学习会 省委教育工委列席旁听工作第三组到会指导 2025.11.05 无锡市惠山区区长赵强一行来南邮调研交流 查看更多 【江苏卫视·江苏新时空】信长星调研高校思政课建设时指出 坚持把立德树人作为根本任务 守正创新推动思政课建设内涵式发展 【新华网客户端】2025年电子信息特色高校发展大会在南京邮电大学举行 【中国教育电视台·全国教育新闻联播】悄悄给饭卡打钱 被这些高校暖到了！ 【江苏教育报】奏响教育强省建设的南邮奋进曲 【光明日报客户端】南京邮电大学：于“热血铸荣光”中厚植青春信仰 【新华社客户端】我国科研人员在无机光伏材料薄膜化领域取得新进展 预告 报道 查看更多 话剧《华罗庚的最后讲演》 2025-11-10 19:00:00 安恒青春剧场 2025年“读经典·颂中华”朗诵比赛 2025-11-07 19:00:00 大学生活动中心安恒青春剧场 美育大讲堂邀您共赏“聆听音乐——为何音乐与音乐何为” 2025-11-10 15:30:00 运满满报告厅（图书馆四楼） 医工交叉院士专家报告会 2025-10-31 09:00:00 行政北楼中天报告厅 人工智能与科研管理专题培训会 2025-10-28 13:30:10 仙林校区中兴报告厅 跨越时空的对话——南京邮电大学红色校史剧《赤子》复演 2025-10-21 19:30:00 大学生活动中心安恒青春剧场 美育大讲堂邀您共赏“《牡丹亭》——从临川笔下到昆曲场上” 2025-10-22 14:00:00 运满满报告厅（图书馆四楼） 2025智能传播与艺术国际研讨会 2025-10-18 14:00:00 计算机学科楼报告厅 查看更多 南邮承办2025江苏高校国际产学研用合作交流周“能源与新材料”分论坛 南邮主办第十七届无线通信与信号处理国际会议 南邮联合主办“中国社会学自主知识体系与中国式现代化”学术座谈会 跨越时空的对话——南邮红色校史剧《赤子》复演 南邮举办第五十六届田径运动会 南邮主办第七届“强邮论坛”暨2025年空地协同低空物流技术学术会议 南邮举办“2025年智能传播与艺术国际研讨会” 南邮举办“热血铸荣光”——弘扬伟大抗战精神 赓续传承红色基因主题文艺展演 南邮举行2025年“邮启新程”主题晚会暨第三十届大学生科技节开幕式 南邮联合主办第18届IEEE毫米波与太赫兹技术联合会议 yyyy MM.dd title text 1 南邮教师获第二届全国电子信息类专业高校教师智慧教学案例竞赛一等奖和最佳创意奖 2 Nature Communications刊发汪联辉教授和王婷教授团队在柔性电化学传感领域的最新研究进展 3 Nature Energy刊发黄维院士辛颢教授铜锌锡硫硒薄膜太阳能最新研究进展 4 南邮教师获全国高校教师教学创新大赛新文科赛道国奖 查看更多 2025 10.29 南邮教师获第二届全国电子信息类专业高校教师智慧教学案例竞赛一等奖和最佳创意奖 2025 10.10 Nature Communications刊发汪联辉教授和王婷教授团队在柔性电化学传感领域的最新研究进展 2025 09.15 Nature Energy刊发黄维院士辛颢教授铜锌锡硫硒薄膜太阳能最新研究进展 2025 08.20 南邮教师获全国高校教师教学创新大赛新文科赛道国奖 查看更多 查看更多 优秀教师 优秀教育工作者 最美辅导员 最美大学生 查看更多 查看更多 微信 视频号 跨越六十载，八旬教授曹伟的求索与守望 美育大讲堂——昆曲《牡丹亭》里的文化传承与美学觉醒 这个系列怎么能少的了南邮？ 在南邮，一起赴“桂”！ 2025.10.31 同心向邮，双倍精彩！ 2025.10.30 青年科学家 | 扎根南邮，他让无线网络更智能！ 2025.10.25 全省体测第一高校的运动会来了！ 2025.10.24 南邮师生热议党的二十届四中全会 科学谋划“十五五”发展新篇章 7 博士后流动站7个 10 一级学科博士学位授权点10个 11 博士学位授权点11个 23 一级学科硕士学位授权点23个 14 硕士专业学位授权点（类别）14个 52 本科专业52个 6 6个学科进入ESI学科排名全球前1% 27 国家一流专业27个 14 14个专业通过国家工程教育专业认证 校历查询 规章制度 服务流程 办公信息 书记信箱 校长信箱 图书馆 艺术馆 南邮校报 VPN登录 实景漫游 校友会 基金会 相关链接 仙林校区地址：南京市仙林大学城文苑路9号 邮编：210023 三牌楼校区地址：南京市新模范马路66号 邮编：210003 锁金村校区地址：南京市龙蟠路177号 邮编：210042 联系电话:（86）-25-85866888 传真:（86）-25-85866999 邮箱:njupt@njupt.edu.cn 苏公网安备32011302320419号 |苏ICP备11073489号-1 版权所有：南京邮电大学","created_at":"2025-11-10T22:54:55.491+08:00","updated_at":"2025-11-10T22:54:55.491+08:00"}]}(FA) ubuntu@W(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ export PAGE_ID=1
curl -X GET "$BASE_URL/web/page/$PAGE_ID" \
  -H "Authorization: Bearer $TOKEN"
{"id":1,"user_id":1,"url":"https://www.njupt.edu.cn/","title":"Example Home","content":"旧版入口 English 电子邮件 首页 南邮概况 \u003e 学校简介 学校章程 南邮精神 校标校训 南邮校史 南邮校歌 现任领导 视频展播 校园实景漫... 校园景色 校区地图 内设机构 \u003e 党政群部门 教学机构 基层党的组... 科研机构 直属单位和... 独立学院 学科建设 科学研究 \u003e 自然科学研... 社会科学研... 高等教育研... 科技基础条... 学术刊物 招生就业 \u003e 本科招生 研究生招生 留学生招生 继续教育招... 就业信息网 信息公开 人才培养 \u003e 本科 研究生 留学生 继续教育 实验教学示... 教学质量监... 智慧校园 通知通告 关闭 NJUPT 查看更多 2024 06.25 信长星书记调研南京邮电大学思政课建设 6月13日，省委书记信长星到南京邮电大学调研思政课建设工作。他指出，要深入学习贯彻习近平总书记关于高校思想政治工作的重要论述和对学校思政课建设的重要指示精神，全面贯彻党的教育方针，落实立德树人根本任务，守正创新推动思政课建设内涵式发展，努力培养更多让党放心、爱国奉献、担当民族复兴重任的时代新人。 因邮电而生，随通信而长，由信息而强，南京邮电大学历史悠久、特色鲜明。信长星走进校史馆，深入了解学校历史沿革、建设发展、人才培养等情况。他指出，南邮在抗日烽火中淬炼而生，在时代大潮中发展壮大，具有鲜明的红色基因、深厚的家国情怀，培养了一大批优秀人才。要加强校史资料的挖掘、整理和研究，结合最新办学实践，因地制宜开展思政教育，激励广大师生发扬优良传统、赓续红色血脉。在与学校辅导员、学生代表交流时，信长星说，培养时代新人，既要提升理论素养，也要涵养人文精神。办好思政课，既要老师用心教，也要学生用心悟。要注重提炼挖掘各门课程的思政教育元素，把做人做事的基本道理，社会主义核心价值观的要求，强国建设、民族复兴的理想责任等融入课程教学之中，把思政小课堂和社会大课堂结合起来，推进课程思政与思政课程同向同行 2025 08.15 南京邮电大学位居2025软科世界大学学术排名全球第376名 8月15日，高等教育评价机构软科发布“2025软科世界大学学术排名”，展示了全球领先的1000所研究型大学，中国内地共有222所大学上榜。南京邮电大学位居全球排名第376名，相较2024年上升82名，位列国内77名。 软科世界大学学术排名于2003年首次发布，是世界范围内首个综合性的全球大学排名，已成为全球最具影响力和权威性的大学排名之一。该排名全部采用国际可比的客观指标和第三方数据，包括获诺贝尔奖和菲尔兹奖的校友和教师数、高被引科学家数、在Nature和Science上发表的论文数、被Web of Science科学引文索引（SCI）和社会科学引文索引（SSCI）收录的论文数、师均学术表现等。软科世界大学学术排名每年排名的全球大学超过2500所，发布最为领先的前1000所大学名单。 点击浏览：2025软科世界大学学术排名完整榜单https://www.shanghairanking.cn/rankings/arwu/2025（撰稿：唐澄澄 初审：胡纵宇 编辑：王存宏 审核：张丰） 2025 03.22 南邮连续五年荣获江苏省综合考核第一等次 3月21日，江苏省委、省政府召开全省2024年度高质量发展总结暨2025年工作推进会议，省委书记信长星出席会议并讲话。会议通报了全省2024年度高质量发展综合考核结果，南邮获评2024年度省属高校高质量发展综合考核第一等次。这也是学校自2020年以来，连续第五次荣膺第一等次。 2024年，学校坚持以习近平新时代中国特色社会主义思想为指导，深入学习贯彻习近平总书记对江苏工作重要讲话精神，坚持社会主义办学方向，落实立德树人根本任务，扎实推进国家“双一流”和江苏高水平大学建设，在全面从严治党和内涵建设方面取得了突出成效。学校扎实开展党纪学习教育，自觉接受省委巡视和审计；牵头重组并成功获批柔性电子全国重点实验室，参与重组并成功获批藏语智能全国重点实验室；新增2个博士学位授权点，6个学科进入ESI全球排名前1%，组建未来信息学科交叉中心，入选省首批学科交叉中心试点建设高校；国家自然科学基金项目连续三年突破百项，国家社科基金项目再创历史新高，引育欧洲科学院院士、国家级人才等7人；开展“AI+创新人才培养”行动，获批省首批集成电路学院、省首批卓越工程师学院，高质量完成本科教育教学审核评估考察。 2025 04.03 国家邮政局局长赵冲久一行来南邮调研 4月2日上午，交通运输部党组成员、国家邮政局党组书记、局长赵冲久来南邮调研。国家邮政局市场监管司司长林虎，人事司司长孙广明，江苏省邮政管理局党组书记、局长、省交通运输厅副厅长蒋波，学校党委书记郭宇锋，党委常委、副校长周亮参加调研和座谈。座谈会由校长叶美兰主持。座谈会现场 郭宇锋书记在致辞中表示，此次调研是对南邮服务邮政事业发展的全面检阅，更是对进一步融入国家邮政战略的激励与鞭策。学校将以此次调研为契机，充分发挥信息科技与邮政物流交叉融合的学科优势，聚焦行业“卡脖子”技术攻关，助力邮政领域核心技术突破；深化协同育人，打造“政产学研用”一体化人才培养高地，为行业高质量发展提供智力支撑；围绕“快递进村”“快递出海”等国家战略需求，以科技赋能行业绿色低碳转型和国际化布局。 座谈会上，现代邮政学院围绕学院概况、党建工作、师资队伍、学科专业、科研创新、人才培养等进行汇报。学院师生代表分别就邮政人才培养、邮政红色基因传承进行发言。 赵冲久局长充分肯定南邮在邮政快递领域为高等教育作出的努力与贡献。他指出，邮政快递业作为国民经济的重要支撑，2024年业务收入已突破1.4万亿，亟需以高质量教育科技人 2025 04.02 南邮4个项目入选教育部2025年度高校思想政治工作质量提升综合改革与精品建设项目 近日，教育部思想政治工作司发布了2025年度高校思想政治工作质量提升综合改革与精品建设项目遴选结果，南邮共有4个项目入选，入选项目数量位列江苏高校第二、省属高校第一，这是学校在教育部高校思想政治工作质量提升综合改革与精品建设项目上的重要突破，也是学校守正创新推动思想政治工作改革发展的重大成果。 教育部高校思想政治工作质量提升综合改革与精品建设项目是不断加强和改进新时代学校思想政治教育、建设高质量思想政治工作体系的重要举措。学校按照教育部、江苏省项目建设要求和标准，结合思政工作实际和育人特色，精心培育校级思想政治工作有关建设项目，切实推动思想政治工作队伍在育人实践中总结经验规律，在改革创新中增强工作质效，形成了一批体系健全、特色鲜明、育人成效明显的具有示范作用的典型做法和经验。 学校将坚持不懈用习近平新时代中国特色社会主义思想铸魂育人，全面落实全国教育大会精神，深入实施新时代立德树人工程，以建好教育部、江苏省、校级项目为抓手，不断推动“时代新人铸魂工程”提质升级，在构建固本铸魂的思想政治教育体系、培养堪当民族复兴重任的时代新人这一使命任务上提供南邮方案、贡献南邮智慧。（撰稿：王波 01 2025.11.04 南邮获“挑战杯”国赛特等奖7项 再捧“优胜杯” 2025.11.10 南邮MBA项目获评江苏省首个QS Stars五星认证 2025.11.07 南邮举办研究生“鼎山研习营” 专题学习党的二十届四中全会精神 2025.11.07 南邮两门课程出海泰国和印尼国家慕课平台 2025.11.06 南邮党委理论学习中心组召开党的二十届四中全会精神专题学习会 省委教育工委列席旁听工作第三组到会指导 2025.11.05 无锡市惠山区区长赵强一行来南邮调研交流 查看更多 【江苏卫视·江苏新时空】信长星调研高校思政课建设时指出 坚持把立德树人作为根本任务 守正创新推动思政课建设内涵式发展 【新华网客户端】2025年电子信息特色高校发展大会在南京邮电大学举行 【中国教育电视台·全国教育新闻联播】悄悄给饭卡打钱 被这些高校暖到了！ 【江苏教育报】奏响教育强省建设的南邮奋进曲 【光明日报客户端】南京邮电大学：于“热血铸荣光”中厚植青春信仰 【新华社客户端】我国科研人员在无机光伏材料薄膜化领域取得新进展 预告 报道 查看更多 话剧《华罗庚的最后讲演》 2025-11-10 19:00:00 安恒青春剧场 2025年“读经典·颂中华”朗诵比赛 2025-11-07 19:00:00 大学生活动中心安恒青春剧场 美育大讲堂邀您共赏“聆听音乐——为何音乐与音乐何为” 2025-11-10 15:30:00 运满满报告厅（图书馆四楼） 医工交叉院士专家报告会 2025-10-31 09:00:00 行政北楼中天报告厅 人工智能与科研管理专题培训会 2025-10-28 13:30:10 仙林校区中兴报告厅 跨越时空的对话——南京邮电大学红色校史剧《赤子》复演 2025-10-21 19:30:00 大学生活动中心安恒青春剧场 美育大讲堂邀您共赏“《牡丹亭》——从临川笔下到昆曲场上” 2025-10-22 14:00:00 运满满报告厅（图书馆四楼） 2025智能传播与艺术国际研讨会 2025-10-18 14:00:00 计算机学科楼报告厅 查看更多 南邮承办2025江苏高校国际产学研用合作交流周“能源与新材料”分论坛 南邮主办第十七届无线通信与信号处理国际会议 南邮联合主办“中国社会学自主知识体系与中国式现代化”学术座谈会 跨越时空的对话——南邮红色校史剧《赤子》复演 南邮举办第五十六届田径运动会 南邮主办第七届“强邮论坛”暨2025年空地协同低空物流技术学术会议 南邮举办“2025年智能传播与艺术国际研讨会” 南邮举办“热血铸荣光”——弘扬伟大抗战精神 赓续传承红色基因主题文艺展演 南邮举行2025年“邮启新程”主题晚会暨第三十届大学生科技节开幕式 南邮联合主办第18届IEEE毫米波与太赫兹技术联合会议 yyyy MM.dd title text 1 南邮教师获第二届全国电子信息类专业高校教师智慧教学案例竞赛一等奖和最佳创意奖 2 Nature Communications刊发汪联辉教授和王婷教授团队在柔性电化学传感领域的最新研究进展 3 Nature Energy刊发黄维院士辛颢教授铜锌锡硫硒薄膜太阳能最新研究进展 4 南邮教师获全国高校教师教学创新大赛新文科赛道国奖 查看更多 2025 10.29 南邮教师获第二届全国电子信息类专业高校教师智慧教学案例竞赛一等奖和最佳创意奖 2025 10.10 Nature Communications刊发汪联辉教授和王婷教授团队在柔性电化学传感领域的最新研究进展 2025 09.15 Nature Energy刊发黄维院士辛颢教授铜锌锡硫硒薄膜太阳能最新研究进展 2025 08.20 南邮教师获全国高校教师教学创新大赛新文科赛道国奖 查看更多 查看更多 优秀教师 优秀教育工作者 最美辅导员 最美大学生 查看更多 查看更多 微信 视频号 跨越六十载，八旬教授曹伟的求索与守望 美育大讲堂——昆曲《牡丹亭》里的文化传承与美学觉醒 这个系列怎么能少的了南邮？ 在南邮，一起赴“桂”！ 2025.10.31 同心向邮，双倍精彩！ 2025.10.30 青年科学家 | 扎根南邮，他让无线网络更智能！ 2025.10.25 全省体测第一高校的运动会来了！ 2025.10.24 南邮师生热议党的二十届四中全会 科学谋划“十五五”发展新篇章 7 博士后流动站7个 10 一级学科博士学位授权点10个 11 博士学位授权点11个 23 一级学科硕士学位授权点23个 14 硕士专业学位授权点（类别）14个 52 本科专业52个 6 6个学科进入ESI学科排名全球前1% 27 国家一流专业27个 14 14个专业通过国家工程教育专业认证 校历查询 规章制度 服务流程 办公信息 书记信箱 校长信箱 图书馆 艺术馆 南邮校报 VPN登录 实景漫游 校友会 基金会 相关链接 仙林校区地址：南京市仙林大学城文苑路9号 邮编：210023 三牌楼校区地址：南京市新模范马路66号 邮编：210003 锁金村校区地址：南京市龙蟠路177号 邮编：210042 联系电话:（86）-25-85866888 传真:（86）-25-85866999 邮箱:njupt@njupt.edu.cn 苏公网安备32011302320419号 |苏ICP备11073489号-1 版权所有：南京邮电大学","created_at":"2025-11-10T22:54:55.491+08:00","updated_at":"2025-11-10T22:54:55.491+08:00"}(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ curl -X PUT "$BAScurl -X PUT "$BASE_URL/web/items/$PAGE_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Example Updated","content":"Hello world"}'
{"id":1,"user_id":1,"url":"https://www.njupt.edu.cn/","title":"Example Updated","content":"Hello world","created_at":"2025-11-10T22:54:55.491+(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ curl -X POST "$BASE_URL/web/search" \/backend/test$ curl -X POST "$BASE_URL/web/search" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"q":"Example","top_k":5}'
{"results":[{"id":1,"url":"https://www.njupt.edu.cn/","title":"Example Updated","snippet":"Hello world","score":80}]}curl -X POST "$BASE_URL/web/search" \nd/test$ curl -X POST "$BASE_URL/web/search" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"urls":["https://www.rfc-editor.org/"],"top_k":3}'
{"results":[{"id":2,"url":"https://www.rfc-editor.org/","title":"» RFC Editor","snippet":"Search RFCs Advanced Search RFC Editor The RFC Series The RFC Series (ISSN 2070-1721) contains technical and organizational documents about the Internet, including the specifications and policy documents produced by five streams: the Intern","score":100}]}(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ curl -X Pcurl -X POST "$BASE_URL/web/ingest" -H "Authorization: Bearer $TOKEN"
curl -X POST "$BASE_URL/web/chunk"  -H "Authorization: Bearer $TOKEN"
{"message":"Web content processed successfully","success":true}{"messag(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ curl -X DELETE "$BASE_URL/web/items/$PAGE_ID" \$ curl -X DELETE "$BASE_URL/web/items/$PAGE_ID" \
  -H "Authorization: Bearer $TOKEN"
(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ curl -X POST "$BASE_URL/chat/sessions" \chat/sessions" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Smoke Test Session"}'
{"id":1,"user_id":1,"title":"Smoke Test Session","created_at":"2025-11-10T22:56:41.944+08:00","updated_at":"2025-11-10T22:56:41.944+08:00"}(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ SESSESSION_ID=$(curl -sS -X POST "$BASE_URL/chat/sessions" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Smoke Test Session"}' | jq -r '.id // .data.id // .session_id // .data.session_id') && echo "$SESSION_ID"
2
(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ SESSION_ID=$(curl -sS -X POST "$BASE_URL/chat/sessions" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Smoke Test Session"}' \
  | python3 - <<'PY'
import sys, json
d=json.load(sys.stdin)
print(d.get("id") or (d.get("data") or {}).get("id") or d.get("session_id") or (d.get("data") or {}).get("session_id") or "")
PY
) && echo "$SESSION_ID"
curl: (23) Failed writing body
Traceback (most recent call last):
  File "<stdin>", line 2, in <module>
  File "/home/ubuntu/miniconda3/envs/FA/lib/python3.10/json/__init__.py", line 293, in load
    return loads(fp.read(),
  File "/home/ubuntu/miniconda3/envs/FA/lib/python3.10/json/__init__.py", line 346, in loads
    return _default_decoder.decode(s)
  File "/home/ubuntu/miniconda3/envs/FA/lib/python3.10/json/decoder.py", line 337, in decode
    obj, end = self.raw_decode(s, idx=_w(s, 0).end())
  File "/home/ubuntu/miniconda3/envs/FA/lib/python3.10/json/decoder.py", line 355, in raw_decode
    raise JSONDecodeError("Expecting value", s, err.value) from None
json.decoder.JSONDecodeError: Expecting value: line 1 column 1 (char 0)
(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ SESSION_ID=$(curl -sS -X POST "$BASE_URL/chat/sessions" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Smoke Test Session"}' | jq -r '.id // .data.id // .session_id // .data.session_id') && echo "$SESSION_ID"
4
(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ curl -X GET "$BASE_URL/chat/sessions" \
  -H "Authorization: Bearer $TOKEN"
[{"id":4,"user_id":1,"title":"Smoke Test Session","created_at":"2025-11-10T22:57:37+08:00","updated_at":"2025-11-10T22:57:37+08:00"},{"id":3,"user_id":1,"title":"Smoke Test Session","created_at":"2025-11-10T22:57:11+08:00","updated_at":"2025-11-10T22:57:11+08:00"},{"id":2,"user_id":1,"title":"Smoke Test Session","created_at":"2025-11-10T22:57:02+08:00","updated_at":"2025-11-10T22:57:02+08:00"},{"id":1,"user_id":1,"title":"Smoke Test Session","created_at":"2025-11-10T22:56:42+08:00","updated_at":"2025-11-10T22:56:42+08:00"}](FA) ubuntu@WFcurl -X POST "$BASE_URL/chat/sessions/$SESSION_ID/messages" \sessions/$SESSION_ID/messages" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role":"user","content":"Hello from curl"}'
{"id":1,"user_id":1,"session_id":4,"content":"Hello from curl","role":"user","created_at":"2025-11-10T22:57:55.642469924+08:00"}(FA) ubuntu@WF(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ curl -X POST "$BASE_URL/chat/messages" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"session_id\":\"$SESSION_ID\",\"content\":\"Second message via alias\"}"
{"id":2,"user_id":1,"session_id":4,"content":"Second message via alias","role":"user","created_at":"2025-11-10T22:58:11.772037186+08:00"}(FA) (FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ curl -X GET "$BASE_URL/chat/sessions/$SESSION_ID/messages" \
  -H "Authorization: Bearer $TOKEN"
[{"id":1,"user_id":1,"session_id":4,"content":"Hello from curl","role":"user","created_at":"2025-11-10T22:57:56+08:00"},{"id":2,"user_id":1,"session_id":4,"content":"Second message via alias","role":"user","created_at":"2025-11-10T22:58:12+08:00"}](FA) ubuntucurl -X POST "$BASE_URL/chat/sessions/$SESSION_ID/complete" \at/sessions/$SESSION_ID/complete" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content":"Briefly echo hello","model":""}'
{"content":"Hello!"}(FA) ubuntu@WFTGF2025:~/Chcurl -N -X POST "$BASE_URL/chat/sessions/$SESSION_ID/stream" \SESSION_ID/stream" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content":"streaming hello","model":""}'
{"delta": "H"}
{"delta": "-e"}
{"delta": "-l"}
{"delta": "-l"}
{"delta": "-o-!"}
{"event": "done"}
(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ curl -N -X POST "$BASE_URL/chat/sessions/$SESSION_ID/stream" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content":"你好","model":""}'
{"delta": "你好"}
{"delta": "！很高兴"}
{"delta": "见到你！"}
{"delta": "有什么"}
{"delta": "我可以帮助你的吗"}
{"delta": "？"}
{"event": "done"}
(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ curl -X DELETE "$BASE_URL/api/auth/delete" \
  -H "Authorization: Bearer $TOKEN"
404 page not found(FA) ubuntu@WFTGF2025:~/ChatWeb/backend/test$ 