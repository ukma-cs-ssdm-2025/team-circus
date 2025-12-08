import {
	FormControl,
	InputLabel,
	MenuItem,
	Select,
	Stack,
} from "@mui/material";
import { CenteredContent, PageCard, PageHeader } from "../components/common";
import { SettingsSection } from "../components/forms";
import type { Language } from "../contexts/LanguageContext";
import { useLanguage } from "../contexts/LanguageContext";
import { type Theme, useTheme } from "../contexts/ThemeContext";
import type { BaseComponentProps } from "../types";

type SettingsProps = BaseComponentProps;

const Settings = ({ className = "" }: SettingsProps) => {
	const { theme, setTheme } = useTheme();
	const { language, setLanguage, t } = useLanguage();

	return (
		<CenteredContent className={className}>
			<PageCard>
				<PageHeader
					title={t("settings.title")}
					subtitle={t("settings.subtitle")}
				/>

				<Stack spacing={4}>
					<SettingsSection title={t("settings.general")}>
						<FormControl fullWidth>
							<InputLabel>{t("settings.theme")}</InputLabel>
							<Select
								value={theme}
								label={t("settings.theme")}
								onChange={(event) => setTheme(event.target.value as Theme)}
							>
								<MenuItem value="light">{t("settings.theme.light")}</MenuItem>
								<MenuItem value="dark">{t("settings.theme.dark")}</MenuItem>
							</Select>
						</FormControl>

						<FormControl fullWidth>
							<InputLabel>{t("settings.language")}</InputLabel>
							<Select
								value={language}
								label={t("settings.language")}
								onChange={(event) =>
									setLanguage(event.target.value as Language)
								}
							>
								<MenuItem value="uk">Українська</MenuItem>
								<MenuItem value="en">English</MenuItem>
							</Select>
						</FormControl>
					</SettingsSection>
				</Stack>
			</PageCard>
		</CenteredContent>
	);
};

export default Settings;
