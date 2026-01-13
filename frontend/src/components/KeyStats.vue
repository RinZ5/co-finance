<script setup lang="ts">
import type { BasicFinancials } from '../types/types.ts'

defineProps<{
  financials: BasicFinancials;
}>();

const formatNumber = (val: number) =>
  new Intl.NumberFormat('en-US', { maximumFractionDigits: 2, notation: "compact", compactDisplay: "short" }).format(val);

const formatCurrency = (val: number) =>
  new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(val);
</script>

<template>
  <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-6 h-full">
    <h2 class="text-lg font-semibold text-slate-800 mb-4 border-b border-gray-100 pb-2">Key Statistics</h2>
    <div class="space-y-4">
      <div class="flex justify-between items-center">
        <span class="text-gray-500 text-sm">Market Cap</span>
        <span class="font-medium text-slate-900">{{ formatNumber(financials.metric.marketCapitalization) }}</span>
      </div>
      <div class="flex justify-between items-center">
        <span class="text-gray-500 text-sm">P/E Ratio</span>
        <span class="font-medium text-slate-900">{{ financials.metric.peBasicExclExtraTTM.toFixed(2) }}</span>
      </div>
      <div class="flex justify-between items-center">
        <span class="text-gray-500 text-sm">Beta</span>
        <span class="font-medium text-slate-900">{{ financials.metric.beta.toFixed(2) }}</span>
      </div>
      <div class="flex justify-between items-center">
        <span class="text-gray-500 text-sm">Div Yield</span>
        <span class="font-medium text-slate-900">{{ financials.metric.dividendYieldIndicatedAnnual.toFixed(2) }}%</span>
      </div>
      <div class="flex justify-between items-center">
        <span class="text-gray-500 text-sm">52W High</span>
        <span class="font-medium text-slate-900">{{ formatCurrency(financials.metric["52WeekHigh"]) }}</span>
      </div>
      <div class="flex justify-between items-center">
        <span class="text-gray-500 text-sm">52W Low</span>
        <span class="font-medium text-slate-900">{{ formatCurrency(financials.metric["52WeekLow"]) }}</span>
      </div>
    </div>
  </div>
</template>
