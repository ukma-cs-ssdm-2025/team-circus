import CssBaseline from "@mui/material/CssBaseline";
import {
	createTheme,
	ThemeProvider as MuiThemeProvider,
} from "@mui/material/styles";
import type { ReactNode } from "react";
import { createContext, useContext, useEffect, useState } from "react";

export type Theme = "light" | "dark";

interface ThemeContextType {
	theme: Theme;
	setTheme: (theme: Theme) => void;
	toggleTheme: () => void;
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

const getStoredTheme = (): Theme => {
	try {
		const stored = localStorage.getItem("mcd_theme");
		if (stored) {
			return stored as Theme;
		}
		if (
			window.matchMedia &&
			window.matchMedia("(prefers-color-scheme: dark)").matches
		) {
			return "dark";
		}
		return "light";
	} catch (error) {
		console.warn("Не вдалося визначити тему із localStorage", error);
		return "light";
	}
};

const setStoredTheme = (theme: Theme): void => {
	try {
		localStorage.setItem("mcd_theme", theme);
	} catch (error) {
		console.warn("Не вдалося зберегти тему до localStorage", error);
	}
};

const createAppTheme = (mode: Theme) =>
	createTheme({
		palette: {
			mode,
			primary: {
				main: "#667eea",
				light: "#9bb5ff",
				dark: "#4c63d2",
				contrastText: "#ffffff",
			},
			secondary: {
				main: "#764ba2",
				light: "#a77bc4",
				dark: "#5a3d7a",
				contrastText: "#ffffff",
			},
			background: {
				default: mode === "light" ? "#f5f7fa" : "#121212",
				paper: mode === "light" ? "#ffffff" : "#1e1e1e",
			},
			text: {
				primary: mode === "light" ? "#2c3e50" : "#ffffff",
				secondary: mode === "light" ? "#7f8c8d" : "#b0b0b0",
			},
		},
		typography: {
			fontFamily: '"Inter", "Roboto", "Helvetica", "Arial", sans-serif',
		},
		shape: {
			borderRadius: 12,
		},
		components: {
			MuiButton: {
				styleOverrides: {
					root: {
						textTransform: "none",
						fontWeight: 600,
						borderRadius: 8,
						padding: "8px 16px",
					},
					contained: {
						boxShadow: "0 4px 15px rgba(102, 126, 234, 0.3)",
						"&:hover": {
							boxShadow: "0 6px 20px rgba(102, 126, 234, 0.4)",
							transform: "translateY(-2px)",
						},
					},
				},
			},
			MuiCard: {
				styleOverrides: {
					root: {
						borderRadius: 16,
						boxShadow: "0 10px 30px rgba(0, 0, 0, 0.1)",
						backdropFilter: "blur(10px)",
					},
				},
			},
			MuiAppBar: {
				styleOverrides: {
					root: {
						backgroundColor:
							mode === "light"
								? "rgba(255, 255, 255, 0.95)"
								: "rgba(30, 30, 30, 0.95)",
						backdropFilter: "blur(10px)",
						boxShadow: "0 2px 20px rgba(0, 0, 0, 0.1)",
					},
				},
			},
			MuiPaper: {
				styleOverrides: {
					root: {
						backgroundImage: "none",
					},
				},
			},
		},
	});

interface ThemeProviderProps {
	children: ReactNode;
}

export const ThemeProvider = ({ children }: ThemeProviderProps) => {
	const [theme, setTheme] = useState<Theme>(getStoredTheme);

	// Зберігаємо тему в localStorage при зміні
	useEffect(() => {
		setStoredTheme(theme);
	}, [theme]);

	// Sync data attributes for CSS modules (light/dark)
	useEffect(() => {
		const root = document.documentElement;
		root.setAttribute("data-theme", theme);
		root.setAttribute("data-mui-color-scheme", theme);
	}, [theme]);

	// Слухач змін теми системи
	useEffect(() => {
		const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
		const handleChange = (e: MediaQueryListEvent) => {
			// Змінюємо тему тільки якщо користувач не встановив власну
			const stored = localStorage.getItem("mcd_theme");
			if (!stored) {
				setTheme(e.matches ? "dark" : "light");
			}
		};

		if (mediaQuery.addEventListener) {
			mediaQuery.addEventListener("change", handleChange);
		} else {
			// Fallback для старих браузерів
			mediaQuery.addListener(handleChange);
		}

		return () => {
			if (mediaQuery.removeEventListener) {
				mediaQuery.removeEventListener("change", handleChange);
			} else {
				mediaQuery.removeListener(handleChange);
			}
		};
	}, []);

	const handleSetTheme = (newTheme: Theme) => {
		setTheme(newTheme);
	};

	const toggleTheme = () => {
		const newTheme = theme === "light" ? "dark" : "light";
		setTheme(newTheme);
		// Очищаємо localStorage при ручній зміні теми
		localStorage.removeItem("mcd_theme");
		localStorage.setItem("mcd_theme", newTheme);
	};

	const muiTheme = createAppTheme(theme);

	return (
		<ThemeContext.Provider
			value={{ theme, setTheme: handleSetTheme, toggleTheme }}
		>
			<MuiThemeProvider theme={muiTheme}>
				<CssBaseline />
				{children}
			</MuiThemeProvider>
		</ThemeContext.Provider>
	);
};

// eslint-disable-next-line react-refresh/only-export-components
export function useTheme() {
	const context = useContext(ThemeContext);
	if (context === undefined) {
		throw new Error("useTheme must be used within a ThemeProvider");
	}
	return context;
}
