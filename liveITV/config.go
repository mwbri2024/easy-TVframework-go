package liveITV

import (
	"fmt"
	"strings"
	"sync"
)

var Enable bool = true
var Debug bool = false

var Enable_Yangshi bool = true
var Enable_Weishi bool = true

var WorkIP_bestzb string = ""
var IPList_bestzb = []IPRecord{
	{Address: "39.135.132.111", Fail: 0},
	{Address: "39.135.97.11", Fail: 0},
	{Address: "39.135.97.33", Fail: 0},
}

var WorkIP_hnbblive string = ""
var IPList_hnbblive = []IPRecord{
	{Address: "39.135.132.111", Fail: 0},
	{Address: "39.135.97.11", Fail: 0},
	{Address: "39.135.97.33", Fail: 0},
}

var WorkIP_fifalive string = ""
var IPList_fifalive = []IPRecord{
	{Address: "39.135.132.111", Fail: 0},
	{Address: "39.135.97.11", Fail: 0},
	{Address: "39.135.97.33", Fail: 0},
}

var ChannelList = []ChannelRecord{
	{Group: "央视", Tvgid: "CCTV1", Contentid: "5000000004000002226", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "央视", Tvgid: "CCTV2", Contentid: "5000000011000031101", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "央视", Tvgid: "CCTV3", Contentid: "5000000004000008883", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "央视", Tvgid: "CCTV4", Contentid: "5000000011000031102", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "央视", Tvgid: "CCTV5", Contentid: "5000000004000008885", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "央视", Tvgid: "CCTV5+", Contentid: "5000000011000031127", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "央视", Tvgid: "CCTV6", Contentid: "5000000004000008886", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "央视", Tvgid: "CCTV7", Contentid: "5000000011000031104", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "央视", Tvgid: "CCTV8", Contentid: "5000000004000008888", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "央视", Tvgid: "CCTV9", Contentid: "5000000011000288020", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "央视", Tvgid: "CCTV10", Contentid: "5000000004000012827", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "央视", Tvgid: "CCTV11", Contentid: "5000000011000031106", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "央视", Tvgid: "CCTV12", Contentid: "5000000011000031107", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "央视", Tvgid: "CCTV13", Contentid: "5000000011000031108", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "央视", Tvgid: "CCTV14", Contentid: "5000000004000006673", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "央视", Tvgid: "CCTV15", Contentid: "5000000011000031109", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "央视", Tvgid: "CCTV16", Contentid: "5000000008000023253", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "央视", Tvgid: "CCTV16-4K", Contentid: "5000000008000023254", Cdn: "bestzb", Quality: "$12M UHD HEVC"},
	{Group: "央视", Tvgid: "CCTV17", Contentid: "5000000011000288014", Cdn: "bestzb", Quality: "$8M FHD"},

	{Group: "卫视", Tvgid: "安徽卫视", Contentid: "5000000004000023002", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "卫视", Tvgid: "北京卫视", Contentid: "5000000004000031556", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "卫视", Tvgid: "东方卫视", Contentid: "5000000004000014098", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "卫视", Tvgid: "东南卫视", Contentid: "5000000004000010584", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "卫视", Tvgid: "甘肃卫视", Contentid: "5000000011000031121", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "卫视", Tvgid: "广东卫视", Contentid: "5000000004000014694", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "卫视", Tvgid: "广西卫视", Contentid: "5000000011000031118", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "卫视", Tvgid: "贵州卫视", Contentid: "5000000004000025843", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "卫视", Tvgid: "海南卫视", Contentid: "5000000004000006211", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "卫视", Tvgid: "河北卫视", Contentid: "5000000006000040016", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "卫视", Tvgid: "河南卫视", Contentid: "5000000011000031119", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "卫视", Tvgid: "黑龙江卫视", Contentid: "5000000004000025203", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "卫视", Tvgid: "湖北卫视", Contentid: "5000000004000014954", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "卫视", Tvgid: "湖南卫视", Contentid: "5000000004000006692", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "卫视", Tvgid: "吉林卫视", Contentid: "5000000011000031117", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "卫视", Tvgid: "江苏卫视", Contentid: "5000000004000019351", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "卫视", Tvgid: "江西卫视", Contentid: "5000000004000011210", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "卫视", Tvgid: "辽宁卫视", Contentid: "5000000004000011671", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "卫视", Tvgid: "三沙卫视", Contentid: "5000000011000288016", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "卫视", Tvgid: "山东卫视", Contentid: "5000000004000020424", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "卫视", Tvgid: "深圳卫视", Contentid: "5000000004000007410", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "卫视", Tvgid: "四川卫视", Contentid: "5000000004000006119", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "卫视", Tvgid: "天津卫视", Contentid: "5000000004000006827", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "卫视", Tvgid: "云南卫视", Contentid: "5000000011000031120", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "卫视", Tvgid: "浙江卫视", Contentid: "5000000004000007275", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "卫视", Tvgid: "重庆卫视", Contentid: "5000000004000025797", Cdn: "bestzb", Quality: "$8M FHD"},
	{Group: "卫视", Tvgid: "内蒙古卫视", Contentid: "7331079758781182663", Cdn: "bestzb", Quality: "$2M SD"},
	{Group: "卫视", Tvgid: "宁夏卫视", Contentid: "5000000006000040022", Cdn: "bestzb", Quality: "$4M SD"},
	{Group: "卫视", Tvgid: "青海卫视", Contentid: "5000000006000040015", Cdn: "bestzb", Quality: "$4M FHD"},
	{Group: "卫视", Tvgid: "山西卫视", Contentid: "5000000006000040023", Cdn: "bestzb", Quality: "$4M SD"},
	{Group: "卫视", Tvgid: "陕西卫视", Contentid: "5000000006000040017", Cdn: "bestzb", Quality: "$4M SD"},
	{Group: "卫视", Tvgid: "西藏卫视", Contentid: "5381760837640571168", Cdn: "bestzb", Quality: "$2M SD"},
	{Group: "卫视", Tvgid: "新疆卫视", Contentid: "5000000006000040018", Cdn: "bestzb", Quality: "$4M SD"},
	{Group: "卫视", Tvgid: "安多卫视", Contentid: "5000000006000022124", Cdn: "bestzb", Quality: "$2M SD"},
	{Group: "卫视", Tvgid: "兵团卫视", Contentid: "5000000006000040020", Cdn: "bestzb", Quality: "$4M SD"},
	{Group: "卫视", Tvgid: "康巴卫视", Contentid: "5000000003000001598", Cdn: "bestzb", Quality: "$2M SD"},
	{Group: "卫视", Tvgid: "厦门卫视", Contentid: "5084893189220221988", Cdn: "bestzb", Quality: "$2M SD"},
	{Group: "卫视", Tvgid: "大湾区卫视", Contentid: "2000000003000000045", Cdn: "hnbblive", Quality: "$4M SD"},
	{Group: "卫视", Tvgid: "山东教育卫视", Contentid: "2000000003000000013", Cdn: "hnbblive", Quality: "$2M SD"},
	{Group: "卫视", Tvgid: "延边卫视", Contentid: "2000000003000000049", Cdn: "hnbblive", Quality: "$4M SD"},

	{Group: "数字", Tvgid: "北京纪实科教", Contentid: "2000000002000000065", Cdn: "hnbblive", Quality: "$8M FHD"},
	{Group: "数字", Tvgid: "上海新闻综合", Contentid: "2000000002000000005", Cdn: "hnbblive", Quality: "$8M FHD"},
	{Group: "数字", Tvgid: "上海都市", Contentid: "2000000002000000012", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "数字", Tvgid: "第一财经", Contentid: "2000000002000000004", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "数字", Tvgid: "东方财经", Contentid: "2000000002000000090", Cdn: "hnbblive", Quality: "$8M FHD"},
	{Group: "数字", Tvgid: "东方影视", Contentid: "2000000002000000013", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "数字", Tvgid: "五星体育", Contentid: "2000000002000000007", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "数字", Tvgid: "广东珠江", Contentid: "2000000003000000033", Cdn: "hnbblive", Quality: "$3M HD"},
	{Group: "数字", Tvgid: "浙江教科影视", Contentid: "2000000004000000001", Cdn: "hnbblive", Quality: "$2M SD"},
	{Group: "数字", Tvgid: "茶频道", Contentid: "2000000002000000070", Cdn: "hnbblive", Quality: "$8M FHD"},
	{Group: "数字", Tvgid: "快乐垂钓", Contentid: "2000000002000000067", Cdn: "hnbblive", Quality: "$8M FHD"},
	{Group: "数字", Tvgid: "金鹰卡通", Contentid: "8919073838942247498", Cdn: "hnbblive", Quality: "$2M SD"},
	{Group: "数字", Tvgid: "家庭理财", Contentid: "2000000002000000064", Cdn: "hnbblive", Quality: "$2M SD"},
	{Group: "数字", Tvgid: "电子竞技", Contentid: "2000000004000000015", Cdn: "hnbblive", Quality: "$3M HD"},
	{Group: "数字", Tvgid: "东方大剧院", Contentid: "2000000004000000004", Cdn: "hnbblive", Quality: "$3M HD"},
	{Group: "数字", Tvgid: "高清娱乐", Contentid: "2000000004000000013", Cdn: "hnbblive", Quality: "$3M HD"},
	{Group: "数字", Tvgid: "精品剧场", Contentid: "2000000004000000002", Cdn: "hnbblive", Quality: "$3M HD"},
	{Group: "数字", Tvgid: "精品综合", Contentid: "2000000003000000008", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "数字", Tvgid: "漫游世界", Contentid: "2000000004000000017", Cdn: "hnbblive", Quality: "$3M HD"},
	{Group: "数字", Tvgid: "欧美影院", Contentid: "2000000004000000005", Cdn: "hnbblive", Quality: "$3M HD"},
	{Group: "数字", Tvgid: "热播精选", Contentid: "2000000003000000016", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "数字", Tvgid: "探索纪录", Contentid: "2000000004000000007", Cdn: "hnbblive", Quality: "$3M HD"},
	{Group: "数字", Tvgid: "少儿动漫", Contentid: "2000000004000000003", Cdn: "hnbblive", Quality: "$3M HD"},
	{Group: "数字", Tvgid: "卡酷少儿", Contentid: "7947519469029992015", Cdn: "hnbblive", Quality: "$2M SD"},
	{Group: "数字", Tvgid: "求索动物", Contentid: "2000000004000000009", Cdn: "hnbblive", Quality: "$8M FHD"},
	{Group: "数字", Tvgid: "求索纪录", Contentid: "2000000004000000010", Cdn: "hnbblive", Quality: "$8M FHD"},
	{Group: "数字", Tvgid: "求索科学", Contentid: "2000000004000000011", Cdn: "hnbblive", Quality: "$8M FHD"},
	{Group: "数字", Tvgid: "求索生活", Contentid: "2000000004000000008", Cdn: "hnbblive", Quality: "$8M FHD"},

	{Group: "NewTV", Tvgid: "爱情喜剧", Contentid: "2000000003000000010", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "NewTV", Tvgid: "超级电视剧", Contentid: "2000000003000000032", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "NewTV", Tvgid: "超级电影", Contentid: "2000000003000000031", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "NewTV", Tvgid: "超级体育", Contentid: "2000000003000000030", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "NewTV", Tvgid: "超级综艺", Contentid: "2000000003000000029", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "NewTV", Tvgid: "潮妈辣婆", Contentid: "2000000003000000018", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "NewTV", Tvgid: "哒啵电竞", Contentid: "2000000003000000066", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "NewTV", Tvgid: "哒啵赛事", Contentid: "2000000003000000040", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "NewTV", Tvgid: "东北热剧", Contentid: "2000000003000000051", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "NewTV", Tvgid: "动作电影", Contentid: "2000000003000000017", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "NewTV", Tvgid: "古装剧场", Contentid: "2000000003000000024", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "NewTV", Tvgid: "黑莓电影", Contentid: "2000000003000000001", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "NewTV", Tvgid: "黑莓动画", Contentid: "2000000003000000002", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "NewTV", Tvgid: "欢乐剧场", Contentid: "2000000003000000050", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "NewTV", Tvgid: "家庭剧场", Contentid: "2000000003000000012", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "NewTV", Tvgid: "金牌综艺", Contentid: "2000000003000000005", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "NewTV", Tvgid: "惊悚悬疑", Contentid: "2000000003000000015", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "NewTV", Tvgid: "精品大剧", Contentid: "2000000003000000020", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "NewTV", Tvgid: "精品纪录", Contentid: "2000000003000000019", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "NewTV", Tvgid: "精品萌宠", Contentid: "2000000003000000067", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "NewTV", Tvgid: "精品体育", Contentid: "2000000003000000021", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "NewTV", Tvgid: "军旅剧场", Contentid: "2000000003000000014", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "NewTV", Tvgid: "军事评论", Contentid: "2000000003000000022", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "NewTV", Tvgid: "魅力潇湘", Contentid: "2000000003000000041", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "NewTV", Tvgid: "农业致富", Contentid: "2000000003000000003", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "NewTV", Tvgid: "武搏世界", Contentid: "2000000003000000007", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "NewTV", Tvgid: "炫舞未来", Contentid: "2000000003000000044", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "NewTV", Tvgid: "怡伴健康", Contentid: "2000000003000000023", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "NewTV", Tvgid: "中国功夫", Contentid: "2000000003000000009", Cdn: "hnbblive", Quality: "$4M FHD"},

	{Group: "SiTV", Tvgid: "动漫秀场", Contentid: "2000000002000000009", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "SiTV", Tvgid: "都市剧场", Contentid: "2000000002000000015", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "SiTV", Tvgid: "法治天地", Contentid: "2000000002000000014", Cdn: "hnbblive", Quality: "$8M FHD"},
	{Group: "SiTV", Tvgid: "欢笑剧场", Contentid: "2000000002000000016", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "SiTV", Tvgid: "金色学堂", Contentid: "2000000002000000061", Cdn: "hnbblive", Quality: "$8M FHD"},
	{Group: "SiTV", Tvgid: "劲爆体育", Contentid: "2000000002000000008", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "SiTV", Tvgid: "乐游", Contentid: "2000000002000000092", Cdn: "hnbblive", Quality: "$8M FHD"},
	{Group: "SiTV", Tvgid: "魅力足球", Contentid: "2000000002000000068", Cdn: "hnbblive", Quality: "$8M FHD"},
	{Group: "SiTV", Tvgid: "生活时尚", Contentid: "2000000002000000006", Cdn: "hnbblive", Quality: "$4M FHD"},
	{Group: "SiTV", Tvgid: "游戏风云", Contentid: "2000000002000000011", Cdn: "hnbblive", Quality: "$4M FHD"},

	{Group: "IHOT", Tvgid: "爱谍战", Contentid: "2000000004000000038", Cdn: "hnbblive", Quality: "$8M FHD"},
	{Group: "IHOT", Tvgid: "爱动漫", Contentid: "2000000004000000059", Cdn: "hnbblive", Quality: "$8M FHD"},
	{Group: "IHOT", Tvgid: "爱科幻", Contentid: "2000000004000000033", Cdn: "hnbblive", Quality: "$8M FHD"},
	{Group: "IHOT", Tvgid: "爱科学", Contentid: "2000000004000000047", Cdn: "hnbblive", Quality: "$8M FHD"},
	{Group: "IHOT", Tvgid: "爱浪漫", Contentid: "2000000004000000035", Cdn: "hnbblive", Quality: "$8M FHD"},
	{Group: "IHOT", Tvgid: "爱历史", Contentid: "2000000004000000046", Cdn: "hnbblive", Quality: "$8M FHD"},
	{Group: "IHOT", Tvgid: "爱旅行", Contentid: "2000000004000000056", Cdn: "hnbblive", Quality: "$8M FHD"},
	{Group: "IHOT", Tvgid: "爱奇谈", Contentid: "2000000004000000058", Cdn: "hnbblive", Quality: "$8M FHD"},
	{Group: "IHOT", Tvgid: "爱赛车", Contentid: "2000000004000000055", Cdn: "hnbblive", Quality: "$8M FHD"},
	{Group: "IHOT", Tvgid: "爱体育", Contentid: "2000000004000000060", Cdn: "hnbblive", Quality: "$8M FHD"},
	{Group: "IHOT", Tvgid: "爱玩具", Contentid: "2000000004000000053", Cdn: "hnbblive", Quality: "$8M FHD"},
	{Group: "IHOT", Tvgid: "爱喜剧", Contentid: "2000000004000000032", Cdn: "hnbblive", Quality: "$8M FHD"},
	{Group: "IHOT", Tvgid: "爱悬疑", Contentid: "2000000004000000036", Cdn: "hnbblive", Quality: "$8M FHD"},
	{Group: "IHOT", Tvgid: "爱幼教", Contentid: "2000000004000000049", Cdn: "hnbblive", Quality: "$8M FHD"},
	{Group: "IHOT", Tvgid: "爱院线", Contentid: "2000000004000000034", Cdn: "hnbblive", Quality: "$8M FHD"},

	{Group: "咪视通", Tvgid: "睛彩广场舞", Contentid: "3000000020000011523", Cdn: "FifastbLive", Quality: "$8M FHD"},
	{Group: "咪视通", Tvgid: "睛彩竞技", Contentid: "3000000020000011528", Cdn: "FifastbLive", Quality: "$8M FHD"},
	{Group: "咪视通", Tvgid: "睛彩篮球", Contentid: "3000000020000011529", Cdn: "FifastbLive", Quality: "$8M FHD"},
	{Group: "咪视通", Tvgid: "睛彩青少", Contentid: "3000000020000011525", Cdn: "FifastbLive", Quality: "$8M FHD"},
	{Group: "咪视通", Tvgid: "3000000001000005308", Contentid: "3000000001000005308", Cdn: "FifastbLive", Quality: "$8M FHD"},
	{Group: "咪视通", Tvgid: "3000000001000005969", Contentid: "3000000001000005969", Cdn: "FifastbLive", Quality: "$8M FHD"},
	{Group: "咪视通", Tvgid: "3000000001000007218", Contentid: "3000000001000007218", Cdn: "FifastbLive", Quality: "$8M FHD"},
	{Group: "咪视通", Tvgid: "3000000001000008001", Contentid: "3000000001000008001", Cdn: "FifastbLive", Quality: "$8M FHD"},
	{Group: "咪视通", Tvgid: "3000000001000008176", Contentid: "3000000001000008176", Cdn: "FifastbLive", Quality: "$8M FHD"},
	{Group: "咪视通", Tvgid: "3000000001000008379", Contentid: "3000000001000008379", Cdn: "FifastbLive", Quality: "$8M FHD"},
	{Group: "咪视通", Tvgid: "3000000001000010129", Contentid: "3000000001000010129", Cdn: "FifastbLive", Quality: "$8M FHD"},
	{Group: "咪视通", Tvgid: "3000000001000010948", Contentid: "3000000001000010948", Cdn: "FifastbLive", Quality: "$8M FHD"},
	{Group: "咪视通", Tvgid: "3000000001000028638", Contentid: "3000000001000028638", Cdn: "FifastbLive", Quality: "$8M FHD"},
	{Group: "咪视通", Tvgid: "3000000001000031494", Contentid: "3000000001000031494", Cdn: "FifastbLive", Quality: "$8M FHD"},
	{Group: "咪视通", Tvgid: "3000000010000000097", Contentid: "3000000010000000097", Cdn: "FifastbLive", Quality: "$5M FHD"},
	{Group: "咪视通", Tvgid: "3000000010000002019", Contentid: "3000000010000002019", Cdn: "FifastbLive", Quality: "$5M FHD"},
	{Group: "咪视通", Tvgid: "3000000010000002809", Contentid: "3000000010000002809", Cdn: "FifastbLive", Quality: "$5M FHD"},
	{Group: "咪视通", Tvgid: "3000000010000003915", Contentid: "3000000010000003915", Cdn: "FifastbLive", Quality: "$5M FHD"},
	{Group: "咪视通", Tvgid: "3000000010000004193", Contentid: "3000000010000004193", Cdn: "FifastbLive", Quality: "$5M FHD"},
	{Group: "咪视通", Tvgid: "3000000010000005837", Contentid: "3000000010000005837", Cdn: "FifastbLive", Quality: "$8M FHD"},
	{Group: "咪视通", Tvgid: "3000000010000006077", Contentid: "3000000010000006077", Cdn: "FifastbLive", Quality: "$5M FHD"},
	{Group: "咪视通", Tvgid: "3000000010000006658", Contentid: "3000000010000006658", Cdn: "FifastbLive", Quality: "$5M FHD"},
	{Group: "咪视通", Tvgid: "3000000010000009788", Contentid: "3000000010000009788", Cdn: "FifastbLive", Quality: "$5M FHD"},
	{Group: "咪视通", Tvgid: "3000000010000010833", Contentid: "3000000010000010833", Cdn: "FifastbLive", Quality: "$5M FHD"},
	{Group: "咪视通", Tvgid: "3000000010000011297", Contentid: "3000000010000011297", Cdn: "FifastbLive", Quality: "$5M FHD"},
	{Group: "咪视通", Tvgid: "3000000010000011518", Contentid: "3000000010000011518", Cdn: "FifastbLive", Quality: "$5M FHD"},
	{Group: "咪视通", Tvgid: "3000000010000012558", Contentid: "3000000010000012558", Cdn: "FifastbLive", Quality: "$5M FHD"},
	{Group: "咪视通", Tvgid: "3000000010000012616", Contentid: "3000000010000012616", Cdn: "FifastbLive", Quality: "$5M FHD"},
	{Group: "咪视通", Tvgid: "3000000010000015470", Contentid: "3000000010000015470", Cdn: "FifastbLive", Quality: "$5M FHD"},
	{Group: "咪视通", Tvgid: "3000000010000015560", Contentid: "3000000010000015560", Cdn: "FifastbLive", Quality: "$5M FHD"},
	{Group: "咪视通", Tvgid: "3000000010000019839", Contentid: "3000000010000019839", Cdn: "FifastbLive", Quality: "$5M FHD"},
	{Group: "咪视通", Tvgid: "3000000010000021904", Contentid: "3000000010000021904", Cdn: "FifastbLive", Quality: "$5M FHD"},
	{Group: "咪视通", Tvgid: "3000000010000023434", Contentid: "3000000010000023434", Cdn: "FifastbLive", Quality: "$5M FHD"},
	{Group: "咪视通", Tvgid: "3000000010000025380", Contentid: "3000000010000025380", Cdn: "FifastbLive", Quality: "$5M FHD"},
	{Group: "咪视通", Tvgid: "3000000010000027691", Contentid: "3000000010000027691", Cdn: "FifastbLive", Quality: "$5M FHD"},
	{Group: "咪视通", Tvgid: "3000000010000031669", Contentid: "3000000010000031669", Cdn: "FifastbLive", Quality: "$5M FHD"},
	{Group: "咪视通", Tvgid: "3000000020000011518", Contentid: "3000000020000011518", Cdn: "FifastbLive", Quality: "$8M FHD"},
	{Group: "咪视通", Tvgid: "3000000020000011519", Contentid: "3000000020000011519", Cdn: "FifastbLive", Quality: "$8M FHD"},
	{Group: "咪视通", Tvgid: "3000000020000011520", Contentid: "3000000020000011520", Cdn: "FifastbLive", Quality: "$8M FHD"},
	{Group: "咪视通", Tvgid: "3000000020000011521", Contentid: "3000000020000011521", Cdn: "FifastbLive", Quality: "$8M FHD"},
	{Group: "咪视通", Tvgid: "3000000020000011522", Contentid: "3000000020000011522", Cdn: "FifastbLive", Quality: "$8M FHD"},
}

// 同步锁
var mu sync.Mutex

// Channel 结构体
type ChannelRecord struct {
	Group     string
	Tvgid     string
	Contentid string
	Cdn       string
	Quality   string
}

// IP 结构体
type IPRecord struct {
	Address string
	Fail    int
}

// 更新 IP 列表
func UpdateIPList(newIPs []IPRecord, ip_list *[]IPRecord) {
	mu.Lock()
	defer mu.Unlock()
	*ip_list = newIPs
}

// 获取 Fail 最小的 IP
func GetBestIP(ip_list []IPRecord) IPRecord {
	mu.Lock()
	defer mu.Unlock()

	bestIP := ip_list[0]
	for _, ip := range ip_list {
		if ip.Fail < bestIP.Fail {
			bestIP = ip
		}
	}

	return bestIP
}

// 失败计数 +1
func IncreaseFail(ipAddress string, ip_list []IPRecord) {
	// 去掉端口号
	addressWithoutPort := strings.Split(ipAddress, ":")[0]

	mu.Lock()
	defer mu.Unlock()

	for i, ip := range ip_list {
		if ip.Address == addressWithoutPort {
			ip_list[i].Fail++
			// fmt.Printf("IP %s 的失败计数增加到 %d\n", addressWithoutPort, ip_list[i].Fail)
			break
		}
	}
}

// 打印 IP 列表
func PrintIPList(ip_list []IPRecord) {
	mu.Lock()
	defer mu.Unlock()

	for _, ip := range ip_list {
		fmt.Printf("  IP: %s, Fail: %d\n", ip.Address, ip.Fail)
	}
}
