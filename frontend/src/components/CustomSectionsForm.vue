<template>
  <div class="section-root">
    <div class="header-row">
      <button class="btn-small" type="button" @click="addSection">
        Add custom section
      </button>
      <span class="hint">Use this for homelab, open source, side projects and more.</span>
    </div>

    <div
      v-for="(section, index) in customSections"
      :key="index"
      class="card"
    >
      <div class="card-header">
        <h3 class="card-title">Section {{ index + 1 }}</h3>
        <button class="icon-button" type="button" @click="removeSection(index)">
          Remove
        </button>
      </div>

      <label class="field">
        <span class="label">Title</span>
        <input
          class="input"
          type="text"
          :value="section.title"
          @input="updateField(index, 'title', $event)"
        />
      </label>

      <label class="field">
        <span class="label">Bullet symbol</span>
        <input
          class="input"
          type="text"
          maxlength="3"
          :value="section.bulletSymbol"
          @input="updateField(index, 'bulletSymbol', $event)"
        />
      </label>

      <div class="bullets-header">
        <span class="label">Items</span>
        <button class="btn-small" type="button" @click="addItem(index)">
          Add item
        </button>
      </div>

      <div
        v-for="(item, itemIndex) in section.items"
        :key="itemIndex"
        class="bullet-row"
      >
        <input
          class="input"
          type="text"
          :value="item"
          @input="updateItem(index, itemIndex, $event)"
        />
        <button class="icon-button" type="button" @click="removeItem(index, itemIndex)">
          ✕
        </button>
      </div>

      <p v-if="section.items.length === 0" class="empty-text">
        This section has no items yet. Click "Add item" to create one.
      </p>
    </div>

    <p v-if="customSections.length === 0" class="empty-text">
      No custom sections yet. Click "Add custom section" to create one.
    </p>
  </div>
</template>

<script setup lang="ts">
import type { CustomSection } from '@/types/resume';

const props = defineProps<{
  customSections: CustomSection[];
}>();

const emit = defineEmits<{
  (e: 'update:customSections', value: CustomSection[]): void;
}>();

function addSection() {
  const next: CustomSection = {
    title: '',
    bulletSymbol: '•',
    items: ['']
  };
  emit('update:customSections', [...props.customSections, next]);
}

function removeSection(index: number) {
  const sections = props.customSections.filter((_, i) => i !== index);
  emit('update:customSections', sections);
}

function updateField(
  index: number,
  key: keyof CustomSection,
  event: Event
) {
  const target = event.target as HTMLInputElement;
  const sections = props.customSections.map((cs, i) =>
    i === index ? { ...cs, [key]: target.value } : cs
  );
  emit('update:customSections', sections);
}

function addItem(index: number) {
  const sections = props.customSections.map((cs, i) =>
    i === index ? { ...cs, items: [...cs.items, ''] } : cs
  );
  emit('update:customSections', sections);
}

function updateItem(
  index: number,
  itemIndex: number,
  event: Event
) {
  const target = event.target as HTMLInputElement;
  const sections = props.customSections.map((cs, i) => {
    if (i !== index) return cs;
    const items = cs.items.map((item, j) =>
      j === itemIndex ? target.value : item
    );
    return { ...cs, items };
  });
  emit('update:customSections', sections);
}

function removeItem(index: number, itemIndex: number) {
  const sections = props.customSections.map((cs, i) => {
    if (i !== index) return cs;
    const items = cs.items.filter((_, j) => j !== itemIndex);
    return { ...cs, items };
  });
  emit('update:customSections', sections);
}
</script>

<style scoped>
.section-root {
  display: grid;
  gap: 0.75rem;
}

.header-row {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.hint {
  font-size: 0.8rem;
  color: #6b7280;
}

.card {
  border-radius: 0.75rem;
  border: 1px solid #e5e7eb;
  padding: 0.75rem 0.9rem;
  background: #f9fafb;
  display: grid;
  gap: 0.6rem;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  margin: 0;
  font-size: 0.95rem;
  color: #111827;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.label {
  font-size: 0.8rem;
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

.bullets-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.bullet-row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 0.5rem;
  align-items: center;
}

.btn-small {
  border-radius: 0.5rem;
  border: 1px solid #d1d5db;
  padding: 0.2rem 0.5rem;
  font-size: 0.8rem;
  background: #ffffff;
  cursor: pointer;
}

.btn-small:hover {
  background: #f3f4f6;
}

.icon-button {
  border: none;
  background: transparent;
  cursor: pointer;
  font-size: 0.9rem;
  padding: 0.2rem;
  color: #9ca3af;
}

.icon-button:hover {
  color: #ef4444;
}

.empty-text {
  font-size: 0.85rem;
  color: #9ca3af;
}
</style>
