import { Link as RouterLink, useNavigate } from "react-router-dom";
import {
	AppBar,
	Toolbar,
	Typography,
	IconButton,
	Box,
	useTheme,
} from "@mui/material";
import {
	Settings as SettingsIcon,
	LightMode as LightModeIcon,
	DarkMode as DarkModeIcon,
	Menu as MenuIcon,
	Logout as LogoutIcon,
} from "@mui/icons-material";
import type { BaseComponentProps } from "../../types";
import { useLanguage } from "../../contexts/LanguageContext";
import { useTheme as useAppTheme } from "../../contexts/ThemeContext";
import { useAuth } from "../../contexts/AuthContextBase";
import { ROUTES } from "../../constants";

interface HeaderProps extends BaseComponentProps {
	onToggleSidebar?: () => void;
}

const Header = ({ className = "", onToggleSidebar }: HeaderProps) => {
	const navigate = useNavigate();
	const theme = useTheme();
	const { t } = useLanguage();
	const { theme: appTheme, toggleTheme } = useAppTheme();
	const { logout, user } = useAuth();

	const handleAccountSettings = () => {
		navigate(ROUTES.SETTINGS);
	};

	const handleLogout = async () => {
		try {
			await logout();
			navigate(ROUTES.LOGIN);
		} catch (error) {
			console.error("Logout failed:", error);
		}
	};

	const commonIconStyles = {
		backgroundColor: theme.palette.mode === "light" ? "#f8fafc" : "#374151",
		border: `1px solid ${theme.palette.mode === "light" ? "#cbd5e1" : "#4b5563"}`,
		color: theme.palette.mode === "light" ? "#475569" : "#d1d5db",
		"&:hover": {
			backgroundColor: theme.palette.mode === "light" ? "#e2e8f0" : "#4b5563",
			borderColor: theme.palette.mode === "light" ? "#94a3b8" : "#6b7280",
			color: theme.palette.mode === "light" ? "#334155" : "#f3f4f6",
			transform: "translateY(-2px)",
			boxShadow: "0 2px 6px rgba(0, 0, 0, 0.15)",
		},
		transition: "all 0.3s ease",
	} as const;

	return (
		<AppBar
			className={className}
			position="sticky"
			elevation={0}
			sx={{
				backgroundColor:
					theme.palette.mode === "light"
						? "rgba(255, 255, 255, 0.95)"
						: "rgba(30, 30, 30, 0.95)",
				backdropFilter: "blur(10px)",
				boxShadow: "0 2px 20px rgba(0, 0, 0, 0.1)",
			}}
		>
			<Toolbar sx={{ display: "flex", justifyContent: "space-between" }}>
				<Box sx={{ display: "flex", alignItems: "center", gap: 1.5 }}>
					{onToggleSidebar && (
						<IconButton
							onClick={onToggleSidebar}
							title={t("sidebar.navigation")}
							aria-label={t("sidebar.navigation")}
							sx={commonIconStyles}
						>
							<MenuIcon />
						</IconButton>
					)}

					<Typography
						variant="h5"
						component={RouterLink}
						to={ROUTES.HOME}
						sx={{
							fontWeight: 700,
							textDecoration: "none",
							color: theme.palette.primary.main,
							cursor: "pointer",
						}}
					>
						MCD
					</Typography>
				</Box>

				<Box sx={{ display: "flex", gap: 1, alignItems: "center" }}>
					{user && (
						<Typography
							variant="body2"
							sx={{
								color: theme.palette.text.secondary,
								mr: 1,
								display: { xs: "none", sm: "block" },
							}}
						>
							{user.login}
						</Typography>
					)}

					<IconButton
						onClick={toggleTheme}
						title={t("header.toggleTheme")}
						sx={commonIconStyles}
					>
						{appTheme === "light" ? <DarkModeIcon /> : <LightModeIcon />}
					</IconButton>

					<IconButton
						onClick={handleAccountSettings}
						title={t("header.settings")}
						sx={commonIconStyles}
					>
						<SettingsIcon />
					</IconButton>

					<IconButton
						onClick={handleLogout}
						title="Logout"
						sx={commonIconStyles}
					>
						<LogoutIcon />
					</IconButton>
				</Box>
			</Toolbar>
		</AppBar>
	);
};

export default Header;
