# Diceware Passphrase Generator

This Diceware Generator aims to provide an easy to use interface for creating diceware passphrases whether you wish to use it as a library for your application or whether you want to quickly generate a passphrase on the command line.

## Security
1. This generator uses a cryptographically secure random number generator to generate passphrases.
2. It is kept up to date with the latest required word count for a secure passphrase.
3. It enforces the latest required character count per passphrase.
4. It restricts the usage of insecure word separators.

## Using as an executable...
In your $GOPATH, run:

```
go install github.com/rosatolen/diceware
```

This will install the diceware generator as the executable `diceware` in your $GOPATH/bin that you can run with:

```
$GOPATH/bin/diceware
```

or if your $GOPATH is a part of your regular path:

```
diceware
```

This will automatically create a diceware passphrase for you with the default settings. If you would like to specify more qualities of your diceware, use `$ $GOPATH/bin/diceware --help` to show your options.

## Using as a Library...
Import the diceware generator by running:

```
go get github.com/rosatolen/diceware
```

and by importing it:

```
import (
    "github.com/rosatolen/diceware/generator"
)
```

Then use either generator.Passphrase or generator.PassphraseWith.

PassphraseWith allows you to specify a wordlist, the number of words, and a different word separator from the default space character (" ") with parameters.

If you specify a wordlist, please make sure that it has the format of a normal diceware wordlist. Meaning:

1. The number of digits in the ID of each diceware word is 5
2. The values of each diceware ID digit are between 1-6, similar to the limit of values when rolling a die
3. Each line containing a word ID and a word should be in the following format: "[5 digit ID][whitespace][word]"

The PassphraseWith method will return error if supplied a word count less than the current secure minimum (6) or supplied a non-space separator that exists as a word in the wordlist or is alphanumeric. Separator constraints are based on the collision issues described here: [http://world.std.com/~reinhold/dicewarefaq.html#spaces](http://world.std.com/~reinhold/dicewarefaq.html#spaces)
