# EsoLambda Language Specification

EsoLambda is a fictional, Turing-complete scripting language inspired by lambda calculus, designed for Elder Scrolls Online (ESO) addons and general computation. It features static typing, expression-oriented syntax, first-class functions, and mutable state for practicality. The syntax uses colons for type-name binding, pipes | for parameter separation, and arrows -> for function bodies.

Since lambda calculus is Turing-complete, EsoLambda achieves Turing completeness through recursion, loops (while/foreach), conditionals, and unbounded data structures (e.g., vectors). It supports arbitrary computation via mutable variables and infinite loops (with care to avoid ESO runtime limits).

## Key Features

- **Static Typing**: Types are declared before variable/parameter names (e.g., int:x).
- **Expression-Oriented**: Everything returns a value; no pure statements.
- **Functions**: First-class, curried by default, support recursion.
- **Data Structures**: Primitives (int, float, bool, string), collections (vec<T>, map<K,V>).
- **Control Flow**: If-else, while, foreach.
- **Mutability**: Variables are immutable by default; use mut keyword for mutation.
- **Built-ins**: ESO-specific (e.g., get-power, cast) and general (print, len).
- **Error Handling**: Runtime errors for type mismatches, out-of-bounds, etc.; no try-catch (panics on error).

## Basic Syntax Rules

- Function definition: returnType:functionName = type:paramName | type:paramName -> body
- Body: Single expression after ->, or indented block (no braces needed, but {} optional for grouping).
- Variable binding: let type:name = value (inside functions).
- Mutation: mut type:name = value for mutable vars; assign with name = newValue.
- Function calls: functionName arg1 | arg2 (pipe for multi-arg) or space for single-arg.
- Strings: "hello" or $"interpolated {var}".
- Lists: vec<int>:numbers = [1 | 2 | 3] (pipe-separated).
- Maps: map<string,int>:dict = {"key": 42 | "other": 99}.
- Operators: + - * / % == != < > <= >= && || ! (usual precedence).
- Comments: # single-line.

## Data Types

- Primitives: int, float, bool, string, void (for side-effects).
- Collections: vec<T> (dynamic array), map<K,V> (hash map).
- Functions: type -> type (e.g., int -> bool).
- Custom: type Name = variant1 | variant2 (simple enums/unions).

## Control Flow

- If: if cond -> thenExpr else elseExpr or multi-line indented.
- While: while cond -> body (loops until false).
- Foreach: foreach (type:item) in collection -> body.
- Recursion: Functions can call themselves.

## ESO-Specific Built-ins

- float:get-power = string:unit | string:powerType -> ... (returns power value).
- bool:is-dead = string:unit -> ...
- void:cast = int:abilityId -> ... (side-effect).
- Constants: string:player = "player", string:target = "reticleover".

## Code Samples

### 1. Hello World (Simple Function)

```
void:main = ->
    print "Hello, ESO World!"
```

**Description**: Defines a void-returning main function that prints a string. Call with main.

**Possible Errors**:
- If print is not defined in the runtime: "Undefined function: print".
- Type mismatch if passing non-string: "Expected string, got int".

### 2. Basic Arithmetic and Conditionals

```
int:add-squared = int:a | int:b ->
    let int:sum = a + b
    sum * sum

void:demo = ->
    let int:result = add-squared 3 | 4   # 49
    if result > 40 ->
        print $"Large: {result}"
    else
        print $"Small: {result}"
```

**Description**: Computes sum and squares it. Uses if-else as an expression (though void here). Interpolation with $"".

**Possible Errors**:
- Overflow on large ints: "Integer overflow".
- Non-int args: "Type mismatch: Expected int, got float".

### 3. Loops and Mutation (Factorial via While)

```
int:factorial = int:n ->
    mut int:result = 1
    mut int:i = 1
    while i <= n ->
        result = result * i
        i = i + 1
    result
```

**Description**: Computes factorial using a mutable loop. Demonstrates mutation with mut and assignment.

**Possible Errors**:
- Negative n: Loops forever if not guarded (runtime hang, no error).
- Mutation on immutable: "Cannot assign to immutable variable".
- Division by zero in body: "Division by zero".

### 4. Recursion (Fibonacci)

```
int:fib = int:n ->
    if n <= 1 ->
        n
    else
        fib (n - 1) | + fib (n - 2)
```

**Description**: Classic recursive Fibonacci. Turing-complete via recursion.

**Possible Errors**:
- Stack overflow on large n: "Recursion depth exceeded".
- Non-int n: "Type mismatch".

### 5. Collections and Foreach

```
vec<string>:filter-even = vec<int>:numbers ->
    mut vec<string>:result = []
    foreach (int:x) in numbers ->
        if x % 2 == 0 ->
            result = result + [$"{x} even"]
    result
```

**Description**: Filters even numbers, converts to strings. Uses vector append with +.

**Possible Errors**:
- Index out-of-bounds: "Vec index out of range".
- Type mismatch in append: "Cannot append string to vec<int>".
- Empty vec len: No error, returns 0.

### 6. ESO Combat Example

```
bool:should-execute = string:unit ->
    float:hp = get-power unit | "health"
    float:max = get-power unit | "health_max"
    hp / max <= 0.20

void:rotation = ->
    while true ->   # Infinite loop for combat pulse
        if should-execute target ->
            cast 123456
        else
            cast 456789
```

**Description**: Checks execute phase and casts abilities in a loop. Turing-complete infinite computation.

**Possible Errors**:
- Undefined ESO built-in: "External function not found".
- Division by zero if max=0: "Division by zero".
- Infinite loop: Runtime timeout in ESO context.

### 7. Higher-Order Functions (Map)

```
vec<int>:map-double = vec<int>:vals | (int -> int):fn ->
    mut vec<int>:result = []
    foreach (int:x) in vals ->
        result = result + [fn x]
    result

void:demo-map = ->
    let vec<int>:doubled = map-double [1 | 2 | 3] | (int:x -> x * 2)
    print doubled   # [2,4,6]
```

**Description**: Applies a function to each element. Shows first-class functions.

**Possible Errors**:
- Fn type mismatch: "Expected int -> int, got int -> float".
- Non-callable: "Attempt to call non-function".

### 8. Error-Prone Example (Division by Zero)

```
int:divide = int:a | int:b ->
    a / b
```

**Description**: Simple division.

**Possible Errors**:
- b=0: "Division by zero".
- Used in loop: Could cause panic mid-computation.

## Possible Errors

EsoLambda errors are runtime panics (no exceptions). Common ones include:

1. **Type Mismatch**: "Expected <type>, got <other>" – E.g., passing float to int param.
2. **Undefined Identifier**: "Undefined variable/function: <name>" – Referencing unbound name.
3. **Division by Zero**: "Division by zero" – In / or %.
4. **Index Out of Bounds**: "Vec/map access out of range" – Bad index/key.
5. **Recursion Depth Exceeded**: "Maximum recursion depth reached" – Deep recursion.
6. **Immutable Assignment**: "Cannot assign to immutable variable" – Mutating non-mut.
7. **Overflow/Underflow**: "Integer overflow" – Exceeding type bounds.
8. **ESO-Specific**: "External API error: <msg>" – E.g., invalid unit tag in get-power.
9. **Syntax Errors** (compile-time if interpreted): "Unexpected token: <token>" – Bad parsing.
10. **Infinite Loop**: No error, but runtime hang/timeout in hosted environments like ESO.

## Implementation Notes

- Interpreter would need to handle static type checking at "compile" time.
- For Turing completeness proof: Can simulate a Turing machine via infinite vec as tape, while loop as step, recursion for states.
- Limitations: No file I/O (except ESO APIs), bounded memory in practice.