<script setup lang="ts">
import { ref, onMounted } from 'vue';
import type { DashboardData } from './types/types.ts';
import NProgress from 'nprogress';
import 'nprogress/nprogress.css';

import StockHeader from './components/StockHeader.vue';
import KeyStats from './components/KeyStats.vue';
import EarningsChart from './components/EarningsChart.vue';
import RecTrends from './components/RecTrends.vue';
import InsiderTable from './components/InsiderTable.vue';

NProgress.configure({ showSpinner: false })

const symbol = ref('AAPL');
const dashboardData = ref<DashboardData | null>(null);
const error = ref<string | null>(null);

const fetchData = async () => {
  NProgress.start();
  error.value = null;

  try {
    const response = await fetch(`/api/dashboard?symbol=${symbol.value}`);

    if (!response.ok) {
      throw new Error(`Error: ${response.status} ${response.statusText}`);
    }

    const data = await response.json();
    dashboardData.value = data;
  } catch (err: any) {
    error.value = err.message || 'Failed to load data';
    console.error(err);
  } finally {
    NProgress.done();
  }
};

onMounted(() => {
  fetchData();
});
</script>

<template>
  <div class="min-h-screen w-full bg-gray-50 text-slate-800 font-sans">

    <div v-if="error" class="flex h-screen items-center justify-center">
      <div class="bg-red-50 text-red-600 px-6 py-4 rounded-xl border border-red-100 text-center">
        <p class="font-bold">Connection Failed</p>
        <p class="text-sm mt-1">{{ error }}</p>
        <button @click="fetchData"
          class="mt-3 px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 text-sm">Retry</button>
      </div>
    </div>

    <div v-else-if="dashboardData" class="w-full h-full p-4 md:p-6 space-y-6 animate-fade-in-up">

      <StockHeader :quote="dashboardData.quote" :profile="dashboardData.financials" />

      <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-6">

        <div class="2xl:col-span-1">
          <KeyStats :financials="dashboardData.financials" />
        </div>

        <div class="md:col-span-2 xl:col-span-2 2xl:col-span-2">
          <EarningsChart :earnings="dashboardData.earnings" />
        </div>

        <div class="md:col-span-1 xl:col-span-1 2xl:col-span-1 flex flex-col gap-6">
          <RecTrends :trend="dashboardData.recommendations" />
          <div class="flex-1 min-h-[300px]">
            <InsiderTable :transactions="dashboardData.insiders" />
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<style scoped>
.animate-fade-in-up {
  animation: fadeInUp 0.5s ease-out;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
