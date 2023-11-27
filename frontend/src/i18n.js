import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import enTranslation from './locales/en.json';
import frTranslation from './locales/fr.json';

i18n.use(initReactI18next).init({
  resources: {
    english: { translation: enTranslation },
    french: { translation: frTranslation }
  },
  lng: 'english', // Set the default language
  fallbackLng: 'english', // Fallback language
  interpolation: {
    escapeValue: false // React already does escaping
  }
});

export default i18n;
