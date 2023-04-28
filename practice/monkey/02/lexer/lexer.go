package lexer // 字句解析機

import (
	"github.com/kuroweb/go-interpreter/practice/monkey/02/token"
)

type Lexer struct {
	input        string
	position     int    // 入力における現在の位置(現在の文字を指し示す)
	readPosition int		// これから読み込む位置(現在の文字の次)
	ch           byte   // 現在検査中の文字
}

// 検査中の文字列に対応するトークンを返却
func (l *Lexer) NextToken() token.Token {
	var tok token.Token // Token構造体を初期化
	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			// 検査中の文字列が英文字のとき
			tok.Literal = l.readIdentifier() // 識別子を読み込む
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			// 検査中の文字列が数字のとき
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			// 検査中の文字列がいずれの条件にも当てはまらない異常値の時
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar() // 検査中の文字列とポジションを更新してポインタを進める
	return tok // 検査中した文字列に対応するトークンを返却
}

// 検査中の文字列とポジションを更新してポインタを進める
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// 次に検査する文字列を取得して返却
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// 識別子を読み込むメソッド
func(l *Lexer) readIdentifier() string {
	position := l.position
	// 検査中の文字列が英文字の場合に、検査中の文字列とポジションを更新してポインタを進める
	for isLetter(l.ch) {
		l.readChar()
	}
	// sliceで識別子の文字列を取得して返却
	// position:   初めに読み込んだ英文字のポジション
	// l.position: 続けて読み込んだ英文字のポジション
	return l.input[position:l.position]
}

// 数字を読み込むメソッド
func (l *Lexer) readNumber() string {
	position := l.position
	// 検査中の文字列が数値の場合に、検査中の文字列とポジションを更新してポインタを進める
	for isDigit(l.ch) {
		l.readChar()
	}
	// sliceで数値を取得して返却
	// position:   初めに読み込んだ数字のポジション
	// l.position: 続けて読み込んだ数字のポジション
	return l.input[position:l.position]
}

// 検査中の文字列が空白な時に、検査中の文字列とポジションを更新してポインタを進める
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// inputの値でLexer構造体を初期化して返却
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// Token構造体を返却する関数
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// 検査中の文字列が英文字なのか？
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// 検査中の文字列が数字なのか？
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
