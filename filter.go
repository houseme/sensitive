package sensitive

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"
)

// Filter 敏感词过滤器
type Filter struct {
	trie       *Trie
	noise      *regexp.Regexp
	buildVer   int64
	updatedVer int64
}

// New 返回一个敏感词过滤器
func New() *Filter {
	return &Filter{
		trie:  NewTrie(),
		noise: regexp.MustCompile(`[\|\s&%$@*]+`),
	}
}

// UpdateNoisePattern 更新去噪模式
func (ft *Filter) UpdateNoisePattern(pattern string) {
	ft.noise = regexp.MustCompile(pattern)
}

// LoadWordDict 加载敏感词字典
func (ft *Filter) LoadWordDict(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	return ft.Load(f)
}

// LoadNetWordDict 加载网络敏感词字典
func (ft *Filter) LoadNetWordDict(url string) error {
	c := http.Client{
		Timeout: 5 * time.Second,
	}
	rsp, err := c.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		_ = rsp.Body.Close()
	}()

	return ft.Load(rsp.Body)
}

// Load common method to add words
func (ft *Filter) Load(rd io.Reader) error {
	buf := bufio.NewReader(rd)
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		ft.AddWord(string(line))
	}

	return nil
}

func (ft *Filter) updateFailureLink() {
	if ft.buildVer != ft.updatedVer {
		// fmt.Println("update failure link")
		ft.trie.BuildFailureLinks()
		ft.buildVer = ft.updatedVer
	}
}

// AddWord 添加敏感词
func (ft *Filter) AddWord(words ...string) {
	ft.trie.Add(words...)
	ft.updatedVer = time.Now().UnixNano()
}

// DelWord 删除敏感词
func (ft *Filter) DelWord(words ...string) {
	ft.trie.Del(words...)
}

// Filter 过滤敏感词
func (ft *Filter) Filter(text string) string {
	ft.updateFailureLink()
	return ft.trie.Filter(text)
}

// Replace 和谐敏感词
func (ft *Filter) Replace(text string, repl rune) string {
	ft.updateFailureLink()
	return ft.trie.Replace(text, repl)
}

// FindIn 检测敏感词
func (ft *Filter) FindIn(text string) (bool, string) {
	ft.updateFailureLink()
	text = ft.RemoveNoise(text)
	return ft.trie.FindIn(text)
}

// FindAll 找到所有匹配词
func (ft *Filter) FindAll(text string) []string {
	ft.updateFailureLink()
	return ft.trie.FindAll(text)
}

// Validate 检测字符串是否合法
func (ft *Filter) Validate(text string) (bool, string) {
	ft.updateFailureLink()
	text = ft.RemoveNoise(text)
	return ft.trie.Validate(text)
}

// RemoveNoise 去除空格等噪音
func (ft *Filter) RemoveNoise(text string) string {
	return ft.noise.ReplaceAllString(text, "")
}

// ValidateWithWildcard 检测字符串是否合法，支持通配符
func (ft *Filter) ValidateWithWildcard(text string, wildcard rune) (bool, string) {
	text = ft.RemoveNoise(text)
	return ft.trie.ValidateWithWildcard(text, wildcard)
}
