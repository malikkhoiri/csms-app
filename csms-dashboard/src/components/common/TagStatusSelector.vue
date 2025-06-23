<template>
  <v-select :items="['Accepted', 'Blocked', 'Expired']" v-model="status" dense />
</template>

<script setup>
import { ref, watch } from 'vue'
const props = defineProps({ status: String })
const status = ref(props.status)
watch(() => props.status, val => status.value = val)
</script> 