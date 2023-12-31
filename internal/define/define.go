package define

import "github.com/golang-jwt/jwt/v5"

type SystemConfig struct {
	Port  string `json:"port"`  // 端口
	Entry string `json:"entry"` // 入口地址
}
type UserBasic struct {
	Name     string `json:"name"`     // 用户名
	Password string `json:"password"` // 密码
}
type UserClaim struct {
	jwt.RegisteredClaims
}

var (
	Key            = []byte("op-panel")
	PID            int
	PageSize       = 20
	ShellDir       = "./shell"
	LogDir         = "./log"
	DefaultWebDir  = "/home/wwwroot/"
	NginxConfigDir = "/home/nginx/conf/"
	MysqlDNS       = "root:123456@tcp(192.168.12.133:3306)"
)
