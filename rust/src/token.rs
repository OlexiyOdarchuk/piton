use std::convert::identity;

#[derive(Debug, Clone, PartialEq)]
pub enum TokenType {
    IllegalType,
    Eof,
    NewLine,
    Indent,
    Dedent,

    Ident(String),
    NumberLiteral(f64),
    StringLiteral(String),
    LeftBracket,
    RightBracket,
    LeftFigureBracket,
    RightFigureBracket,
    Comma,

    Assign,
    Plus,
    Minus,
    Multiply,
    Divide,
    GreaterThan,
    LessThan,
    LeftParen,
    RightParen,
    Colon,
    Dot,

    Equal,
    NotEqual,
    GreaterThanEqual,
    LessThanEqual,
    Not,
    And,
    Or,

    Functia,
    Drukuvaty,
    Vvid,
    Yaksho,
    Inackshe,
    Korin,
    Loh10,
    Abs,
    Synus,
    ArkSynus,
    Kosynus,
    ArkKosynus,
    Stupin,
    Vernuty,
    Poky,
    Vykorystaty,
}

impl TokenType {
    pub fn from_ident(ident: &str) -> Self {
        match ident {
            "functia" => TokenType::Functia,
            "drukuvaty" => TokenType::Drukuvaty,
            "vvid" => TokenType::Vvid,
            "yaksho" => TokenType::Yaksho,
            "inackshe" => TokenType::Inackshe,
            "korin" => TokenType::Korin,
            "loh10" => TokenType::Loh10,
            "abs" => TokenType::Abs,
            "synus" => TokenType::Synus,
            "kosynus" => TokenType::Kosynus,
            "arksynus" => TokenType::ArkSynus,
            "arkkosynus" => TokenType::ArkKosynus,
            "stupin" => TokenType::Stupin,
            "vernuty" => TokenType::Vernuty,
            "poky" => TokenType::Poky,
            "vykorystaty" => TokenType::Vykorystaty,
            "ta" => TokenType::And,
            "abo" => TokenType::Or,
            "ne" => TokenType::Not,
            ident => TokenType::Ident(String::from(ident)),
        }
    }
}

#[derive(Debug, Clone, PartialEq)]
pub struct Token {
    pub token_type: TokenType,
    pub line_num: u32,
}

impl Token {
    pub fn new(token_type: TokenType, line_num: u32) -> Self {
        Token {
            token_type,
            line_num,
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    #[test]
    fn keywords() {
        assert_eq!(TokenType::from_ident("functia"), TokenType::Functia);
        assert_eq!(TokenType::from_ident("drukuvaty"), TokenType::Drukuvaty);
        assert_eq!(TokenType::from_ident("vvid"), TokenType::Vvid);
        assert_eq!(TokenType::from_ident("yaksho"), TokenType::Yaksho);
        assert_eq!(TokenType::from_ident("inackshe"), TokenType::Inackshe);
        assert_eq!(TokenType::from_ident("korin"), TokenType::Korin);
        assert_eq!(TokenType::from_ident("loh10"), TokenType::Loh10);
        assert_eq!(TokenType::from_ident("abs"), TokenType::Abs);
        assert_eq!(TokenType::from_ident("synus"), TokenType::Synus);
        assert_eq!(TokenType::from_ident("arksynus"), TokenType::ArkSynus);
        assert_eq!(TokenType::from_ident("kosynus"), TokenType::Kosynus);
        assert_eq!(TokenType::from_ident("arkkosynus"), TokenType::ArkKosynus);
        assert_eq!(TokenType::from_ident("stupin"), TokenType::Stupin);
        assert_eq!(TokenType::from_ident("vernuty"), TokenType::Vernuty);
        assert_eq!(TokenType::from_ident("poky"), TokenType::Poky);
        assert_eq!(TokenType::from_ident("vykorystaty"), TokenType::Vykorystaty);
        assert_eq!(TokenType::from_ident("ta"), TokenType::And);
        assert_eq!(TokenType::from_ident("abo"), TokenType::Or);
        assert_eq!(TokenType::from_ident("ne"), TokenType::Not);
        assert_eq!(
            TokenType::from_ident("blablabla"),
            TokenType::Ident(String::from("blablabla"))
        );
    }
}
