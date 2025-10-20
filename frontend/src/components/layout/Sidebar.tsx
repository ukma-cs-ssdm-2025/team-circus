import {
  Box,
  List,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Typography,
  useTheme,
  Drawer,
  IconButton,
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
  open: boolean;
  onClose: () => void;
}

const Sidebar = ({ className = "", open, onClose }: SidebarProps) => {
  const theme = useTheme();
  const navigate = useNavigate();
  const location = useLocation();
  const { t } = useLanguage();

  const isActive = (path: string) => location.pathname === path;

  const handleNavigate = (path: string) => {
    navigate(path);
    onClose();
  };

  return (
    <Drawer
      anchor="left"
      open={open}
      onClose={onClose}
      ModalProps={{ keepMounted: true }}
      PaperProps={{
        sx: {
          width: { xs: "78vw", sm: 320 },
          maxWidth: 360,
          backgroundColor:
            theme.palette.mode === "light"
              ? "rgba(255, 255, 255, 0.95)"
              : "rgba(18, 18, 18, 0.95)",
          backdropFilter: "blur(12px)",
          display: "flex",
        },
      }}
    >
      <Box
        className={className}
        component="nav"
        sx={{
          width: "100%",
          height: "100%",
          display: "flex",
          flexDirection: "column",
          px: { xs: 3, md: 6 },
          py: { xs: 4, md: 6 },
          gap: 2,
        }}
      >
        <Box
          sx={{
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
          }}
        >
          <Typography variant="h5" sx={{ fontWeight: 700 }}>
            {t("sidebar.navigation")}
          </Typography>
          <IconButton
            onClick={onClose}
            title="Close navigation"
            aria-label="Close navigation"
            sx={{
              backgroundColor:
                theme.palette.mode === "light"
                  ? "#f8fafc"
                  : "#374151",
              border: `1px solid ${theme.palette.mode === "light" ? "#cbd5e1" : "#4b5563"}`,
              color:
                theme.palette.mode === "light"
                  ? "#475569"
                  : "#d1d5db",
              "&:hover": {
                backgroundColor:
                  theme.palette.mode === "light"
                    ? "#e2e8f0"
                    : "#4b5563",
                borderColor:
                  theme.palette.mode === "light"
                    ? "#94a3b8"
                    : "#6b7280",
                color:
                  theme.palette.mode === "light"
                    ? "#334155"
                    : "#f3f4f6",
                transform: "translateY(-2px)",
                boxShadow: "0 2px 6px rgba(0, 0, 0, 0.15)",
              },
              transition: "all 0.3s ease",
            }}
          >
            <CloseIcon />
          </IconButton>
        </Box>

        <List sx={{ flexGrow: 1 }}>
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
