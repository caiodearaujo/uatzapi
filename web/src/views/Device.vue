<template>
  <v-container>
    <v-system-bar color="primary" dark>
      <strong>{{ device.push_name || device.business_name }}</strong>
      <v-icon left class="ml-2">mdi-cellphone</v-icon>{{ formatPhoneNumber(device.phone_number) }}
    <v-col class="text-right">
      Webhook: 
      <span :class="{'text-grey': !device.webhook}">
        {{ device.webhook || 'Inativo' }}
      </span>
    </v-col>
</v-system-bar>

    <!-- Loader -->
    <v-row v-if="loading" class="justify-center">
      <v-col cols="12" class="text-center">
        <v-progress-circular indeterminate color="primary"></v-progress-circular>
        <p>Carregando...</p>
      </v-col>
    </v-row>

    <v-row>
      <v-col>
        <v-img
          src="/form.svg"
          alt="Header Image"
          height="200px"
        ></v-img>
        <h1 class="text-center mt-4">Cadastro de Webhook</h1>

        <v-form @submit.prevent="submitWebhook">
          <v-text-field
            v-model="webhookUrl"
            label="URL do Webhook"
            :rules="[rules.required]"
            required
          ></v-text-field>

          <v-btn type="submit" color="primary">Salvar Webhook</v-btn>
        </v-form>

        <v-snackbar v-model="snackbarVisible" timeout="3000">
          Webhook salvo com sucesso!
          <template v-slot:action="{ attrs }">
            <v-btn color="blue" text v-bind="attrs" @click="snackbarVisible = false">Fechar</v-btn>
          </template>
        </v-snackbar>

        <br />
        <v-divider></v-divider>
        <br />

        <!-- Loader -->
        <v-row v-if="loading" class="justify-center">
          <v-progress-circular indeterminate color="primary"></v-progress-circular>
        </v-row>

        <v-data-table
          v-else
          :headers="headers"
          :items="webhooks"
          class="elevation-1"
          :pagination="{ rowsPerPage: -1 }" 
          :no-data-text="webhooks.length === 0 ? 'Nenhum webhook cadastrado' : 'Carregando...'"
          hide-default-footer
        >
          <template v-slot:top>
            <v-toolbar>
              <v-toolbar-title>Lista de Webhooks</v-toolbar-title>
            </v-toolbar>
          </template>
          <template v-slot:item.actions="{ item }">
            <v-btn
              v-if="item.active"
              @click="deactivateWebhook(item.id)"
              color="red"
            >
              Desativar
            </v-btn>
            <v-btn
              v-else
              color="grey"
              :disabled="true"
            >
              Desativado
            </v-btn>
          </template>
        </v-data-table>
      </v-col>
    </v-row>
    <!-- Snackbar -->
    <v-snackbar v-model="snackbar" :timeout="2000" location="end" variant="tonal" target="cursor">
      {{ snackbarMessage }}
    </v-snackbar>
  </v-container>
</template>

<script>
import axios from 'axios';
import { API_BASE_URL, API_BASE_ENDPOINTS, API_KEY_TOKEN } from '@/api.config';

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
      webhookUrl: '',
      webhooks: [],
      rules: {
        required: value => !!value || 'Campo obrigatório.',
      },
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
        .get(`${API_BASE_URL}${API_BASE_ENDPOINTS.DEVICES}/${deviceId}`, {
          headers: {
            'X-Api-Key': API_KEY_TOKEN,
          },
        })
        .then((response) => {
          this.device = response.data;
          this.loading = false; // Desativa o loader quando os dados são carregados
        })
        .catch((error) => {
          console.error('Erro ao buscar dados do dispositivo:', error);
          this.loading = false; // Desativa o loader mesmo se ocorrer erro
        });
    },
    async fetchWebhook() {
      const deviceId = parseInt(this.$route.params.id);
      this.loading = true; // Inicie o carregamento
      try {
        const response = await axios.get(`${API_BASE_URL}/webhook/${deviceId}`, {
          headers: {
            'X-Api-Key': API_KEY_TOKEN,
          },
        });
        if (response.status === 200) {
          this.webhookUrl = response.data.webhook_url;
        }
      } catch (error) {
        if (error.response.status !== 204) {
          console.error('Erro ao buscar webhook:', error);
        }
      } finally {
        this.loading = false; // Encerre o carregamento
      }
    },
    async submitWebhook() {
      const deviceId = parseInt(this.$route.params.id);
      this.loading = true; // Inicie o carregamento
      try {
        const response = await axios.post(`${API_BASE_URL}/webhook`, {
          device_id: deviceId,
          webhook_url: this.webhookUrl,
        }, {
          headers: {
            'X-Api-Key': API_KEY_TOKEN,
          },
        });
        if (response.status === 200) {
          this.snackbarVisible = true; // Exibe o snackbar
          this.fetchAllWebhooks();
        }
      } catch (error) {
        console.error('Erro ao salvar webhook:', error);
      } finally {
        this.loading = false; // Encerre o carregamento
      }
    },
    async fetchAllWebhooks() {
      const deviceId = parseInt(this.$route.params.id);
      this.loading = true; // Inicie o carregamento
      try {
        const response = await axios.get(`${API_BASE_URL}/webhook/${deviceId}/all`, {
          headers: {
            'X-Api-Key': API_KEY_TOKEN,
          },
        });
        console.log(response.data);
        this.webhooks = response.data || [];
      } catch (error) {
        console.error('Erro ao listar webhooks:', error);
      } finally {
        this.loading = false; // Encerre o carregamento
      }
    },
    async deactivateWebhook(webhookId) {
      const deviceId = this.$route.params.id;
      this.loading = true; // Inicie o carregamento
      try {
        await axios.delete(`${API_BASE_URL}/webhook/${webhookId}`, {
          data: {
            device_id: deviceId,
            webhook_url: this.webhookUrl,
          },
          headers: {
            'X-Api-Key': API_KEY_TOKEN,
          },
        });
        this.fetchAllWebhooks();
      } catch (error) {
        console.error('Erro ao desativar webhook:', error);
      } finally {
        this.loading = false; // Encerre o carregamento
      }
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
    formatPhoneNumber(phoneNumber) {
      // Remove any non-numeric characters
      const cleanedNumber = phoneNumber.replace(/\D/g, '');

      // Check if the number starts with the Brazilian country code '55' and has the correct length
      if (cleanedNumber.length === 13 && cleanedNumber.startsWith('55')) {
        const countryCode = '+55';
        const areaCode = cleanedNumber.substring(2, 4);      // Extract the area code (DDD)
        const firstDigit = cleanedNumber.substring(4, 5);    // Extract the first digit (usually 9 for mobile phones)
        const firstPart = cleanedNumber.substring(5, 9);     // Extract the first part of the phone number
        const secondPart = cleanedNumber.substring(9, 13);   // Extract the second part of the phone number

        // Return the formatted phone number in the international format
        return `${countryCode} (${areaCode}) ${firstDigit} ${firstPart}-${secondPart}`;
      }

      // Return a message for invalid numbers
      return phoneNumber;
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
