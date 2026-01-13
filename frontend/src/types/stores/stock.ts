import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useStockStore = defineStore('stock', () => {
  const symbol = ref('AAPL');

  function setSymbol(newSymbol: string) {
    if (newSymbol && typeof newSymbol === 'string') {
      symbol.value = newSymbol.toUpperCase();
    }
  }

  return { symbol, setSymbol };
});
