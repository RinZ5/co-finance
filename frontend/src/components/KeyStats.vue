<script setup lang="ts">
import { computed } from 'vue';
import BaseCard from './base/BaseCard.vue';
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
  <BaseCard class="h-full">
    <h2 class="section-header border-b border-border-light pb-2">Key Statistics</h2>
    <div class="space-y-4">
      <div class="flex justify-between items-center">
        <span class="text-on-surface-tertiary text-sm">Market Cap</span>
        <span class="font-medium text-on-surface">{{ formatNumber(financials.metric.marketCapitalization) }}</span>
      </div>
      <div class="flex justify-between items-center">
        <span class="text-on-surface-tertiary text-sm">P/E Ratio</span>
        <span class="font-medium text-on-surface">{{ financials.metric.peBasicExclExtraTTM.toFixed(2) }}</span>
      </div>
      <div class="flex justify-between items-center">
        <span class="text-on-surface-tertiary text-sm">Beta</span>
        <span class="font-medium text-on-surface">{{ financials.metric.beta.toFixed(2) }}</span>
      </div>
      <div class="flex justify-between items-center">
        <span class="text-on-surface-tertiary text-sm">Div Yield</span>
        <span class="font-medium text-on-surface">{{ financials.metric.dividendYieldIndicatedAnnual.toFixed(2) }}%</span>
      </div>
      <div class="flex justify-between items-center">
        <span class="text-on-surface-tertiary text-sm">52W High</span>
        <span class="font-medium text-on-surface">{{ formatCurrency(financials.metric["52WeekHigh"]) }}</span>
      </div>
      <div class="flex justify-between items-center">
        <span class="text-on-surface-tertiary text-sm">52W Low</span>
        <span class="font-medium text-on-surface">{{ formatCurrency(financials.metric["52WeekLow"]) }}</span>
      </div>

      <div class="mt-4 pt-4 border-t border-border-light">
        <h3 class="text-xs font-semibold text-on-surface-tertiary uppercase tracking-wider mb-2">
          Financials (Annual)
        </h3>

        <div class="divide-y divide-border-light">
          <div v-if="annualMetrics.currentRatio" class="flex justify-between items-center py-3">
            <span class="text-on-surface-tertiary text-sm">Current Ratio</span>
            <span class="font-medium text-on-surface">{{ annualMetrics.currentRatio.toFixed(2) }}</span>
          </div>
          <div v-if="annualMetrics.salesPerShare" class="flex justify-between items-center py-3">
            <span class="text-on-surface-tertiary text-sm">Sales/Share</span>
            <span class="font-medium text-on-surface">{{ formatCurrency(annualMetrics.salesPerShare) }}</span>
          </div>
          <div v-if="annualMetrics.netMargin" class="flex justify-between items-center py-3">
            <span class="text-on-surface-tertiary text-sm">Net Margin</span>
            <span class="font-medium text-on-surface">{{ (annualMetrics.netMargin * 100).toFixed(2) }}%</span>
          </div>
        </div>
      </div>

    </div>
  </BaseCard>
</template>
