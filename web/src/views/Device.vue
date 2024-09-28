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

    <!-- Card com Texto -->
    <v-router-view></v-router-view>
  </v-container>
</template>