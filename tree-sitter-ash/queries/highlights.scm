(comment) @comment
(string) @string
(integer) @number
(float) @number.float
(boolean) @boolean
(typed_constant) @constant
(function_definition name: (identifier) @function)
(call_expression function: (identifier) @function.call)
(variable_declaration name: (identifier) @variable)
(parameter name: (identifier) @variable.parameter)
[
  "if"
  "else"
  "while"
  "return"
  "import"
] @keyword
