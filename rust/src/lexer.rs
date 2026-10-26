use crate::token;
use crate::token::Token;
use crate::token::TokenType;

pub fn tokenize(input: &str) -> Vec<Token> {
    let mut tokens: Vec<Token> = Vec::new();
    let mut indents: Vec<usize> = vec![0];
    let mut line_num: u32 = 0;
    for mut line in input.lines() {
        line_num += 1;
        let mut indent = 0;
        for ch in line.chars() {
            match ch {
                '#' => {
                    line = "";
                    break;
                }
                ' ' => indent += 1,
                '\t' => indent += 4,
                _ => break,
            }
        }
        let trimmed = line.trim();
        if trimmed.is_empty() {
            continue;
        }

        if indent > indents[indents.len() - 1] {
            indents.push(indent);
            tokens.push(Token::new(TokenType::Indent, line_num));
        } else if indent < indents[indents.len() - 1] {
            while indents.len() > 1 && indent < indents[indents.len() - 1] {
                indents.pop();
                tokens.push(Token::new(TokenType::Dedent, line_num))
            }
        }

        let line_tokens = tokenize_line(line, line_num);
        tokens.extend(line_tokens);
        tokens.push(Token::new(TokenType::NewLine, line_num))
    }
    while indents.len() > 1 {
        indents.pop();
        tokens.push(Token::new(TokenType::Dedent, line_num));
    }

    tokens.push(Token::new(TokenType::Eof, line_num));
    tokens
}

fn tokenize_line(line: &str, line_num: u32) -> Vec<Token> {
    let mut tokens = Vec::new();
    let mut chars = line.chars().peekable();
    while let Some(c) = chars.next() {
        match c {
            ' ' | '\r' | '\t' => (),
            '#' => break,
            '(' => tokens.push(Token::new(TokenType::LeftParen, line_num)),
            ')' => tokens.push(Token::new(TokenType::RightParen, line_num)),
            '[' => tokens.push(Token::new(TokenType::LeftBracket, line_num)),
            ']' => tokens.push(Token::new(TokenType::RightBracket, line_num)),
            '{' => tokens.push(Token::new(TokenType::LeftFigureBracket, line_num)),
            '}' => tokens.push(Token::new(TokenType::RightFigureBracket, line_num)),
            ',' => tokens.push(Token::new(TokenType::Comma, line_num)),
            '.' => tokens.push(Token::new(TokenType::Dot, line_num)),
            '+' => tokens.push(Token::new(TokenType::Plus, line_num)),
            '-' => tokens.push(Token::new(TokenType::Minus, line_num)),
            '*' => tokens.push(Token::new(TokenType::Multiply, line_num)),
            '/' => tokens.push(Token::new(TokenType::Divide, line_num)),
            ':' => tokens.push(Token::new(TokenType::Colon, line_num)),

            '=' => {
                if let Some(&'=') = chars.peek() {
                    chars.next();
                    tokens.push(Token::new(TokenType::Equal, line_num));
                } else {
                    tokens.push(Token::new(TokenType::Assign, line_num))
                }
            }
            '>' => {
                if let Some(&'=') = chars.peek() {
                    chars.next();
                    tokens.push(Token::new(TokenType::GreaterThanEqual, line_num));
                } else {
                    tokens.push(Token::new(TokenType::GreaterThan, line_num))
                }
            }
            '<' => {
                if let Some(&'=') = chars.peek() {
                    chars.next();
                    tokens.push(Token::new(TokenType::LessThanEqual, line_num));
                } else {
                    tokens.push(Token::new(TokenType::LessThan, line_num))
                }
            }
            '!' => {
                if let Some(&'=') = chars.peek() {
                    chars.next();
                    tokens.push(Token::new(TokenType::NotEqual, line_num));
                } else {
                    tokens.push(Token::new(TokenType::IllegalType, line_num))
                }
            }
            '"' => {
                let mut string_val = String::new();
                while let Some(&next_char) = chars.peek() {
                    if next_char == '"' {
                        chars.next();
                        break;
                    }
                    string_val.push(chars.next().unwrap());
                }
                tokens.push(Token::new(TokenType::StringLiteral(string_val), line_num));
            }
            ch if ch.is_ascii_digit() => {
                let mut num_val = String::from(ch);
                while let Some(&next_num) = chars.peek() {
                    if next_num.is_ascii_digit() || next_num == '.' {
                        num_val.push(chars.next().unwrap());
                    } else {
                        break;
                    }
                }
                tokens.push(Token::new(
                    TokenType::NumberLiteral(num_val.parse().unwrap()),
                    line_num,
                ))
            }
            ch if ch.is_alphabetic() || ch == '_' => {
                let mut ident = String::from(ch);
                while let Some(&next_ch) = chars.peek() {
                    if next_ch.is_alphabetic() || next_ch == '_' {
                        ident.push(chars.next().unwrap());
                    } else {
                        break;
                    }
                }
                tokens.push(Token {
                    token_type: TokenType::from_ident(&ident),
                    line_num,
                })
            }
            _ => tokens.push(Token::new(TokenType::IllegalType, line_num)),
        }
    }
    tokens
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::token::TokenType::*;

    fn get_types(input: &str) -> Vec<TokenType> {
        tokenize(input).into_iter().map(|t| t.token_type).collect()
    }

    #[test]
    fn test_single_and_double_char_operators() {
        let input = "= == ! != < <= > >= + - * / : , . ( ) [ ] { }";
        let expected = vec![
            Assign,
            Equal,
            IllegalType,
            NotEqual,
            LessThan,
            LessThanEqual,
            GreaterThan,
            GreaterThanEqual,
            Plus,
            Minus,
            Multiply,
            Divide,
            Colon,
            Comma,
            Dot,
            LeftParen,
            RightParen,
            LeftBracket,
            RightBracket,
            LeftFigureBracket,
            RightFigureBracket,
            NewLine,
            Eof,
        ];

        assert_eq!(get_types(input), expected);
    }

    #[test]
    fn test_literals_and_identifiers() {
        let input = "x = 42\nname = \"piton\"";
        let tokens = tokenize(input);

        assert_eq!(tokens[0].token_type, TokenType::from_ident("x"));
        assert_eq!(tokens[1].token_type, Assign);
        assert_eq!(tokens[2].token_type, NumberLiteral(42.0));
        assert_eq!(tokens[3].token_type, NewLine);

        assert_eq!(tokens[4].token_type, TokenType::from_ident("name"));
        assert_eq!(tokens[5].token_type, Assign);
        assert_eq!(tokens[6].token_type, StringLiteral("piton".to_string()));
        assert_eq!(tokens[7].token_type, NewLine);
        assert_eq!(tokens[8].token_type, Eof);
    }

    #[test]
    fn test_indentation_and_dedentation() {
        let input = "def main():\n    x = 1\n        y = 2\nz = 3";
        let types = get_types(input);

        let expected = vec![
            TokenType::from_ident("def"),
            TokenType::from_ident("main"),
            LeftParen,
            RightParen,
            Colon,
            NewLine,
            Indent,
            TokenType::from_ident("x"),
            Assign,
            NumberLiteral(1.0),
            NewLine,
            Indent,
            TokenType::from_ident("y"),
            Assign,
            NumberLiteral(2.0),
            NewLine,
            Dedent,
            Dedent,
            TokenType::from_ident("z"),
            Assign,
            NumberLiteral(3.0),
            NewLine,
            Eof,
        ];

        assert_eq!(types, expected);
    }

    #[test]
    fn test_comments_and_empty_lines() {
        let input = "# Це коментар\nx = 1 # Внутрішній коментар\n\n  # Ще коментар\ny = 2";
        let types = get_types(input);

        let expected = vec![
            TokenType::from_ident("x"),
            Assign,
            NumberLiteral(1.0),
            NewLine,
            TokenType::from_ident("y"),
            Assign,
            NumberLiteral(2.0),
            NewLine,
            Eof,
        ];

        assert_eq!(types, expected);
    }

    #[test]
    fn test_line_number_tracking() {
        let input = "a\n    b\nc";
        let tokens = tokenize(input);

        assert_eq!(tokens[0].line_num, 1); // 'a'
        assert_eq!(tokens[1].line_num, 1); // NewLine (1)
        assert_eq!(tokens[2].line_num, 2); // Indent
        assert_eq!(tokens[3].line_num, 2); // 'b'
        assert_eq!(tokens[4].line_num, 2); // NewLine (2)
        assert_eq!(tokens[5].line_num, 3); // Dedent
        assert_eq!(tokens[6].line_num, 3); // 'c'
    }
}
