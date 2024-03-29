<template>
  <div id="app">
    <h1>To-Do List</h1>
    <!-- Login and Signup Components -->
    <login v-if="!isLoggedIn" @login-submit="handleLogin"></login>
  <signup v-if="!isLoggedIn" @signup-submit="handleSignup"></signup>
    <to-do-form @todo-added="addToDo"></to-do-form>

    <h2 id="list-summary">{{listSummary}}</h2>


    <ul aria-labelledby="list-summary" class="stack-large">
      <li v-for="item in ToDoItems" :key="item.id">
        <to-do-item :label="item.label" :done="item.done" :id="item.id" @checkbox-changed="updateDoneStatus(item.id)" @item-deleted="deleteToDo(item.id)" @item-edited="editToDo(item.id, $event)"></to-do-item>
      </li>
    </ul>
  </div>
</template>

<script>
import toDoItem from "./components/ToDoItem.vue"
import toDoForm from "./components/ToDoForm.vue"
import Login from "./components/LoginForm.vue"
import Signup from "./components/SignupForm.vue"
import axios from 'axios';

export default {
  
  name: "app",
  components: {
    toDoItem,
    toDoForm,
    Login,
    Signup
  },
  computed:{
    listSummary(){
      const numberFinishedItems=this.ToDoItems.filter((item)=>item.done).length
      return `${numberFinishedItems} out of ${this.ToDoItems.length} items completed`
    }
  },
  data(){
    return{
      ToDoItems:[
    ],
    jwtToken: '', // JWT Token
    isLoggedIn: false, // Add this line
    };

  },
  
  methods:{

    async handleLogin(credentials) {
      try {
        const response = await axios.post('http://localhost:8080/login', credentials);
        
        this.jwtToken = response.data.token;
        console.log(this.jwtToken);
        this.isLoggedIn = true;
        this.fetchToDoItems();
      } catch (error) {
        console.error(error);

        //
      }
    },

    async handleSignup(userInfo) {
      try {
        await axios.post('http://localhost:8080/signup', userInfo);
        console.log("The user has successfuly logged in")
        this.isLoggedIn = true;
        this.fetchToDoItems();
      } catch (error) {
        console.error(error);
      }
    },



    async fetchToDoItems() {
      if (!this.jwtToken) return; // Ensure JWT token is available
      try {

        const response = await axios.get('http://localhost:8080/todo', {
          headers: {
            'Authorization': `Bearer ${this.jwtToken}`
          }
        });
        this.ToDoItems = response.data;
      } catch (error) {
        console.error(error);
      }
    },

    async addToDo(toDoLabel) {
      if (!this.jwtToken) return;
      try {
        const response = await axios.post('http://localhost:8080/todo', { label: toDoLabel, done: false }, {
          headers: {
            'Authorization': `Bearer ${this.jwtToken}`
          }
        });
        this.ToDoItems.push(response.data);
      } catch (error) {
        console.error(error);
      }
    },

    async updateDoneStatus(toDoId) {
      if (!this.jwtToken) return;
      const toDoToUpdate = this.ToDoItems.find((item) => item.id === toDoId);
      toDoToUpdate.done = !toDoToUpdate.done;
      try {
        await axios.put(`http://localhost:8080/todo/${toDoId}`, toDoToUpdate, {
          headers: {
            'Authorization': `Bearer ${this.jwtToken}`
          }
        });
      } catch (error) {
        console.error(error);
      }
    },

  async deleteToDo(toDoId) {
      if (!this.jwtToken) return;
      try {
        await axios.delete(`http://localhost:8080/todo/${toDoId}`, {
          headers: {
            'Authorization': `Bearer ${this.jwtToken}`
          }
        });
        const itemIndex = this.ToDoItems.findIndex((item) => item.id === toDoId);
        this.ToDoItems.splice(itemIndex, 1);
      } catch (error) {
        console.error(error);
      }
    },

    async editToDo(toDoId, newLabel) {
      if (!this.jwtToken) return;
      const toDoToEdit = this.ToDoItems.find((item) => item.id === toDoId);
      toDoToEdit.label = newLabel;
      try {
        await axios.put(`http://localhost:8080/todo/${toDoId}`, toDoToEdit, {
          headers: {
            'Authorization': `Bearer ${this.jwtToken}`
          }
        });
      } catch (error) {
        console.error(error);
      }
    }
  }
};

</script>

<style>
  /* Global styles */
  .btn {
    padding: 0.8rem 1rem 0.7rem;
    border: 0.2rem solid #4d4d4d;
    cursor: pointer;
    text-transform: capitalize;
  }
  .btn__danger {
    color: #fff;
    background-color: #ca3c3c;
    border-color: #bd2130;
  }
  .btn__filter {
    border-color: lightgrey;
  }
  .btn__danger:focus {
    outline-color: #c82333;
  }
  .btn__primary {
    color: #fff;
    background-color: #000;
  }
  .btn-group {
    display: flex;
    justify-content: space-between;
  }
  .btn-group > * {
    flex: 1 1 auto;
  }
  .btn-group > * + * {
    margin-left: 0.8rem;
  }
  .label-wrapper {
    margin: 0;
    flex: 0 0 100%;
    text-align: center;
  }
  [class*="__lg"] {
    display: inline-block;
    width: 100%;
    font-size: 1.9rem;
  }
  [class*="__lg"]:not(:last-child) {
    margin-bottom: 1rem;
  }
  @media screen and (min-width: 620px) {
    [class*="__lg"] {
      font-size: 2.4rem;
    }
  }
  .visually-hidden {
    position: absolute;
    height: 1px;
    width: 1px;
    overflow: hidden;
    clip: rect(1px 1px 1px 1px);
    clip: rect(1px, 1px, 1px, 1px);
    clip-path: rect(1px, 1px, 1px, 1px);
    white-space: nowrap;
  }
  [class*="stack"] > * {
    margin-top: 0;
    margin-bottom: 0;
  }
  .stack-small > * + * {
    margin-top: 1.25rem;
  }
  .stack-large > * + * {
    margin-top: 2.5rem;
  }
  @media screen and (min-width: 550px) {
    .stack-small > * + * {
      margin-top: 1.4rem;
    }
    .stack-large > * + * {
      margin-top: 2.8rem;
    }
  }
  /* End global styles */
  #app {
    background: #fff;
    margin: 2rem 0 4rem 0;
    padding: 1rem;
    padding-top: 0;
    position: relative;
    box-shadow:
      0 2px 4px 0 rgb(0 0 0 / 20%),
      0 2.5rem 5rem 0 rgb(0 0 0 / 10%);
  }
  @media screen and (min-width: 550px) {
    #app {
      padding: 4rem;
    }
  }
  #app > * {
    max-width: 50rem;
    margin-left: auto;
    margin-right: auto;
  }
  #app > form {
    max-width: 100%;
  }
  #app h1 {
    display: block;
    min-width: 100%;
    width: 100%;
    text-align: center;
    margin: 0;
    margin-bottom: 1rem;
  }
</style>
