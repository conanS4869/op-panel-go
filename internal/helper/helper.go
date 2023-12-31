package helper

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/satori/go.uuid"
	"log"
	"math/rand"
	"op-panel-go/define"
	"os"
	"os/exec"
	"path"
	"syscall"
	"time"
)

func RandomString(n int) string {
	s := "0123456789AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
	ans := make([]byte, 0, n)
	for i := 0; i < n; i++ {
		ans = append(ans, s[rand.Intn(len(s))])
	}
	return string(ans)
}
func GenerateToken() (string, error) {
	tokenStruct := jwt.NewWithClaims(jwt.SigningMethodHS256, &define.UserClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(time.Hour * 24 * 30),
			},
		},
	})
	token, err := tokenStruct.SignedString(define.Key)
	if err != nil {
		return "", err
	}
	return token, nil
}

func ParseToken(token string) error {
	claims, err := jwt.ParseWithClaims(token, &define.UserClaim{}, func(token *jwt.Token) (interface{}, error) {
		return define.Key, nil
	})
	if err != nil {
		return err
	}
	if !claims.Valid {
		return errors.New("error Token")
	}
	return nil
}
func If(bo bool, a, b interface{}) interface{} {
	if bo {
		return a
	}
	return b
}

func RunShell(shellPath, logPath string) {
	// 分配权限 0777
	cmdChmod := exec.Command("sh", "-c", "chmod +x "+shellPath)
	var errChmod bytes.Buffer
	cmdChmod.Stderr = &errChmod
	if err := cmdChmod.Run(); err != nil {
		log.Println("[CHMOD ERROR] : " + err.Error())
	}
	// 打印当前时间 & 追加日志
	f, err := os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND, 0666)
	if errors.Is(err, os.ErrNotExist) {
		os.MkdirAll(path.Dir(logPath), 0777)
		f, err = os.Create(logPath)
		if err != nil {
			log.Fatalln("[CREATE ERROR] : " + err.Error())
		}
	}
	w := bufio.NewWriter(f)
	w.WriteString(time.Now().Format("2006-01-02 15:04:05") + "\n")
	w.Flush()

	// 运行 & 打印追加日志
	cmdShell := exec.Command("sh", "-c", "nohup "+shellPath+" >> "+logPath+" 2>&1 &")
	var outShell, errShell bytes.Buffer
	cmdShell.Stdout = &outShell
	cmdShell.Stderr = &errShell
	if err := cmdShell.Run(); err != nil {
		log.Println("[SHELL ERROR] : "+err.Error()+" ErrShell : ", errShell.String())
	}
	log.Println(outShell.String())
}

func GetUUID() string {
	return uuid.NewV4().String()
}
func SendSIGINT() {
	process, err := os.FindProcess(define.PID)
	log.Println("---------------------", define.PID)
	if err != nil {
		fmt.Println("无法找到子进程:", err)
		os.Exit(1)
	}
	err = process.Signal(syscall.SIGINT)
	if err != nil {
		fmt.Println("发送信号给子进程失败:", err)
		os.Exit(1)
	}
}

func TouchFile(filePath, data string) error {
	f, err := os.Create(filePath)
	if errors.Is(err, os.ErrNotExist) {
		os.MkdirAll(path.Dir(filePath), 0777)
		f, err = os.Create(filePath)
		if err != nil {
			log.Println("[CREATE FILE ERROR] : " + err.Error())
			return err
		}
	}
	w := bufio.NewWriter(f)
	w.WriteString(data)
	w.Flush()
	return nil
}
