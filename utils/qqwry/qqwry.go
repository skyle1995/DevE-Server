package qqwry

import (
	"bytes"
	_ "embed"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"strings"
	"sync"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

//go:embed qqwry.dat
var ipData []byte

var (
	data    []byte
	dataLen uint32
	ipCache = &sync.Map{}
)

const (
	indexLen      = 7
	redirectMode1 = 0x01
	redirectMode2 = 0x02
)

type IPInfo struct {
	City string // 城市
	ISP  string // 运营商
}

// 初始化函数，自动加载嵌入的IP数据库
func init() {
	LoadData(ipData)
}

// byte3ToUInt32 将3字节转换为uint32
func byte3ToUInt32(data []byte) uint32 {
	i := uint32(data[0]) & 0xff
	i |= (uint32(data[1]) << 8) & 0xff00
	i |= (uint32(data[2]) << 16) & 0xff0000
	return i
}

// gb18030Decode 将GB18030编码转换为UTF-8
func gb18030Decode(src []byte) string {
	in := bytes.NewReader(src)
	out := transform.NewReader(in, simplifiedchinese.GB18030.NewDecoder())
	d, _ := io.ReadAll(out)
	return string(d)
}

// QueryIP 查询IP信息
// 返回城市和运营商信息
func QueryIP(queryIP string) (city string, isp string, err error) {
	// 先从缓存中查询
	if v, ok := ipCache.Load(queryIP); ok {
		cacheInfo := v.(IPInfo)
		return cacheInfo.City, cacheInfo.ISP, nil
	}

	// 解析IP地址
	ip := net.ParseIP(queryIP).To4()
	if ip == nil {
		return "", "", errors.New("IP地址不是有效的IPv4地址")
	}

	// 转换为uint32
	ip32 := binary.BigEndian.Uint32(ip)

	// 获取索引范围
	posA := binary.LittleEndian.Uint32(data[:4])
	posZ := binary.LittleEndian.Uint32(data[4:8])

	// 二分查找
	var offset uint32 = 0
	for {
		mid := posA + (((posZ-posA)/indexLen)>>1)*indexLen
		buf := data[mid : mid+indexLen]
		_ip := binary.LittleEndian.Uint32(buf[:4])

		if posZ-posA == indexLen {
			offset = byte3ToUInt32(buf[4:])
			buf = data[mid+indexLen : mid+indexLen+indexLen]
			if ip32 < binary.LittleEndian.Uint32(buf[:4]) {
				break
			} else {
				offset = 0
				break
			}
		}

		if _ip > ip32 {
			posZ = mid
		} else if _ip < ip32 {
			posA = mid
		} else if _ip == ip32 {
			offset = byte3ToUInt32(buf[4:])
			break
		}
	}

	if offset <= 0 {
		return "", "", errors.New("未找到IP信息")
	}

	// 解析地址信息
	posM := offset + 4
	mode := data[posM]
	var ispPos uint32

	switch mode {
	case redirectMode1:
		posC := byte3ToUInt32(data[posM+1 : posM+4])
		mode = data[posC]
		posCA := posC
		if mode == redirectMode2 {
			posCA = byte3ToUInt32(data[posC+1 : posC+4])
			posC += 4
		}
		for i := posCA; i < dataLen; i++ {
			if data[i] == 0 {
				city = string(data[posCA:i])
				break
			}
		}
		if mode != redirectMode2 {
			posC += uint32(len(city) + 1)
		}
		ispPos = posC
	case redirectMode2:
		posCA := byte3ToUInt32(data[posM+1 : posM+4])
		for i := posCA; i < dataLen; i++ {
			if data[i] == 0 {
				city = string(data[posCA:i])
				break
			}
		}
		ispPos = offset + 8
	default:
		posCA := offset + 4
		for i := posCA; i < dataLen; i++ {
			if data[i] == 0 {
				city = string(data[posCA:i])
				break
			}
		}
		ispPos = offset + uint32(5+len(city))
	}

	// 转换城市编码
	if city != "" {
		city = strings.TrimSpace(gb18030Decode([]byte(city)))
	}

	// 解析ISP信息
	ispMode := data[ispPos]
	if ispMode == redirectMode1 || ispMode == redirectMode2 {
		ispPos = byte3ToUInt32(data[ispPos+1 : ispPos+4])
	}

	if ispPos > 0 {
		for i := ispPos; i < dataLen; i++ {
			if data[i] == 0 {
				isp = string(data[ispPos:i])
				if isp != "" {
					if strings.Contains(isp, "CZ88.NET") {
						isp = ""
					} else {
						isp = strings.TrimSpace(gb18030Decode([]byte(isp)))
					}
				}
				break
			}
		}
	}

	// 存入缓存
	ipCache.Store(queryIP, IPInfo{City: city, ISP: isp})
	return city, isp, nil
}

// GetIPLocation 获取IP地址的位置信息
// 返回格式：城市 运营商
func GetIPLocation(ip string) string {
	city, isp, err := QueryIP(ip)
	if err != nil {
		return ""
	}
	return city + " " + isp
}

// GetIPCity 获取IP地址的城市信息
func GetIPCity(ip string) (string, error) {
	city, _, err := QueryIP(ip)
	if err != nil {
		return "", err
	}
	return city, nil
}

// GetIPISP 获取IP地址的运营商信息
func GetIPISP(ip string) (string, error) {
	_, isp, err := QueryIP(ip)
	if err != nil {
		return "", err
	}
	return isp, nil
}

// LoadData 从内存加载IP数据库
func LoadData(database []byte) {
	data = database
	dataLen = uint32(len(data))
}