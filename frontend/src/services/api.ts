import axios from 'axios';
import type { StockQuote } from '../types/types';

const API_BASE_URL = import.meta.env.VITE_BACKEND_URL || '';

const apiClient = axios.create({
  baseURL: `${API_BASE_URL}/api`,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000,
});

export default {
  getQuote(symbol: string) {
    return apiClient.get<StockQuote>('/quote', { params: { symbol } });
  },

  getDashboard(symbol: string) {
    return apiClient.get('/dashboard', { params: { symbol } });
  },

  getCompanyNews(symbol: string, from: string, to: string) {
    return apiClient.get('/company-news', {
      params: { symbol, from, to }
    });
  },

  getMarketStatus(exchange: string) {
    return apiClient.get('market-status', {
      params: { exchange }
    })
  },

  get(endpoint: string, params: any = {}) {
    return apiClient.get(endpoint, { params });
  }
};
