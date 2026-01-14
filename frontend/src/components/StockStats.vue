<script setup lang="ts">
import type { StockQuote } from '../types/types.ts'

const props = defineProps<{
  quote: StockQuote;
}>();

const formatCurrency = (val: number) =>
  new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(val);
</script>

<template>
  <div class="md:hidden mt-4 pt-4 border-t border-gray-100
           grid grid-cols-2 gap-x-6 gap-y-3 text-xs">
    <div class="flex justify-between">
      <span class="text-gray-500">Open</span>
      <span class="font-medium text-slate-800">
        {{ formatCurrency(quote.o) }}
      </span>
    </div>

    <div class="flex justify-between">
      <span class="text-gray-500">High</span>
      <span class="font-medium text-slate-800">
        {{ formatCurrency(quote.h) }}
      </span>
    </div>

    <div class="flex justify-between">
      <span class="text-gray-500">Low</span>
      <span class="font-medium text-slate-800">
        {{ formatCurrency(quote.l) }}
      </span>
    </div>

    <div class="flex justify-between">
      <span class="text-gray-500">Prev Close</span>
      <span class="font-medium text-slate-800">
        {{ formatCurrency(quote.pc) }}
      </span>
    </div>
  </div>
</template>