import { createContext, useContext, useState, useEffect } from 'react';
import type { ReactNode } from 'react';

export type Language = 'uk' | 'en';

interface LanguageContextType {
  language: Language;
  setLanguage: (lang: Language) => void;
  t: (key: string) => string;
}

const LanguageContext = createContext<LanguageContextType | undefined>(undefined);

const getStoredLanguage = (): Language => {
  try {
    const stored = localStorage.getItem('mcd_language');
    return (stored as Language) || 'uk';
  } catch (error) {
    console.warn('Не вдалося завантажити мову з localStorage', error);
    return 'uk';
  }
};

const setStoredLanguage = (language: Language): void => {
  try {
    localStorage.setItem('mcd_language', language);
  } catch (error) {
    console.warn('Не вдалося зберегти мову до localStorage', error);
  }
};

const translations = {
  uk: {
    // Header
    'header.settings': 'Налаштування акаунту',
    'header.toggleTheme': 'Перемкнути тему',

    // Home page
    'home.createDocument': 'Створити документ',
    'home.createGroup': 'Створити групу',
    'home.createDocumentAlert': 'Створити документ - функція в розробці',
    'home.createGroupAlert': 'Створити групу - функція в розробці',

    // Settings page
    'settings.title': 'Налаштування акаунту',
    'settings.subtitle': 'Керуйте своїми налаштуваннями та преференціями',
    'settings.general': 'Загальні налаштування',
    'settings.theme': 'Тема',
    'settings.language': 'Мова',
    'settings.notifications': 'Сповіщення',
    'settings.documents': 'Робота з документами',
    'settings.notificationsLabel': 'Отримувати сповіщення',
    'settings.autoSaveLabel': 'Автоматичне збереження',
    'settings.save': 'Зберегти зміни',
    'settings.reset': 'Скинути до стандартних',
    'settings.theme.light': 'Світла',
    'settings.theme.dark': 'Темна',

    // Footer
    'footer.description': 'Ваш надійний помічник для управління документами',
    'footer.navigation': 'Навігація',
    'footer.home': 'Головна',
    'footer.documents': 'Документи',
    'footer.groups': 'Групи',
    'footer.settings': 'Налаштування',
    'footer.copyright': 'Всі права захищені.',

    // Sidebar
    'sidebar.navigation': 'Навігація',
    'sidebar.home': 'Головна',
    'sidebar.documents': 'Документи',
    'sidebar.settings': 'Налаштування',
    'sidebar.groups': 'Групи',
    'sidebar.viewGroups': 'Переглянути групи',
    'sidebar.createGroup': 'Створити групу',

    // 404 page
    'notFound.title': 'Сторінку не знайдено',
    'notFound.message': 'Вибачте, але сторінка, яку ви шукаєте, не існує або була переміщена.',
    'notFound.home': 'Повернутися на головну',
    'notFound.back': 'Назад',

    // Groups page
    'groups.title': 'Групи документів',
    'groups.subtitle': 'Переглядайте створені групи та переходьте до їх документів',
    'groups.error': 'Не вдалося завантажити групи.',
    'groups.refresh': 'Спробувати знову',
    'groups.empty': 'Поки немає жодної групи.',
    'groups.createdAt': 'Створено',
    'groups.totalLabel': 'Всього груп:',

    // Documents page
    'documents.title': 'Документи',
    'documents.subtitle': 'Ознайомтеся з документами та знаходьте потрібні за кілька секунд',
    'documents.filterGroup': 'Група',
    'documents.filterAll': 'Усі групи',
    'documents.searchPlaceholder': 'Пошук документів',
    'documents.error': 'Не вдалося завантажити документи.',
    'documents.refresh': 'Спробувати знову',
    'documents.empty': 'Документів не знайдено.',
    'documents.groupUnknown': 'Без групи',
    'documents.noContent': 'Без вмісту',
    'documents.createdAt': 'Створено',

    // Document editor page
    'documentEditor.subtitle': 'Редагуйте документ та переглядайте зміни в реальному часі',
    'documentEditor.fallbackTitle': 'Документ',
    'documentEditor.nameLabel': 'Назва документа',
    'documentEditor.namePlaceholder': 'Введіть назву',
    'documentEditor.nameRequired': 'Назва є обовʼязковою',
    'documentEditor.contentLabel': 'Вміст документа',
    'documentEditor.contentPlaceholder': 'Напишіть текст у форматі Markdown...',
    'documentEditor.previewTitle': 'Попередній перегляд',
    'documentEditor.previewEmpty': 'Почніть вводити текст, щоб побачити попередній перегляд.',
    'documentEditor.createdAtLabel': 'Створено',
    'documentEditor.saveButton': 'Зберегти зміни',
    'documentEditor.savingButton': 'Збереження...',
    'documentEditor.saveSuccess': 'Документ успішно збережено.',
    'documentEditor.saveError': 'Не вдалося зберегти документ. Спробуйте ще раз.',
    'documentEditor.backToList': 'Повернутися до документів',
    'documentEditor.loadError': 'Не вдалося завантажити документ.',
    'documentEditor.notFound': 'Документ не знайдено або він був видалений.',
  },
  en: {
    // Header
    'header.settings': 'Account Settings',
    'header.toggleTheme': 'Toggle Theme',

    // Home page
    'home.createDocument': 'Create Document',
    'home.createGroup': 'Create Group',
    'home.createDocumentAlert': 'Create Document - feature in development',
    'home.createGroupAlert': 'Create Group - feature in development',

    // Settings page
    'settings.title': 'Account Settings',
    'settings.subtitle': 'Manage your settings and preferences',
    'settings.general': 'General Settings',
    'settings.theme': 'Theme',
    'settings.language': 'Language',
    'settings.notifications': 'Notifications',
    'settings.documents': 'Document Work',
    'settings.notificationsLabel': 'Receive notifications',
    'settings.autoSaveLabel': 'Auto-save',
    'settings.save': 'Save Changes',
    'settings.reset': 'Reset to Default',
    'settings.theme.light': 'Light',
    'settings.theme.dark': 'Dark',

    // Footer
    'footer.description': 'Your reliable assistant for document management',
    'footer.navigation': 'Navigation',
    'footer.home': 'Home',
    'footer.documents': 'Documents',
    'footer.groups': 'Groups',
    'footer.settings': 'Settings',
    'footer.copyright': 'All rights reserved.',

    // Sidebar
    'sidebar.navigation': 'Navigation',
    'sidebar.home': 'Home',
    'sidebar.documents': 'Documents',
    'sidebar.settings': 'Settings',
    'sidebar.groups': 'Groups',
    'sidebar.viewGroups': 'View Groups',
    'sidebar.createGroup': 'Create Group',

    // 404 page
    'notFound.title': 'Page Not Found',
    'notFound.message': 'Sorry, but the page you are looking for does not exist or has been moved.',
    'notFound.home': 'Return to Home',
    'notFound.back': 'Back',

    // Groups page
    'groups.title': 'Document Groups',
    'groups.subtitle': 'Browse created groups and navigate to their documents',
    'groups.error': 'Failed to load groups.',
    'groups.refresh': 'Try again',
    'groups.empty': 'No groups yet.',
    'groups.createdAt': 'Created',
    'groups.totalLabel': 'Total groups:',

    // Documents page
    'documents.title': 'Documents',
    'documents.subtitle': 'Browse documents and find what you need in seconds',
    'documents.filterGroup': 'Group',
    'documents.filterAll': 'All groups',
    'documents.searchPlaceholder': 'Search documents',
    'documents.error': 'Failed to load documents.',
    'documents.refresh': 'Try again',
    'documents.empty': 'No documents found.',
    'documents.groupUnknown': 'No group',
    'documents.noContent': 'No content',
    'documents.createdAt': 'Created',

    // Document editor page
    'documentEditor.subtitle': 'Edit the document and preview changes in real time',
    'documentEditor.fallbackTitle': 'Document',
    'documentEditor.nameLabel': 'Document name',
    'documentEditor.namePlaceholder': 'Enter a name',
    'documentEditor.nameRequired': 'Name is required',
    'documentEditor.contentLabel': 'Document content',
    'documentEditor.contentPlaceholder': 'Write your Markdown content here...',
    'documentEditor.previewTitle': 'Live preview',
    'documentEditor.previewEmpty': 'Start typing to see a formatted preview.',
    'documentEditor.createdAtLabel': 'Created',
    'documentEditor.saveButton': 'Save changes',
    'documentEditor.savingButton': 'Saving...',
    'documentEditor.saveSuccess': 'Document saved successfully.',
    'documentEditor.saveError': 'We could not save the document. Please try again.',
    'documentEditor.backToList': 'Back to documents',
    'documentEditor.loadError': 'Failed to load the document.',
    'documentEditor.notFound': 'The document was not found or has been removed.',
  }
};

interface LanguageProviderProps {
  children: ReactNode;
}

export const LanguageProvider = ({ children }: LanguageProviderProps) => {
  const [language, setLanguage] = useState<Language>(getStoredLanguage);

  useEffect(() => {
    setStoredLanguage(language);
  }, [language]);

  const handleSetLanguage = (lang: Language) => {
    setLanguage(lang);
  };

  const t = (key: string): string => {
    return translations[language][key as keyof typeof translations[typeof language]] || key;
  };

  return (
    <LanguageContext.Provider value={{ language, setLanguage: handleSetLanguage, t }}>
      {children}
    </LanguageContext.Provider>
  );
};

// eslint-disable-next-line react-refresh/only-export-components
export function useLanguage() {
  const context = useContext(LanguageContext);
  if (context === undefined) {
    throw new Error('useLanguage must be used within a LanguageProvider');
  }
  return context;
}
