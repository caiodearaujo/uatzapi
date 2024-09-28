<template>
  <v-container>
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
        <br/>
        <v-divider></v-divider>
        <br/>
        <v-data-table
          :headers="headers"
          :items="webhooks"
          class="elevation-1"
          :pagination="{ rowsPerPage: -1 }" 
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
  </v-container>
</template>

<script>
import axios from 'axios';
import { API_BASE_URL } from '@/api.config';

export default {
  data() {
    return {
      webhookUrl: '',
      webhooks: [],
      snackbarVisible: false,
      headers: [
        { title: 'Webhook URL', key: 'webhook_url' },
        { title: 'Data de Criação', key: 'timestamp' },
        { title: 'Ações', key: 'actions', sortable: false },
      ],
      rules: {
        required: value => !!value || 'Campo obrigatório.',
      },
    };
  },
  methods: {
    async fetchWebhook() {
      const deviceId = parseInt(this.$route.params.id);

      try {
        const response = await axios.get(`${API_BASE_URL}/webhook/${deviceId}`);
        if (response.status === 200) {
          this.webhookUrl = response.data.webhook_url;
        }
      } catch (error) {
        if (error.response.status !== 204) {
          console.error('Erro ao buscar webhook:', error);
        }
      }
    },
    async submitWebhook() {
      const deviceId = parseInt(this.$route.params.id);

      try {
        const response = await axios.post(`${API_BASE_URL}/webhook`, {
          device_id: deviceId,
          webhook_url: this.webhookUrl,
        });
        if (response.status === 200) {
          this.snackbarVisible = true; // Exibe o snackbar
          this.fetchAllWebhooks();
        }
      } catch (error) {
        console.error('Erro ao salvar webhook:', error);
      }
    },
    async fetchAllWebhooks() {
      const deviceId = parseInt(this.$route.params.id);
      try {
        const response = await axios.get(`${API_BASE_URL}/webhook/${deviceId}/all`);
        this.webhooks = response.data;
      } catch (error) {
        console.error('Erro ao listar webhooks:', error);
      }
    },
    async deactivateWebhook(webhookId) {
      const deviceId = this.$route.params.id;
      try {
        await axios.delete(`${API_BASE_URL}/webhook/${webhookId}`, {
          data: {
            device_id: deviceId,
            webhook_url: this.webhookUrl, // Use o URL do webhook específico se necessário
          },
        });
        this.fetchAllWebhooks();
      } catch (error) {
        console.error('Erro ao desativar webhook:', error);
      }
    },
  },
  mounted() {
    this.fetchWebhook();
    this.fetchAllWebhooks();
  },
};
</script>

<style scoped>
.v-container {
  max-width: 800px;
  margin: 0 auto;
}
</style>
