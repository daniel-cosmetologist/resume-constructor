<template>
  <div class="app-root">
    <header class="app-header">
      <h1 class="app-title">Resume Builder</h1>
      <p class="app-subtitle">
        Fill the form, see live PDF preview, download your resume.
      </p>
    </header>

    <main class="app-main">
      <RouterView />
    </main>
  </div>
</template>

<script setup lang="ts">
import { RouterView } from 'vue-router';
</script>

<style scoped>
.app-root {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f5f5;
  color: #222;
  font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
}

.app-header {
  padding: 1.5rem 2rem;
  background: #111827;
  color: #f9fafb;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.app-title {
  margin: 0;
  font-size: 1.6rem;
}

.app-subtitle {
  margin: 0.3rem 0 0;
  opacity: 0.8;
  font-size: 0.95rem;
}

.app-main {
  flex: 1;
  padding: 1.5rem;
}
</style>
