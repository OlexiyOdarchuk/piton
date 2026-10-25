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
            ' ' => (),
            '\r' => (),
            '\t' => (),
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
                if let Some(next_c) = chars.peek()
                    && next_c == &'='
                {
                    tokens.push(Token::new(TokenType::Equal, line_num));
                    continue;
                } else {
                    tokens.push(Token::new(TokenType::Assign, line_num))
                }
            }
            '>' => {
                if let Some(next_c) = chars.peek()
                    && next_c == &'='
                {
                    tokens.push(Token::new(TokenType::GreaterThanEqual, line_num));
                    continue;
                } else {
                    tokens.push(Token::new(TokenType::GreaterThan, line_num))
                }
            }
            '<' => {
                if let Some(next_c) = chars.peek()
                    && next_c == &'='
                {
                    tokens.push(Token::new(TokenType::LessThanEqual, line_num));
                    continue;
                } else {
                    tokens.push(Token::new(TokenType::LessThan, line_num))
                }
            }
            '!' => {
                if let Some(next_c) = chars.peek()
                    && next_c == &'='
                {
                    tokens.push(Token::new(TokenType::NotEqual, line_num));
                    continue;
                } else {
                    tokens.push(Token::new(TokenType::IllegalType, line_num))
                }
            }
            _ => tokens.push(Token::new(TokenType::IllegalType, line_num)),
        }
    }
    tokens
}

#[cfg(test)]
mod tests {
    use super::*;
    use TokenType::*;

    fn types(src: &str) -> Vec<TokenType> {
        tokenize(src).iter().map(|t| t.token_type.clone()).collect()
    }

    #[test]
    fn indents() {
        assert_eq!(
            types("a\n    b\nc"),
            [NewLine, Indent, NewLine, Dedent, NewLine, Eof]
        );
    }
}
