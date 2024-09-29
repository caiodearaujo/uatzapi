<template>
  <v-container>
    <!-- Header com Imagem e Título -->
    <v-row>
      <v-col>
        <v-img
          src="/devices.svg"
          alt="Header Image"
          height="200px"
        ></v-img>
        <h1 class="text-center mt-4">Dispositivo conectado</h1>
      </v-col>
    </v-row>

    <!-- Loader -->
    <v-row v-if="loading" class="justify-center">
      <v-col cols="12" class="text-center">
        <v-progress-circular indeterminate color="primary"></v-progress-circular>
        <p>Carregando...</p>
      </v-col>
    </v-row>

    <!-- Conteúdo Principal - Apenas exibe se não estiver carregando -->
    <v-row v-if="!loading">
      <!-- Card Principal com informações do device -->
      <v-col>
        <v-card class="mx-auto" max-width="600" elevation="6">
          <v-row>
            <v-col cols="12" md="4">
              <v-img :src="device.profile_picture" aspect-ratio="1" class="rounded-circle"></v-img>
            </v-col>
            <v-col cols="12" md="8">
              <v-card-title>{{ device.push_name || device.business_name }}</v-card-title>
              <v-card-subtitle>{{ device.phone_number }}</v-card-subtitle>
              <v-card-text>
                Webhook:
                <span :class="{'text-grey': !device.webhook}">
                  {{ device.webhook || 'Inativo' }}
                </span>
              </v-card-text>
            </v-col>
          </v-row>
        </v-card>
      </v-col>
    </v-row>

    <!-- Tabela com os contatos -->
    <v-row v-if="!loading">
      <v-col cols="12">
        <h2 class="text-center">Contatos</h2>
        <v-card class="mx-auto" max-width="800" outlined>
          <v-data-table
            :headers="headers"
            :items="device.contacts"
            :search="search"
            :no-data-text="device.contacts.length === 0 ? 'Nenhum contato encontrado' : 'Carregando...'"
            :sort-by="['name']"
            :items-per-page="15"
            :loading="loading"
            class="elevation-1"
          >
            <template v-slot:top>
              <v-card flat>
                <v-card-title class="d-flex align-center pe-2">
                  <v-icon icon="mdi-contacts"></v-icon> &nbsp;
                  Pesquisar por nome
                  <v-spacer></v-spacer>
                  <v-text-field
                    v-model="search"
                    density="compact"
                    label="Pesquisar"
                    prepend-inner-icon="mdi-magnify"
                    variant="solo-filled"
                    flat
                    hide-details
                    single-line
                  ></v-text-field>
                </v-card-title>
              </v-card>
            </template>

            <template v-slot:item.profile_picture="{ item }">
              <v-img
                :src="item.profile_picture || 'https://via.placeholder.com/100'"
                aspect-ratio="1"
                class="rounded-circle"
                max-width="100px"
              ></v-img>
            </template>

            <!-- Adicionando botão de copiar no campo número -->
            <template v-slot:item.number="{ item }">
              <v-row>
                <v-col cols="8" justify="right">
                  {{ item.number }}
                  <v-btn icon="mdi-content-copy" @click="copyToClipboard(item.number)" size="x-small" variant="plain"></v-btn>
                </v-col>
              </v-row>
            </template>
          </v-data-table>
        </v-card>
      </v-col>
    </v-row>

    <!-- Snackbar -->
    <v-snackbar v-model="snackbar" :timeout="2000" :transition="zoom" location="end" variant="tonal" target="cursor">
      {{ snackbarMessage }}
    </v-snackbar>
  </v-container>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      device: {
        profile_picture: '',
        push_name: '',
        business_name: '',
        phone_number: '',
        webhook: '',
        contacts: [],
      },
      loading: true, // Estado de carregamento
      search: '', // Valor de busca
      snackbar: false, // Controle do Snackbar
      snackbarMessage: '', // Mensagem do Snackbar
      headers: [
        { text: 'Foto', value: 'profile_picture', sortable: false },
        { text: 'Nome', value: 'name' },
        { text: 'Número', value: 'number' },
      ], // Cabeçalhos da tabela
    };
  },
  created() {
    this.fetchDeviceData();
  },
  methods: {
    fetchDeviceData() {
      const deviceId = parseInt(this.$route.params.id);
      axios
        .get(`http://localhost:8080/device/${deviceId}`)
        .then((response) => {
          this.device = response.data;
          this.loading = false; // Desativa o loader quando os dados são carregados
        })
        .catch((error) => {
          console.error('Erro ao buscar dados do dispositivo:', error);
          this.loading = false; // Desativa o loader mesmo se ocorrer erro
        });
    },
    // Função para copiar o número para o clipboard e exibir o Snackbar
    copyToClipboard(text) {
      navigator.clipboard.writeText(text).then(
        () => {
          this.snackbarMessage = 'Número copiado!';
          this.snackbar = true;
        },
        (error) => {
          console.error('Erro ao copiar o texto: ', error);
          this.snackbarMessage = 'Erro ao copiar número';
          this.snackbar = true;
        }
      );
    },
  },
};
</script>

<style>
.rounded-circle {
  border-radius: 50%;
  margin: 5%
}
.text-grey {
  color: grey;
}
</style>
