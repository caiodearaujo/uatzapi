<script setup lang="ts">
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { API_BASE_URL, API_BASE_ENDPOINTS, API_KEY_TOKEN } from '@/api.config';
import Device from '@/types/device';

const devices = ref<Device[]>([]);
const loading = ref(true); // Variável para controlar o estado de loading
const noDevices = ref(false); // Variável para controlar se não há dispositivos

onMounted(async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}${API_BASE_ENDPOINTS.DEVICES}`, {
      headers: {
        'X-Api-Key': API_KEY_TOKEN
      },
    });

    if (response.data === null) {
      response.data = [];
    }
    devices.value = Array.isArray(response.data)
      ? response.data.map((deviceData: any) => new Device(deviceData))
      : [];

    if (devices.value.length === 0) {
      noDevices.value = true; // Atualiza o estado quando não há dispositivos
    }
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

        <!-- Loader enquanto os dados estão sendo carregados -->
        <v-container v-if="loading" class="d-flex justify-center align-center" style="height: 400px;">
          <v-progress-circular indeterminate color="primary"></v-progress-circular>
        </v-container>

        <!-- Exibe o v-banner quando não há dispositivos -->
        <v-container v-else-if="noDevices">
          <v-banner color="warning" icon="mdi-alert">
            <template v-slot:actions>
              <v-btn @click="() => $router.push({ name: 'Connect' })">
                Conectar novo dispositivo
              </v-btn>
            </template>
            Nenhum dispositivo foi conectado ainda.
          </v-banner>
        </v-container>

        <!-- Lista de dispositivos -->
        <v-container v-else>
          <v-row>
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
                    <v-chip color="secondary" label>{{ item.contacts }} Contatos</v-chip>
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
  transition: box-shadow 0.3s, transform 0.3s;
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

.v-enter, .v-leave-to {
  opacity: 1;
}
</style>
