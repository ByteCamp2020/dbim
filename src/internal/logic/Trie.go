package logic

const (
	INIT_TRIE_CHILDREN_NUM = 200 // > 26
)

type trieNode struct {
	isEndOfWord bool
	children map[rune]*trieNode
}

func newtrieNode() *trieNode {
	return &trieNode{
		isEndOfWord: false,
		children:    make(map[rune]*trieNode, INIT_TRIE_CHILDREN_NUM),
	}
}

// Match index object
type matchIndex struct {
	start int // start index
	end   int // end index
}

// Construct from scratch
func newMatchIndex(start, end int) *matchIndex {
	return &matchIndex{
		start: start,
		end:   end,
	}
}

// Construct from existing match index object
func buildMatchIndex(obj *matchIndex) *matchIndex {
	return &matchIndex{
		start: obj.start,
		end:   obj.end,
	}
}

// dfa util
type DFAUtil struct {
	// The root node
	root *trieNode
}

func (this *DFAUtil) insertWord(word []rune) {
	currNode := this.root
	for _, c := range word {
		if cildNode, exist := currNode.children[c]; !exist {
			cildNode = newtrieNode()
			currNode.children[c] = cildNode
			currNode = cildNode
		} else {
			currNode = cildNode
		}
	}

	currNode.isEndOfWord = true
}

// Check if there is any word in the trie that starts with the given prefix.
func (this *DFAUtil) startsWith(prefix []rune) bool {
	currNode := this.root
	for _, c := range prefix {
		if cildNode, exist := currNode.children[c]; !exist {
			return false
		} else {
			currNode = cildNode
		}
	}

	return true
}

// Searc and make sure if a word is existed in the underlying trie.
func (this *DFAUtil) searchWord(word []rune) bool {
	currNode := this.root
	for _, c := range word {
		if cildNode, exist := currNode.children[c]; !exist {
			return false
		} else {
			currNode = cildNode
		}
	}

	return currNode.isEndOfWord
}

func (this *DFAUtil) searchSentence(sentence string) bool {
	start, end := 0, 1
	sentenceRuneList := []rune(sentence)

	// Iterate the sentence from the beginning to the end.
	startsWith := false
	for end <= len(sentenceRuneList) {
		// Check if a sensitive word starts with word range from [start:end)
		// We find the longest possible path
		// Then we check any sub word is the sensitive word from long to short
		if this.startsWith(sentenceRuneList[start:end]) {
			startsWith = true
			end += 1
		} else {
			if startsWith == true {
				// Check any sub word is the sensitive word from long to short
				for index := end - 1; index > start; index-- {
					if this.searchWord(sentenceRuneList[start:index]) {
						return false
					}
				}
			}
			start, end = end-1, end+1
			startsWith = false
		}
	}

	// If finishing not because of unmatching, but reaching the end, we need to
	// check if the previous startsWith is true or not.
	// If it's true, we need to check if there is any candidate?
	if startsWith {
		for index := end - 1; index > start; index-- {
			if this.searchWord(sentenceRuneList[start:index]) {
				return false
			}
		}
	}

	return true
}

//
func (this *DFAUtil) CheckSentence(sentence string) bool {
	return this.searchSentence(sentence)
}

// Create new DfaUtil object
// wordList:word list
func NewDFAUtil(wordList []string) *DFAUtil {
	this := &DFAUtil{
		root: newtrieNode(),
	}

	for _, word := range wordList {
		wordRuneList := []rune(word)
		if len(wordRuneList) > 0 {
			this.insertWord(wordRuneList)
		}
	}

	return this
}