<script setup lang="ts">
import { computed } from 'vue';
import type { StockQuote } from '../types/types.ts'

const props = defineProps<{
  quote: StockQuote;
  variant?: 'mobile' | 'desktop';
}>();

const formatCurrency = (val: number) =>
  new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(val);

const isPositive = computed(() => props.quote.d >= 0);

const priceClasses = computed(() => {
  const base = props.variant === 'mobile' 
    ? 'text-3xl' 
    : 'text-4xl font-bold tracking-tight text-slate-900';
  return base;
});

const changeClasses = computed(() => {
  const base = props.variant === 'mobile'
    ? 'flex items-center justify-end gap-2 text-sm font-medium'
    : 'flex items-center justify-end gap-2 text-lg font-medium';
  return `${base} ${isPositive.value ? 'text-green-600' : 'text-red-600'}`;
});

const percentBadgeClasses = computed(() => {
  return `text-sm px-2 py-0.5 rounded-full bg-opacity-10 ${
    isPositive.value 
      ? 'bg-green-100 text-green-700' 
      : 'bg-red-100 text-red-700'
  }`;
});
</script>

<template>
  <div :class="variant === 'mobile' ? 'md:hidden text-right' : 'hidden md:block text-right'">
    <div :class="priceClasses">
      {{ formatCurrency(quote.c) }}
    </div>

    <div :class="changeClasses">
      <span>{{ isPositive ? '+' : '' }}{{ quote.d.toFixed(2) }}</span>
      
      <span v-if="variant === 'desktop'" :class="percentBadgeClasses">
        {{ isPositive ? '+' : '' }}{{ quote.dp.toFixed(2) }}%
      </span>
      
      <span v-else>
        ({{ isPositive ? '+' : '' }}{{ quote.dp.toFixed(2) }}%)
      </span>
    </div>
  </div>
</template>