<script setup lang="ts">
import type { StockQuote, BasicFinancials, MarketStatus } from '../types/types.ts'
import StockSymbol from './StockSymbol.vue'
import StockPrice from './StockPrice.vue'
import StockStats from './StockStats.vue'
import PriceRange from './PriceRange.vue'

const props = defineProps<{
  quote: StockQuote;
  profile: BasicFinancials;
  marketStatus?: MarketStatus
}>();
</script>

<template>
  <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-6
              grid grid-cols-1 md:grid-cols-12 gap-6 items-center">

    <div class="md:col-span-4 flex flex-col gap-4">
      <div class="flex justify-between gap-4">
        <StockSymbol :symbol="profile.symbol" :market-status="marketStatus" />
        <div class="md:hidden">
          <StockPrice :quote="quote" variant="mobile" />
        </div>
      </div>
      <StockStats :quote="quote" />
    </div>

    <div class="md:col-span-8 hidden md:flex items-center gap-3 pl-6 border-l border-gray-100 h-full">

      <div class="grow min-w-0">
        <PriceRange :quote="quote" />
      </div>

      <div class="shrink-0 text-right">
        <StockPrice :quote="quote" variant="desktop" />
      </div>

    </div>

  </div>
</template>
