package generator

import (
	"strings"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type GoDicephraseSuite struct{}

var _ = Suite(&GoDicephraseSuite{})

const testCount = 8
const wordlist = "../wordlists/original-diceware.wordlist.asc"

func isUnique(element string, set []string) bool {
	count := 0
	for _, val := range set {
		if val == element {
			count++
		}
	}
	return count == 1
}

func (s *GoDicephraseSuite) Test_eachWordIsUnique(c *C) {
	passphrase, _ := Passphrase(wordlist)

	words := strings.Split(passphrase, " ")
	for _, v := range words {
		c.Assert(isUnique(v, words), Equals, true)
	}
}

func (s *GoDicephraseSuite) Test_eachPassphraseIsDifferent(c *C) {
	passphrase1, _ := Passphrase(wordlist)
	passphrase2, _ := Passphrase(wordlist)

	c.Assert(passphrase1, Not(Equals), passphrase2)
}

func (s *GoDicephraseSuite) Test_defaultPassphraseHasRecommendedWordCount(c *C) {
	passphrase, _ := Passphrase(wordlist)

	c.Assert(strings.Split(passphrase, " "), HasLen, WordCount)
}

func (s *GoDicephraseSuite) Test_createPassphraseWith8Words(c *C) {
	passphrase, _ := PassphraseWith(wordlist, " ", testCount)

	c.Assert(strings.Split(passphrase, " "), HasLen, testCount)
}

func (s *GoDicephraseSuite) Test_createPassphraseWithUnderscoreSeparation(c *C) {
	passphrase, _ := PassphraseWith(wordlist, "_", testCount)

	c.Assert(strings.Split(passphrase, "_"), HasLen, testCount)
}

func (s *GoDicephraseSuite) Test_cannotCreateDicephraseWithInsecureWordCount(c *C) {
	lowWordCount := 4
	passphrase, err := PassphraseWith(wordlist, " ", lowWordCount)
	c.Assert(passphrase, Equals, "")
	c.Assert(err, ErrorMatches, "diceware: cannot create passphrase: word count requested is 4; must be over 6 words")
}

func (s *GoDicephraseSuite) Test_cannotCreateDicephraseWithoutSeparator(c *C) {
	passphrase, err := PassphraseWith(wordlist, "", testCount)
	c.Assert(passphrase, Equals, "")
	c.Assert(err, ErrorMatches, "diceware: cannot create passphrase: separator cannot be an empty string")
}

func (s *GoDicephraseSuite) Test_separatorCannotExistAsWordInWordlist(c *C) {
	diceword := "?"
	passphrase, err := PassphraseWith(wordlist, diceword, testCount)
	c.Assert(passphrase, Equals, "")
	c.Assert(err, ErrorMatches, "diceware: cannot create passphrase: separator \\? cannot be a word in the wordlist.*")
}

func (s *GoDicephraseSuite) Test_cannotCreateDicephraseWithAlphnumericSeparator(c *C) {
	passphrase, err := PassphraseWith(wordlist, "ab12", testCount)
	c.Assert(passphrase, Equals, "")
	c.Assert(err, ErrorMatches, "diceware: cannot create passphrase: separator ab12 cannot contain alphanumeric characters")
}

func (s *GoDicephraseSuite) Test_generatesErrorForBadWordlistPath(c *C) {
	passphrase, err := PassphraseWith("/bad/file/location", "ab12", testCount)
	c.Assert(passphrase, Equals, "")
	c.Assert(err, ErrorMatches, "diceware: cannot read wordlist file /bad/file/location:.*")
}

func (s *GoDicephraseSuite) Test_generatesErrorWhenItCannotFindWordForAnID(c *C) {
	passphrase, err := PassphraseWith("../wordlists/empty-wordlist", Separator, testCount)
	c.Assert(passphrase, Equals, "")
	c.Assert(err, ErrorMatches, "diceware: cannot choose random word with wordlist ../wordlists/empty-wordlist: cannot find word with the ID:.*")
}

func (s *GoDicephraseSuite) Test_newPasswordGeneratedWhenCharReqNotMet(c *C) {
	badPass := "aa bb cc dd ee ff"
	passphrase, err := checkCharReq(badPass, wordlist, Separator, testCount)
	c.Assert(err, Equals, nil)
	c.Assert(passphrase, Not(Equals), badPass)
}

func (s *GoDicephraseSuite) Test_currentPasswordIsKeptWhenCharReqIsMet(c *C) {
	goodPass := "dow kraft freon plat fear dr"
	passphrase, err := checkCharReq(goodPass, wordlist, Separator, testCount)
	c.Assert(err, Equals, nil)
	c.Assert(passphrase, Equals, goodPass)
}

func (s *GoDicephraseSuite) Test_sixteenCharPassphraseDoesNotMeetMinimumCharReq(c *C) {
	sixteenChars := "1234567890123456"
	val := meetsCharReq(sixteenChars)
	c.Assert(val, Equals, false)
}

func (s *GoDicephraseSuite) Test_seventeenCharPassphraseDoesNotMeetMinimumCharReq(c *C) {
	seventeenChars := "12345678901234567"
	val := meetsCharReq(seventeenChars)
	c.Assert(val, Equals, false)
}

func (s *GoDicephraseSuite) Test_eighteenCharPassphraseMeetsMinimumCharReq(c *C) {
	eighteenChars := "123456789012345678"
	val := meetsCharReq(eighteenChars)
	c.Assert(val, Equals, true)
}
