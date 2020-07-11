package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/zachdeibert/protomux/config/ast"
	"github.com/zachdeibert/protomux/config/common"
	"github.com/zachdeibert/protomux/config/lexer"
	"github.com/zachdeibert/protomux/config/tokenizer"
)

func splitGroups(args []string) []ArgumentGroup {
	start := 0
	groups := []ArgumentGroup{}
	locations := Locations(args)
	for i, arg := range args {
		if len(arg) >= 2 && arg[0] == '-' && i != start {
			groups = append(groups, ArgumentGroup{
				Arguments: args[start:i],
				Locations: locations[start:i],
			})
			start = i
		}
	}
	if start < len(args) {
		groups = append(groups, ArgumentGroup{
			Arguments: args[start:],
			Locations: locations[start:],
		})
	}
	return groups
}

// Parse the command-line arguments into an AST
func Parse(args []string, flags []CommandFlag) (*ast.AST, error) {
	ret := &ast.AST{
		Blocks:     []ast.Block{},
		Parameters: []ast.Parameter{},
	}
	keyLocations := map[string]common.Location{}
	res := map[string]ArgumentGroup{}
	for _, group := range splitGroups(args) {
		var flag *CommandFlag = nil
		if group.Arguments[0][1] == '-' {
			for _, f := range flags {
				if f.LongFlag == group.Arguments[0][2:] {
					flag = &f
					break
				}
			}
		} else {
			for _, f := range flags {
				if f.ShortFlag == group.Arguments[0][1:] {
					flag = &f
					break
				}
			}
		}
		if flag == nil {
			return nil, ErrorUnrecognizedFlag(group.Arguments[0], group.Locations[0])
		}
		if len(group.Arguments) < flag.NumArgs+1 {
			return nil, ErrorMissingArgs(group.Arguments[0], common.Merge(group.Locations))
		}
		switch flag.NumArgs {
		case 0:
			if len(group.Arguments) > 1 {
				return nil, ErrorTooManyArgs(group.Arguments[0], common.Merge(group.Locations[1:]))
			}
			res[flag.TreePath] = ArgumentGroup{
				Arguments: []string{"true"},
				Locations: group.Locations,
			}
			keyLocations[flag.TreePath] = group.Locations[0]
			break
		case 1:
			res[flag.TreePath] = ArgumentGroup{
				Arguments: group.Arguments[1:],
				Locations: group.Locations[1:],
			}
			keyLocations[flag.TreePath] = group.Locations[0]
			break
		case 2:
			res[group.Arguments[1]] = ArgumentGroup{
				Arguments: group.Arguments[2:],
				Locations: group.Locations[2:],
			}
			keyLocations[group.Arguments[1]] = group.Locations[1]
			break
		default:
			panic("Missing case")
		}
	}
	for k, v := range res {
		loc := keyLocations[k]
		if k[0] != '/' {
			return nil, ErrorTreePathFormat("Tree path missing leading '/'", common.Location{
				FileName:  loc.FileName,
				Line:      []byte(k),
				LineNo:    loc.LineNo,
				CharStart: 0,
				CharLen:   1,
			})
		}
		if len(k) < 2 {
			return nil, ErrorTreePathFormat("Tree path too short", common.Location{
				FileName:  loc.FileName,
				Line:      []byte(k),
				LineNo:    loc.LineNo,
				CharStart: len(k),
				CharLen:   1,
			})
		}
		parts := strings.Split(k[1:], "/")
		tree := ret
		for i, w := range parts {
			if i == len(parts)-1 {
				str := fmt.Sprintf("%s [ %s ]", w, strings.Join(v.Arguments, " , "))
				tokenStream := tokenizer.CreateTokenReader(strings.NewReader(str), "cmdline")
				lexemeStream := lexer.CreateLexemeReader(tokenStream)
				first, err := lexemeStream.Next()
				if err != nil {
					return nil, err
				}
				second, err := lexemeStream.Next()
				if err != nil {
					return nil, err
				}
				param, _, err := ast.ParseParameter(lexemeStream, *first, *second)
				if err != nil {
					origErr := err
					str = fmt.Sprintf("%s [ \"%s\" ]", w, strings.Join(v.Arguments, "\" , \""))
					tokenStream = tokenizer.CreateTokenReader(strings.NewReader(str), "cmdline")
					lexemeStream = lexer.CreateLexemeReader(tokenStream)
					first, err = lexemeStream.Next()
					if err != nil {
						return nil, err
					}
					second, err = lexemeStream.Next()
					if err != nil {
						return nil, err
					}
					param, _, err = ast.ParseParameter(lexemeStream, *first, *second)
					if err != nil {
						return nil, origErr
					}
				}
				tree.Parameters = append(tree.Parameters, *param)
			} else {
				idx := 0
				if w[len(w)-1] == ']' {
					split := strings.SplitAfter(w, "[")
					if len(split) != 2 {
						return nil, ErrorTreePathFormat("Tree path component missing '['", common.Location{
							FileName:  loc.FileName,
							Line:      []byte(w),
							LineNo:    loc.LineNo,
							CharStart: len(w) - 1,
							CharLen:   1,
						})
					}
					val, err := strconv.ParseInt(split[1][:len(split[1])-1], 10, 32)
					if err != nil {
						return nil, ErrorTreePathFormat(fmt.Sprintf("Tree path contains invalid index '%s': %s", split[1][:len(split[1])-1], err.Error()), common.Location{
							FileName:  loc.FileName,
							Line:      []byte(w),
							LineNo:    loc.LineNo,
							CharStart: len(split[0]) + 1,
							CharLen:   len(split[1]) - 1,
						})
					}
					idx = int(val)
					w = split[0][:len(split[0])-1]
				}
				for i, block := range tree.Blocks {
					if block.Name == w {
						if idx == 0 {
							tree = &tree.Blocks[i].Children
							idx--
							break
						} else {
							idx--
						}
					}
				}
				if idx >= 0 {
					for ; idx >= 0; idx-- {
						tree.Blocks = append(tree.Blocks, ast.Block{
							Name:     w,
							Location: loc,
							Children: ast.AST{
								Blocks:     []ast.Block{},
								Parameters: []ast.Parameter{},
							},
						})
					}
					tree = &tree.Blocks[len(tree.Blocks)-1].Children
				}
			}
		}
	}
	return ret, nil
}
