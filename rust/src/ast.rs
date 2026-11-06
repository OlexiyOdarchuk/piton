#[derive(Debug, Clone, PartialEq)]
pub(crate) struct ElseIf {
    pub(crate) condition: Expr,
    pub(crate) body: Vec<Stmt>,
}

#[derive(Debug, Clone, Copy, PartialEq)]
pub(crate) enum BinOperator {
    Plus,
    Minus,
    Multiply,
    Divide,
    GreaterThan,
    LessThan,

    Equal,
    NotEqual,
    GreaterThanEqual,
    LessThanEqual,
    And,
    Or,

    Stupin,
}

#[derive(Debug, Clone, Copy, PartialEq)]
pub(crate) enum UnOperator {
    Minus,
    Kosynus,
    Synus,
    ArkKosynus,
    ArkSynus,
    Abs,
    Korin,
    Loh10,
    Not,
}

#[derive(Debug, Clone, PartialEq)]
pub(crate) enum Expr {
    Prefix {
        operator: UnOperator,
        right: Box<Expr>,
    },
    Infix {
        operator: BinOperator,
        left: Box<Expr>,
        right: Box<Expr>,
    },
    Call {
        name: String,
        receiver: Option<Box<Expr>>,
        args: Vec<Expr>,
    },
    Index {
        arr: Box<Expr>,
        index: Box<Expr>,
    },
    Slice {
        left: Box<Expr>,
        start: Option<Box<Expr>>,
        end: Option<Box<Expr>>,
    },
    Selector {
        right: String,
        left: Box<Expr>,
    },
    Number(f64),
    Str(String),
    Ident(String),
    Spysok(Vec<Expr>),
    Slovnyk(Vec<(Expr, Expr)>),
}

#[derive(Debug, Clone, PartialEq)]
pub(crate) enum Stmt {
    Import {
        file_name: Expr,
    },
    FuncDef {
        name: String,
        params: Vec<String>,
        body: Vec<Stmt>,
        module: Option<String>,
    },
    Drukuvaty {
        expr: Expr,
    },
    VarDec {
        name: String,
        expr: Expr,
    },
    Vvid {
        name: String,
    },
    Assign {
        target: Expr,
        expr: Expr,
    },
    If {
        condition: Expr,
        body: Vec<Stmt>,
        else_if: Vec<ElseIf>,
        else_body: Vec<Stmt>,
    },
    Vernuty {
        expr: Expr,
    },
    Expr {
        expr: Expr,
    },
    Poky {
        condition: Expr,
        body: Vec<Stmt>,
    },
}
