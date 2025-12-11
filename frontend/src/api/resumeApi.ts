import type { ResumeRequest } from '@/types/resume';

const API_URL = '/api/v1/resume/pdf';

export async function generateResumePdf(data: ResumeRequest): Promise<Blob> {
  const response = await fetch(API_URL, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  });

  if (!response.ok) {
    let message = `Request failed with status ${response.status}`;
    try {
      const errorBody = (await response.json()) as { message?: string };
      if (errorBody && typeof errorBody.message === 'string' && errorBody.message.length > 0) {
        message = errorBody.message;
      }
    } catch {
      // тело ошибки не JSON, оставляем дефолтное сообщение
    }
    throw new Error(message);
  }

  const blob = await response.blob();
  return blob;
}
