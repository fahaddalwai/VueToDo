<template>
  <div id="app">
    <h1>To-Do List</h1>
    <to-do-form @todo-added="addToDo"></to-do-form>
    <ul>
      <li v-for="item in ToDoItems" :key="item.id"> //For each item in the the array its making an instance of todo
        <to-do-item :label="item.label" :done="item.done" :id="item.id"></to-do-item>
      </li>
    </ul>
  </div>
</template>

<script>
import toDoItem from "./components/ToDoItem.vue"
import toDoForm from "./components/ToDoForm.vue"

export default {
  
  name: "app",
  components: {
    toDoItem,
    toDoForm
  },
  
  data(){
    return{
      ToDoItems:[
        {id:crypto.randomUUID(),label:"Learn Vue",done:true},
        {id:crypto.randomUUID(),label:"Learn Go",done:false},
        {id:crypto.randomUUID(),label:"Make loging authentication page",done:false},
        {id:crypto.randomUUID(),label:"Finish Leetcode questions",done:false},
      ]
    };

  },
  methods:{
  addToDo(toDoLabel){
    this.ToDoItems.push({id:crypto.randomUUID(), label: toDoLabel, done: false});
  }
}
  
};

</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
