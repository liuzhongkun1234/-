package service

import (
	"encoding/json"
	"math"
	"math/rand"
	"sync"
	"time"
)

// 定价json数组
type DeviceData struct {
	DeviceID      string  `json:"device_id"`     // 设备编号
	TemperatureF4 float64 `json:"temperaturef4"` // F4压差值 (纯数字)
	PressureF7    float64 `json:"pressuref7"`    // F7压力值
	HuoXingTan    float64 `json:"huoxingtan"`    // 活性炭（模拟）
	F9YaCha       float64 `json:"f9yacha"`       // F9压差值
	// 换热器1数据
	FuReWenDuCurr1           float64 `json:"furewenducurr1"`           // 换热器1温度
	HuanReQiYaChaCurr1       float64 `json:"huanreqiyachacurr1"`       // 换热器压差1
	ZhuHuoQiYaChaCurr1       float64 `json:"zhuhuoqiyachacurr1"`       // 住火器压差1
	HuanReQiJinKouWenDuCurr1 float64 `json:"huanreqijinkouwenducurr1"` // 换热器进口温度1
	HuanReQiChuKouWenDuCurr1 float64 `json:"huanreqichukouwenducurr1"` // 换热器出口温度1
	// 换热器2数据
	HuanReQiChuKouWenDuCurr2 float64 `json:"huanreqichukouwenducurr2"` // 换热器出口温度2
	HuanReQiJinKouWenDuCurr2 float64 `json:"huanreqijinkouwenducurr2"` // 换热器进口温度2
	FanYingWenDuCurr2        float64 `json:"fanyingwenducurr2"`        // 反应温度2
	JiaReWenDuCurr2          float64 `json:"jiarewenducurr2"`          // 加热温度2
	CuiHuaJinKouWenDuCurr2   float64 `json:"cuihuajinkouwenducurr2"`   // 催化进口温度2
	CuiHuaWenDuCurr2         float64 `json:"cuihuawenducurr2"`         // 催化温度2
	JiaReWenDuCurr3          float64 `json:"jiarewenducurr3"`          // 加热温度3
	//转轮模块
	LunZhuanDinjiPinLV   float64 `json:"lunzhuandingjipinlv"`  // 转轮电机频率
	XiFuQvYaCha          float64 `json:"xifuqyacha"`           // 吸附区压差
	TuoFuQvYaCha         float64 `json:"tuofuqyacha"`          // 脱附区压差
	XiFuQvJinKouWenDu    float64 `json:"xifuqvjinkouwendu"`    // 吸附区进口温度
	XiFuQvChuKouWenDu    float64 `json:"xifuqvchukouwendu"`    // 吸附区出口温度
	PenLinShuiYali       float64 `json:"penlinshuiyali"`       // 喷淋水压力
	TuoFuChuKouWenDu     float64 `json:"tuofuchukouwendu"`     // 脱附区出口温度
	TuoFuJinKouWenDu     float64 `json:"tuofujinkouwendu"`     // 脱附区进口温度
	LengQueQvChUKouWenDu float64 `json:"lengqueqvchukouwendu"` // 冷却区出口温度
	Timestamp            int64   `json:"timestamp"`            // 毫秒级时间戳
}

var (
	mu sync.Mutex // 互斥锁：防止多个客户端同时读写导致数据混乱
	//前处理模拟值
	currentTempF4  float64 = 25.33 // F4压差值 (纯数字)
	currentPresF7  float64 = 25.23 // F7压力值
	HuoXingTanCurr float64 = 50.23 // 活性炭（模拟）
	F9YaChaCurr    float64 = 45.23 // F9压差值
	//轮转点击模块
	LunZhuanDinjiPinLVCurr   float64 = 5.0   // 转轮电机频率
	XiFuQvYaChaCurr          float64 = 41    // 吸附区压差
	TuoFuQvYaChaCurr         float64 = 15.0  // 脱附区压差
	XiFuQvJinKouWenDuCurr    float64 = 23.47 // 吸附区进口温度
	XiFuQvChuKouWenDuCurr    float64 = 43.18 // 吸附区出口温度
	PenLinShuiYaliCurr       float64 = 1.0   // 喷淋水压力
	TuoFuChuKouWenDuCurr     float64 = 53.69 // 脱附区出口温度
	TuoFuJinKouWenDuCurr     float64 = 225.0 // 脱附区进口温度
	LengQueQvChUKouWenDuCurr float64 = 20.0  // 冷却区出口温度
	//换热器1数据
	FuReWenDuCurr1           float64 = 220.0 // 换热器1温度
	HuanReQiYaChaCurr1       float64 = -6.0  // 换热器压差1
	ZhuHuoQiYaChaCurr1       float64 = 0.2   // 住火器压差1
	HuanReQiJinKouWenDuCurr1 float64 = 120.0 // 换热器进口温度1
	HuanReQiChuKouWenDuCurr1 float64 = 225.0 // 换热器出口温度1
	//换热器2数据
	HuanReQiChuKouWenDuCurr2 float64 = 103.0 // 换热器出口温度2
	HuanReQiJinKouWenDuCurr2 float64 = 255.0 // 换热器进口温度2
	FanYingWenDuCurr2        float64 = 430.0 // 反应温度2
	JiaReWenDuCurr2          float64 = 460.0 // 加热温度2
	CuiHuaJinKouWenDuCurr2   float64 = 450.0 // 催化进口温度2
	CuiHuaWenDuCurr2         float64 = 400.0 // 催化温度2
	JiaReWenDuCurr3          float64 = 480.0 // 加热温度3
)

func clampFloat(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// ProcessCommand 处理客户端发来的指令
func ProcessCommand(cmd string) string {
	if cmd != "获取全部数据" {
		return `{"error": "未知指令，请发送：获取全部数据"}`
	}
	mu.Lock()
	defer mu.Unlock()
	//模拟温度数据变换
	//F4压差值在18.35到26.56之间随机波动，步长不超过0.4
	deltaT := (rand.Float64() - 0.5) * 0.4
	currentTempF4 += deltaT
	minF4, maxF4 := 18.35, 26.56
	if currentTempF4 > maxF4 {
		// 超出上限时，在上限附近随机反弹
		currentTempF4 = maxF4 - rand.Float64()*0.1
	} else if currentTempF4 < minF4 {
		// 超出下限时，在下限附近随机反弹
		currentTempF4 = minF4 + rand.Float64()*0.1
	}

	//F7压力值在25.34到35.23之间随机波动，步长不超过0.4
	deltaP := (rand.Float64() - 0.5) * 0.4
	currentPresF7 += deltaP
	minF7, maxF7 := 25.34, 35.23
	if currentPresF7 > maxF7 {
		currentPresF7 = maxF7 - rand.Float64()*0.1
	} else if currentPresF7 < minF7 {
		currentPresF7 = minF7 + rand.Float64()*0.1
	}

	// 活性炭在50.23到56.32之间随机波动，步长不超过0.4
	HuoXingTanCurr += (rand.Float64() - 0.5) * 0.4
	minHuo, maxHuo := 50.23, 56.32
	if HuoXingTanCurr > maxHuo {
		HuoXingTanCurr = maxHuo - rand.Float64()*0.1
	} else if HuoXingTanCurr < minHuo {
		HuoXingTanCurr = minHuo + rand.Float64()*0.1
	}

	// F9压差值在45.23到55.23之间随机波动，步长不超过0.4
	F9YaChaCurr += (rand.Float64() - 0.5) * 0.4
	minF9, maxF9 := 45.23, 55.23
	if F9YaChaCurr > maxF9 {
		F9YaChaCurr = maxF9 - rand.Float64()*0.1
	} else if F9YaChaCurr < minF9 {
		F9YaChaCurr = minF9 + rand.Float64()*0.1
	}
	//沸石轮转点击模块
	// 转轮电机频率在5到6之间随机波动，步长不超过0.5
	LunZhuanDinjiPinLVCurr += (rand.Float64() - 0.5) * 1.0
	if LunZhuanDinjiPinLVCurr > 6 {
		LunZhuanDinjiPinLVCurr = 6 - rand.Float64()*0.5
	} else if LunZhuanDinjiPinLVCurr < 5 {
		LunZhuanDinjiPinLVCurr = 5 + rand.Float64()*0.5
	}
	// 吸附区压差在-12到1之间随机波动，步长不超过0.5
	XiFuQvYaChaCurr += (rand.Float64() - 0.5) * 1.0
	if XiFuQvYaChaCurr > -12 {
		XiFuQvYaChaCurr = -12 - rand.Float64()*0.5
	} else if XiFuQvYaChaCurr < 1 {
		XiFuQvYaChaCurr = 1 + rand.Float64()*0.5
	}
	// 脱附区压差在-12到1之间随机波动，步长不超过0.5
	TuoFuQvYaChaCurr += (rand.Float64() - 0.5) * 1.0
	if TuoFuQvYaChaCurr > -12 {
		TuoFuQvYaChaCurr = -12 - rand.Float64()*0.5
	} else if TuoFuQvYaChaCurr < 1 {
		TuoFuQvYaChaCurr = 1 + rand.Float64()*0.5
	}
	// 吸附区进口温度每次变化不超过10
	XiFuQvJinKouWenDuCurr += (rand.Float64() - 0.5) * 20
	if XiFuQvJinKouWenDuCurr > 33.47 {
		XiFuQvJinKouWenDuCurr = 33.47 - rand.Float64()*0.5
	} else if XiFuQvJinKouWenDuCurr < 13.47 {
		XiFuQvJinKouWenDuCurr = 13.47 + rand.Float64()*0.5
	}
	// 吸附区出口温度每次变化不超过10
	XiFuQvChuKouWenDuCurr += (rand.Float64() - 0.5) * 20
	if XiFuQvChuKouWenDuCurr > 53.18 {
		XiFuQvChuKouWenDuCurr = 53.18 - rand.Float64()*0.5
	} else if XiFuQvChuKouWenDuCurr < 33.18 {
		XiFuQvChuKouWenDuCurr = 33.18 + rand.Float64()*0.5
	}
	// 喷淋水压力每次变化不超过10
	PenLinShuiYaliCurr += (rand.Float64() - 0.5) * 20
	if PenLinShuiYaliCurr > 11 {
		PenLinShuiYaliCurr = 11 - rand.Float64()*0.5
	} else if PenLinShuiYaliCurr < 0 {
		PenLinShuiYaliCurr = rand.Float64() * 0.5
	}
	// 脱附区出口温度每次变化不超过10
	TuoFuChuKouWenDuCurr += (rand.Float64() - 0.5) * 20
	if TuoFuChuKouWenDuCurr > 63.69 {
		TuoFuChuKouWenDuCurr = 63.69 - rand.Float64()*0.5
	} else if TuoFuChuKouWenDuCurr < 43.69 {
		TuoFuChuKouWenDuCurr = 43.69 + rand.Float64()*0.5
	}
	// 脱附区进口温度每次变化不超过10
	TuoFuJinKouWenDuCurr += (rand.Float64() - 0.5) * 20
	if TuoFuJinKouWenDuCurr > 235.0 {
		TuoFuJinKouWenDuCurr = 235.0 - rand.Float64()*0.5
	} else if TuoFuJinKouWenDuCurr < 215.0 {
		TuoFuJinKouWenDuCurr = 215.0 + rand.Float64()*0.5
	}
	// 冷却区出口温度每次变化不超过10
	LengQueQvChUKouWenDuCurr += (rand.Float64() - 0.5) * 20
	if LengQueQvChUKouWenDuCurr > 30.0 {
		LengQueQvChUKouWenDuCurr = 30.0 - rand.Float64()*0.5
	} else if LengQueQvChUKouWenDuCurr < 10.0 {
		LengQueQvChUKouWenDuCurr = 10.0 + rand.Float64()*0.5
	}

	// 换热器波动，步长不超过15
	FuReWenDuCurr1 += (rand.Float64() - 0.5) * 30
	FuReWenDuCurr1 = clampFloat(FuReWenDuCurr1, 205.0, 235.0)
	HuanReQiYaChaCurr1 += (rand.Float64() - 0.5) * 30
	HuanReQiYaChaCurr1 = clampFloat(HuanReQiYaChaCurr1, -21.0, 9.0)
	ZhuHuoQiYaChaCurr1 += (rand.Float64() - 0.5) * 30
	ZhuHuoQiYaChaCurr1 = clampFloat(ZhuHuoQiYaChaCurr1, -14.8, 15.2)
	HuanReQiJinKouWenDuCurr1 += (rand.Float64() - 0.5) * 30
	HuanReQiJinKouWenDuCurr1 = clampFloat(HuanReQiJinKouWenDuCurr1, 105.0, 135.0)
	HuanReQiChuKouWenDuCurr1 += (rand.Float64() - 0.5) * 30
	HuanReQiChuKouWenDuCurr1 = clampFloat(HuanReQiChuKouWenDuCurr1, 210.0, 240.0)

	HuanReQiChuKouWenDuCurr2 += (rand.Float64() - 0.5) * 30
	HuanReQiChuKouWenDuCurr2 = clampFloat(HuanReQiChuKouWenDuCurr2, 88.0, 118.0)
	HuanReQiJinKouWenDuCurr2 += (rand.Float64() - 0.5) * 30
	HuanReQiJinKouWenDuCurr2 = clampFloat(HuanReQiJinKouWenDuCurr2, 240.0, 270.0)
	FanYingWenDuCurr2 += (rand.Float64() - 0.5) * 30
	FanYingWenDuCurr2 = clampFloat(FanYingWenDuCurr2, 415.0, 445.0)
	JiaReWenDuCurr2 += (rand.Float64() - 0.5) * 30
	JiaReWenDuCurr2 = clampFloat(JiaReWenDuCurr2, 445.0, 475.0)
	CuiHuaJinKouWenDuCurr2 += (rand.Float64() - 0.5) * 30
	CuiHuaJinKouWenDuCurr2 = clampFloat(CuiHuaJinKouWenDuCurr2, 435.0, 465.0)
	CuiHuaWenDuCurr2 += (rand.Float64() - 0.5) * 30
	CuiHuaWenDuCurr2 = clampFloat(CuiHuaWenDuCurr2, 385.0, 415.0)
	JiaReWenDuCurr3 += (rand.Float64() - 0.5) * 30
	JiaReWenDuCurr3 = clampFloat(JiaReWenDuCurr3, 465.0, 495.0)

	// 结构花数据
	data := DeviceData{
		DeviceID:      "ShengXiao_ShangHaiRuiNian_001",
		TemperatureF4: math.Round(currentTempF4*100) / 100,
		PressureF7:    math.Round(currentPresF7*100) / 100,
		HuoXingTan:    math.Round(HuoXingTanCurr*100) / 100,
		F9YaCha:       math.Round(F9YaChaCurr*100) / 100,
		//轮转点击模块
		LunZhuanDinjiPinLV:       math.Round(LunZhuanDinjiPinLVCurr*100) / 100,
		XiFuQvYaCha:              math.Round(XiFuQvYaChaCurr*100) / 100, // 吸附区压差
		XiFuQvJinKouWenDu:        math.Round(XiFuQvJinKouWenDuCurr*100) / 100,
		XiFuQvChuKouWenDu:        math.Round(XiFuQvChuKouWenDuCurr*100) / 100,
		PenLinShuiYali:           math.Round(PenLinShuiYaliCurr*100) / 100,
		TuoFuChuKouWenDu:         math.Round(TuoFuChuKouWenDuCurr*100) / 100,
		TuoFuJinKouWenDu:         math.Round(TuoFuJinKouWenDuCurr*100) / 100,
		LengQueQvChUKouWenDu:     math.Round(LengQueQvChUKouWenDuCurr*100) / 100,
		FuReWenDuCurr1:           math.Round(FuReWenDuCurr1*100) / 100,
		HuanReQiYaChaCurr1:       math.Round(HuanReQiYaChaCurr1*100) / 100,
		ZhuHuoQiYaChaCurr1:       math.Round(ZhuHuoQiYaChaCurr1*100) / 100,
		HuanReQiJinKouWenDuCurr1: math.Round(HuanReQiJinKouWenDuCurr1*100) / 100,
		HuanReQiChuKouWenDuCurr1: math.Round(HuanReQiChuKouWenDuCurr1*100) / 100,
		HuanReQiChuKouWenDuCurr2: math.Round(HuanReQiChuKouWenDuCurr2*100) / 100,
		HuanReQiJinKouWenDuCurr2: math.Round(HuanReQiJinKouWenDuCurr2*100) / 100,
		FanYingWenDuCurr2:        math.Round(FanYingWenDuCurr2*100) / 100,
		JiaReWenDuCurr2:          math.Round(JiaReWenDuCurr2*100) / 100,
		CuiHuaJinKouWenDuCurr2:   math.Round(CuiHuaJinKouWenDuCurr2*100) / 100,
		CuiHuaWenDuCurr2:         math.Round(CuiHuaWenDuCurr2*100) / 100,
		JiaReWenDuCurr3:          math.Round(JiaReWenDuCurr3*100) / 100,
		//换热器
		Timestamp:    time.Now().UnixMilli(),
		TuoFuQvYaCha: math.Round(TuoFuQvYaChaCurr*100) / 100,
	}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return `{"error": "数据序列化失败"}`
	}
	return string(jsonBytes)
}
