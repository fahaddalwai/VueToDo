<template>
    <div class="login-container">
      <h2>Login</h2>
      <form @submit.prevent="submitLogin">
        <div>
          <label for="username">Username:</label>
          <input type="text" id="username" v-model="credentials.username">
        </div>
        <div>
          <label for="password">Password:</label>
          <input type="password" id="password" v-model="credentials.password">
        </div>
        <button type="submit">Login</button>
      </form>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        credentials: {
          username: '',
          password: ''
        }
      };
    },
    methods: {
      submitLogin() {
        this.$emit('login-submit', this.credentials);
      }
    }
  };
  </script>
  
  <style scoped>
    /* Add your styling here */
  </style>
  