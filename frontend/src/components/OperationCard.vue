<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'

const props = defineProps({
  operationData: {
    type: Object,
    required: true,
  },
})

const emit = defineEmits(['update:dates', 'update:numbers', 'update:comment'])

// Reactive data
const selectedNumbers = ref<string>('')
const selectedDates = ref<string>('')
const comment = ref<string>('')

// Watchers to sync with parent data
watch(
  () => props.operationData.numbers,
  (newVal) => {
    selectedNumbers.value = Array.isArray(newVal) ? newVal.join(', ') : ''
  },
  { immediate: true }
)

watch(
  () => props.operationData.dates,
  (newVal) => {
    selectedDates.value = Array.isArray(newVal) ? newVal.join(', ') : ''
  },
  { immediate: true }
)

watch(
  () => props.operationData.comment,
  (newVal) => {
    comment.value = newVal || ''
  },
  { immediate: true }
)

// Methods
const onNumberChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  const value = target.value
  const numbersArray = value.split(',').map((n) => n.trim()).filter((n) => n)
  emit('update:numbers', numbersArray)
}

const onDateChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  const value = target.value
  const datesArray = value.split(',').map((d) => d.trim()).filter((d) => d)
  emit('update:dates', datesArray)
}

const onCommentChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  emit('update:comment', target.value)
}
</script>

<template>
  <Card class="operation-card">
    <CardHeader class="operation-header">
      <CardTitle class="operation-title text-sm">{{ operationData.title }}</CardTitle>
    </CardHeader>

    <CardContent class="operation-content space-y-3 p-4">
      <div class="operation-field space-y-1">
        <Label class="text-xs font-medium text-muted-foreground">Tooth Numbers</Label>
        <Input
          :value="selectedNumbers"
          placeholder="e.g., 1.1, 1.2, 2.3"
          @input="onNumberChange"
          class="h-9"
        />
      </div>

      <div class="operation-field space-y-1">
        <Label class="text-xs font-medium text-muted-foreground">Dates</Label>
        <Input
          :value="selectedDates"
          placeholder="e.g., 2024-01-15, 2024-02-20"
          @input="onDateChange"
          class="h-9"
        />
      </div>

      <div class="operation-field space-y-1">
        <Label class="text-xs font-medium text-muted-foreground">Comment</Label>
        <Input
          :value="comment"
          placeholder="Add notes..."
          @input="onCommentChange"
          class="h-9"
        />
      </div>
    </CardContent>
  </Card>
</template>

<style scoped>
.operation-card {
  border: 1px solid hsl(var(--border));
  border-radius: 0.5rem;
  background: hsl(var(--card));
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.operation-header {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 0.75rem;
  border-bottom: 1px solid hsl(var(--border));
  background-color: hsl(var(--muted));
}

.operation-title {
  font-weight: 600;
  color: hsl(var(--foreground));
}
</style>
