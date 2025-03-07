package cache

import (
	"strconv"
	"sync"
	"time"
)

type MemCacheUtil struct {
	Data                  map[string]string
	DataRWLock            sync.RWMutex
	CaptchaCacheMaxNumber int
}

func NewMemCacheUtil(captchaCacheMaxNumber int) *MemCacheUtil {
	return &MemCacheUtil{
		Data:                  make(map[string]string),
		CaptchaCacheMaxNumber: captchaCacheMaxNumber,
	}
}

func (l *MemCacheUtil) Exists(key string) bool {
	l.DataRWLock.RLock()
	timeVal := l.Data[key+"_HoldTime"]
	cacheHoldTime, err := strconv.ParseInt(timeVal, 10, 64)
	l.DataRWLock.RUnlock()

	if err != nil {
		return false
	}

	if cacheHoldTime == 0 {
		return true
	}

	if cacheHoldTime < time.Now().Unix() {
		l.Delete(key)
		return false
	}
	return true
}

func (l *MemCacheUtil) Get(key string) string {

	if l.Exists(key) {
		l.DataRWLock.RLock()
		val := l.Data[key]
		l.DataRWLock.RUnlock()

		return val
	}

	return ""
}

func (l *MemCacheUtil) Set(key string, val string, expiresInSeconds int) {

	//设置阈值，达到即clear缓存
	if len(l.Data) >= l.CaptchaCacheMaxNumber*2 {
		//		log.Println("CACHE_MAP达到阈值，clear map")
		l.Clear()
	}

	l.DataRWLock.Lock()
	l.Data[key] = val
	if expiresInSeconds > 0 {
		// 缓存失效时间
		nowTime := time.Now().Unix() + int64(expiresInSeconds)
		l.Data[key+"_HoldTime"] = strconv.FormatInt(nowTime, 10)
	} else {
		l.Data[key+"_HoldTime"] = strconv.FormatInt(0, 10)
	}

	l.DataRWLock.Unlock()
}

func (l *MemCacheUtil) Delete(key string) {
	l.DataRWLock.Lock()
	defer l.DataRWLock.Unlock()
	delete(l.Data, key)
	delete(l.Data, key+"_HoldTime")
}

func (l *MemCacheUtil) Clear() {
	for key, _ := range l.Data {
		l.Delete(key)
	}
}
