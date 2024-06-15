package tokens

import (
	"fmt"

	"github.com/brandonksides/crazy-types/models"
)

type TokenType int

const (
	// Operators, Grouping, and Punctuation
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_SQUARE_BRACKET
	RIGHT_SQUARE_BRACKET
	COMMA
	NEWLINE
	DOT
	MINUS
	PLUS
	SLASH
	STAR
	TILDE
	AMPERSAND
	VERTICAL_BAR
	EQUAL
	LEFT_ANGLE_BRACKET
	RIGHT_ANGLE_BRACKET

	// Literals
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
)

func (t TokenType) IsInfixOperator() bool {
	return t == PLUS ||
		t == SLASH ||
		t == MINUS ||
		t == STAR ||
		t == VERTICAL_BAR ||
		t == EQUAL
}

func (t TokenType) IsPrefixOperator() bool {
	return t == MINUS ||
		t == TILDE
}

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
	"~":    TILDE,
	"&":    AMPERSAND,
	"|":    VERTICAL_BAR,
	"=":    EQUAL,
	"<":    LEFT_ANGLE_BRACKET,
	">":    RIGHT_ANGLE_BRACKET,
	"[":    LEFT_SQUARE_BRACKET,
	"]":    RIGHT_SQUARE_BRACKET,
	"let":  LET,
	"in":   IN,
	"if":   IF,
	"for":  FOR,
	"then": THEN,
	"else": ELSE,
	"func": FUNC,
}

type Token struct {
	Type           TokenType
	Value          string
	SourceLocation models.SourceLocation
}

func Tokenize(lines []string) ([]Token, *models.InterpreterError) {
	toks := make([]Token, 0)

	for lineNumber, line := range lines {
		lineToks, err := tokenizeLine(line, lineNumber)
		if err != nil {
			return nil, err
		}

		toks = append(toks, lineToks...)
	}

	return toks, nil
}

func (t Token) IsInfixOperator() bool {
	return t.Type.IsInfixOperator()
}

func (t Token) IsPrefixOperator() bool {
	return t.Type.IsPrefixOperator()
}

func tokenizeLine(line string, lineNumber int) ([]Token, *models.InterpreterError) {
	toks := make([]Token, 0)
	col := 0
	for col < len(line) {
		char := line[col]
		if char == ' ' || char == '\t' {
			col++
			continue
		} else if tokType, ok := tokMap[string(char)]; ok {
			toks = append(toks, Token{
				Type:  tokType,
				Value: string(char),
				SourceLocation: models.SourceLocation{
					LineNumber:   lineNumber,
					ColumnNumber: col,
				},
			})
			col++
		} else if char == '"' {
			strTok, length, err := tokenizeString(line[col:])
			if err != nil {
				return nil, &models.InterpreterError{
					Err: err,
					SourceLocation: models.SourceLocation{
						LineNumber:   lineNumber,
						ColumnNumber: col + length,
					},
				}
			}
			strTok.SourceLocation = models.SourceLocation{
				LineNumber:   lineNumber,
				ColumnNumber: col,
			}
			col += length
			toks = append(toks, strTok)
		} else if char == '-' || (char >= '0' && char <= '9') {
			numTok, length, err := tokenizeNumber(line[col:])
			if err != nil {
				return nil, &models.InterpreterError{
					Err: err,
					SourceLocation: models.SourceLocation{
						LineNumber:   lineNumber,
						ColumnNumber: col + length,
					},
				}
			}
			numTok.SourceLocation = models.SourceLocation{
				LineNumber:   lineNumber,
				ColumnNumber: col,
			}
			col += length
			toks = append(toks, numTok)
		} else if char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char == '_' {
			idTok, length, err := tokenizeOther(line[col:])
			if err != nil {
				return nil, &models.InterpreterError{
					Err: err,
					SourceLocation: models.SourceLocation{
						LineNumber:   lineNumber,
						ColumnNumber: col + length,
					},
				}
			}
			idTok.SourceLocation = models.SourceLocation{
				LineNumber:   lineNumber,
				ColumnNumber: col,
			}
			col += length
			toks = append(toks, idTok)
		} else {
			return nil, &models.InterpreterError{
				Err: fmt.Errorf("unexpected character %c", char),
				SourceLocation: models.SourceLocation{
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
	for col < len(line) {
		if line[col] == '"' {
			return Token{
				Type:  STRING,
				Value: line[1:col],
			}, col + 1, nil
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
