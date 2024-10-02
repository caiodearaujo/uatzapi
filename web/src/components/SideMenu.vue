<script lang="ts">
import { defineComponent, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';

export default defineComponent({
  setup() {
    const route = useRoute();
    const router = useRouter(); // Para redirecionar após logout

    const items = ref<any[]>([
      { title: 'Dashboard', icon: 'mdi-view-dashboard', to: '/dashboard' },
      { title: 'Conectar', icon: 'mdi-qrcode-plus', to: '/connect' },
    ]);

    const extraItems = ref<any[]>([]);
    const extraMenuName = ref<string>('');

    // Watch the route to update the selected item
    watch(route, (newRoute) => {
      extraItems.value = Array.isArray(newRoute.meta.extraMenuItems) 
        ? newRoute.meta.extraMenuItems 
        : [];
      extraMenuName.value = typeof newRoute.name === 'string' ? newRoute.name : '';
    });

    // Função de logout
    const logout = () => {
      // Limpa os dados do localStorage
      localStorage.removeItem('loginTime');
      localStorage.removeItem('isLogged');
      
      // Redireciona para a página de login
      router.push({ name: 'Login' });
    };

    return { items, extraItems, extraMenuName, route, logout };
  }
});
</script>

<template>
  <v-navigation-drawer floating
        permanent
        app>
    <v-img
      src="@/assets/logo_uatz.svg"
      height="80"
    />
  <v-divider/>
    <v-list density="compact" nav>
      <v-list-item v-for="item in items" :key="item.title" :to="item.to" :prepend-icon="item.icon" :value="item.title" :title="item.title"> </v-list-item>
      <v-divider/>
      <template v-if="extraItems.length > 0">
        <v-list-subheader>{{ extraMenuName }}</v-list-subheader>
        <v-list-item style="padding-left: 33px;" v-for="item in extraItems" :key="item.title" :to="{ name: item.to.name, params: { id: route.params.id } }"  :prepend-icon="item.icon" :value="item.title" :title="item.title"> </v-list-item>
      </template>
      <v-divider/>
      <!-- Botão de logout com a ação associada -->
      <v-list-item prepend-icon="mdi-logout" title="Logout" value="logout" @click="logout"></v-list-item>

    </v-list>
  </v-navigation-drawer>
</template>
