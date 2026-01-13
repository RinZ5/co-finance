import { defineStore } from 'pinia';
import { ref } from 'vue';

const formatDate = (date: Date) => date.toISOString().split('T')[0];

export const useStockStore = defineStore('stock', () => {
  const symbol = ref('AAPL');

  const today = new Date();
  const threeDaysAgo = new Date();
  threeDaysAgo.setDate(today.getDate() - 3);

  const newsFrom = ref(formatDate(threeDaysAgo));
  const newsTo = ref(formatDate(today));

  function setSymbol(newSymbol: string) {
    if (newSymbol && typeof newSymbol === 'string') {
      symbol.value = newSymbol.toUpperCase();
    }
  }

  function setDateRange(from: string, to: string) {
    newsFrom.value = from;
    newsTo.value = to;
  }

  return { symbol, newsFrom, newsTo, setSymbol, setDateRange };
});
