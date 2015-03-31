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

// The primary method for generating a new Concordance struct
// scanner :: The source of the input
// caseSensitive :: a true value treats differently cased words as different words
// topWords :: Specifies the maximum length of the MostUsed array. A value <= 0
// will return them all
func NewConcordance(scanner *bufio.Scanner, caseSensitive bool, topWords int) *Concordance {
	c := &Concordance{}
	c.Counts, c.Total = WordCount(scanner, caseSensitive)
	c.Unique = len(c.Counts)
	c.process()
	c.TruncateTopWords(topWords)

	return c
}

// Takes a scanner, runs through it and counts unqiue words and their
// number of occurrences.
// caseSensitive :: a true value treats differently cased words as different words
// a false value results in all words being downcased before counting
func WordCount(scanner *bufio.Scanner, caseSensitive bool) (map[string]int, int) {
	// Set the scanner to break on words and not lines
	scanner.Split(bufio.ScanWords)
	m := make(map[string]int, 4096)
	total := 0
	for scanner.Scan() {
		var word string
		if caseSensitive {
			word = ScrubWord(scanner.Text())
		} else {
			word = ScrubWord(strings.ToLower(scanner.Text()))
		}

		m[word]++
		total++
	}

	// Unparseable words are stored as an empty string so we remove that
	delete(m, "")
	return m, total
}

// Takes a word token and strips non alphabetic characters from the beginning
// and end of the word. Any nonalphabetic characters in the middle of the word
// are ignored
func ScrubWord(s string) string {
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

// A wrapper function that runs all the process functions needed to generate
// the concordance and other word count stats
func (c *Concordance) process() {
	c.LengthHistogram = make([]int, 64)
	c.MostUsed = make([]WordTuple, 0, c.Unique)

	// Init both MostUsed and LengthHistogram in one pass of the Counts map
	for k, v := range c.Counts {
		// Most Used
		c.MostUsed = append(c.MostUsed, WordTuple{Word: k, Count: v})

		// Length Histogram - increment the counter where i is the word length
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

// Removes trailing 0 values from the histogram slice
func (c *Concordance) trimHist() {
	max := 0
	for i, v := range c.LengthHistogram {
		if v > 0 {
			max = i
		}
	}
	c.LengthHistogram = c.LengthHistogram[:max+1]
}

// Truncates the top words array to the value specified by maxWords
func (c *Concordance) TruncateTopWords(n int) {
	if n > 0 && len(c.MostUsed) > n {
		c.MostUsed = c.MostUsed[:n]
	}
}
