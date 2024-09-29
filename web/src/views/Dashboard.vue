<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios';
import { API_BASE_URL, API_ENDPOINTS } from '@/api.config';
import Device from '@/types/device';

const devices = ref<Device[]>([]);
const loading = ref(true); // Adiciona uma variável para controlar o loading

onMounted(async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}${API_ENDPOINTS.DEVICES}`);
    devices.value = response.data.map((deviceData: any) => new Device(deviceData));
  } catch (error) {
    console.error(error);
  } finally {
    loading.value = false; // Finaliza o loading
  }
});
</script>

<template>
  <v-app>
    <v-main>
      <v-container>
        <v-row>
          <v-col>
            <v-img
              src="/social.svg"
              alt="Header Image"
              height="200px"
            ></v-img>
            <h1 class="text-center mt-4">Meus dispositivos</h1>
          </v-col>
        </v-row>
        <v-divider/>

        <!-- Loader -->
        <v-container v-if="loading" class="d-flex justify-center align-center" style="height: 400px;">
          <v-progress-circular indeterminate color="primary"></v-progress-circular>
        </v-container>

        <v-container v-else>
          <v-row >
            <v-col cols="12" md="6" v-for="item in devices" :key="item.id">
              <v-card
                class="mx-auto hover-card"
                prepend-icon="mdi-cellphone"
                :subtitle="item.number"
                :title="item.pushName"
                @click="() => $router.push({ name: 'Device', params: { id: item.id } })" 
                transition="fade-transition"
              >
                <v-container>
                  <div class="d-flex justify ga-2">
                    <v-chip color="secondary" label>{{item.contacts}} Contatos</v-chip>
                  </div>
                </v-container>
              </v-card>
            </v-col>
          </v-row>
        </v-container>
      </v-container>
    </v-main>
  </v-app>
</template>

<style scoped>
.hover-card {
  transition: box-shadow 0.3s, transform 0.3s; /* Efeito de hover */
}

.hover-card:hover {
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
  transform: scale(1.02);
}

.fade-transition {
  transition: opacity 0.5s ease;
}

.v-enter-active, .v-leave-active {
  opacity: 0;
}

.v-enter, .v-leave-to /* .v-leave-active in <2.1.8 */ {
  opacity: 1;
}
</style>
