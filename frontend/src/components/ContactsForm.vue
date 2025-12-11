<template>
  <div class="form-grid">
    <label class="field">
      <span class="label">Email</span>
      <input
        class="input"
        type="email"
        :value="localContacts.email"
        @input="updateField('email', $event)"
        autocomplete="off"
      />
    </label>

    <label class="field">
      <span class="label">Phone</span>
      <input
        class="input"
        type="tel"
        :value="localContacts.phone"
        @input="updateField('phone', $event)"
        autocomplete="off"
      />
    </label>

    <label class="field full-width">
      <span class="label">Location</span>
      <input
        class="input"
        type="text"
        :value="localContacts.location"
        @input="updateField('location', $event)"
        autocomplete="off"
      />
    </label>

    <div class="links-header">
      <span class="label">Links</span>
      <button class="btn-small" type="button" @click="addLink">
        Add link
      </button>
    </div>

    <div v-for="(link, index) in localContacts.links" :key="index" class="link-row">
      <input
        class="input"
        type="text"
        :value="link.label"
        @input="updateLink(index, 'label', $event)"
        placeholder="Label (GitHub, LinkedIn...)"
        autocomplete="off"
      />
      <input
        class="input"
        type="url"
        :value="link.url"
        @input="updateLink(index, 'url', $event)"
        placeholder="https://"
        autocomplete="off"
      />
      <button class="icon-button" type="button" @click="removeLink(index)">
        ✕
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { Contacts } from '@/types/resume';

const props = defineProps<{
  contacts: Contacts;
}>();

const emit = defineEmits<{
  (e: 'update:contacts', value: Contacts): void;
}>();

const localContacts = computed<Contacts>({
  get: () => props.contacts,
  set: (value) => emit('update:contacts', value)
});

function updateField(key: 'email' | 'phone' | 'location', event: Event) {
  const target = event.target as HTMLInputElement;
  localContacts.value = { ...localContacts.value, [key]: target.value };
}

function addLink() {
  const links = [...localContacts.value.links, { label: '', url: '' }];
  localContacts.value = { ...localContacts.value, links };
}

function updateLink(index: number, key: 'label' | 'url', event: Event) {
  const target = event.target as HTMLInputElement;
  const links = localContacts.value.links.map((link, i) =>
    i === index ? { ...link, [key]: target.value } : link
  );
  localContacts.value = { ...localContacts.value, links };
}

function removeLink(index: number) {
  const links = localContacts.value.links.filter((_, i) => i !== index);
  localContacts.value = { ...localContacts.value, links };
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

.full-width {
  grid-column: 1 / -1;
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

.links-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.link-row {
  display: grid;
  grid-template-columns: minmax(0, 1.2fr) minmax(0, 2fr) auto;
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
</style>
