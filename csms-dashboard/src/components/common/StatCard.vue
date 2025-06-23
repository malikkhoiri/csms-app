<template>
  <v-card :color="color || 'primary'" class="pa-4" dark>
    <v-row align="center">
      <v-col cols="3" class="d-flex align-center justify-center">
        <v-icon size="36">{{ icon }}</v-icon>
      </v-col>
      <v-col cols="9">
        <div class="text-h6">{{ title }}</div>
        <div class="text-h4 font-weight-bold">{{ value }}</div>
      </v-col>
    </v-row>
  </v-card>
</template>

<script setup>
const props = defineProps({
  title: String,
  value: [String, Number],
  icon: String,
  color: String,
})
</script> 