import { useEffect } from "react";
import { THEME } from "../constants";
import type { Theme } from "../types";
import { useLocalStorage } from "./useLocalStorage";

export function useTheme() {
	const [theme, setTheme] = useLocalStorage<Theme>(
		"mcd_theme",
		THEME.LIGHT as Theme,
	);

	const toggleTheme = () => {
		setTheme((prev) => (prev === THEME.LIGHT ? THEME.DARK : THEME.LIGHT));
	};

	useEffect(() => {
		document.documentElement.setAttribute("data-theme", theme);
	}, [theme]);

	return { theme, toggleTheme };
}
