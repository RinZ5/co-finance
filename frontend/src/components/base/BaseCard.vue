<template>
  <div :class="cardClasses">
    <slot />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  variant?: 'default' | 'hover' | 'stat'
  padding?: 'none' | 'sm' | 'md' | 'lg'
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
  padding: 'lg'
})

const cardClasses = computed(() => [
  'rounded-xl border shadow-md',
  {
    'bg-surface border-border-light': props.variant === 'default',
    'card-hover': props.variant === 'hover',
    'card-stat': props.variant === 'stat',
    'p-lg': props.padding === 'lg',
    'p-md': props.padding === 'md',
    'p-sm': props.padding === 'sm',
    'p-0': props.padding === 'none'
  }
])
</script>