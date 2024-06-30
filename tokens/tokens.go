package tokens

import (
	"fmt"

	"github.com/brandonksides/grundfunken/models"
)

type TokenType int

const (
	// Grouping and Punctuation
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_SQUARE_BRACKET
	RIGHT_SQUARE_BRACKET
	COMMA
	NEWLINE
	LEFT_ANGLE_BRACKET
	RIGHT_ANGLE_BRACKET
	LEFT_SQUIGGLY_BRACKET
	RIGHT_SQUIGGLY_BRACKET
	COLON

	// Operators
	MINUS
	PLUS
	SLASH
	STAR
	EQUAL
	PERCENT
	DOT

	// Values
	IDENTIFIER
	STRING
	NUMBER

	// Keywords
	LET
	IN
	IF
	FOR
	THEN
	ELSE
	FUNC
	AND
	OR
	NOT
	IS
)

var tokMap = map[string]TokenType{
	"(":    LEFT_PAREN,
	")":    RIGHT_PAREN,
	",":    COMMA,
	"\n":   NEWLINE,
	".":    DOT,
	"-":    MINUS,
	"+":    PLUS,
	"/":    SLASH,
	"*":    STAR,
	"=":    EQUAL,
	"<":    LEFT_ANGLE_BRACKET,
	">":    RIGHT_ANGLE_BRACKET,
	"[":    LEFT_SQUARE_BRACKET,
	"]":    RIGHT_SQUARE_BRACKET,
	"{":    LEFT_SQUIGGLY_BRACKET,
	"}":    RIGHT_SQUIGGLY_BRACKET,
	"%":    PERCENT,
	":":    COLON,
	"let":  LET,
	"in":   IN,
	"if":   IF,
	"for":  FOR,
	"then": THEN,
	"else": ELSE,
	"func": FUNC,
	"and":  AND,
	"or":   OR,
	"not":  NOT,
	"is":   IS,
}

type Token struct {
	Type           TokenType
	Value          string
	SourceLocation models.SourceLocation
}

type TokenStack struct {
	toks   []Token
	curLoc models.SourceLocation
}

func (stack *TokenStack) CurrentSourceLocation() models.SourceLocation {
	return stack.curLoc
}

// Pop removes and returns the next token in the stack
func (stack *TokenStack) Pop() (Token, error) {
	if len(stack.toks) == 0 {
		return Token{}, fmt.Errorf("expected token")
	}

	this, rest := stack.toks[0], stack.toks[1:]

	if len(rest) == 0 {
		stack.curLoc = models.SourceLocation{
			LineNumber:   this.SourceLocation.LineNumber + len(this.Value),
			ColumnNumber: this.SourceLocation.ColumnNumber,
			File:         this.SourceLocation.File,
		}
	} else {
		stack.curLoc = rest[0].SourceLocation
	}

	return Token{}, nil
}

// Peek returns the next token in the stack without removing it
func (stack *TokenStack) Peek() (Token, bool) {
	if len(stack.toks) == 0 {
		return Token{}, false
	}

	return stack.toks[0], true
}

func Tokenize(filename string, lines []string) (*TokenStack, *models.InterpreterError) {
	toks := make([]Token, 0)

	for lineNumber, line := range lines {
		lineToks, err := tokenizeLine(filename, line, lineNumber)
		if err != nil {
			return nil, err
		}

		toks = append(toks, lineToks...)
	}

	return &TokenStack{
		toks: toks,
		curLoc: models.SourceLocation{
			File: filename,
		},
	}, nil
}

func tokenizeLine(file string, line string, lineNumber int) ([]Token, *models.InterpreterError) {
	toks := make([]Token, 0)
	col := 0
	for col < len(line) {
		char := line[col]
		if char == ' ' || char == '\t' {
			col++
			continue
		} else if char == '/' {
			if col+1 < len(line) && line[col+1] == '/' {
				break
			}
		}
		if tokType, ok := tokMap[string(char)]; ok {
			toks = append(toks, Token{
				Type:  tokType,
				Value: string(char),
				SourceLocation: models.SourceLocation{
					File:         file,
					LineNumber:   lineNumber,
					ColumnNumber: col,
				},
			})
			col++
		} else if char == '"' {
			strTok, length, err := tokenizeString(line[col:])
			if err != nil {
				return nil, &models.InterpreterError{
					Underlying: err,
					SourceLocation: models.SourceLocation{
						File:         file,
						LineNumber:   lineNumber,
						ColumnNumber: col + length,
					},
				}
			}
			strTok.SourceLocation = models.SourceLocation{
				File:         file,
				LineNumber:   lineNumber,
				ColumnNumber: col,
			}
			col += length
			toks = append(toks, strTok)
		} else if char == '-' || (char >= '0' && char <= '9') {
			numTok, length, err := tokenizeNumber(line[col:])
			if err != nil {
				return nil, &models.InterpreterError{
					Underlying: err,
					SourceLocation: models.SourceLocation{
						File:         file,
						LineNumber:   lineNumber,
						ColumnNumber: col + length,
					},
				}
			}
			numTok.SourceLocation = models.SourceLocation{
				File:         file,
				LineNumber:   lineNumber,
				ColumnNumber: col,
			}
			col += length
			toks = append(toks, numTok)
		} else if char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char == '_' {
			idTok, length, err := tokenizeOther(line[col:])
			if err != nil {
				return nil, &models.InterpreterError{
					Underlying: err,
					SourceLocation: models.SourceLocation{
						File:         file,
						LineNumber:   lineNumber,
						ColumnNumber: col + length,
					},
				}
			}
			idTok.SourceLocation = models.SourceLocation{
				File:         file,
				LineNumber:   lineNumber,
				ColumnNumber: col,
			}
			col += length
			toks = append(toks, idTok)
		} else {
			return nil, &models.InterpreterError{
				Underlying: fmt.Errorf("unexpected character %c", char),
				SourceLocation: models.SourceLocation{
					File:         file,
					LineNumber:   lineNumber,
					ColumnNumber: col,
				},
			}
		}
	}

	return toks, nil
}

func tokenizeString(line string) (Token, int, error) {
	col := 1

	lineTokVal := ""
	for col < len(line) {
		if line[col] == '"' {
			return Token{
				Type:  STRING,
				Value: lineTokVal,
			}, col + 1, nil
		} else if line[col] == '\\' {
			col++
			if col >= len(line) {
				break
			}
			escaped := line[col]
			switch escaped {
			case 'n':
				lineTokVal += "\n"
			case 't':
				lineTokVal += "\t"
			case 'r':
				lineTokVal += "\r"
			case '\\':
				lineTokVal += "\\"
			case '"':
				lineTokVal += "\""
			default:
				return Token{}, col, fmt.Errorf("unexpected escape character %c", escaped)
			}
		} else {
			lineTokVal += string(line[col])
		}
		col++
	}

	return Token{}, col - 1, fmt.Errorf("unterminated string")
}

func tokenizeNumber(line string) (Token, int, error) {
	col := 0
	for col < len(line) {
		if (line[col] < '0' || line[col] > '9') && line[col] != '.' {
			if line[col] >= 'a' && line[col] <= 'z' || line[col] >= 'A' && line[col] <= 'Z' || line[col] == '_' {
				return Token{}, col, fmt.Errorf("unexpected character %c", line[col])
			}
			return Token{
				Type:  NUMBER,
				Value: line[:col],
			}, col, nil
		}
		col++
	}

	return Token{
		Type:  NUMBER,
		Value: line[:col],
	}, col, nil
}

func tokenizeOther(line string) (Token, int, error) {
	col := 0
	for col < len(line) {
		if (line[col] < 'a' || line[col] > 'z') && (line[col] < 'A' || line[col] > 'Z') && line[col] != '_' && (line[col] < '0' || line[col] > '9') {
			break
		}
		col++
	}

	word := line[:col]
	if tokType, ok := tokMap[word]; ok {
		return Token{
			Type:  tokType,
			Value: word,
		}, col, nil
	}

	return Token{
		Type:  IDENTIFIER,
		Value: word,
	}, col, nil
}
