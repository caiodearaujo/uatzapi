const API_BASE_ENDPOINTS = {
  CONNECT: '/connect',
  DEVICES: '/device',
};

let API_BASE_URL = '';

if (import.meta.env.VITE_ENVIRONMENT === 'production') {
  API_BASE_URL = 'http://localhost:8080';
} else {
    API_BASE_URL = 'http://localhost:8080';
}

const API_KEY_TOKEN = import.meta.env.VITE_API_KEY_TOKEN;

export { API_BASE_ENDPOINTS, API_BASE_URL, API_KEY_TOKEN };