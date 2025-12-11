<template>
  <div class="preview-root">
    <div v-if="error" class="alert alert-error">
      {{ error }}
    </div>

    <div v-else-if="loading" class="preview-placeholder">
      Generating PDF preview...
    </div>

    <div v-else-if="pdfUrl" class="preview-container">
      <object
        :data="pdfUrl"
        type="application/pdf"
        class="preview-object"
      >
        <p>
          PDF preview is not available in this browser.
          You can still download the file using the button on the left.
        </p>
      </object>
    </div>

    <div v-else class="preview-placeholder">
      Start filling the form to see a live PDF preview here.
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  pdfUrl: string | null;
  loading: boolean;
  error: string | null;
}>();
</script>

<style scoped>
.preview-root {
  min-height: 300px;
}

.preview-container {
  border-radius: 0.75rem;
  border: 1px solid #e5e7eb;
  overflow: hidden;
  background: #f9fafb;
  height: 600px;
  max-height: 80vh;
}

.preview-object {
  width: 100%;
  height: 100%;
  border: none;
}

.preview-placeholder {
  border-radius: 0.75rem;
  border: 1px dashed #d1d5db;
  padding: 1.25rem;
  text-align: center;
  font-size: 0.9rem;
  color: #6b7280;
  background: #f9fafb;
}

.alert {
  border-radius: 0.5rem;
  padding: 0.6rem 0.75rem;
  font-size: 0.85rem;
  margin-bottom: 0.6rem;
}

.alert-error {
  background: #fef2f2;
  border: 1px solid #fecaca;
  color: #b91c1c;
}
</style>
