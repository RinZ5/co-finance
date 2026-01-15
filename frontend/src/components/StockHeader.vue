<script setup lang="ts">
import type { StockQuote, BasicFinancials, MarketStatus } from '../types/types.ts'
import BaseCard from './base/BaseCard.vue'
import StockSymbol from './StockSymbol.vue'
import StockPrice from './StockPrice.vue'
import StockStats from './StockStats.vue'
import PriceRange from './PriceRange.vue'

const props = defineProps<{
  quote: StockQuote;
  profile: BasicFinancials;
  marketStatus?: MarketStatus;
  symbol: string;
}>();
</script>

<template>
  <BaseCard class="grid grid-cols-1 md:grid-cols-12 gap-6 items-center">

    <div class="md:col-span-4 flex flex-col gap-4">
      <div class="flex justify-between gap-4">
        <StockSymbol :symbol="profile.symbol" :market-status="marketStatus" />
        <div class="md:hidden">
          <StockPrice :quote="quote" :symbol="symbol" variant="mobile" />
        </div>
      </div>
      <StockStats :quote="quote" />
    </div>

    <div class="md:col-span-8 hidden md:flex items-stretch gap-6 pl-6 border-l border-border-light h-full">

      <div class="grow min-w-0 pr-40">
        <PriceRange :quote="quote" />
      </div>

      <div class="shrink-0 text-right">
        <StockPrice :quote="quote" :symbol="symbol" variant="desktop" />
      </div>

    </div>

  </BaseCard>
</template>
