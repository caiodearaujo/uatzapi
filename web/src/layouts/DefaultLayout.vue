<script lang="ts">
import { defineComponent, computed } from 'vue';
import { useRoute } from 'vue-router';
import SideMenu from '../components/SideMenu.vue';

export default defineComponent({
  components: { SideMenu },
  setup() {
    const route = useRoute();

    const showMenu = computed(() => {
      // Verifica se o nome da rota atual é diferente de 'Login'
      return route.name !== 'Login';
    });

    return { showMenu };
  },
});
</script>

<template>
  <v-app>
    <SideMenu v-if="showMenu" />
    <v-main>
      <router-view />
    </v-main>
  </v-app>
</template>
