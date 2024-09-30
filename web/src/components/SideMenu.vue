<script lang="ts">
import { defineComponent, ref, watch } from 'vue';
import { useRoute } from 'vue-router';

export default defineComponent({
  setup() {
    const route = useRoute();

    const items = ref<any[]>([
      { title: 'Dashboard', icon: 'mdi-view-dashboard', to: '/dashboard' },
      { title: 'Conectar', icon: 'mdi-qrcode-plus', to: '/connect' },
      { title: 'Configurações', icon: 'mdi-cog', to: '/settings' },
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

    return { items, extraItems, extraMenuName, route };
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
      <v-list-item prepend-icon="mdi-logout" title="Logout" value="logout"></v-list-item>

    </v-list>
  </v-navigation-drawer>
</template>
