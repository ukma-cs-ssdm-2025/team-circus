import { createContext, useContext, useState, useEffect } from "react";
import type { ReactNode } from "react";

export type Language = "uk" | "en";

interface LanguageContextType {
  language: Language;
  setLanguage: (lang: Language) => void;
  t: (key: string) => string;
}

const LanguageContext = createContext<LanguageContextType | undefined>(
  undefined,
);

const getStoredLanguage = (): Language => {
  try {
    const stored = localStorage.getItem("mcd_language");
    return (stored as Language) || "uk";
  } catch (error) {
    console.warn("Не вдалося завантажити мову з localStorage", error);
    return "uk";
  }
};

const setStoredLanguage = (language: Language): void => {
  try {
    localStorage.setItem("mcd_language", language);
  } catch (error) {
    console.warn("Не вдалося зберегти мову до localStorage", error);
  }
};

const translations = {
  uk: {
    // Header
    "header.settings": "Налаштування акаунту",
    "header.toggleTheme": "Перемкнути тему",

    // Home page
    "home.createDocument": "Створити документ",
    "home.createGroup": "Створити групу",
    "home.createDocumentAlert": "Створити документ - функція в розробці",
    "home.createGroupAlert": "Створити групу - функція в розробці",

    // Settings page
    "settings.title": "Налаштування акаунту",
    "settings.subtitle": "Керуйте своїми налаштуваннями та преференціями",
    "settings.general": "Загальні налаштування",
    "settings.theme": "Тема",
    "settings.language": "Мова",
    "settings.notifications": "Сповіщення",
    "settings.documents": "Робота з документами",
    "settings.notificationsLabel": "Отримувати сповіщення",
    "settings.autoSaveLabel": "Автоматичне збереження",
    "settings.save": "Зберегти зміни",
    "settings.reset": "Скинути до стандартних",
    "settings.theme.light": "Світла",
    "settings.theme.dark": "Темна",

    // Footer
    "footer.description": "Ваш надійний помічник для управління документами",
    "footer.navigation": "Навігація",
    "footer.home": "Головна",
    "footer.documents": "Документи",
    "footer.groups": "Групи",
    "footer.settings": "Налаштування",
    "footer.copyright": "Всі права захищені.",

    // Sidebar
    "sidebar.navigation": "Навігація",
    "sidebar.home": "Головна",
    "sidebar.documents": "Документи",
    "sidebar.settings": "Налаштування",
    "sidebar.groups": "Групи",
    "sidebar.close": "Закрити меню",
    "sidebar.viewGroups": "Переглянути групи",
    "sidebar.createGroup": "Створити групу",

    // 404 page
    "notFound.title": "Сторінку не знайдено",
    "notFound.message":
      "Вибачте, але сторінка, яку ви шукаєте, не існує або була переміщена.",
    "notFound.home": "Повернутися на головну",
    "notFound.back": "Назад",

    // Groups page
    "groups.title": "Групи документів",
    "groups.subtitle":
      "Переглядайте створені групи та переходьте до їх документів",
    "groups.error": "Не вдалося завантажити групи.",
    "groups.refresh": "Спробувати знову",
    "groups.empty": "Поки немає жодної групи.",
    "groups.createdAt": "Створено",
    "groups.totalLabel": "Всього груп:",
    "groups.createButton": "Створити групу",
    "groups.createDialogTitle": "Створення групи",
    "groups.editDialogTitle": "Редагування групи",
    "groups.createConfirm": "Створити",
    "groups.updateConfirm": "Зберегти",
    "groups.cancel": "Скасувати",
    "groups.nameLabel": "Назва групи",
    "groups.namePlaceholder": "Введіть назву групи",
    "groups.nameHelper": "Назва є обовʼязковою",
    "groups.editLabel": "Редагувати",
    "groups.deleteLabel": "Видалити",
    "groups.createSuccess": "Групу успішно створено.",
    "groups.updateSuccess": "Групу успішно оновлено.",
    "groups.deleteSuccess": "Групу успішно видалено.",
    "groups.createError": "Не вдалося створити групу.",
    "groups.updateError": "Не вдалося оновити групу.",
    "groups.deleteError": "Не вдалося видалити групу.",
    "groups.deleteConfirmTitle": "Видалення групи",
    "groups.deleteConfirmDescription":
      'Ви впевнені, що хочете видалити групу "{name}"?',
    "groups.deleteConfirmAccept": "Видалити",
    "groups.manageMembersButton": "Керувати учасниками",
    "groups.membersTitle": 'Учасники групи "{name}"',
    "groups.membersTitleFallback": "Учасники групи",
    "groups.membersSubtitle":
      "Додайте авторів та оглядачів, щоб працювати разом.",
    "groups.membersEmpty": "У цій групі ще немає учасників.",
    "groups.membersAddUser": "Додати користувача",
    "groups.membersAddUserPlaceholder": "Оберіть користувача зі списку",
    "groups.membersRoleLabel": "Роль",
    "groups.membersAddButton": "Додати",
    "groups.membersAddValidation":
      "Оберіть користувача, якого хочете додати.",
    "groups.membersUsersError":
      "Не вдалося завантажити список користувачів.",
    "groups.membersActionError":
      "Не вдалося виконати дію. Спробуйте ще раз.",
    "groups.membersRemoveTooltip": "Видалити з групи",
    "groups.membersAddSuccess": "Користувача додано до групи.",
    "groups.membersAddError": "Не вдалося додати користувача.",
    "groups.membersUpdateSuccess": "Роль учасника оновлено.",
    "groups.membersUpdateError": "Не вдалося оновити роль учасника.",
    "groups.membersRemoveSuccess": "Учасника видалено з групи.",
    "groups.membersRemoveError": "Не вдалося видалити учасника.",
    "groups.membersForbidden":
      "Лише автор групи може керувати учасниками.",
    "groups.roleLabel": "Роль",
    "groups.role.author": "Автор",
    "groups.role.coauthor": "Співавтор",
    "groups.role.reviewer": "Оглядач",

    // Documents page
    "documents.title": "Документи",
    "documents.subtitle":
      "Ознайомтеся з документами та знаходьте потрібні за кілька секунд",
    "documents.filterGroup": "Група",
    "documents.filterAll": "Усі групи",
    "documents.searchPlaceholder": "Пошук документів",
    "documents.error": "Не вдалося завантажити документи.",
    "documents.refresh": "Спробувати знову",
    "documents.empty": "Документів не знайдено.",
    "documents.groupUnknown": "Без групи",
    "documents.noContent": "Без вмісту",
    "documents.createdAt": "Створено",
    "documents.createButton": "Створити документ",
    "documents.noGroupsHint": "Створіть групу, щоб додавати документи.",
    "documents.createDialogTitle": "Створення документа",
    "documents.createConfirm": "Створити",
    "documents.defaultContent": "# Новий документ",
    "documents.cancel": "Скасувати",
    "documents.nameLabel": "Назва документа",
    "documents.namePlaceholder": "Введіть назву документа",
    "documents.nameHelper": "Назва є обовʼязковою",
    "documents.groupLabel": "Група",
    "documents.contentLabel": "Вміст документа",
    "documents.contentPlaceholder": "Напишіть текст у форматі Markdown...",
    "documents.contentHelper": "Вміст є обовʼязковим",
    "documents.createSuccess": "Документ успішно створено.",
    "documents.createError": "Не вдалося створити документ.",
    "documents.editLabel": "Редагувати",
    "documents.deleteLabel": "Видалити",
    "documents.deleteConfirmTitle": "Видалення документа",
    "documents.deleteConfirmDescription":
      'Ви впевнені, що хочете видалити документ "{name}"? Це дію неможливо скасувати.',
    "documents.deleteConfirmAccept": "Видалити",
    "documents.deleteSuccess": "Документ успішно видалено.",
    "documents.deleteError": "Не вдалося видалити документ.",
    "documents.noEditableGroups":
      "Немає груп, де ви можете створювати документи.",
    "documents.deleteForbidden":
      "У вас немає прав видаляти цей документ.",

    // Document editor page
    "documentEditor.subtitle":
      "Редагуйте документ та переглядайте зміни в реальному часі",
    "documentEditor.fallbackTitle": "Документ",
    "documentEditor.nameLabel": "Назва документа",
    "documentEditor.namePlaceholder": "Введіть назву",
    "documentEditor.nameRequired": "Назва є обовʼязковою",
    "documentEditor.contentLabel": "Вміст документа",
    "documentEditor.contentPlaceholder": "Напишіть текст у форматі Markdown...",
    "documentEditor.previewTitle": "Попередній перегляд",
    "documentEditor.previewEmpty":
      "Почніть вводити текст, щоб побачити попередній перегляд.",
    "documentEditor.createdAtLabel": "Створено",
    "documentEditor.saveButton": "Зберегти зміни",
    "documentEditor.savingButton": "Збереження...",
    "documentEditor.saveSuccess": "Документ успішно збережено.",
    "documentEditor.saveError":
      "Не вдалося зберегти документ. Спробуйте ще раз.",
    "documentEditor.backToList": "Повернутися до документів",
    "documentEditor.loadError": "Не вдалося завантажити документ.",
    "documentEditor.notFound": "Документ не знайдено або він був видалений.",
    "documentEditor.readOnlyNotice":
      "Ви маєте права лише на перегляд цього документа.",
  },
  en: {
    // Header
    "header.settings": "Account Settings",
    "header.toggleTheme": "Toggle Theme",

    // Home page
    "home.createDocument": "Create Document",
    "home.createGroup": "Create Group",
    "home.createDocumentAlert": "Create Document - feature in development",
    "home.createGroupAlert": "Create Group - feature in development",

    // Settings page
    "settings.title": "Account Settings",
    "settings.subtitle": "Manage your settings and preferences",
    "settings.general": "General Settings",
    "settings.theme": "Theme",
    "settings.language": "Language",
    "settings.notifications": "Notifications",
    "settings.documents": "Document Work",
    "settings.notificationsLabel": "Receive notifications",
    "settings.autoSaveLabel": "Auto-save",
    "settings.save": "Save Changes",
    "settings.reset": "Reset to Default",
    "settings.theme.light": "Light",
    "settings.theme.dark": "Dark",

    // Footer
    "footer.description": "Your reliable assistant for document management",
    "footer.navigation": "Navigation",
    "footer.home": "Home",
    "footer.documents": "Documents",
    "footer.groups": "Groups",
    "footer.settings": "Settings",
    "footer.copyright": "All rights reserved.",

    // Sidebar
    "sidebar.navigation": "Navigation",
    "sidebar.home": "Home",
    "sidebar.documents": "Documents",
    "sidebar.settings": "Settings",
    "sidebar.groups": "Groups",
    "sidebar.close": "Close menu",
    "sidebar.viewGroups": "View Groups",
    "sidebar.createGroup": "Create Group",

    // 404 page
    "notFound.title": "Page Not Found",
    "notFound.message":
      "Sorry, but the page you are looking for does not exist or has been moved.",
    "notFound.home": "Return to Home",
    "notFound.back": "Back",

    // Groups page
    "groups.title": "Document Groups",
    "groups.subtitle": "Browse created groups and navigate to their documents",
    "groups.error": "Failed to load groups.",
    "groups.refresh": "Try again",
    "groups.empty": "No groups yet.",
    "groups.createdAt": "Created",
    "groups.totalLabel": "Total groups:",
    "groups.createButton": "Create Group",
    "groups.createDialogTitle": "Create Group",
    "groups.editDialogTitle": "Edit Group",
    "groups.createConfirm": "Create",
    "groups.updateConfirm": "Save",
    "groups.cancel": "Cancel",
    "groups.nameLabel": "Group name",
    "groups.namePlaceholder": "Enter group name",
    "groups.nameHelper": "Name is required",
    "groups.editLabel": "Edit",
    "groups.deleteLabel": "Delete",
    "groups.createSuccess": "Group created successfully.",
    "groups.updateSuccess": "Group updated successfully.",
    "groups.deleteSuccess": "Group deleted successfully.",
    "groups.createError": "Failed to create group.",
    "groups.updateError": "Failed to update group.",
    "groups.deleteError": "Failed to delete group.",
    "groups.deleteConfirmTitle": "Delete group",
    "groups.deleteConfirmDescription":
      'Are you sure you want to delete the group "{name}"?',
    "groups.deleteConfirmAccept": "Delete",
    "groups.manageMembersButton": "Manage members",
    "groups.membersTitle": 'Members of "{name}"',
    "groups.membersTitleFallback": "Group members",
    "groups.membersSubtitle":
      "Add collaborators to work together on documents.",
    "groups.membersEmpty": "This group has no members yet.",
    "groups.membersAddUser": "Add user",
    "groups.membersAddUserPlaceholder": "Select a user from the list",
    "groups.membersRoleLabel": "Role",
    "groups.membersAddButton": "Add",
    "groups.membersAddValidation": "Select a user to add to the group.",
    "groups.membersUsersError": "Failed to load users list.",
    "groups.membersActionError": "Action failed. Please try again.",
    "groups.membersRemoveTooltip": "Remove from group",
    "groups.membersAddSuccess": "User added to the group.",
    "groups.membersAddError": "Failed to add user to the group.",
    "groups.membersUpdateSuccess": "Member role updated.",
    "groups.membersUpdateError": "Failed to update member role.",
    "groups.membersRemoveSuccess": "Member removed from the group.",
    "groups.membersRemoveError": "Failed to remove member from the group.",
    "groups.membersForbidden": "Only the group author can manage members.",
    "groups.roleLabel": "Role",
    "groups.role.author": "Author",
    "groups.role.coauthor": "Co-author",
    "groups.role.reviewer": "Reviewer",

    // Documents page
    "documents.title": "Documents",
    "documents.subtitle": "Browse documents and find what you need in seconds",
    "documents.filterGroup": "Group",
    "documents.filterAll": "All groups",
    "documents.searchPlaceholder": "Search documents",
    "documents.error": "Failed to load documents.",
    "documents.refresh": "Try again",
    "documents.empty": "No documents found.",
    "documents.groupUnknown": "No group",
    "documents.noContent": "No content",
    "documents.createdAt": "Created",
    "documents.createButton": "Create Document",
    "documents.noGroupsHint": "Create a group before adding documents.",
    "documents.createDialogTitle": "Create Document",
    "documents.createConfirm": "Create",
    "documents.defaultContent": "# New document",
    "documents.cancel": "Cancel",
    "documents.nameLabel": "Document name",
    "documents.namePlaceholder": "Enter document name",
    "documents.nameHelper": "Name is required",
    "documents.groupLabel": "Group",
    "documents.contentLabel": "Document content",
    "documents.contentPlaceholder": "Write your Markdown content here...",
    "documents.contentHelper": "Content is required",
    "documents.createSuccess": "Document created successfully.",
    "documents.createError": "Failed to create document.",
    "documents.editLabel": "Edit",
    "documents.deleteLabel": "Delete",
    "documents.deleteConfirmTitle": "Delete document",
    "documents.deleteConfirmDescription":
      'Are you sure you want to delete the document "{name}"? This action cannot be undone.',
    "documents.deleteConfirmAccept": "Delete",
    "documents.deleteSuccess": "Document deleted successfully.",
    "documents.deleteError": "Failed to delete document.",
    "documents.noEditableGroups":
      "No groups are available where you can create documents.",
    "documents.deleteForbidden": "You cannot delete this document.",

    // Document editor page
    "documentEditor.subtitle":
      "Edit the document and preview changes in real time",
    "documentEditor.fallbackTitle": "Document",
    "documentEditor.nameLabel": "Document name",
    "documentEditor.namePlaceholder": "Enter a name",
    "documentEditor.nameRequired": "Name is required",
    "documentEditor.contentLabel": "Document content",
    "documentEditor.contentPlaceholder": "Write your Markdown content here...",
    "documentEditor.previewTitle": "Live preview",
    "documentEditor.previewEmpty": "Start typing to see a formatted preview.",
    "documentEditor.createdAtLabel": "Created",
    "documentEditor.saveButton": "Save changes",
    "documentEditor.savingButton": "Saving...",
    "documentEditor.saveSuccess": "Document saved successfully.",
    "documentEditor.saveError":
      "We could not save the document. Please try again.",
    "documentEditor.backToList": "Back to documents",
    "documentEditor.loadError": "Failed to load the document.",
    "documentEditor.notFound":
      "The document was not found or has been removed.",
    "documentEditor.readOnlyNotice":
      "You only have view access to this document.",
  },
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
    return (
      translations[language][
        key as keyof (typeof translations)[typeof language]
      ] || key
    );
  };

  return (
    <LanguageContext.Provider
      value={{ language, setLanguage: handleSetLanguage, t }}
    >
      {children}
    </LanguageContext.Provider>
  );
};

// eslint-disable-next-line react-refresh/only-export-components
export function useLanguage() {
  const context = useContext(LanguageContext);
  if (context === undefined) {
    throw new Error("useLanguage must be used within a LanguageProvider");
  }
  return context;
}
