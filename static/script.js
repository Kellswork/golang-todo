var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (_) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
// create document query constants
var localhostAddress = "http://localhost:9000/todo";
var submitButton = document.querySelector("#submit");
var newTodoInput = document.querySelector("#new-todo input");
var todoListContainer = document.querySelector("#todos");
var isEditingtask = false; // editTask is for tracking when we are editing a tasks
var isComplete = false;
var todoID = "";
// create async func to fetch todos
function getTodos() {
    return __awaiter(this, void 0, void 0, function () {
        var response, responseData;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, fetch(localhostAddress)];
                case 1:
                    response = _a.sent();
                    return [4 /*yield*/, response.json()];
                case 2:
                    responseData = _a.sent();
                    console.log(responseData);
                    // return just the data frm the response
                    return [2 /*return*/, responseData.data];
            }
        });
    });
}
// async func to add a todo task
function createTodo(data) {
    return __awaiter(this, void 0, void 0, function () {
        var response, result, error_1;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    _a.trys.push([0, 3, , 4]);
                    return [4 /*yield*/, fetch(localhostAddress, {
                            method: "POST",
                            headers: {
                                "Content-Type": "application/json"
                            },
                            body: JSON.stringify(data)
                        })];
                case 1:
                    response = _a.sent();
                    return [4 /*yield*/, response.json()];
                case 2:
                    result = _a.sent();
                    console.log("Success:", result);
                    return [3 /*break*/, 4];
                case 3:
                    error_1 = _a.sent();
                    console.error("Error:", error_1);
                    return [3 /*break*/, 4];
                case 4: return [2 /*return*/];
            }
        });
    });
}
function updateTodo(id, data) {
    return __awaiter(this, void 0, void 0, function () {
        var response, result, error_2;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    _a.trys.push([0, 3, , 4]);
                    return [4 /*yield*/, fetch(localhostAddress + "/" + id, {
                            method: "PUT",
                            headers: {
                                "Content-Type": "application/json"
                            },
                            body: JSON.stringify(data)
                        })];
                case 1:
                    response = _a.sent();
                    return [4 /*yield*/, response.json()];
                case 2:
                    result = _a.sent();
                    console.log("Success:", result);
                    return [3 /*break*/, 4];
                case 3:
                    error_2 = _a.sent();
                    console.error("Error:", error_2);
                    return [3 /*break*/, 4];
                case 4: return [2 /*return*/];
            }
        });
    });
}
function deleteTodo(id) {
    return __awaiter(this, void 0, void 0, function () {
        var response, result, error_3;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    _a.trys.push([0, 3, , 4]);
                    return [4 /*yield*/, fetch(localhostAddress + "/" + id, {
                            method: "Delete"
                        })];
                case 1:
                    response = _a.sent();
                    return [4 /*yield*/, response.json()];
                case 2:
                    result = _a.sent();
                    console.log("Success:", result);
                    return [3 /*break*/, 4];
                case 3:
                    error_3 = _a.sent();
                    console.error("Error:", error_3);
                    return [3 /*break*/, 4];
                case 4: return [2 /*return*/];
            }
        });
    });
}
// submit todo event handler
function AddTask() {
    return __awaiter(this, void 0, void 0, function () {
        var value, data;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    console.log("add", isEditingtask);
                    value = newTodoInput.value;
                    data = { title: value };
                    if (!isEditingtask) return [3 /*break*/, 2];
                    return [4 /*yield*/, createTodo(data)];
                case 1:
                    _a.sent();
                    _a.label = 2;
                case 2: return [4 /*yield*/, loadTodos()];
                case 3:
                    _a.sent();
                    newTodoInput.value = "";
                    return [2 /*return*/];
            }
        });
    });
}
// function to display the toodlist saved in the db
function loadTodos() {
    var _a, _b;
    return __awaiter(this, void 0, void 0, function () {
        var todoList, deleteButton, _loop_1, _i, deleteButton_1, button, editButton, _loop_2, _c, editButton_1, button;
        return __generator(this, function (_d) {
            switch (_d.label) {
                case 0: return [4 /*yield*/, getTodos()];
                case 1:
                    todoList = _d.sent();
                    // if (todoListContainer === null) return
                    todoListContainer.innerHTML = "";
                    console.log("loading", todoList);
                    // show a message if there is no todo in the databse
                    if (todoList.length == 0) {
                        todoListContainer.innerHTML += " \n  <div class=\"todo\" style= \"display : " + (todoList.length === 0 ? "block" : "none") + ";\"> \n      <span id=\"no-todo\">You do not have any tasks \n      </span>\n    </div>";
                    }
                    else {
                        // loop through the list and display each todo as inner html
                        todoList.forEach(function (todo) {
                            // if (todoListContainer === null) return
                            todoListContainer.innerHTML += " \n  <div class=\"todo " + (todo.completed ? "completed" : "") + "\" > \n      <span id=\"todoname\" data-iscomplete =" + todo.completed + ">" + todo.title + "</span>\n      <div class=\"actions\">\n       <button data-id=" + todo.id + " class=\"edit\">\n       <i class=\"fas fa-edit\"></i>\n       </button>\n       <button data-id=" + todo.id + " class=\"delete\">\n       <i class=\"far fa-trash-alt\"></i>\n       </button>\n       <div>\n      \n    </div>";
                        });
                    }
                    deleteButton = Array.from(document.querySelectorAll(".delete"));
                    _loop_1 = function (button) {
                        button.onclick = function () {
                            return __awaiter(this, void 0, void 0, function () {
                                var todoID;
                                return __generator(this, function (_a) {
                                    switch (_a.label) {
                                        case 0:
                                            todoID = button.getAttribute("data-id") || "";
                                            return [4 /*yield*/, deleteTodo(todoID)];
                                        case 1:
                                            _a.sent();
                                            return [4 /*yield*/, loadTodos()];
                                        case 2:
                                            _a.sent();
                                            return [2 /*return*/];
                                    }
                                });
                            });
                        };
                    };
                    // fethc all the delete button, loop through it and attyach an onclick l;istner
                    for (_i = 0, deleteButton_1 = deleteButton; _i < deleteButton_1.length; _i++) {
                        button = deleteButton_1[_i];
                        _loop_1(button);
                    }
                    editButton = Array.from(document.querySelectorAll(".edit"));
                    _loop_2 = function (button) {
                        console.log("button1", button);
                        var parent_1 = (_b = (_a = button.parentNode) === null || _a === void 0 ? void 0 : _a.parentNode) === null || _b === void 0 ? void 0 : _b.children;
                        var todoName = (parent_1 === null || parent_1 === void 0 ? void 0 : parent_1[0]) || undefined;
                        var title = todoName.innerText;
                        button.onclick = function () {
                            todoID = button.getAttribute("data-id") || "";
                            // isComplete = todoName.getAttribute("data-iscomplete");
                            // if(newTodoInput === null) return;
                            newTodoInput.value = title;
                            isEditingtask = true;
                            // if edit task is true then we load the updateTasks() function when the user clicksd the submit button
                            console.log("edit tsk", isEditingtask);
                            console.log("id", todoID, isComplete);
                        };
                    };
                    // this func is responsible for loading the text in the input box when a user clicks the edit button
                    for (_c = 0, editButton_1 = editButton; _c < editButton_1.length; _c++) {
                        button = editButton_1[_c];
                        _loop_2(button);
                    }
                    console.log("loaded");
                    return [2 /*return*/];
            }
        });
    });
}
// call laod todos when you laod the page
loadTodos();
// this function is very sismilar to createTodo(), in the future, we will create a function that can handle both situations.
function updateTasks() {
    return __awaiter(this, void 0, void 0, function () {
        var value, data;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    value = newTodoInput.value;
                    isComplete = isComplete === true;
                    data = { title: value, completed: isComplete };
                    console.log(todoID, data);
                    if (!isEditingtask) return [3 /*break*/, 2];
                    return [4 /*yield*/, updateTodo(todoID, data)];
                case 1:
                    _a.sent();
                    _a.label = 2;
                case 2: return [4 /*yield*/, loadTodos()];
                case 3:
                    _a.sent();
                    newTodoInput.value = "";
                    todoID = "";
                    isComplete = false;
                    console.log("edit", isEditingtask);
                    isEditingtask = false;
                    return [2 /*return*/];
            }
        });
    });
}
// call the submit buttoon event listner for either updating or creating a new todo task.
submitButton.addEventListener("click", function () {
    isEditingtask ? updateTasks() : AddTask();
    console.log("last", isEditingtask);
});
