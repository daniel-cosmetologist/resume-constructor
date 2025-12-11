<template>
  <div class="section-root">
    <div class="header-row">
      <button class="btn-small" type="button" @click="addEducation">
        Add education
      </button>
      <span class="hint">Add your education history.</span>
    </div>

    <div
      v-for="(ed, index) in education"
      :key="index"
      class="card"
    >
      <div class="card-header">
        <h3 class="card-title">Entry {{ index + 1 }}</h3>
        <button class="icon-button" type="button" @click="removeEducation(index)">
          Remove
        </button>
      </div>

      <div class="grid-2">
        <label class="field">
          <span class="label">Institution</span>
          <input
            class="input"
            type="text"
            :value="ed.institution"
            @input="updateField(index, 'institution', $event)"
          />
        </label>

        <label class="field">
          <span class="label">Degree</span>
          <input
            class="input"
            type="text"
            :value="ed.degree"
            @input="updateField(index, 'degree', $event)"
          />
        </label>
      </div>

      <div class="grid-2">
        <label class="field">
          <span class="label">Location</span>
          <input
            class="input"
            type="text"
            :value="ed.location"
            @input="updateField(index, 'location', $event)"
          />
        </label>

        <div class="grid-2-inner">
          <label class="field">
            <span class="label">Start date (YYYY-MM)</span>
            <input
              class="input"
              type="text"
              :value="ed.startDate"
              @input="updateField(index, 'startDate', $event)"
            />
          </label>
          <label class="field">
            <span class="label">End date (YYYY-MM)</span>
            <input
              class="input"
              type="text"
              :value="ed.endDate"
              @input="updateField(index, 'endDate', $event)"
            />
          </label>
        </div>
      </div>

      <label class="field">
        <span class="label">Details</span>
        <textarea
          class="textarea"
          rows="2"
          :value="ed.details"
          @input="updateField(index, 'details', $event)"
        />
      </label>
    </div>

    <p v-if="education.length === 0" class="empty-text">
      No education entries yet. Click "Add education" to create one.
    </p>
  </div>
</template>

<script setup lang="ts">
import type { EducationEntry } from '@/types/resume';

const props = defineProps<{
  education: EducationEntry[];
}>();

const emit = defineEmits<{
  (e: 'update:education', value: EducationEntry[]): void;
}>();

function addEducation() {
  const next: EducationEntry = {
    institution: '',
    degree: '',
    location: '',
    startDate: '',
    endDate: '',
    details: ''
  };
  emit('update:education', [...props.education, next]);
}

function removeEducation(index: number) {
  const items = props.education.filter((_, i) => i !== index);
  emit('update:education', items);
}

function updateField(
  index: number,
  key: keyof EducationEntry,
  event: Event
) {
  const target = event.target as HTMLInputElement | HTMLTextAreaElement;
  const items = props.education.map((ed, i) =>
    i === index ? { ...ed, [key]: target.value } : ed
  );
  emit('update:education', items);
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

.grid-2 {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
  gap: 0.6rem;
}

.grid-2-inner {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
  gap: 0.6rem;
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

.input,
.textarea {
  border-radius: 0.5rem;
  border: 1px solid #d1d5db;
  padding: 0.45rem 0.6rem;
  font-size: 0.9rem;
}

.input:focus,
.textarea:focus {
  outline: none;
  border-color: #2563eb;
  box-shadow: 0 0 0 1px #2563eb22;
}

.textarea {
  resize: vertical;
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
