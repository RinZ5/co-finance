<script setup lang="ts">
import { ref, computed } from 'vue';
import type { StockQuote, BasicFinancials } from '../types/types.ts'

const props = defineProps<{
  quote: StockQuote;
  profile: BasicFinancials;
}>();

const formatCurrency = (val: number) =>
  new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(val);

const isPositive = computed(() => props.quote.d >= 0);
</script>

<template>
  <div
    class="bg-white rounded-xl shadow-sm border border-gray-100 p-6 flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
    <div>
      <div class="flex items-center gap-2">
        <h1 class="text-3xl font-bold text-slate-900">{{ profile.symbol }}</h1>
        <span class="px-2 py-0.5 rounded text-xs font-semibold bg-gray-200 text-gray-600">US Market</span>
      </div>
      <p class="text-gray-500 text-sm mt-1">Real-time Data</p>
    </div>

    <div class="text-right">
      <div class="text-4xl font-bold tracking-tight text-slate-900">
        {{ formatCurrency(quote.c) }}
      </div>
      <div
        :class="['flex items-center justify-end gap-2 text-lg font-medium', isPositive ? 'text-green-600' : 'text-red-600']">
        <span>{{ isPositive ? '+' : '' }}{{ quote.d.toFixed(2) }}</span>
        <span class="text-sm bg-opacity-10 px-2 py-0.5 rounded-full"
          :class="isPositive ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'">
          {{ isPositive ? '+' : '' }}{{ quote.dp.toFixed(2) }}%
        </span>
      </div>
    </div>
  </div>
</template>
