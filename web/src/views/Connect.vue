<script setup lang="ts">
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { API_BASE_URL, API_BASE_ENDPOINTS, API_KEY_TOKEN } from '@/api.config';
const qrCode = ref('');
const loading = ref(true);
const maxRetries = 20; // Número máximo de requisições
let retryCount = 0;

const fetchQRCode = async () => {
  try {
    loading.value = true; // Inicia o loader
    const response = await axios.get(`${API_BASE_URL}${API_BASE_ENDPOINTS.CONNECT}`, {
      headers: {
        'X-Api-Key': API_KEY_TOKEN
      },
    });
    qrCode.value = response.data.qrCode; // Recebe o QR code da API
  } catch (error) {
    console.error(error);
  } finally {
    loading.value = false; // Para o loader
  }
};

onMounted(() => {
  fetchQRCode(); // Primeira requisição

  const intervalId = setInterval(() => {
    retryCount += 1;
    if (retryCount >= maxRetries) {
      clearInterval(intervalId); // Para as requisições após 20 tentativas
    } else {
      fetchQRCode(); // Faz a requisição a cada 20 segundos
    }
  }, 20000); // 20 segundos
});
</script>


<template>
  <v-app>
    <v-main>
      <v-container>
        <v-row>
          <v-col>
            <v-img
              src="/connect.svg"
              alt="Header Image"
              height="200px"
            ></v-img>
            <h1 class="text-center mt-4">Conectar novo dispositivo</h1>
          </v-col>
        </v-row>
        <v-divider />
        <v-row dense>
          <v-col cols="12" md="6">
            <v-list lines="one">
              <v-list-item>
                <v-list-item-content>
                  <v-list-item-title>
                    1. Abra o WhatsApp no seu celular
                  </v-list-item-title>
                </v-list-item-content>
              </v-list-item>
              <v-list-item>
                <v-list-item-content>
                  <v-list-item-title>
                    2. Toque em <b>Mais opções</b>
                      <img src="/icons/conn_moreopt.svg" height="24" width="24" style="display: inline-block; margin: 0 5px;" />
                      no Android <br/> 
                      ou em <b>Configurações</b>
                      <img src="/icons/conn_conf.svg" height="24" width="24" style="display: inline-block; margin: 0 5px;" />
                      no iPhone
                    </v-list-item-title>
                </v-list-item-content>
              </v-list-item>
              <v-list-item>
                <v-list-item-content>
                  <v-list-item-title>
                    3. Toque em <b>Dispositivos conectados</b> e, <br/>em seguida, em <b>Conectar dispositivo.</b>
                  </v-list-item-title>
                </v-list-item-content>
              </v-list-item>
              <v-list-item>
                <v-list-item-content>
                  <v-list-item-title
                    >4. Aponte seu celular para esta tela para escanear o QR
                    Code.</v-list-item-title
                  >
                </v-list-item-content>
              </v-list-item>
            </v-list>
          </v-col>
          <v-col cols="12" md="6" class="d-flex justify-center fill-height">
            <v-progress-circular v-if="loading" indeterminate color="primary"></v-progress-circular>
            <img v-else :src="'data:image/png;base64,' + qrCode" alt="QR Code" />
          </v-col>
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<style scoped>
.fill-height {
  
  height: 100%;
  min-height: 100%;
}
</style>