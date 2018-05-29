// Package generator implements a cryptographically strong
// passphrase generator in the style of diceware
package generator

import (
	cr "crypto/rand"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

const maxDieRolls = 5

// DefaultWordCount is the required number of dicewords in a secure passphrase
// Date Established: March 5, 2014
// Ref: https://diceware.blogspot.com/2014/03/time-to-add-word.html
const DefaultWordCount = 6

// DefaultSeparator is the recommended word separator for dice words
const DefaultSeparator = " "

// charCount is the required number of characters in a diceware passphrase
// Ref: http://world.std.com/~reinhold/dicewarefaq.html#14characters
const charCount = 17

var meetsCharReq = func(pass string) bool { return strings.Count(pass, "")-1 > charCount }

// Passphrase generates a diceware passphrase with 6 words and be at least
// 17 characters long. If there is an error, it returns an empty string and
// the error. Errors can happen while opening the WordList file, the word IDs
// do not follow dice limitation, causing errors while finding a word based on
// a randomly generated id, or when it is unable to get the required entropy from
// the system's random number generator.
func Passphrase(path string) (string, error) {
	return PassphraseWith(path, DefaultSeparator, DefaultWordCount)
}

// PassphraseWith generates a diceware passphrase according to a supplied
// word count, separator. All generated phrases will be at least 17 characters
// long. If there is an error, it returns an empty string and the error.
// Errors can happen while opening the WordList file, the word IDs
// do not follow dice limitation, causing errors while finding a word based on
// a randomly generated id, or when it is unable to get the required entropy from
// the system's random number generator.
// In addition, PassphraseWith will reject the following parameters:
// * A word count less than the minimum recommended for security. As of
//   3/5/2014, the minimum is 6 words.
// * An empty string separator
// * Separators that appear as a tword in the wordlist
// * Separators that are alphanumeric
// DefaultSeparator constraints are based on the collision issues described here:
// http://world.std.com/~reinhold/dicewarefaq.html#spaces
func PassphraseWith(path, separator string, wordCount int) (string, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("diceware: cannot read wordlist file %v: %s", path, err)
	}
	wordlist := string(bytes)
	if err := checkParams(path, separator, wordlist, wordCount); err != nil {
		return "", fmt.Errorf("diceware: cannot create passphrase: %s", err)
	}
	words := []string{}
	for i := 0; i < wordCount; i++ {
		word, err := randWord(wordlist)
		if err != nil {
			return "", fmt.Errorf("diceware: cannot choose random word with wordlist %s: %s", path, err)
		}
		words = append(words, word)
	}
	return checkCharReq(strings.Join(words, separator), path, separator, wordCount)
}

func checkCharReq(pass, path, separator string, wordCount int) (string, error) {
	if !meetsCharReq(pass) {
		return PassphraseWith(path, separator, wordCount)
	}
	return pass, nil
}

var alum = regexp.MustCompile("[[:alnum:]]")

func checkParams(path, separator, wordlist string, wc int) error {
	switch {
	case separator == "":
		return errors.New("separator cannot be an empty string")
	case isInWordlist(separator, wordlist):
		return fmt.Errorf("separator %s cannot be a word in the wordlist %s", separator, path)
	case alum.MatchString(separator):
		return fmt.Errorf("separator %s cannot contain alphanumeric characters", separator)
	case wc < DefaultWordCount:
		return fmt.Errorf("word count requested is %v; must be over %v words", wc, DefaultWordCount)
	default:
		return nil
	}
}

func isInWordlist(word, wordlist string) bool {
	if word := find(word, wordlist); word == "" {
		return false
	}
	return true
}

func randWord(wordlist string) (string, error) {
	wordID, err := randID()
	if err != nil {
		return "", fmt.Errorf("cannot generate random ID: %s", err)
	}
	if word := find(wordID, wordlist); word != "" {
		return word, nil
	}
	return "", fmt.Errorf("cannot find word with the ID: %v", wordID)
}

var idWordLine = regexp.MustCompile(`\d\s\.*`)

func find(tofind, wordlist string) string {
	for _, l := range strings.Split(wordlist, "\n") {
		if idWordLine.MatchString(l) && strings.Contains(l, tofind) {
			return strings.Trim(strings.Fields(l)[1], " ")
		}
	}
	return ""
}

func randID() (string, error) {
	digits := []string{}
	for i := 0; i < maxDieRolls; i++ {
		die := &die{maxRoll: 6, minRoll: 1}
		digit, err := die.rollDie()
		if err != nil {
			return "", fmt.Errorf("cannot generate digit number %d of %d digit word ID: %s", i, maxDieRolls, err)
		}
		digits = append(digits, digit)
	}
	return strings.Join(digits, ""), nil
}

type die struct {
	maxRoll int
	minRoll int
}

func (d *die) rollDie() (string, error) {
	for {
		r, err := cr.Int(cr.Reader, big.NewInt(int64(d.maxRoll+1)))
		if err != nil {
			return "", fmt.Errorf("cannot source enough entropy: %s", err)
		}
		rollVal := r.Int64()
		if rollVal > 0 {
			base := 10
			return strconv.FormatInt(rollVal, base), nil
		}
	}
}
