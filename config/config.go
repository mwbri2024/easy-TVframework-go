package config

// 基础配置
const Version = "20250331.0000"
const Port = "8123"
const Debug bool = false

// 自动升级配置
const Update bool = false
const VersionCheckURL = "https://abc.com/app/latest_version.json"
const BinaryDownloadBaseURL = "https://abc.com/app/"

// 运行配置
var Token string = ""
var Province string = ""
var Operator string = ""

// 省份列表
var ProvinceList = map[string]string{
	"北京":  "Live",
	"天津":  "Live",
	"河北":  "Live",
	"山西":  "Live",
	"内蒙":  "Live",
	"辽宁":  "Live",
	"吉林":  "Live",
	"黑龙江": "Live",
	"上海":  "Live",
	"江苏":  "Live",
	"浙江":  "Live",
	"安徽":  "Live",
	"福建":  "Live",
	"江西":  "Live",
	"山东":  "Live",
	"河南":  "Live",
	"湖北":  "Live",
	"湖南":  "Live",
	"广东":  "Live",
	"广西":  "Live",
	"海南":  "Live",
	"重庆":  "Live",
	"四川":  "Live",
	"贵州":  "Live",
	"云南":  "Live",
	"西藏":  "Live",
	"陕西":  "Live",
	"甘肃":  "Live",
	"青海":  "Live",
	"宁夏":  "Live",
	"新疆":  "Live",
	"香港":  "Live",
	"澳门":  "Live",
	"台湾":  "Live",
}

// 运营商列表
var OperatorList = map[string]string{
	"电信": "中国电信",
	"联通": "中国联通",
	"移动": "中国移动",
}
