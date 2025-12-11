export interface Link {
  label: string;
  url: string;
}

export interface Contacts {
  email: string;
  phone: string;
  location: string;
  links: Link[];
}

export interface ExperienceEntry {
  company: string;
  position: string;
  location: string;
  startDate: string;
  endDate: string;
  description: string;
  bullets: string[];
}

export interface EducationEntry {
  institution: string;
  degree: string;
  location: string;
  startDate: string;
  endDate: string;
  details: string;
}

export interface CustomSection {
  title: string;
  bulletSymbol: string;
  items: string[];
}

export interface Photo {
  mimeType: string;
  data: string;
}

export interface ResumeRequest {
  fullName: string;
  position: string;
  summary: string;
  contacts: Contacts;
  skills: string[];
  experience: ExperienceEntry[];
  education: EducationEntry[];
  customSections: CustomSection[];
  photo: Photo | null;
}

export function createEmptyResume(): ResumeRequest {
  return {
    fullName: '',
    position: '',
    summary: '',
    contacts: {
      email: '',
      phone: '',
      location: '',
      links: []
    },
    skills: [],
    experience: [],
    education: [],
    customSections: [],
    photo: null
  };
}
