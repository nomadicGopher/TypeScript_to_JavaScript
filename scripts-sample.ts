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
