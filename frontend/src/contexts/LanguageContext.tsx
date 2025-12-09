import type { ReactNode } from "react";
import { createContext, useContext, useEffect, useState } from "react";
import { en } from "../locales/en";
import { uk } from "../locales/uk";

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

const translations = { uk, en };

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
