// create document query constants
let localhostAddress = "http://localhost:9000/todo"
let submitButton = document.querySelector('#submit');
let newTodoInput = document.querySelector('#new-todo input');
let todoListContainer = document.querySelector('#todos')
let isEditingtask = false // editTask is for tracking when we are editing a tasks
let isComplete = ''
let todoID = ''

// create async func to fetch todos
async function getTodos() {
 const response = await fetch(localhostAddress);
 const responseData = await response.json();
 console.log(responseData);
 // return just the data frm the response
 return responseData.data
}

// async func to add a todo task
async function createTodo(data) {
 try {
  // send POST request with user oinput as the req body
  const response = await fetch(localhostAddress, {
   method: "POST",
   headers: {
    "Content-Type": "application/json",
   },
   body: JSON.stringify(data),
  });

  const result = await response.json();
  console.log("Success:", result);
 } catch (error) {
  console.error("Error:", error);
 }
}

async function updateTodo(id, data) {
 try {
  // async func to send PUT request to update data iwth the id
  const response = await fetch(`${localhostAddress}/${id}`, {
   method: "PUT",
   headers: {
    "Content-Type": "application/json",
   },
   body: JSON.stringify(data),
  });
  const result = await response.json();
  console.log("Success:", result);
 } catch (error) {
  console.error("Error:", error);
 }
}

async function deleteTodo(id) {
 try {
  const response = await fetch(`${localhostAddress}/${id}`, { method: "Delete" });
  const result = await response.json();
  console.log("Success:", result);
 } catch (error) {
  console.error("Error:", error);
 }
}


// submit todo event handler
async function AddTask() {
 {
  console.log('add', isEditingtask);

  // when the user clicks submit, call the createTodo
  // fetch all todos again to add the just created todo
  let value = newTodoInput.value;
  const data = { title: value }
  if (isEditingtask) u
  await createTodo(data)
  await loadTodos();
  newTodoInput.value = ''
  console.log('edit tsk', isEditingtask);

 }
}



// function to display the toodlist saved in the db
async function loadTodos() {
 const todoList = await getTodos();
 todoListContainer.innerHTML = '';
 console.log('loading', todoList)
 // show a message if there is no todo in the databse
 if (todoList.length == 0) {
  todoListContainer.innerHTML += ` 
  <div class="todo" style= "display : ${todoList.length === 0 ? 'block' : 'none'};"> 
      <span id="no-todo">You do not have any tasks 
      </span>
    </div>`;
 } else {
  // loop through the list and display each todo as inner html
  todoList.forEach(
   todo => {
    document.querySelector('#todos').innerHTML += ` 
  <div class="todo ${todo.completed ? "completed" : ''}" > 
      <span id="todoname" data-iscomplete =${todo.completed}>${todo.title}</span>
      <div class="actions">
       <button data-id=${todo.id} class="edit">
       <i class="fas fa-edit"></i>
       </button>
       <button data-id=${todo.id} class="delete">
       <i class="far fa-trash-alt"></i>
       </button>
       <div>
      
    </div>`
   });

 }

 // this func is inside load todo because the html is omnmly created after the todo loads.
 let deleteButton = await document.querySelectorAll('.delete');

 // fethc all the delete button, loop through it and attyach an onclick l;istner
 for (let button of deleteButton) {
  button.onclick = async function () {
   // get the todo id using the html attrubite stuff
   const todoID = button.getAttribute('data-id')
   await deleteTodo(todoID);
   await loadTodos();
  }
 }


// edit button func is similar to delete.
 let editButton = await document.querySelectorAll('.edit');

 // this func is responsible for loading the text in the input box when a user clicks the edit button
 for (let button of editButton) {
  console.log('button1', button)
  const parent = button.parentNode.parentNode.children;
  const todoName = parent[0];
  const title = todoName.innerText;

  button.onclick = function () {
   todoID = button.getAttribute('data-id');
   isComplete = todoName.getAttribute('data-iscomplete');
   newTodoInput.value = title
   isEditingtask = true
   // if edit task is true then we load the updateTasks() function when the user clicksd the submit button
   console.log('edit tsk', isEditingtask);
   console.log('id', todoID, isComplete)
  }
 }


 console.log('loaded')
}
// call laod todos when you laod the page
loadTodos();

// this function is very sismilar to createTodo(), in the future, we will create a function that can handle both situations.
async function updateTasks() {
 let value = newTodoInput.value;
 isComplete = isComplete === true;

 const data = { title: value, completed: isComplete }
 console.log(todoID, data)
 if (isEditingtask)
  await updateTodo(todoID, data)
 await loadTodos();
 newTodoInput.value = ''
 todoID = ''
 isComplete = ''
 console.log('edit', isEditingtask);
 isEditingtask = false
}
// call the submit buttoon event listner for either updating or creating a new todo task.
submitButton.addEventListener('click', () => {
 isEditingtask ? updateTasks() : AddTask()
 console.log('last', isEditingtask)
})