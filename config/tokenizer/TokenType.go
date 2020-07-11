package tokenizer

// TokenType represents the type of a token
type TokenType int

const (
	// KeyToken represents a string without quotes
	KeyToken TokenType = iota
	// IntToken represents a token that is the format of an integer
	IntToken TokenType = iota
	// StringToken represents a string with quotes
	StringToken TokenType = iota
	// OpenBraceSymbol is a '{'
	OpenBraceSymbol TokenType = iota
	// CloseBraceSymbol is a '}'
	CloseBraceSymbol TokenType = iota
	// OpenBracketSymbol is a '['
	OpenBracketSymbol TokenType = iota
	// CloseBracketSymbol is a ']'
	CloseBracketSymbol TokenType = iota
	// ColonSymbol is a ':'
	ColonSymbol TokenType = iota
	// DotSymbol is a '.'
	DotSymbol TokenType = iota
	// CommaSymbol is a ','
	CommaSymbol TokenType = iota
	// LineFeedToken represents the end of a line
	LineFeedToken TokenType = iota
)

const (
	// WhitespaceToken represents whitespace
	WhitespaceToken TokenType = -1 - iota
	// CommentToken represents the start of a comment
	CommentToken TokenType = -1 - iota
	// InvalidToken represents a token that is not recognized
	InvalidToken TokenType = -1 - iota
)

// TokenLookup contains the type of token based on the first character's ASCII value
var TokenLookup = [256]TokenType{
	// <NUL> <SOH> <STX> <ETX> <EOT> <ENQ> <ACK> <BEL>
	InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken,
	// <BS> <TAB> <LF> <VT> <FF> <CR> <SO> <SI>
	InvalidToken, WhitespaceToken, WhitespaceToken, InvalidToken, InvalidToken, WhitespaceToken, InvalidToken, InvalidToken,
	// <DLE> <DC1> <DC2> <DC3> <DC4> <NAK> <SYN> <ETB>
	InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken,
	// <CAN> <EM> <SUB> <ESC> <FS> <GS> <RS> <US>
	InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken,
	// ' ' ! " # $ % & '
	WhitespaceToken, InvalidToken, StringToken, CommentToken, InvalidToken, InvalidToken, InvalidToken, StringToken,
	// ( ) * + , - . /
	InvalidToken, InvalidToken, InvalidToken, InvalidToken, CommaSymbol, InvalidToken, DotSymbol, InvalidToken,
	// 0 1 2 3 4 5 6 7
	IntToken, IntToken, IntToken, IntToken, IntToken, IntToken, IntToken, IntToken,
	// 8 9 : ; < = > ?
	IntToken, IntToken, ColonSymbol, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken,
	// @ A B C D E F G
	InvalidToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken,
	// H I J K L M N O
	KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken,
	// P Q R S T U V W
	KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken,
	// X Y Z [ \ ] ^ _
	KeyToken, KeyToken, KeyToken, OpenBracketSymbol, InvalidToken, CloseBracketSymbol, InvalidToken, InvalidToken,
	// ` a b c d e f g
	InvalidToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken,
	// h i j k l m n o
	KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken,
	// p q r s t u v w
	KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken,
	// x y z { | } ~ <DEL>
	KeyToken, KeyToken, KeyToken, OpenBraceSymbol, InvalidToken, CloseBraceSymbol, InvalidToken, InvalidToken,
	// Ç ü é â ä à å ç
	KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken,
	// ê ë è ï î ì Ä Å
	KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken,
	// É æ Æ ô ö ò û ù
	KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken,
	// ÿ Ö Ü ¢ £ ¥ ₧ ƒ
	KeyToken, KeyToken, KeyToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken,
	// á í ó ú ñ Ñ ª º
	KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, InvalidToken, InvalidToken,
	// ¿ ⌐ ¬ ½ ¼ ¡ « »
	InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken,
	// ░ ▒ ▓ │ ┤ ╡ ╢ ╖
	InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken,
	// ╕ ╣ ║ ╗ ╝ ╜ ╛ ┐
	InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken,
	// └ └ ┬ ├ ─ ┼ ╞ ╟
	InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken,
	// ╚ ╔ ╩ ╦ ╠ ═ ╬ ╧
	InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken,
	// ╨ ╤ ╥ ╙ ╘ ╒ ╓ ╫
	InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken,
	// ╪ ┘ ┌ █ ▄ ▌ ▐ ▀
	InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken,
	// α ß Γ π Σ σ µ τ
	KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken, KeyToken,
	// Φ Θ Ω δ ∞ φ ε ∩
	KeyToken, KeyToken, KeyToken, KeyToken, InvalidToken, KeyToken, KeyToken, KeyToken,
	// ≡ ± ≥ ≤ ⌠ ⌡ ÷ ≈
	InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken,
	// ° ∙ · √ ⁿ ² ■ <NBSP>
	InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, InvalidToken, WhitespaceToken,
}
