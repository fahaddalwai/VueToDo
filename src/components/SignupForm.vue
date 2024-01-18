<template>
    <div class="signup-container">
      <h2>Signup</h2>
      <form @submit.prevent="submitSignup">
        <div>
          <label for="username">Username:</label>
          <input type="text" id="username" v-model="userInfo.username">
        </div>
        
        <div>
          <label for="password">Password:</label>
          <input type="password" id="password" v-model="userInfo.password">
        </div>
        <button type="submit">Signup</button>
      </form>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        userInfo: {
          username: '',
          email: '',
          password: ''
        }
      };
    },
    methods: {
      submitSignup() {
        this.$emit('signup-submit', this.userInfo);
      }
    }
  };
  </script>
  
  <style scoped>
    /* Add your styling here */
  </style>
  