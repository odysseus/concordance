# Concordance

A Go library for generating word counts and stats about plaintext files.

## The Money Methods

**Generating the Full Concordance**

The input data needs to be wrapped in a `bufio.Scanner`, after that specify whether you want case sensitive matching in `caseSensitive` and the maximum number of most used words to return using `topWords`

```go
Function:
func NewConcordance

Arguments:
(scanner *bufio.Scanner, caseSensitive bool, topWords int)

Returns:
*Concordance
```

**Word Count Only**

The word count method can be called separately and returns a map of the words and their counts, as well as the total number of words counted (non-unique, for unique words take the `len` of the map)

```go
Function:
func WordCount

Arguments:
(scanner *bufio.Scanner, caseSensitive bool)

Returns:
(map[string]int, int)
```

**Word Scrubbing Only**

`scrubWord` can be called by itself to remove leading and trailing punctuation that is not part of the word. Note that the input to this function *must* be a token because it does none of the word splitting itself. To generate tokens from an arbitrary text data source the best approach is to wrap it in an `io.Reader` then wrap that in a `bufio.Scanner` and set the scanners split function to words using `scanner.Split(bufio.ScanWords))` (where `scanner` is the var name of the scanner).
