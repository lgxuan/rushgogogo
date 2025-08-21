package filter

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"regexp"
	handlercontext "rushgogogo/internal/handlerContext"
	"rushgogogo/pkgs/utils"
	"strings"
	"sync"
)

var (
	seen               = make(map[string]struct{})
	seenMutex          sync.RWMutex
	seenSensitiveInfo  = make(map[string]map[string]string)
	sensitiveInfoMutex sync.RWMutex
	workerPool         chan struct{}
	unusefulPath       = []string{
		".ico",
		".png",
		".jpg",
		".gif",
		".jpeg",
		".webp",
		".svg",
		".woff",
		".woff2",
		".ttf",
		".eot",
		".mp4",
		".mp3",
		".avi",
		".mov",
		".wmv",
		".flv",
		".pdf",
		".zip",
		".rar",
		".7z",
		".tar",
		".gz",
		".bz2",
	}
)

func init() {
	maxWorkers := getMaxWorkers()
	workerPool = make(chan struct{}, maxWorkers)
	fmt.Printf("init worker pool: %d\n", maxWorkers)
}

func getMaxWorkers() int {
	defaultWorkers := 10
	return defaultWorkers
}

type InformationFilter struct {
	filterName string
	filterType string
	filterReg  *regexp.Regexp
	source     string
	level      string
	enabled    bool
}

func (f *InformationFilter) Filter(html string, source string) string {
	if f.filterType == "reg" {
		information := f.filterReg.FindAllString(html, -1)
		for _, info := range information {
			// 全局去重检查，避免不同过滤器重复打印相同内容
			sensitiveInfoMutex.Lock()
			if sourceMap, exists := seenSensitiveInfo[info]; exists {
				// 如果这个敏感信息在同一个来源中已经被发现过，跳过
				if _, sourceExists := sourceMap[source]; sourceExists {
					sensitiveInfoMutex.Unlock()
					continue
				}
			} else {
				// 初始化这个敏感信息来源map
				seenSensitiveInfo[info] = make(map[string]string)
			}
			// 记录这个敏感信息已经被发现
			seenSensitiveInfo[info][source] = f.filterName
			sensitiveInfoMutex.Unlock()

			inf := fmt.Sprintf("%s found information %s : %s", source, f.filterName, info)
			utils.Log(inf, f.level)
		}
	}
	return html
}
func NewInformationFilter(filterName string, filterType string, filterReg *regexp.Regexp, source string, level string, enabled bool) InformationFilter {
	return InformationFilter{
		filterName: filterName,
		filterType: filterType,
		filterReg:  filterReg,
		source:     source,
		level:      level,
		enabled:    enabled,
	}
}

// FilterResponse 异步过滤响应
func FilterResponse(w http.Response, filters []InformationFilter) *http.Response {

	// 先读取响应体
	data, err := io.ReadAll(w.Body)
	if err != nil {
		return &w
	}
	key := sha256.Sum256(data)
	keyHex := hex.EncodeToString(key[:])

	// 使用读写锁防止重复处理
	seenMutex.RLock()
	if _, ok := seen[keyHex]; ok {
		seenMutex.RUnlock()
		return &w
	}
	seenMutex.RUnlock()

	// 获取写锁并再次检查（双重检查锁定模式）
	seenMutex.Lock()
	if _, ok := seen[keyHex]; ok {
		seenMutex.Unlock()
		return &w
	}
	seen[keyHex] = struct{}{}
	seenMutex.Unlock()
	// 创建新的响应体
	w.Body = io.NopCloser(bytes.NewBuffer(data))

	// 使用goroutine池限制并发数量
	select {
	case workerPool <- struct{}{}: // 获取一个工作槽位
		// 异步处理，不影响响应速度
		go func() {
			defer func() { <-workerPool }() // 处理完成后释放工作槽位

			// 检查是否为无用路径
			for _, i := range unusefulPath {
				if strings.HasSuffix(w.Request.URL.Path, i) {
					// 结束异步处理
					return
				}
			}
			// 转换为UTF-8字符串
			html, err := utils.ConvertToUTF8(data)
			if err != nil {
				// 如果转换失败，使用原始数据
				html = string(data)
			}

			// 应用过滤器
			for _, filter := range filters {
				if !filter.enabled {
					continue
				}
				filter.Filter(html, w.Request.URL.String())
			}
		}()
	default:
	}
	return &w
}

// FilterWithContext 使用handlercontext进行过滤
func FilterWithContext(proxyCtx *handlercontext.HandlerContext) {
	// 这里可以使用handlercontext包的功能
	ctx := proxyCtx.GetProxyCtx()
	_ = ctx // 避免未使用变量警告
}
