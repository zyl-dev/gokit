package masking

import "regexp"

// LogMaskingConfig 打印配置时日志脱敏
func LogMaskingConfig(a string) string {
	a = RegExpReplaceAllString(a, `"appkey":"(.*?)",`, `"appkey":"",`)
	a = RegExpReplaceAllString(a, `"dsn":"(.*?)@`, `"dsn":"`)
	a = RegExpReplaceAllString(a, `"sentrydsn":"(.*?)",`, `"sentrydsn":"","`)
	a = RegExpReplaceAllString(a, `"coinmarketmapapikey":"(.*?)",`, `"coinmarketmapapikey":"","`)
	a = RegExpReplaceAllString(a, `"password":"(.*?)",`, `"password":"","`)
	a = RegExpReplaceAllString(a, `"pass":"(.*?)",`, `"pass":"","`)
	a = RegExpReplaceAllString(a, `"appids":{(.*?)}`, `"appids":""`)
	a = RegExpReplaceAllString(a, `"Addrs":\["(.*?)"`, `"Addrs":[""`)
	return a
}

// RegExpReplaceAllString a 代表原始字符串，str 代表正则表达式，repl 最后替换的结果
func RegExpReplaceAllString(a, str, repl string) string {
	re := regexp.MustCompile(str)
	s := re.ReplaceAllString(a, repl)
	return s
}
