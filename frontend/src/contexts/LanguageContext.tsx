import { createContext, useContext, useState, useEffect } from 'react';
import type { ReactNode } from 'react';

export type Language = 'uk' | 'en';

interface LanguageContextType {
  language: Language;
  setLanguage: (lang: Language) => void;
  t: (key: string) => string;
}

const LanguageContext = createContext<LanguageContextType | undefined>(undefined);

// Функції для роботи з localStorage
const getStoredLanguage = (): Language => {
  try {
    const stored = localStorage.getItem('mcd_language');
    return (stored as Language) || 'uk';
  } catch {
    return 'uk';
  }
};

const setStoredLanguage = (language: Language): void => {
  try {
    localStorage.setItem('mcd_language', language);
  } catch {
    // Ігноруємо помилки localStorage
  }
};

// Тексти для перекладів
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
    'footer.navigation': 'Навігація',
    'footer.support': 'Підтримка',
    'footer.home': 'Головна',
    'footer.documents': 'Документи',
    'footer.groups': 'Групи',
    'footer.settings': 'Налаштування',
    'footer.help': 'Допомога',
    'footer.contact': 'Контакти',
    'footer.privacy': 'Конфіденційність',
    'footer.copyright': 'Всі права захищені.',
    
    // 404 page
    'notFound.title': 'Сторінку не знайдено',
    'notFound.message': 'Вибачте, але сторінка, яку ви шукаєте, не існує або була переміщена.',
    'notFound.home': 'Повернутися на головну',
    'notFound.back': 'Назад',
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
    'footer.navigation': 'Navigation',
    'footer.support': 'Support',
    'footer.home': 'Home',
    'footer.documents': 'Documents',
    'footer.groups': 'Groups',
    'footer.settings': 'Settings',
    'footer.help': 'Help',
    'footer.contact': 'Contact',
    'footer.privacy': 'Privacy',
    'footer.copyright': 'All rights reserved.',
    
    // 404 page
    'notFound.title': 'Page Not Found',
    'notFound.message': 'Sorry, but the page you are looking for does not exist or has been moved.',
    'notFound.home': 'Return to Home',
    'notFound.back': 'Back',
  }
};

interface LanguageProviderProps {
  children: ReactNode;
}

export const LanguageProvider = ({ children }: LanguageProviderProps) => {
  const [language, setLanguage] = useState<Language>(getStoredLanguage);

  // Зберігаємо мову в localStorage при зміні
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

export const useLanguage = () => {
  const context = useContext(LanguageContext);
  if (context === undefined) {
    throw new Error('useLanguage must be used within a LanguageProvider');
  }
  return context;
};
