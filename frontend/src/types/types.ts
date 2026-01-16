export interface StockQuote {
  c: number;
  d: number;
  dp: number;
  h: number;
  l: number;
  o: number;
  pc: number;
  t: number;
}

export interface BasicFinancials {
  metric: {
    peBasicExclExtraTTM: number;
    marketCapitalization: number;
    "52WeekHigh": number;
    "52WeekLow": number;
    dividendYieldIndicatedAnnual: number;
    beta: number;
  };
  Series?: {
    Annual: {
      currentRatio?: Array<{ period: string; v: number }>;
      salesPerShare?: Array<{ period: string; v: number }>;
      netMargin?: Array<{ period: string; v: number }>;
    };
  };
  metricType: string;
  symbol: string;
}

export interface EarningsSurprise {
  actual: number;
  estimate: number;
  period: string;
  quarter: number;
  year: number;
  surprise: number;
  surprisePercent: number;
  symbol: string;
}

export interface RecommendationTrend {
  buy: number;
  hold: number;
  period: string;
  sell: number;
  strongBuy: number;
  strongSell: number;
  symbol: string;
}

export interface InsiderTransaction {
  name: string;
  share: number;
  change: number;
  filingDate: string;
  transactionDate: string;
  transactionPrice: number;
  symbol: string;
}

export interface DashboardData {
  quote: StockQuote;
  financials: BasicFinancials;
  earnings: EarningsSurprise[];
  recommendations: RecommendationTrend[];
  insiders: InsiderTransaction[];
}

export interface CompanyNews {
  category: string;
  datetime: number;
  headline: string;
  id: number;
  image: string;
  related: string;
  source: string;
  summary: string;
  url: string;
}

export interface MarketStatus {
  exchange: string
  holiday: string | null
  isOpen: boolean
  session: 'regular' | 'pre-market' | 'post-market'
  t: number
  timezone: string
}
