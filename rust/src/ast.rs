#[derive(Debug, Clone, Copy)]
enum BinOperator {
    Assign,
    Plus,
    Minus,
    Multiply,
    Divide,
    GreaterThan,
    LessThan,
    Dot,

    Equal,
    NotEqual,
    GreaterThanEqual,
    LessThanEqual,
    And,
    Or,

    Stupin,
}

#[derive(Debug, Clone, Copy)]
enum UnOperator {
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

enum Expr {
    Prefix {
        op: UnOperator,
        right: Box<Expr>,
    },
    Infix {
        operator: BinOperator,
        left: Box<Expr>,
        right: Box<Expr>,
    },
    Call {
        name: String,
        receiver: Box<Expr>,
        args: Vec<Expr>,
    },
    Index {
        arr: Box<Expr>,
        index: Box<Expr>,
    },
    Spysok {
        left: Box<Expr>,
        start: Box<Expr>,
        end: Box<Expr>,
    },
    Selector {
        right: String,
        letf: Box<Expr>,
    },
}

enum Stmt {
    Import {
        file_name: Expr,
    },
    FuncDef {
        name: String,
        params: Vec<String>,
        body: Vec<Stmt>,
        module: String,
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
        else_if: Vec<Stmt>,
        else_body: Vec<Stmt>,
    },
    ElseIf {
        condition: Expr,
        body: Vec<Stmt>,
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
