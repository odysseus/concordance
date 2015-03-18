package concordance

import (
	"bufio"
	"fmt"
	"sort"
	"strings"
)

type WordTuple struct {
	Word  string
	Count int
}

func (t *WordTuple) String() string {
	return fmt.Sprintf("%v: %v", t.Word, t.Count)
}

type ByCount []WordTuple

func (s ByCount) Len() int {
	return len(s)
}

func (s ByCount) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByCount) Less(i, j int) bool {
	return s[i].Count < s[j].Count
}

type Concordance struct {
	Counts          map[string]int
	Total           int
	Unique          int
	MostUsed        ByCount
	LengthHistogram []int
}

func NewConcordance(scanner *bufio.Scanner, caseSensitive bool, topwords int) *Concordance {
	c := &Concordance{}
	c.Counts, c.Total = WordCount(scanner, caseSensitive)
	c.Unique = len(c.Counts)
	c.process()
	c.TruncateTopWords(topwords)

	return c
}

// Takes a scanner, runs through it and counts unqiue words and their
// number of occurrences. When caseSensitive is false all strings will
// be downcased before being counted
func WordCount(scanner *bufio.Scanner, caseSensitive bool) (map[string]int, int) {
	scanner.Split(bufio.ScanWords)
	m := make(map[string]int, 4096)
	total := 0
	for scanner.Scan() {
		var word string
		if caseSensitive {
			word = scrubWord(scanner.Text())
		} else {
			word = scrubWord(strings.ToLower(scanner.Text()))
		}

		m[word]++
		total++
	}

	// Unparseable words are stored as an empty string so we remove that
	delete(m, "")
	return m, total
}

// Takes a word token and strips non alphabetic characters from the
// beginning and end of the word
func scrubWord(s string) string {
	minAlpha := 0
	maxAlpha := 0
	anyAlpha := false
	i := 0

	// Find the first alphabetic character
	for i < len(s) {
		if alphaChar(s[i]) {
			anyAlpha = true
			minAlpha = i
			break
		}
		i++
	}

	// Find the last alphabetic character
	for i < len(s) {
		if alphaChar(s[i]) {
			anyAlpha = true
			maxAlpha = i
		}
		i++
	}

	if anyAlpha {
		return s[minAlpha : maxAlpha+1]
	} else {
		return ""
	}
}

// Returns true if the character is an alphabetic character
func alphaChar(r uint8) bool {
	return inRange(r, 65, 90) || inRange(r, 97, 122)
}

// Returns true if n is within the range lo..hi inclusive
func inRange(n, lo, hi uint8) bool {
	return n >= lo && n <= hi
}

func (c *Concordance) process() {
	c.LengthHistogram = make([]int, 64)
	c.MostUsed = make([]WordTuple, 0, c.Unique)

	// Init both MostUsed and LengthHistogram in one pass of the Counts map
	for k, v := range c.Counts {
		// Most Used
		c.MostUsed = append(c.MostUsed, WordTuple{Word: k, Count: v})

		// Length Histogram
		wordlen := len(k)
		// Resize the length histogram if it's too short
		if wordlen >= len(c.LengthHistogram) {
			newlen := 2 * len(c.LengthHistogram)
			for newlen < wordlen {
				newlen *= 2
			}
			newHist := make([]int, 2*len(c.LengthHistogram))
			// Copy over the old values
			for i, v := range c.LengthHistogram {
				newHist[i] = v
			}
			c.LengthHistogram = newHist
			c.LengthHistogram[wordlen]++
		} else {
			c.LengthHistogram[wordlen]++
		}
	}
	sort.Sort(sort.Reverse(ByCount(c.MostUsed)))
	c.trimHist()
}

func (c *Concordance) trimHist() {
	max := 0
	for i, v := range c.LengthHistogram {
		if v > 0 {
			max = i
		}
	}
	c.LengthHistogram = c.LengthHistogram[:max+1]
}

func (c *Concordance) TruncateTopWords(n int) {
	c.MostUsed = c.MostUsed[:n]
}
