<script setup lang="ts">
import { computed, ref } from 'vue';
import { useStockStore } from '../stores/stock';
import { useWebSocketStore } from '../stores/websocketStore';

const store = useStockStore();
const wsStore = useWebSocketStore();
const rawQuery = ref('')

const searchQuery = computed({
  get: () => rawQuery.value,
  set: (val: string) => {
    rawQuery.value = val.toUpperCase()
  }
})

const handleSubmit = () => {
  if (searchQuery.value.trim()) {
    const symbol = searchQuery.value

    store.setSymbol(symbol)
    wsStore.subscribe(symbol)

    searchQuery.value = '';
  }
};
</script>

<template>
  <form @submit.prevent="handleSubmit" class="relative group w-full max-w-xs">
    <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
      <svg class="h-4 w-4 text-gray-400 group-focus-within:text-blue-500 transition-colors" fill="none"
        viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
      </svg>
    </div>
    <input v-model="searchQuery" type="text" placeholder="Search ticker (e.g. NVDA)..."
      class="w-full bg-white border border-gray-200 text-gray-700 text-sm rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 block pl-10 p-2.5 shadow-sm transition-all outline-none">
  </form>
</template>
