/// <reference types="tree-sitter-cli/dsl" />
// @ts-check

const PREC = {
  ASSIGN: 1,
  OR: 2,
  AND: 3,
  EQUALITY: 4,
  COMPARE: 5,
  ADD: 6,
  MULTIPLY: 7,
  UNARY: 8,
  CALL: 9,
};

module.exports = grammar({
  name: "ash",

  extras: $ => [/[ \t\r\n]+/, $.comment],
  word: $ => $.identifier,

  rules: {
    source_file: $ => repeat($._statement),

    _statement: $ => choice(
      $.import_statement,
      $.variable_declaration,
      $.function_definition,
      $.if_statement,
      $.while_statement,
      $.return_statement,
      $.expression_statement,
      $.block
    ),

    comment: _ => token(choice(
      seq("//", /.*/),
      seq("/*", /[^*]*\*+([^/*][^*]*\*+)*/, "/")
    )),

    import_statement: $ => seq("import", choice($.string, $.angle_import), ";"),
    angle_import: _ => token(seq("<", /[^>]+/, ">")),

    block: $ => seq("{", repeat($._statement), "}"),

    variable_declaration: $ => seq(
      field("type", $._type),
      field("name", $.identifier),
      optional(seq("=", field("value", $._expression))),
      ";"
    ),

    function_definition: $ => seq(
      field("return_type", $._type),
      field("name", $.identifier),
      "(", optional(commaSep1($.parameter)), ")",
      field("body", $.block)
    ),

    parameter: $ => seq(field("type", $._type), field("name", $.identifier)),

    if_statement: $ => prec.right(seq(
      "if", "(", field("condition", $._expression), ")",
      field("consequence", $._statement),
      optional(seq("else", field("alternative", $._statement)))
    )),

    while_statement: $ => seq("while", "(", field("condition", $._expression), ")", field("body", $._statement)),
    return_statement: $ => seq("return", optional($._expression), ";"),
    expression_statement: $ => seq($._expression, ";"),

    _type: $ => choice($._base_type, $.aggregate_type),
    _base_type: $ => choice($.primitive_type, $.identifier),
    primitive_type: _ => choice(
      "void", "boolean", "int", "float", "string", "buffer",
      "item", "effect", "skill", "location", "monster", "class", "stat",
      "familiar", "slot", "element", "coinmaster"
    ),
    aggregate_type: $ => seq($._base_type, "[", optional($._base_type), "]"),

    _expression: $ => choice(
      $.assignment_expression,
      $.binary_expression,
      $.unary_expression,
      $.call_expression,
      $.parenthesized_expression,
      $.typed_constant,
      $.string,
      $.float,
      $.integer,
      $.boolean,
      $.identifier
    ),

    assignment_expression: $ => prec.right(PREC.ASSIGN, seq(
      field("left", $.identifier),
      field("operator", choice("=", "+=", "-=", "*=", "/=")),
      field("right", $._expression)
    )),

    binary_expression: $ => choice(
      prec.left(PREC.OR, seq($._expression, "||", $._expression)),
      prec.left(PREC.AND, seq($._expression, "&&", $._expression)),
      prec.left(PREC.EQUALITY, seq($._expression, choice("==", "!="), $._expression)),
      prec.left(PREC.COMPARE, seq($._expression, choice("<", "<=", ">", ">="), $._expression)),
      prec.left(PREC.ADD, seq($._expression, choice("+", "-"), $._expression)),
      prec.left(PREC.MULTIPLY, seq($._expression, choice("*", "/", "%"), $._expression))
    ),

    unary_expression: $ => prec(PREC.UNARY, seq(choice("!", "-", "+"), $._expression)),

    call_expression: $ => prec(PREC.CALL, seq(
      field("function", $.identifier),
      "(", optional(commaSep1($._expression)), ")"
    )),

    parenthesized_expression: $ => seq("(", $._expression, ")"),
    typed_constant: _ => token(seq("$", /[A-Za-z_][A-Za-z0-9_]*/, "[", /[^\]]*/, "]")),
    string: _ => token(seq('"', repeat(choice(/[^"\\\n]/, /\\./)), '"')),
    float: _ => token(/\d+\.\d+/),
    integer: _ => token(/\d+/),
    boolean: _ => choice("true", "false"),
    identifier: _ => /[A-Za-z_][A-Za-z0-9_]*/,
  }
});

function commaSep1(rule) {
  return seq(rule, repeat(seq(",", rule)));
}
