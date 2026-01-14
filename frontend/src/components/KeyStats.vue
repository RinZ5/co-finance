<script setup lang="ts">
import { computed } from 'vue';
import type { BasicFinancials } from '../types/types.ts'

const props = defineProps<{
  financials: BasicFinancials;
}>();

const formatNumber = (val: number) =>
  new Intl.NumberFormat('en-US', { maximumFractionDigits: 2, notation: "compact", compactDisplay: "short" }).format(val);

const formatCurrency = (val: number) =>
  new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(val);

const annualMetrics = computed(() => {
  const annual = props.financials.Series?.Annual;

  return {
    currentRatio: annual?.currentRatio?.[0]?.v,
    salesPerShare: annual?.salesPerShare?.[0]?.v,
    netMargin: annual?.netMargin?.[0]?.v
  };
});
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

      <div class="mt-4 pt-4 border-t border-gray-100">
        <h3 class="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-2">
          Financials (Annual)
        </h3>

        <div class="divide-y divide-gray-100">
          <div v-if="annualMetrics.currentRatio" class="flex justify-between items-center py-3">
            <span class="text-gray-500 text-sm">Current Ratio</span>
            <span class="font-medium text-slate-900">{{ annualMetrics.currentRatio.toFixed(2) }}</span>
          </div>
          <div v-if="annualMetrics.salesPerShare" class="flex justify-between items-center py-3">
            <span class="text-gray-500 text-sm">Sales/Share</span>
            <span class="font-medium text-slate-900">{{ formatCurrency(annualMetrics.salesPerShare) }}</span>
          </div>
          <div v-if="annualMetrics.netMargin" class="flex justify-between items-center py-3">
            <span class="text-gray-500 text-sm">Net Margin</span>
            <span class="font-medium text-slate-900">{{ (annualMetrics.netMargin * 100).toFixed(2) }}%</span>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>
