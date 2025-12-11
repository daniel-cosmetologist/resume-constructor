<template>
  <div class="field">
    <span class="label">Summary</span>
    <textarea
      class="textarea"
      :value="summary"
      @input="onSummaryChange"
      rows="4"
    />
    <p class="hint">Short professional summary. It will appear at the top of the resume.</p>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  summary: string;
}>();

const emit = defineEmits<{
  (e: 'update:summary', value: string): void;
}>();

function onSummaryChange(event: Event) {
  const target = event.target as HTMLTextAreaElement;
  emit('update:summary', target.value);
}
</script>

<style scoped>
.field {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.label {
  font-size: 0.85rem;
  color: #4b5563;
}

.textarea {
  border-radius: 0.5rem;
  border: 1px solid #d1d5db;
  padding: 0.5rem 0.6rem;
  font-size: 0.9rem;
  resize: vertical;
}

.textarea:focus {
  outline: none;
  border-color: #2563eb;
  box-shadow: 0 0 0 1px #2563eb22;
}

.hint {
  font-size: 0.8rem;
  color: #6b7280;
}
</style>
