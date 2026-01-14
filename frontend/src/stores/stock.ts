import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useStockStore = defineStore('stock', () => {
  const symbol = ref('AAPL');

  const today = new Date();
  const threeDaysAgo = new Date();
  threeDaysAgo.setDate(today.getDate() - 3);

  const newsFrom = ref<Date>(threeDaysAgo);
  const newsTo = ref<Date>(today);

  function setSymbol(newSymbol: string) {
    if (newSymbol && typeof newSymbol === 'string') {
      symbol.value = newSymbol.toUpperCase();
    }
  }

  function setDateRange(from: Date, to: Date) {
    newsFrom.value = from;
    newsTo.value = to;
  }

  return { symbol, newsFrom, newsTo, setSymbol, setDateRange };
});
