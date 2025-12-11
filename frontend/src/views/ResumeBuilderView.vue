<template>
  <div class="builder-layout">
    <section class="builder-left">
      <div class="builder-section">
        <h2>Personal information</h2>
        <PersonalInfoForm
          v-model:fullName="resume.fullName"
          v-model:position="resume.position"
        />
      </div>

      <div class="builder-section">
        <h2>Summary</h2>
        <SummaryForm v-model:summary="resume.summary" />
      </div>

      <div class="builder-section">
        <h2>Contacts</h2>
        <ContactsForm v-model:contacts="resume.contacts" />
      </div>

      <div class="builder-section">
        <h2>Skills</h2>
        <SkillsForm v-model:skills="resume.skills" />
      </div>

      <div class="builder-section">
        <h2>Experience</h2>
        <ExperienceForm v-model:experience="resume.experience" />
      </div>

      <div class="builder-section">
        <h2>Education</h2>
        <EducationForm v-model:education="resume.education" />
      </div>

      <div class="builder-section">
        <h2>Custom sections</h2>
        <CustomSectionsForm v-model:customSections="resume.customSections" />
      </div>

      <div class="builder-section">
        <h2>Photo</h2>
        <PhotoUpload v-model:photo="resume.photo" />
      </div>

      <div class="builder-actions">
        <button class="btn btn-primary" type="button" @click="downloadPdf" :disabled="isPreviewLoading">
          Download PDF
        </button>
        <button class="btn" type="button" @click="updatePreview" :disabled="isPreviewLoading">
          Regenerate preview
        </button>
        <span v-if="isPreviewLoading" class="status-text">Generating PDF...</span>
        <span v-else-if="lastUpdated" class="status-text">Last updated: {{ lastUpdated }}</span>
      </div>
    </section>

    <section class="builder-right">
      <h2 class="preview-title">PDF Preview</h2>
      <PdfPreview :pdf-url="pdfUrl" :loading="isPreviewLoading" :error="errorMessage" />
    </section>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, watch, onBeforeUnmount } from 'vue';
import PersonalInfoForm from '@/components/PersonalInfoForm.vue';
import SummaryForm from '@/components/SummaryForm.vue';
import ContactsForm from '@/components/ContactsForm.vue';
import SkillsForm from '@/components/SkillsForm.vue';
import ExperienceForm from '@/components/ExperienceForm.vue';
import EducationForm from '@/components/EducationForm.vue';
import CustomSectionsForm from '@/components/CustomSectionsForm.vue';
import PhotoUpload from '@/components/PhotoUpload.vue';
import PdfPreview from '@/components/PdfPreview.vue';
import { createEmptyResume, type ResumeRequest } from '@/types/resume';
import { generateResumePdf } from '@/api/resumeApi';

const resume = reactive<ResumeRequest>(createEmptyResume());

const pdfUrl = ref<string | null>(null);
const isPreviewLoading = ref(false);
const errorMessage = ref<string | null>(null);
const lastUpdated = ref<string | null>(null);

let previewTimeoutId: number | undefined;

async function updatePreview() {
  if (isPreviewLoading.value) {
    // Дадим закончить текущему запросу, чтобы не плодить лишние
  }

  isPreviewLoading.value = true;
  errorMessage.value = null;

  try {
    const blob = await generateResumePdf(resume as ResumeRequest);
    if (pdfUrl.value) {
      URL.revokeObjectURL(pdfUrl.value);
    }
    pdfUrl.value = URL.createObjectURL(blob);
    lastUpdated.value = new Date().toLocaleTimeString();
  } catch (err) {
    const message = err instanceof Error ? err.message : 'Failed to generate PDF';
    errorMessage.value = message;
  } finally {
    isPreviewLoading.value = false;
  }
}

async function downloadPdf() {
  isPreviewLoading.value = true;
  errorMessage.value = null;

  try {
    const blob = await generateResumePdf(resume as ResumeRequest);
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'resume.pdf';
    document.body.appendChild(a);
    a.click();
    a.remove();
    URL.revokeObjectURL(url);
    lastUpdated.value = new Date().toLocaleTimeString();
  } catch (err) {
    const message = err instanceof Error ? err.message : 'Failed to download PDF';
    errorMessage.value = message;
  } finally {
    isPreviewLoading.value = false;
  }
}

function schedulePreview() {
  if (previewTimeoutId !== undefined) {
    window.clearTimeout(previewTimeoutId);
  }
  previewTimeoutId = window.setTimeout(() => {
    updatePreview();
  }, 1000);
}

watch(
  () => resume,
  () => {
    schedulePreview();
  },
  { deep: true }
);

onBeforeUnmount(() => {
  if (previewTimeoutId !== undefined) {
    window.clearTimeout(previewTimeoutId);
  }
  if (pdfUrl.value) {
    URL.revokeObjectURL(pdfUrl.value);
  }
});
</script>

<style scoped>
.builder-layout {
  display: grid;
  grid-template-columns: minmax(0, 2.1fr) minmax(0, 2fr);
  gap: 1.5rem;
  align-items: flex-start;
}

.builder-left,
.builder-right {
  background: #ffffff;
  border-radius: 0.75rem;
  padding: 1.25rem 1.5rem;
  box-shadow: 0 2px 6px rgba(15, 23, 42, 0.08);
}

.builder-section + .builder-section {
  margin-top: 1.25rem;
}

.builder-section h2 {
  margin: 0 0 0.5rem;
  font-size: 1rem;
  color: #111827;
}

.builder-actions {
  margin-top: 1.25rem;
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 0.75rem;
}

.preview-title {
  margin: 0 0 0.75rem;
  font-size: 1rem;
  color: #111827;
}

.btn {
  border-radius: 0.5rem;
  border: 1px solid #d1d5db;
  padding: 0.45rem 0.9rem;
  font-size: 0.9rem;
  background: #ffffff;
  cursor: pointer;
  transition: background 0.1s ease, border-color 0.1s ease;
}

.btn:hover {
  background: #f3f4f6;
}

.btn-primary {
  border-color: #2563eb;
  background: #2563eb;
  color: #ffffff;
}

.btn-primary:hover {
  background: #1d4ed8;
}

.btn:disabled {
  opacity: 0.6;
  cursor: default;
}

.status-text {
  font-size: 0.85rem;
  color: #4b5563;
}

@media (max-width: 1024px) {
  .builder-layout {
    grid-template-columns: minmax(0, 1fr);
  }
}
</style>
