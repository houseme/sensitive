package sensitive

// Trie 短语组成的 Trie 树。
type Trie struct {
	Root *Node
}

// Node Trie 树上的一个节点。
type Node struct {
	isRootNode bool
	isPathEnd  bool
	Character  rune
	Children   map[rune]*Node
	Failure    *Node
	Parent     *Node
	depth      int
}

// BuildFailureLinks 更新 Aho-Corasick 的失败表
func (tree *Trie) BuildFailureLinks() {
	for node := range tree.bfs() {
		pointer := node.Parent
		var link *Node
		for link == nil {
			if pointer.IsRootNode() {
				link = pointer
				break
			}
			link = pointer.Failure.Children[node.Character]
			pointer = pointer.Failure
		}
		// fmt.Printf("%s[%d] link to %s[%d] \n", string(node.Character), node.depth, string(link.Character), link.depth)
		node.Failure = link
	}
	// fmt.Println("finish build failure link")
}

// bfs Breadth First Search
func (tree *Trie) bfs() <-chan *Node {
	ch := make(chan *Node)
	go func() {
		queue := new(LinkList)
		for _, child := range tree.Root.Children {
			queue.Push(child)
		}

		for !queue.Empty() {
			n := queue.Pop().(*Node)
			ch <- n
			for _, child := range n.Children {
				queue.Push(child)
			}
		}

		close(ch)
	}()
	return ch
}

// NewTrie 新建一棵 Trie
func NewTrie() *Trie {
	return &Trie{
		Root: NewRootNode(0),
	}
}

// Add 添加若干个词
func (tree *Trie) Add(words ...string) {
	for _, word := range words {
		tree.add(word)
	}
}

func (tree *Trie) add(word string) {
	current := tree.Root
	runes := []rune(word)
	for position := 0; position < len(runes); position++ {
		r := runes[position]
		if next, ok := current.Children[r]; ok {
			current = next
		} else {
			newNode := NewNode(r)
			newNode.depth = current.depth + 1
			newNode.Parent = current
			current.Children[r] = newNode
			current = newNode
		}
	}
	current.isPathEnd = true
}

func (tree *Trie) Del(words ...string) {
	for _, word := range words {
		tree.del(word)
	}
}

func (tree *Trie) del(word string) {
	current := tree.Root
	runes := []rune(word)
	for position := 0; position < len(runes); position++ {
		r := runes[position]
		if next, ok := current.Children[r]; !ok {
			return
		} else {
			current = next
		}

		if position == len(runes)-1 {
			current.SoftDel()
		}
	}
}

// Replace 词语替换
func (tree *Trie) Replace(text string, character rune) string {
	var (
		node  = tree.Root
		next  *Node
		runes = []rune(text)
	)

	var ac = new(ac)
	for position := 0; position < len(runes); position++ {
		next = ac.next(node, runes[position])
		if next == nil {
			next = ac.fail(node, runes[position])
		}

		node = next
		ac.replace(node, runes, position, character)
	}

	return string(runes)
}

// Filter 直接过滤掉字符串中的敏感词
func (tree *Trie) Filter(text string) string {
	var (
		parent      = tree.Root
		current     *Node
		left        = 0
		found       bool
		runes       = []rune(text)
		length      = len(runes)
		resultRunes = make([]rune, 0, length)
	)

	for position := 0; position < length; position++ {
		current, found = parent.Children[runes[position]]

		if !found || (!current.IsPathEnd() && position == length-1) {
			resultRunes = append(resultRunes, runes[left])
			parent = tree.Root
			position = left
			left++
			continue
		}

		if current.IsPathEnd() {
			left = position + 1
			parent = tree.Root
		} else {
			parent = current
		}
	}

	resultRunes = append(resultRunes, runes[left:]...)
	return string(resultRunes)
}

// Validate 验证字符串是否合法，如不合法则返回 false 和检测到
// 的第一个敏感词
func (tree *Trie) Validate(text string) (bool, string) {
	const EMPTY = ""
	var (
		node  = tree.Root
		next  *Node
		runes = []rune(text)
	)

	var ac = new(ac)
	for position := 0; position < len(runes); position++ {
		next = ac.next(node, runes[position])
		if next == nil {
			next = ac.fail(node, runes[position])
		}

		node = next
		if first := ac.firstOutput(node, runes, position); len(first) > 0 {
			return false, first
		}
	}

	return true, EMPTY
}

// ValidateWithWildcard 验证字符串是否合法，如不合法则返回 false 和检测到
func (tree *Trie) ValidateWithWildcard(text string, wildcard rune) (bool, string) {
	runes := []rune(text)
	for curl := 0; curl < len(runes); curl++ {
		patter := ""
		parent := tree.Root
		if tree.dfs(runes, parent, curl, wildcard, "", &patter) {
			return false, patter
		}
	}
	return true, ""
}

func (tree *Trie) dfs(runes []rune, parent *Node, curl int, wildcard rune, str string, patter *string) bool {
	if parent == nil {
		return false
	}
	if parent.IsPathEnd() {
		*patter = str
		return true
	}
	if curl >= len(runes) {
		return false
	}

	// 匹配到了
	if current, found := parent.Children[runes[curl]]; found {
		if is1 := tree.dfs(runes, current, curl+1, wildcard, str+string(runes[curl]), patter); is1 {
			return true
		}
	}

	// 先看有没有*
	if current1, found1 := parent.Children[wildcard]; found1 {
		if is2 := tree.dfs(runes, current1, curl+1, wildcard, str+string(wildcard), patter); is2 {
			return true
		}

		if current2, found2 := current1.Children[runes[curl]]; found2 {
			if is3 := tree.dfs(runes, current2, curl+1, wildcard, str+string(wildcard)+string(runes[curl]), patter); is3 {
				return true
			}
		}
	}
	return false
}

// FindIn 判断 text 中是否含有词库中的词
func (tree *Trie) FindIn(text string) (bool, string) {
	validated, first := tree.Validate(text)
	return !validated, first
}

// FindAll 找有所有包含在词库中的词
func (tree *Trie) FindAll(text string) []string {
	var (
		node  = tree.Root
		next  *Node
		runes = []rune(text)
	)

	var ac = new(ac)
	for position := 0; position < len(runes); position++ {
		next = ac.next(node, runes[position])
		if next == nil {
			next = ac.fail(node, runes[position])
		}

		node = next
		ac.output(node, runes, position)
	}

	return ac.results
}

// NewNode 新建子节点
func NewNode(character rune) *Node {
	return &Node{
		Character: character,
		Children:  make(map[rune]*Node, 0),
	}
}

// NewRootNode 新建根节点
func NewRootNode(character rune) *Node {
	root := &Node{
		isRootNode: true,
		Character:  character,
		Children:   make(map[rune]*Node, 0),
		depth:      0,
	}

	root.Failure = root

	return root
}

// IsLeafNode 判断是否叶子节点
func (node *Node) IsLeafNode() bool {
	return len(node.Children) == 0
}

// IsRootNode 判断是否为根节点
func (node *Node) IsRootNode() bool {
	return node.isRootNode
}

// IsPathEnd 判断是否为某个路径的结束
func (node *Node) IsPathEnd() bool {
	return node.isPathEnd
}

// SoftDel 置软删除状态
func (node *Node) SoftDel() {
	node.isPathEnd = false
}
