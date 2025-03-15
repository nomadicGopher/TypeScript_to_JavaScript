## Instructions
1. Download your OS's version of the program from [Releases](https://github.com/nomadicGopher/TypeScript_to_JavaScript/releases).
    * **TypeScript_to_JavaScript.exe**: Windows
    * **TypeScript_to_JavaScript**: Linux
2. Add a scripts.ts file in the same directory, or remember to set the -file command-line-argument if it is outside the program's directory.
3. Run the program with the following potential flags: `TypeScript_to_JavaScript.exe -minify -stream=#.##`

```
-file string
    Path to the TypeScript file. (default "scripts.ts")
-minify
    Minify the JavaScript output.
-stream int
    File streaming minimum threshold in megabytes. (default 10)
```

## Sample TypeScript content
```typescript
// Define an interface for a Person
interface Person {
  firstName: string;
  lastName: string;
  age: number;
  greet(): string;
}

// Create a class that implements the Person interface
class Student implements Person {
  firstName: string;
  lastName: string;
  age: number;
  studentId: number;

  constructor(firstName: string, lastName: string, age: number, studentId: number) {
    this.firstName = firstName;
    this.lastName = lastName;
    this.age = age;
    this.studentId = studentId;
  }

  // Implement the greet method
  greet(): string {
    return `Hello, my name is ${this.firstName} ${this.lastName} and I am ${this.age} years old.`;
  }

  // Additional method to get the student ID
  getStudentId(): number {
    return this.studentId;
  }
}

// Create an instance of the Student class
const student = new Student("John", "Doe", 20, 12345);

// Call the greet method
console.log(student.greet());

// Call the getStudentId method
console.log(`My student ID is ${student.getStudentId()}.`);
```