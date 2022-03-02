package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Conf struct {
	Addr            string
	ShutdownTimeout time.Duration
	Router          func(router *gin.Engine)
	Wg              *sync.WaitGroup
}

func InitHttp(conf Conf, middleware ...gin.HandlerFunc) {
	engine := gin.New()

	// sfsfsdsf
	engine.Use(middleware...)

	//注册路由
	conf.Router(engine)

	//监听地址
	l, err := net.Listen("tcp4", conf.Addr)
	if err != nil {
		panic("init http server listen fail :" + err.Error())
	}

	//server配置
	server := &http.Server{
		Addr:         conf.Addr,
		Handler:      engine,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		conf.Wg.Wait()
		ctx, cancel := context.WithTimeout(context.Background(), conf.ShutdownTimeout)
		defer cancel()
		err = server.Shutdown(ctx)
		fmt.Println("----- serve shutdown", time.Now().Format("2006-01-02 15:04:05.999999999"), err)
	}()

	//监听端口
	err = server.Serve(l)
	fmt.Println("----- serve close", time.Now().Format("2006-01-02 15:04:05.999999999"), err)

}

type respEntity struct {
	Code RespCode    `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func RespSuccess(c *gin.Context, data interface{}) {
	c.JSON(200, respEntity{
		Code: RespCodeSuccess,
		Msg:  RespMsg[RespCodeSuccess],
		Data: data,
	})
}

func RespError(c *gin.Context, code RespCode, data interface{}) {
	c.JSON(200, respEntity{
		Code: code,
		Msg:  RespMsg[code],
		Data: data,
	})
}

func BuildParams(params map[string]string) string {

	vals := url.Values{}
	for k, v := range params {
		if v == "" {
			continue
		}
		vals.Add(k, v)
	}

	return vals.Encode()
}

// 发送普通请求
func SendRequest(method, url string, param []byte, header map[string]string) ([]byte, error) {

	client := &http.Client{}
	//req, err := http.NewRequest(method, url, strings.NewReader(postParam))
	req, err := http.NewRequest(method, url, bytes.NewReader(param))
	if err != nil {
		return nil, err
	}
	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if header != nil {
		for key, value := range header {
			req.Header.Set(key, value)
		}
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New("error reponse code:" + string(res.StatusCode))
	}

	content, err := ioutil.ReadAll(res.Body)
	return content, err
}
