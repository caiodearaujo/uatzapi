<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios';
import { API_BASE_URL, API_ENDPOINTS } from '@/api.config';
import Device from '@/types/device';

const devices = ref<Device[]>([]);

onMounted(async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}${API_ENDPOINTS.DEVICES}`);
    devices.value = response.data.map((deviceData: any) => new Device(deviceData));
  } catch (error) {
    console.error(error);
  }
});
</script>

<template>
  <v-app>
    <v-main>
      <v-container>
        <h1>Dashboard</h1>
        <v-divider/>
        <v-container>
          <v-row dense>
    <v-col cols="12" md="6" v-for="item in devices" :key="item.id">
      <v-card
        class="mx-auto"
        prepend-icon="mdi-cellphone"
        :subtitle="item.number"
        :title="item.pushName"
      >
        <v-container>
              <div class="d-flex justify ga-2">
              <v-chip color="secondary" label>{{item.contacts}} Contatos</v-chip>
              </div>
        </v-container>

        <v-card-actions>
          <v-btn variant="tonal" prepend-icon="mdi-magnify" :to="{name: 'Device', params: { id: item.id }}">Detalhes</v-btn>
        </v-card-actions>
      </v-card>
    </v-col>

    </v-row>
        </v-container>
      </v-container>
      </v-main>
  </v-app>
</template>
