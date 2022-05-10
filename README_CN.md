# KMS 密钥鉴权

## 第三方鉴权配置方法

鉴权配置以 json 数组方式配置，json 对象如下所示。

```json
{
    "name":      "参数名称",
    "type":      "参数类型",
    "data":      "参数内容",
    "in":        "鉴权返回值所存放位置",
    "captureIn": "鉴权参数捕获位置",      // 仅在 cookie 中生效
    "from":      "请求参数"              // 仅在 cookie 中生效 
}
```

- type
    - string
    - number
    - boolean
    - keyid
    - keysecret
    - signature
    - expire
    - method
    - authurl
    - cookie
    - signcmd

- in | captureIn | from
    - body
    - query
    - header
    - cookie

- method
    - GET
    - POST
    - PUT

### signature
- 配置格式: cmd1 para1 para2... | cmd2 para1 para2 ... | ...
- base64 std/url encode/decode
- hex encode/decode
- sha1 <SECRET_KEY>
- sha256 <SECRET_KEY>
- append begin/end <APPEND_STR>
- url path/query
- md5 [<SALT>]
- crc32 [IEEE/CASTAGNOLI]
- crc64 [ISO/ECMA]
- aes CBC pkcs5/pkcs7 decode/encode <SECRET_KEY>
- rsa decode/encode <SECRET_KEYS>
- sort json/xml/query same/snake/gonic desc/asc

#### 示例
```
sort query gonic asc|append begin GET\n/iaas/\n|sha256 <SECRET_KEY>|base64 std encode
```
如上所示签名方法处理步骤：
- 用驼峰命名法将参数按升序排序，输出成query格式，如action.1=foo&action.2=bar&age=18&name=bob
- 在头部追加字符串GET\n/iaas/\n，如GET\n/iaas/\naction.1=foo&action.2=bar&age=18&name=bob
- 用<SECRET_KEY>将以上所得字符串进行HMAC sha256计算
- 将以上所得hash值进行base64 std编码，即为最终用于服务器校验的Signature字符串

#### 完整鉴权方法配置示例
如下配置表示采用`sort query gonic asc|append begin GET\n/iaas/\n|sha256 <SECRET_KEY>|base64 std encode`签名方法，
输出access_key_id和signature两个参数的签名配置。
```
[
	{
		"type":"string",
		"name":"signcmd",
		"data":"sort query gonic asc|append begin GET\n/iaas/\n|sha256 <SECRET_KEY>|base64 std encode"
	},
	{
		"name":"access_key_id",
		"type":"keyid",
		"in":"query"
	},
	{
		"name":"signature",
		"type":"signature",
		"in":"body"	
	}
]
```

### cookie
该鉴权方法通过调用外部系统接口实现鉴权代理，配置项中主要以 http 请求参数作为基础。

#### 示例
```json
[
	{
		"name": "userName",
		"type": "keyid",
		"from": "body"
	},
	{
		"name": "password",
        "type": "keysecret",
        "from": "body"
	},
    {
        "name": "constField",
        "type": "string",
        "data": "content",
        "from": "body"
	},
    {
        "type": "authurl",
        "data": "https://localhost/login"
	},
    {
		"type": "method",
		"data": "POST"
	},	
	{
		"name": "token",
		"type": "string",
		"captureIn": "body",
		"in": "body"
	},
	{
		"type": "expire",
		"data": "5"
	}
]
```
上述示例中，定义了三个鉴权请求参数[userName, password, constField]，这三个参数都将通过body传输给鉴权API[https://localhost/login]。从该API返回参数中获取到token并返回。该配置项中还包含expire类型，这会将最终所需要的返回参数缓存起来以便下一次访问的时候不用再去调用API。