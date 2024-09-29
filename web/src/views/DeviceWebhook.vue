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
  </v-container>
</template>

<script>
import axios from 'axios';
import { API_BASE_URL, API_KEY_TOKEN} from '@/api.config';

export default {
  data() {
    return {
      webhookUrl: '',
      webhooks: [],
      snackbarVisible: false,
      loading: false,  // Adicione este estado
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
        this.webhooks = response.data;
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
