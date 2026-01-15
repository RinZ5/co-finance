<script setup lang="ts">
import type { StockQuote } from '../types/types.ts'

const props = defineProps<{
  quote: StockQuote;
}>();

const formatCurrency = (val: number) =>
  new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(val);
</script>

<template>
  <div class="md:hidden mt-4 pt-4 border-t border-border-light
           grid grid-cols-2 gap-x-6 gap-y-3 text-xs">
    <div class="flex justify-between">
      <span class="text-on-surface-tertiary">Open</span>
      <span class="font-medium text-on-surface">
        {{ formatCurrency(quote.o) }}
      </span>
    </div>

    <div class="flex justify-between">
      <span class="text-on-surface-tertiary">High</span>
      <span class="font-medium text-on-surface">
        {{ formatCurrency(quote.h) }}
      </span>
    </div>

    <div class="flex justify-between">
      <span class="text-on-surface-tertiary">Low</span>
      <span class="font-medium text-on-surface">
        {{ formatCurrency(quote.l) }}
      </span>
    </div>

    <div class="flex justify-between">
      <span class="text-on-surface-tertiary">Prev Close</span>
      <span class="font-medium text-on-surface">
        {{ formatCurrency(quote.pc) }}
      </span>
    </div>
  </div>
</template>