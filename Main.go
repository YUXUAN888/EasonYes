package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

const HELLODLLPATH = "shushu.dll"

var (
	MediaDll  = syscall.NewLazyDLL(HELLODLLPATH)
	startPlay = MediaDll.NewProc("startPlay") // 开始播放mp3
	//stopPlay     = MediaDll.NewProc("stopPlay") // 停止播放mp3
	continuePlay = MediaDll.NewProc("continuePlay") //继续播放mp3
	suspendPlay  = MediaDll.NewProc("suspendPlay")  // 暂停播放mp3
)

func main() {
	print(" _____                   __   __        " + "\n")
	time.Sleep(time.Second * 1)
	print("| ____|__ _ ___  ___  _ _\\ \\ / /__  ___ " + "\n")
	time.Sleep(time.Second * 1)
	print("|  _| / _` / __|/ _ \\| '_ \\ V / _ \\/ __|" + "\n")
	time.Sleep(time.Second * 1)
	print("| |__| (_| \\__ \\ (_) | | | | |  __/\\__ \\" + "\n")
	time.Sleep(time.Second * 1)
	print("|_____\\__,_|___/\\___/|_| |_|_|\\___||___/" + "\n")
	time.Sleep(time.Second * 1)
	print("----------------欢迎来到 EasonYes ----------------\n")
	print("----------------请键入以下内容----------------\n")
	print("----------------'-h'：获取帮助----------------\n")
	//var a string
	for true {
		print("EasonYes> ")
		reader := bufio.NewReader(os.Stdin)
		res, _, err := reader.ReadLine()
		if nil != err {
			fmt.Println("reader.ReadLine() error:", err)
		}
		GetScan(string(res))
	}
}
func GetScan(ds string) {
	if ds == "-h" {
		print("-本软件是 Golang 开发的，快速的，方便的，跨平台的陈奕迅音乐软件。-\n" + "-在这里所有Eason的歌曲，专辑，全部免费下载，且最高音质-\n")
		print("作者) Yu_Xuan\n")
		print("版本) EasonYes 1.0\n")
		print("1) -x 获取 FSM Server 上无损的Eason歌曲\n")
		print("2) -d 下载歌曲，用法：-d 歌曲名\n")
		print("3) -p 播放歌曲，用法：-p 歌曲名（注意必须先下载再播放）\n")
	} else if ds == "-x" {
		print("------------下面是在 FSM Server 获取到的Eason歌曲（若需要下载，请键入'-d 歌曲名'；若需要播放，请键入'-p 歌曲名'）------------\n")
		url := "http://101.43.69.145:666/Go/EasonYes/MusicList.json"
		//将要读取的网站放入到get方法中
		resp, err := http.Get(url)
		if err != nil {
			panic(err.Error())
		}
		//读取数据
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("读取数据错误:", err)
			return
		}
		//关闭链接
		defer resp.Body.Close()
		//fmt.Println(string(bytes))
		var kkk Eason
		err = json.Unmarshal([]byte(bytes), &kkk)
		if err != nil {
			fmt.Println("反序列化失败:", err)
		}
		for i := 0; i < len(kkk.Yeah); i++ {
			print(kkk.Yeah[i] + "\n")
		}
		print("-------------------------------------------\n")
	} else if strings.HasPrefix(ds, "-d") {
		//下载
		array := strings.Fields(ds)
		if array == nil {
			print("错误的命令\n")
			return
		}
		music := array[1]
		print("-开始下载 " + music + " , 将保存到程序目录的'EasonYes'文件夹-\n")
		print("-正在下载中......-\n")
		dir, _ := os.Getwd()
		DownloadFile("http://101.43.69.145:1469/file/EasonMusicList/"+music+".mp3", dir+"/EasonYes/"+music+".mp3")
		print("-下载完成!-\n")
	} else if strings.HasPrefix(ds, "-p") {
		// 播放
		array := strings.Fields(ds)
		if array == nil {
			print("错误的命令\n")
			return
		}
		dir, _ := os.Getwd()
		music := array[1]
		print("-开始准备播放音乐 " + music + " -\n-若要停止播放，请键入'-n'-\n")
		if Exists(dir + "/EasonYes/" + music + ".mp3") {
			go (func() {
				_, _, _ = startPlay.Call(strPtr(dir + "/EasonYes/" + music + ".mp3"))
			})()
		} else {
			print("-你还没下载 " + music + " 请先下载再播放-\n")
		}
	} else if ds == "-n" {
		print("-已停止播放音乐-\n-若要继续播放，请键入'-c'-\n")
		_, _, _ = suspendPlay.Call()
	} else if ds == "-c" {
		print("-已继续播放音乐-\n-若要停止播放，请键入'-n'-\n")
		_, _, _ = continuePlay.Call()
	} else {
		print("-没有此命令-\n")
		return
	}
}

func Exists(path string) bool {

	_, err := os.Stat(path) //os.Stat获取文件信息

	if err != nil {

		if os.IsExist(err) {

			return true

		}

		return false

	}

	return true

}

func strPtr(s string) uintptr {
	news, _ := syscall.BytePtrFromString(s)
	return uintptr(unsafe.Pointer(news))
}

func DownloadFile(url, path string) (pathRes, downloadRes bool) {
	var res *http.Response
	tt := strings.Split(path, "/")
	pathDir := strings.TrimSuffix(path, tt[len(tt)-1])
	pathDir = strings.TrimSuffix(pathDir, "/")
	err := os.MkdirAll(pathDir, 0777)
	if err != nil {
		fmt.Printf("mkdir file failed:[%s], error:[%v], path:[%s]\n", pathDir, err, path)
		return false, false
	} else {
		fmt.Printf("mkdir file success:[%s], error:[%v], path:[%s]\n", pathDir, err, path)
	}
	buffer, err := os.Create(path)
	if err != nil {
		fmt.Printf("open file failed, error:%v, path:[%s]\n", err, path)
		return false, false
	} else {
		fmt.Printf("open file success, error:%v, path:[%s]\n", err, path)
	}
	defer buffer.Close()
	//flag 0:work; 1:fail; 2:over
	flag := 0
	for i := 0; i < 3; i++ {
		for {
			//res, err = http.Get(url)
			client := &http.Client{}
			req, err2 := http.NewRequest("GET", url, nil)
			if err2 != nil {
				fmt.Printf("make new http request [%s] error: [%v]\n", url, err2)
				return true, false
			}
			currentDate := time.Now().Format("Mon, 02 Jan 2006 15:04:05 GMT")

			req.Header.Add("Date", currentDate)
			//req.Header.Add("Authorization", auth)
			//req.Header.Add("Content-Type","video/mp4")
			res, err2 = client.Do(req)
			if err2 != nil {
				fmt.Printf("send http req [%s] error: [%v]\n", url, err2)
				return true, false
			}

			if res.StatusCode != 200 {
				fmt.Printf("download failed:[%s], StatusCode:[%d],url:[%s]\n", path, res.StatusCode, url)
				res.Body.Close()
				if i >= 2 {
					return true, false
				} else {
					flag = 1
					break
				}
			}
			flag = 0
			break
		}
		if flag == 1 {
			time.Sleep(1 * time.Second)
			continue
		}
		buf := make([]byte, 102400)

		for {
			n, err := res.Body.Read(buf)
			if err != nil && err != io.EOF {
				fmt.Printf("Download failed:[%s], error:%v, path:[%s]\n", url, err, path)
				res.Body.Close()
				if i >= 2 {
					return true, false
				} else {
					flag = 1
					break
				}
			}
			buffer.Write(buf[:n])
			if err == io.EOF {
				res.Body.Close()
				flag = 2
				break
			}
		}
		if flag == 1 {
			time.Sleep(1 * time.Second)
			continue
		} else if flag == 2 {
			break
		}
	}
	fmt.Printf("Download success:[%s], path:[%s]\n", url, path)
	return true, true
}

type Eason struct {
	Yeah []string `json:"Yeah"`
}
