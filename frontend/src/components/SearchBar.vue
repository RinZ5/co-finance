<script setup lang="ts">
import { computed, ref } from 'vue';
import { useStockStore } from '../stores/stock';
import { useWebSocketStore } from '../stores/websocketStore';

const store = useStockStore();
const wsStore = useWebSocketStore();

const rawQuery = ref('');
const error = ref('');

const isValidStockSymbol = (symbol: string): boolean => {
  const cleanSymbol = symbol.trim().toUpperCase();
  return /^[A-Z]{1,5}$/.test(cleanSymbol);
};

const searchQuery = computed({
  get: () => rawQuery.value,
  set: (val: string) => {
    rawQuery.value = val.toUpperCase();
    error.value = '';
  }
});

const placeholder = computed(() => {
  return error.value || 'Search ticker (e.g. NVDA)';
});

const handleSubmit = () => {
  const symbol = searchQuery.value.trim();

  if (!symbol) {
    error.value = 'Please enter a stock symbol';
    return;
  }

  if (!isValidStockSymbol(symbol)) {
    error.value = 'Invalid stock symbol. Must be 1–5 letters';
    return;
  }

  store.setSymbol(symbol);
  wsStore.subscribe(symbol);
  searchQuery.value = '';
  error.value = '';
};
</script>

<template>
  <form @submit.prevent="handleSubmit" class="relative group ">
    <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
      <svg class="h-4 w-4 text-on-surface-tertiary transition-colors" fill="none" viewBox="0 0 24 24"
        stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
      </svg>
    </div>

    <input v-model="searchQuery" type="text" :placeholder="placeholder"
      class="w-full bg-surface border text-on-surface text-sm rounded-lg block pl-10 p-2.5 shadow-sm transition-all outline-none focus:outline-none focus:ring-0"
      :class="error ? 'border-error' : 'border-border'" />
  </form>
</template>
