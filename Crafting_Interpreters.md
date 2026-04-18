![Paths of Implementation Mountain](https://www.craftinginterpreters.com/image/a-map-of-the-territory/mountain.png)
1. Paths of Implementation Mountain
   1. Source Code
      1. ![string](https://www.craftinginterpreters.com/image/a-map-of-the-territory/string.png)
   2. Scanning
      1. First step is scanning, also known as lexing/lexical analysis
      2. A scanner takins in a linear stream of charecters and chuncks them together into a series of something more aking to "words".
      3. Each of these "words" is called a token.
   3. Tokens
      1. ![tokens](https://www.craftinginterpreters.com/image/a-map-of-the-territory/tokens.png)
         1. Single Chars (, . : ;)
         2. Mutliple Chars
            1. Numbers (123)
            2. String Literals ("hi!")
            3. Identifiers(min) 
         3. Without meaning / ignored
            1. Comments (//)
            2. Whitespace ( )
   4. Parsing
      1. This is where syntax gets a grammer (larger expressions and statements out of smaller parts)
      2. A Parser takes flat sequence of tokens -> builds tree (nested nature of the grammer)
      3. Names for these trees
         1. Parse Tree
         2. Syntax Tree
         3. Abstract Syntax Tree (ASTs)
         4. Tree's
         5. ![AST](https://www.craftinginterpreters.com/image/a-map-of-the-territory/ast.png)
   5. Syntax Tree
      1. Transpiling
      2. High-Level Language
      3. Intermediate Representation
         1. Optimization
   6. Code Generation
      1. Bytecode
      2. Machine Code

- Static Analysis
  - **Binding**
  - **Resolution**
    - For each **identifier** find where it is defined and wire the two together. This is where **scope** matters.
    - If statically type this is when we type check
  - Semantic Insight is stored as 
    - attributes on the syntax tree itself, extra fields in the nodes that arn't initialzed during parsing but get filled in later.
    - Lookup Table, keys are identifiers (names of vars and declarations) in that case we call it a **symbol table** with the values tell us what that identifier refers to.

> Everything up to here is called the frontend

- Intermediate Representations (IR)
  - The compiler can be thought of as a pipeline. Where each stage takes the user code and organises it better for the next stage.
  - The frontend of the pipeline is specific to the source langugage. 
  - The backend is concerned with the final architecture where the program will run. 
  - In the middle the code may be stored in some intermediate representation (IR)
  - IR acts as an interface between these two languages.
  - IR's allow writing one front end for each source language. And then one back end for each target architecture.
  - Then you can mix and match the frontend and the backend for every combination.

- Optimization
  - Once we understand what the user's program means, we are free to swap it out with a different program that has the same semantics but implements them more efficiently, we can optimize that.
  - Contant Folding, if an expression always evaluates to the eact same value, we can evalute at compile time and replace the code for the expression with its result. 
  - Keywords
    - Constant Propogation
    - Common Subexpression Elimination
    - Loop Invariant code motion
    - global value numbering
    - strength reduction
    - scaler replacement of aggregates
    - dead code elimination
    - loop unrolling
- Code Generation
  - after all optimizations, the last step is to convert into machine code.

- Virtual Machine
  - if the compiler produces bytecode
    - write a mini compiler for each target architecture
    - Write a **virtual machine(VM)**, a program that emulates a hypothetical chip. 

- Runtime
  - Services the language provides while the program is running.
    - Garbage Collector
    - Instance of : Tests what kind of object you have.

- Single-pass Compilers
  - Interleaved parsing, analysis, code generation.
  - No intermediate data structures to store global information about the program, and you don’t revisit any previously parsed part of the code.

- Tree-walk Interpreters
  - Execution of code right after parsing it to an AST.
  - The interpreter traverses the syntax tree one branch and leaf at a time.
  
- Transpilers
  - What if you treated some other source language as if it were an intermediate representation?
  - You write a front end for your language.
  - Then in the back you produce a string of valid source code for some other language that's about as high level as yours.
  - Then use the existing compliation tools for THAT language to be able to execute.

- Just in time Compiliation (JIT)
  - The fastest way to execute code is by compiling it, but you may not know what architecture your end user machine supports.
  - On the end users machine, when the program is loaded, either from source, or platfrom independent bytecode, you compile it to native code for the architecture their computer supports.

- Compilers and Interpreters
  - Compiling is an implementation technique that involces translating a source langauge to some other, usuaully lower level form. Bytecode / Machine Code. Transpiling from one high level to one high level you are compiling too.
  - compiler means it trnslates source code to some other from but does not execute it.
  - Interpreter means it takes in source code and executes it immediately, it runs programs from source.
  - ![venn](https://www.craftinginterpreters.com/image/a-map-of-the-territory/venn.png)

## The golox language

- // --> Comment
- print "Hello, World!" --> built in statement, not a library function
- Dynamic Typing
  - Varibles can store vlaues of any type, and a single varible can even store value of different types at different times.
  - Invalid Operations on values of the wrong type trigger an error and runtime.
- Automatic Memory Management
  - Two common types
    - Reference Counting
    - Tracing Garbage Collection
- Data types
  - Booleans --> True, False
  - Numbers --> Double-precision floating point (basic int, decimal literals)
  - Strings --> "I am a string", "123", ""
  - Nil --> no value
- Expressions
  - Arithmetic
    -  add + me;
    -  subtract - me;
    -  multiply * me;
    -  divide / me;
    - the subexpressions on either side of the operator are **operands**
    - because there are two of them, these are called binary operators. Because the operator is fixed in the middle of the operands, AKA infix operators
      - as opposed to prefix where the operator comes before the operands
      - postfix where it comes after 
    - `-` this operator can also be used to negate a number
      - `-negateMe;`
    - these operators work on numbers and raise an exeption for other types, the exeption is the + and - can also pass it two strings to concatenate them.
- Comparison and equality
  - more operators that always return a Boolean result.
  - `less < than;`
  - `lessThan <= orEqual;`
  - `greater > than;`
  - `greaterThan >= orEqual;`
  - we can test two values of any kind for equality or inequality
    - `1 == 2; // false.`
    - `"cat" != "dog"; // true.`
    - `314 == "pi"; // false.`
    - `123 == "123" // false. `
- Logical Operators
  - `!` not operator, returns false if its operand is true, and vice versa
  - `and` expression determines if two values are BOTH true. it returns the left operand if its false, or the right operand otherwise.
    - `true and false; // false.`
    - `true and true; // true.`
  - `or` expression determines if EITHER of two values (or both) are true. it returns the left operand if it is true and the right operand otherwise.
    - `false or false; // false`
    - `true or false; // true`
  - `and` and `or` are like control flow structures because the **short-curcuit** since the `and` returns the left operand if it false and does not even evaluate the right on in that case. if the left operand of an or is true the right is skipped.

- Precedence and grouping
  - All operators have the same precedence and associativity that you’d expect coming from C.

- Statements
  - expressions main job is to produce a value
  - statements job is to produce and effect
  - since by definition, statements dont evaluate to a value, to be useful they have to otherwise change the world, by moidfying state, reading input, producing output.
  - ex. `print`
    - print statement evaluates a sinle expression and dispalys the result to the user.
  - an expression followed by a semicolon promotes the expression to a statement. this is called an expression statement.
  - to pack a seris of statement where a single one is expected you can wrap them up in a block.
  - ` {print "one statement"; print "Two Statements.";}`
  - blocks also affect scoping

- Variables
  - you declare varibles using `var` statements.
  - `var imAVarible = "here is my value";`
  - `var iAmNil;` defaults to `nil`

- Control Flow
  - `if` statement executes one of two statements based on some condition
    - `if (condition) { print "yes";} else {print "no";}`
  - `while` loop executes the body repeatedly as long as the condition expression evaluates to true
    - `var a = 1; while (a < 10) {print a; a = a + 1;}`
  - `for` loops
    - `for (var a = 1; a < 10; a = a + 1) {print a;}`

- Functions
  - `makeBreakfast(bacon, egss, toast);`
  - `makeBreakfast();`
  - to define a function
    - `fun printSum(a,b){ print a + b;}`
  - an argument is an actual value you pass to a function when you call it.
  - a parameter is a varible that holds the value of the argument isnide the body of the function.
  - The body of a function is always a block, inside it you can return a value using a `return` statement.
  - If exection reaches the end of the block without hitting a `return` it implicitly returns `nil`

- Closures
  - functions are first class in golox, which means they are real values that you can get a reference to, sotre in variables, pass around, etc...
  - You can declare local functions inside another function.
  - the function has to "hold on" to references to any surrounding varibles that it uses so that they stay around even after the outer funciton has returned, we call functions that do this `closures`



```c++
fun addPair(a, b) {
    return a + b;
} 
fun identity(a){
    return a;
}

print identity(addPair) (1, 2); // Prints "3"
```
- Classes
  - for a dynamically typed language, we need some way of defining compound data types to bundle blobs of stuf together 

- Classes or prototypes
  - when it comes to objects there are actually tqo approaches to them
  - classes
    - came first, think (C++, Java, C#, etc...)
    - two core concepts
      - instances 
        - store the state for ecah object and have a reference to the instances class
        - To call a method on an instance, there is always a level of indirection.
        - You look up the instance’s class and then you find the method there
        - In a statically typed language like C++, method lookup typically happens at compile time based on the static type of the instance, giving you static dispatch.
        - In contrast, dynamic dispatch looks up the class of the actual instance object at runtime. 
        - This is how virtual methods in statically typed languages and all methods in a dynamically typed language like golox work.
        - ![Class Lookup](https://www.craftinginterpreters.com/image/the-lox-language/class-lookup.png)
      - classes
        - contain the methods and inheritance chain.  
  - prototypes 
    - Javascript
    - merge the two concepts from classes
    - there are only objects, no classes
    - each indiviudal object amy contain state and methods.
    - Objects can directly inherit from each other, or delegate to in prototype lingo.
    - ![Prototype Lookup](https://www.craftinginterpreters.com/image/the-lox-language/prototype-lookup.png)

- Classes in golox
```c++
class Breakfast {
    cook() {
        print "Eggs a-frying'!";
    }

    serve(who) {
        print "Enjoy your breakfast, " + who + ".";
    }
}
``` 
  - The body of a class contains its methods
  - They look like function delcerations without the `fun` keyword
  - When the class declaration is executed, golox creates a class object and stores that in variable named after the class.
  - Just like functions, classes are first class in golox. 
  - `var someVarible = Breakfast;`
  - `someFunction(Breakfast);`
  - in golox the class itself is a factory function for instances
  - `var breakfast = Breakfast();`
  - `print breakfast  // "Breakfast instance"`

- Instantiation and initialization
  - golox lets you freely add properties onto objects.
  - `breakfast.meat = "sausage";`
  - `breakfast.bread = "sourdough";` 
  - assigning to a field creates it if it does not already exist
  - if you want to access a field or method on the current object within a method you use good old `this`
  - if your class has a method named `init()` it is called automatically when the object is contructed.
  - Any parameters passed to the class are forwarded to its initalizer

- Inheritance
  - golox supports single inheritance
  - when you declare a class, you can specify a class that it inherits from using a less-than (<) operator.
  - to call methods on the class that is inherited we use the `super` keyword.

- The Standard Library
  - the set of functionality that is implemented directly in the interpreter and that all user-defined behavior is built on top of.