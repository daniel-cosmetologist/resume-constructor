<template>
  <div class="form-grid">
    <label class="field">
      <span class="label">Full name</span>
      <input
        class="input"
        type="text"
        :value="fullName"
        @input="onFullNameChange($event)"
        autocomplete="off"
      />
    </label>

    <label class="field">
      <span class="label">Position</span>
      <input
        class="input"
        type="text"
        :value="position"
        @input="onPositionChange($event)"
        autocomplete="off"
      />
    </label>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  fullName: string;
  position: string;
}>();

const emit = defineEmits<{
  (e: 'update:fullName', value: string): void;
  (e: 'update:position', value: string): void;
}>();

function onFullNameChange(event: Event) {
  const target = event.target as HTMLInputElement;
  emit('update:fullName', target.value);
}

function onPositionChange(event: Event) {
  const target = event.target as HTMLInputElement;
  emit('update:position', target.value);
}
</script>

<style scoped>
.form-grid {
  display: grid;
  gap: 0.75rem;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.label {
  font-size: 0.85rem;
  color: #4b5563;
}

.input {
  border-radius: 0.5rem;
  border: 1px solid #d1d5db;
  padding: 0.45rem 0.6rem;
  font-size: 0.9rem;
}

.input:focus {
  outline: none;
  border-color: #2563eb;
  box-shadow: 0 0 0 1px #2563eb22;
}
</style>
