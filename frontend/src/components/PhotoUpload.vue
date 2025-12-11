<template>
  <div class="photo-root">
    <div class="photo-row">
      <input
        ref="fileInput"
        type="file"
        class="file-input"
        accept="image/*"
        @change="onFileChange"
      />
      <button class="btn-small" type="button" @click="triggerSelect">
        Choose photo
      </button>
      <button
        v-if="photo"
        class="btn-small btn-secondary"
        type="button"
        @click="clearPhoto"
      >
        Remove photo
      </button>
    </div>

    <p class="hint">
      Any image will be automatically resized and cropped to 3×4, up to 2&nbsp;MB.
    </p>

    <div v-if="photoPreviewUrl" class="preview">
      <img :src="photoPreviewUrl" alt="Photo preview" class="preview-image" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch, onBeforeUnmount } from 'vue';
import type { Photo } from '@/types/resume';

const props = defineProps<{
  photo: Photo | null;
}>();

const emit = defineEmits<{
  (e: 'update:photo', value: Photo | null): void;
}>();

const fileInput = ref<HTMLInputElement | null>(null);
const photoPreviewUrl = ref<string | null>(null);

const localPhoto = computed({
  get: () => props.photo,
  set: (value: Photo | null) => emit('update:photo', value)
});

function triggerSelect() {
  fileInput.value?.click();
}

function clearPhoto() {
  localPhoto.value = null;
  if (photoPreviewUrl.value) {
    URL.revokeObjectURL(photoPreviewUrl.value);
    photoPreviewUrl.value = null;
  }
  if (fileInput.value) {
    fileInput.value.value = '';
  }
}

function onFileChange(event: Event) {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  if (!file) {
    return;
  }

  const reader = new FileReader();
  reader.onload = () => {
    const result = reader.result;
    if (typeof result !== 'string') {
      return;
    }

    const commaIndex = result.indexOf(',');
    if (commaIndex === -1) {
      return;
    }

    const meta = result.substring(0, commaIndex);
    const data = result.substring(commaIndex + 1);

    let mimeType = file.type;
    const start = meta.indexOf(':');
    const end = meta.indexOf(';');
    if (start !== -1 && end !== -1 && end > start) {
      mimeType = meta.substring(start + 1, end);
    }

    localPhoto.value = {
      mimeType,
      data
    };

    if (photoPreviewUrl.value) {
      URL.revokeObjectURL(photoPreviewUrl.value);
    }
    photoPreviewUrl.value = URL.createObjectURL(file);
  };
  reader.readAsDataURL(file);
}

watch(
  () => props.photo,
  (newPhoto) => {
    if (!newPhoto) {
      if (photoPreviewUrl.value) {
        URL.revokeObjectURL(photoPreviewUrl.value);
        photoPreviewUrl.value = null;
      }
    }
  }
);

onBeforeUnmount(() => {
  if (photoPreviewUrl.value) {
    URL.revokeObjectURL(photoPreviewUrl.value);
  }
});
</script>

<style scoped>
.photo-root {
  display: grid;
  gap: 0.5rem;
}

.photo-row {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  align-items: center;
}

.file-input {
  display: none;
}

.btn-small {
  border-radius: 0.5rem;
  border: 1px solid #d1d5db;
  padding: 0.25rem 0.6rem;
  font-size: 0.8rem;
  background: #ffffff;
  cursor: pointer;
}

.btn-small:hover {
  background: #f3f4f6;
}

.btn-secondary {
  border-color: #ef4444;
  color: #ef4444;
}

.btn-secondary:hover {
  background: #fee2e2;
}

.hint {
  font-size: 0.8rem;
  color: #6b7280;
}

.preview {
  margin-top: 0.25rem;
}

.preview-image {
  width: 90px;
  height: 120px;
  object-fit: cover;
  border-radius: 0.5rem;
  border: 1px solid #d1d5db;
}
</style>
