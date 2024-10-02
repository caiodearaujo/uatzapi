<template>
  <v-container class="login-container">
    <v-card class="pa-5" elevation="2" max-width="400">
      <v-img
        src="@/assets/logo_uatz.svg"
        height="100"
        class="mb-4"
        contain
      ></v-img>
      <v-form @submit.prevent="login">
        <v-text-field
          v-model="username"
          label="Username"
          required
          outlined
        ></v-text-field>
        <v-text-field
          v-model="password"
          label="Password"
          type="password"
          required
          outlined
        ></v-text-field>
        <v-btn type="submit" color="primary" class="mt-4" block>Login</v-btn>
      </v-form>
    </v-card>
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
      const expectedUsername = import.meta.env.VITE_USER_USERNAME || 'admin';
      const expectedPassword = import.meta.env.VITE_USER_PASSWORD || 'admin';

      if (this.username === expectedUsername && this.password === expectedPassword) {
        // Login bem-sucedido
        localStorage.setItem('loginTime', Date.now());
        localStorage.setItem('isLogged', 'true');
        this.$router.push({ name: 'Dashboard' });
      } else {
        alert('Usuário ou senha inválidos');
      }
    },
  },
};
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}

.v-card {
  width: 100%; /* Largura total */
  max-width: 400px; /* Limitar a largura máxima */
  margin: 0 auto; /* Centralizar horizontalmente */
}
</style>
