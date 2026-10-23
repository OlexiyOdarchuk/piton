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
            tokens.push(Token::new(TokenType::Indent, "", line_num));
        } else if indent < indents[indents.len() - 1] {
            while indents.len() > 1 && indent < indents[indents.len() - 1] {
                indents.pop();
                tokens.push(Token::new(TokenType::Dedent, "", line_num))
            }
        }

        let line_tokens = tokenize_line(line, line_num);
        tokens.extend(line_tokens);
        tokens.push(Token::new(TokenType::NewLine, "\n", line_num))
    }
    while indents.len() > 1 {
        indents.pop();
        tokens.push(Token::new(TokenType::Dedent, "", line_num));
    }

    tokens.push(Token::new(TokenType::Eof, "", line_num));
    tokens
}

fn tokenize_line(line: &str, line_num: u32) -> Vec<Token> {
    let mut tokens = Vec::new();
    let mut chars = line.chars().peekable();
    while let Some(c) = chars.next() {
        // TODO: operators, two-char operators (== != <= >=), numbers, idents, strings
        match c {
            _ => {}
        }
    }
    tokens
}

#[cfg(test)]
mod tests {
    use super::*;
    use TokenType::*;

    fn types(src: &str) -> Vec<TokenType> {
        tokenize(src).iter().map(|t| t.token_type).collect()
    }

    #[test]
    fn indents() {
        assert_eq!(
            types("a\n    b\nc"),
            [NewLine, Indent, NewLine, Dedent, NewLine, Eof]
        );
    }
}
