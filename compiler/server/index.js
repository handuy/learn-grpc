function map(f, a) {
  let result = []; // Create a new Array
  let i; // Declare variable
  for (i = 0; i != a.length; i++)
    result[i] = f(a[i]);
  return result;
}
const f = function(x) {
   return x * x * x; 
}
let numbers = [1, 3, 5, 77, 111];
let cube = map(f,numbers);
console.log(cube);