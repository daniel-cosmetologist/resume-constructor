<template>
  <div class="skills-root">
    <div class="skills-header">
      <button class="btn-small" type="button" @click="addSkill">Add skill</button>
      <span class="hint">One skill per line.</span>
    </div>

    <div v-for="(skill, index) in skills" :key="index" class="skill-row">
      <input
        class="input"
        type="text"
        :value="skill"
        @input="updateSkill(index, $event)"
        autocomplete="off"
      />
      <button class="icon-button" type="button" @click="removeSkill(index)">
        ✕
      </button>
    </div>

    <p v-if="skills.length === 0" class="empty-text">
      No skills yet. Click "Add skill" to start.
    </p>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  skills: string[];
}>();

const emit = defineEmits<{
  (e: 'update:skills', value: string[]): void;
}>();

function addSkill() {
  emit('update:skills', [...props.skills, '']);
}

function updateSkill(index: number, event: Event) {
  const target = event.target as HTMLInputElement;
  const skills = props.skills.map((s, i) => (i === index ? target.value : s));
  emit('update:skills', skills);
}

function removeSkill(index: number) {
  const skills = props.skills.filter((_, i) => i !== index);
  emit('update:skills', skills);
}
</script>

<style scoped>
.skills-root {
  display: grid;
  gap: 0.5rem;
}

.skills-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.hint {
  font-size: 0.8rem;
  color: #6b7280;
}

.skill-row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 0.5rem;
  align-items: center;
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
