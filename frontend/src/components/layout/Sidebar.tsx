import {
	Box,
	Drawer,
	List,
	ListItemButton,
	ListItemIcon,
	ListItemText,
	Typography,
	IconButton,
	useTheme,
} from "@mui/material";
import {
	Group as GroupIcon,
	Home as HomeIcon,
	Description as DescriptionIcon,
	Close as CloseIcon,
} from "@mui/icons-material";
import { useNavigate, useLocation } from "react-router-dom";
import { useLanguage } from "../../contexts/LanguageContext";
import type { BaseComponentProps } from "../../types";
import { ROUTES } from "../../constants";

interface SidebarProps extends BaseComponentProps {
	open?: boolean;
	onClose?: () => void;
}

const Sidebar = ({ className = "", open = false, onClose }: SidebarProps) => {
	const theme = useTheme();
	const navigate = useNavigate();
	const location = useLocation();
	const { t } = useLanguage();

	const handleClose = () => {
		if (onClose) {
			onClose();
		}
	};

	const handleNavigate = (path: string) => {
		navigate(path);
		handleClose();
	};

	const isActive = (path: string) => location.pathname === path;

	return (
		<Drawer
			anchor="left"
			open={open}
			onClose={handleClose}
			ModalProps={{ keepMounted: true }}
			sx={{
				"& .MuiDrawer-paper": {
					width: 280,
					backgroundColor:
						theme.palette.mode === "light"
							? "rgba(255, 255, 255, 0.95)"
							: "rgba(30, 30, 30, 0.95)",
					borderRight: `1px solid ${theme.palette.mode === "dark" ? "rgba(255, 255, 255, 0.2)" : theme.palette.divider}`,
					backdropFilter: "blur(10px)",
				},
			}}
		>
			<Box
				component="nav"
				className={className}
				sx={{
					width: "100%",
					height: "100%",
					p: 2,
				}}
			>
				<Box
					sx={{
						display: "flex",
						alignItems: "center",
						justifyContent: "space-between",
						px: 1,
						pb: 1,
					}}
				>
					<Typography variant="h6" sx={{ fontWeight: 700 }}>
						{t("sidebar.navigation")}
					</Typography>
					<IconButton
						onClick={handleClose}
						aria-label={t("sidebar.close")}
						title={t("sidebar.close")}
						size="small"
						sx={{
							color: theme.palette.text.secondary,
							"&:hover": { color: theme.palette.text.primary },
						}}
					>
						<CloseIcon />
					</IconButton>
				</Box>
				<List>
					<ListItemButton
						selected={isActive(ROUTES.HOME)}
						onClick={() => handleNavigate(ROUTES.HOME)}
					>
						<ListItemIcon>
							<HomeIcon />
						</ListItemIcon>
						<ListItemText primary={t("sidebar.home")} />
					</ListItemButton>
					<ListItemButton
						selected={isActive(ROUTES.DOCUMENTS)}
						onClick={() => handleNavigate(ROUTES.DOCUMENTS)}
					>
						<ListItemIcon>
							<DescriptionIcon />
						</ListItemIcon>
						<ListItemText primary={t("sidebar.documents")} />
					</ListItemButton>
					<ListItemButton
						selected={isActive(ROUTES.GROUPS)}
						onClick={() => handleNavigate(ROUTES.GROUPS)}
					>
						<ListItemIcon>
							<GroupIcon />
						</ListItemIcon>
						<ListItemText primary={t("sidebar.groups")} />
					</ListItemButton>
				</List>
			</Box>
		</Drawer>
	);
};

export default Sidebar;
