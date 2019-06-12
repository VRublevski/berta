# An interpreter for Berta programming language.

This repo contains source code for a small programming language. 
The language is dynamically typed and the supports  procedural paradygm. 
The builtin types are booleans, ints, doubles, strings and arrays. 
There is no a builtin support for compound types, e.g. structs, but this functionality can be implemented via closures, functions are first class objects in Berta language.
Some examples of Berta language: 


```
>> var array = [1, 2, 3, 4, 5]
>> var sum = fun(ar, n){ 
      var s = 0; 
      for(var i = 0; i < n; i = i + 1){ 
          s = s + array[i] 
      } 
      return s; 
}
>> sum(array, 5)
15

>> var makePerson = fun(name, age){ 
      var dispatch = fun(field){ 
          if (field == "name"){ 
              return name
          } else { 
              return age 
          } 
      } 
      return dispatch; 
}
>> var man = makePerson("Bruce, 34)
>> man("name")
"Bruce"
>> man("age")
34

```

### Download and Install

To install the interpreter you can simply clone the repositary, 
open the repositary folder in your terminal and run

    go install ./main 
