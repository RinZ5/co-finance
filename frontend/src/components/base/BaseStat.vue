<template>
  <div class="stat-card">
    <div v-if="header" class="stat-header">
      {{ header }}
    </div>
    <div :class="valueClasses">
      <slot name="value">{{ formattedValue }}</slot>
    </div>
    <div v-if="$slots.change" :class="changeClasses">
      <slot name="change" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  header?: string
  value?: string | number
  trend?: 'up' | 'down' | 'neutral'
  formatValue?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  trend: 'neutral',
  formatValue: true
})

const formattedValue = computed(() => {
  if (!props.formatValue || props.value === undefined) return ''
  
  if (typeof props.value === 'number') {
    return new Intl.NumberFormat('en-US', {
      maximumFractionDigits: 2,
      minimumFractionDigits: props.value % 1 !== 0 ? 2 : 0
    }).format(props.value)
  }
  
  return props.value
})

const valueClasses = computed(() => [
  'stat-value',
  {
    'price-up': props.trend === 'up',
    'price-down': props.trend === 'down',
    'price-neutral': props.trend === 'neutral'
  }
])

const changeClasses = computed(() => [
  'text-sm mt-1',
  {
    'trend-up': props.trend === 'up',
    'trend-down': props.trend === 'down',
    'trend-neutral': props.trend === 'neutral'
  }
])
</script>