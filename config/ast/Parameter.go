package ast

import (
	"fmt"
	"strings"

	"github.com/zachdeibert/protomux/config/common"
	"github.com/zachdeibert/protomux/config/lexer"
)

// Parameter represents a parameter set on a block in the AST
type Parameter struct {
	Name     string
	Location common.Location
	Type     ParameterType
	Values   []interface{}
}

func parseStringParameter(lexeme lexer.Lexeme) (interface{}, error) {
	if lexeme.Type != lexer.StringLexeme {
		return nil, ErrorParameterArrayType(lexeme, "string")
	}
	return &StringParameterData{
		Value: lexeme.StringValue,
	}, nil
}

func parseConnectionParameter(lexeme lexer.Lexeme) (interface{}, error) {
	if lexeme.Type != lexer.ConnectionLexeme {
		return nil, ErrorParameterArrayType(lexeme, "connection")
	}
	return &ConnectionParameterData{
		Host: lexeme.StringValue,
		Port: lexeme.IntValue,
	}, nil
}

func parseBooleanParameter(lexeme lexer.Lexeme) (interface{}, error) {
	if lexeme.Type != lexer.KeywordLexeme {
		return nil, ErrorParameterArrayType(lexeme, "boolean")
	}
	var val bool
	if lexeme.StringValue == "true" {
		val = true
	} else if lexeme.StringValue == "false" {
		val = false
	} else {
		return nil, ErrorBooleanParse(lexeme)
	}
	return &BooleanParameterData{
		Value: val,
	}, nil
}

// ParseParameter parses a Parameter
func ParseParameter(stream *lexer.LexemeReader, first lexer.Lexeme, second lexer.Lexeme) (*Parameter, *lexer.Lexeme, error) {
	param := &Parameter{
		Name:     first.StringValue,
		Location: common.Merge([]common.Location{first.Location, second.Location}),
	}
	switch second.Type {
	case lexer.ArrayStartLexeme:
		param.Values = []interface{}{}
		locations := []common.Location{param.Location}
		wasValue := false
		var parser func(lexer.Lexeme) (interface{}, error) = nil
		for lexeme, err := stream.Next(); ; lexeme, err = stream.Next() {
			if err != nil {
				return nil, nil, err
			}
			if lexeme == nil {
				return nil, nil, ErrorUnexpectedArrayEOF(param)
			}
			switch lexeme.Type {
			case lexer.ArrayEndLexeme:
				if !wasValue {
					return nil, nil, ErrorUnexpectedArrayLexeme(*lexeme)
				}
				param.Location = common.Merge(locations)
				if lexeme, err = stream.Next(); err != nil {
					return nil, nil, err
				} else if lexeme == nil {
					return nil, nil, ErrorUnexpectedArrayEOF(param)
				} else if lexeme.Type != lexer.LineFeedLexeme {
					return nil, nil, ErrorUnexpectedArrayLexeme(*lexeme)
				}
				return param, nil, nil
			case lexer.ArraySeparatorLexeme:
				if !wasValue {
					return nil, nil, ErrorUnexpectedArrayLexeme(*lexeme)
				}
				wasValue = false
				break
			case lexer.LineFeedLexeme:
				break
			default:
				if wasValue {
					return nil, nil, ErrorUnexpectedArrayLexeme(*lexeme)
				}
				wasValue = true
				if parser == nil {
					switch lexeme.Type {
					case lexer.StringLexeme:
						param.Type = StringParameter
						parser = parseStringParameter
						break
					case lexer.ConnectionLexeme:
						param.Type = ConnectionParameter
						parser = parseConnectionParameter
						break
					case lexer.KeywordLexeme:
						param.Type = BooleanParameter
						parser = parseBooleanParameter
						break
					case lexer.BlockStartLexeme, lexer.BlockEndLexeme, lexer.ArrayStartLexeme:
						break
					default:
						panic("Missing case")
					}
				}
				val, err := parser(*lexeme)
				if err != nil {
					return nil, nil, err
				}
				locations = append(locations, lexeme.Location)
				param.Values = append(param.Values, val)
				break
			}
		}
	case lexer.ConnectionLexeme:
		val, err := parseConnectionParameter(second)
		if err != nil {
			return nil, nil, err
		}
		param.Type = ConnectionParameter
		param.Values = []interface{}{val}
		break
	case lexer.StringLexeme:
		val, err := parseStringParameter(second)
		if err != nil {
			return nil, nil, err
		}
		param.Type = StringParameter
		param.Values = []interface{}{val}
		break
	case lexer.LineFeedLexeme:
		return nil, nil, ErrorSingleLexemeLine(first)
	case lexer.KeywordLexeme, lexer.BlockStartLexeme, lexer.BlockEndLexeme, lexer.ArrayEndLexeme, lexer.ArraySeparatorLexeme:
		return nil, nil, ErrorUnexpectedParameterLexeme(second)
	default:
		panic("Missing case")
	}
	return param, nil, nil
}

// Merge another Parameter onto this Parameter
func (p *Parameter) Merge(overwrite Parameter) {
	*p = overwrite
}

func (p Parameter) String() string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("Parameter '%s'", p.Name))
	for i, v := range p.Values {
		var tree rune
		if i < len(p.Values)-1 {
			tree = '\u251C'
		} else {
			tree = '\u2514'
		}
		buf.WriteString(fmt.Sprintf("\n %c\u2500'%s'", tree, v))
	}
	return buf.String()
}
