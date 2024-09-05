<template>
  <v-container>
    <v-form @submit.prevent="login">
      <v-text-field v-model="username" label="Username" required></v-text-field>
      <v-text-field v-model="password" label="Password" type="password" required></v-text-field>
      <v-btn type="submit" color="primary">Login</v-btn>
    </v-form>
  </v-container>
</template>

<script>
export default {
  data() {
    return {
      username: '',
      password: '',
    };
  },
  methods: {
    login() {
      const expectedUsername = "admin";
      const expectedPassword = "admin";

      if (this.username === expectedUsername && this.password === expectedPassword) {
        // Login successful
        localStorage.setItem('loginTime', Date.now());
        localStorage.setItem('isLogged', 'true');
        this.$router.push({ name: 'Dashboard' });
      } else {
        alert('Invalid username or password');
      }
    },
  },
};
</script>
