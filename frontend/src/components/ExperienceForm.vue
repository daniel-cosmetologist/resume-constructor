<template>
  <div class="section-root">
    <div class="header-row">
      <button class="btn-small" type="button" @click="addExperience">
        Add experience
      </button>
      <span class="hint">Add your work experience entries.</span>
    </div>

    <div
      v-for="(exp, index) in experience"
      :key="index"
      class="card"
    >
      <div class="card-header">
        <h3 class="card-title">Position {{ index + 1 }}</h3>
        <button class="icon-button" type="button" @click="removeExperience(index)">
          Remove
        </button>
      </div>

      <div class="grid-2">
        <label class="field">
          <span class="label">Company</span>
          <input
            class="input"
            type="text"
            :value="exp.company"
            @input="updateField(index, 'company', $event)"
          />
        </label>

        <label class="field">
          <span class="label">Position</span>
          <input
            class="input"
            type="text"
            :value="exp.position"
            @input="updateField(index, 'position', $event)"
          />
        </label>
      </div>

      <div class="grid-2">
        <label class="field">
          <span class="label">Location</span>
          <input
            class="input"
            type="text"
            :value="exp.location"
            @input="updateField(index, 'location', $event)"
          />
        </label>

        <div class="grid-2-inner">
          <label class="field">
            <span class="label">Start date (YYYY-MM)</span>
            <input
              class="input"
              type="text"
              :value="exp.startDate"
              @input="updateField(index, 'startDate', $event)"
            />
          </label>
          <label class="field">
            <span class="label">End date (YYYY-MM or empty)</span>
            <input
              class="input"
              type="text"
              :value="exp.endDate"
              @input="updateField(index, 'endDate', $event)"
            />
          </label>
        </div>
      </div>

      <label class="field">
        <span class="label">Description</span>
        <textarea
          class="textarea"
          rows="3"
          :value="exp.description"
          @input="updateField(index, 'description', $event)"
        />
      </label>

      <div class="bullets-header">
        <span class="label">Bullet points</span>
        <button class="btn-small" type="button" @click="addBullet(index)">
          Add bullet
        </button>
      </div>

      <div
        v-for="(bullet, bulletIndex) in exp.bullets"
        :key="bulletIndex"
        class="bullet-row"
      >
        <input
          class="input"
          type="text"
          :value="bullet"
          @input="updateBullet(index, bulletIndex, $event)"
        />
        <button class="icon-button" type="button" @click="removeBullet(index, bulletIndex)">
          ✕
        </button>
      </div>
    </div>

    <p v-if="experience.length === 0" class="empty-text">
      No experience entries yet. Click "Add experience" to create one.
    </p>
  </div>
</template>

<script setup lang="ts">
import type { ExperienceEntry } from '@/types/resume';

const props = defineProps<{
  experience: ExperienceEntry[];
}>();

const emit = defineEmits<{
  (e: 'update:experience', value: ExperienceEntry[]): void;
}>();

function addExperience() {
  const next: ExperienceEntry = {
    company: '',
    position: '',
    location: '',
    startDate: '',
    endDate: '',
    description: '',
    bullets: []
  };
  emit('update:experience', [...props.experience, next]);
}

function removeExperience(index: number) {
  const items = props.experience.filter((_, i) => i !== index);
  emit('update:experience', items);
}

function updateField(
  index: number,
  key: keyof ExperienceEntry,
  event: Event
) {
  const target = event.target as HTMLInputElement | HTMLTextAreaElement;
  const items = props.experience.map((exp, i) =>
    i === index ? { ...exp, [key]: target.value } : exp
  );
  emit('update:experience', items);
}

function addBullet(index: number) {
  const items = props.experience.map((exp, i) =>
    i === index ? { ...exp, bullets: [...exp.bullets, ''] } : exp
  );
  emit('update:experience', items);
}

function updateBullet(
  index: number,
  bulletIndex: number,
  event: Event
) {
  const target = event.target as HTMLInputElement;
  const items = props.experience.map((exp, i) => {
    if (i !== index) return exp;
    const bullets = exp.bullets.map((b, j) =>
      j === bulletIndex ? target.value : b
    );
    return { ...exp, bullets };
  });
  emit('update:experience', items);
}

function removeBullet(index: number, bulletIndex: number) {
  const items = props.experience.map((exp, i) => {
    if (i !== index) return exp;
    const bullets = exp.bullets.filter((_, j) => j !== bulletIndex);
    return { ...exp, bullets };
  });
  emit('update:experience', items);
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
